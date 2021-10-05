package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/kthomas/go.uuid"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"github.com/unibrightio/proxy-api/common"
	"github.com/unibrightio/proxy-api/cron"
	"github.com/unibrightio/proxy-api/dbutil"
	"github.com/unibrightio/proxy-api/helpers"
	"github.com/unibrightio/proxy-api/httpd/handler"
	"github.com/unibrightio/proxy-api/logger"
	"github.com/unibrightio/proxy-api/messaging"
	"github.com/unibrightio/proxy-api/types"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/unibrightio/proxy-api/httpd/docs"
)

// @title Baseledger Proxy API documentation
// @version 1.0.0
// @host localhost:8081
// @securityDefinitions.basic BasicAuth
func main() {
	setupViper()
	logger.SetupLogger()
	setupDb()
	cron.StartCron()
	subscribeToWorkgroupMessages()

	r := gin.Default()
	r.Use(helpers.CORSMiddleware())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/trustmeshes", basicAuth, handler.GetTrustmeshesHandler())
	r.GET("/trustmeshes/:id", basicAuth, handler.GetTrustmeshHandler())
	r.POST("/suggestion", basicAuth, handler.CreateInitialSuggestionRequestHandler())
	r.POST("/feedback", basicAuth, handler.CreateSynchronizationFeedbackHandler())
	r.GET("/sunburst/:txId", basicAuth, handler.GetSunburstHandler())
	r.POST("send_offchain_message", basicAuth, handler.SendOffchainMessageHandler())
	r.GET("/organization", basicAuth, handler.GetOrganizationsHandler())
	r.POST("/organization", basicAuth, handler.CreateOrganizationHandler())
	r.DELETE("/organization/:id", basicAuth, handler.DeleteOrganizationHandler())
	r.GET("/workgroup", basicAuth, handler.GetWorkgroupsHandler())
	r.POST("/workgroup", basicAuth, handler.CreateWorkgroupHandler())
	r.DELETE("/workgroup/:id", basicAuth, handler.DeleteWorkgroupHandler())
	r.GET("/participation", basicAuth, handler.GetWorkgroupMemberHandler())
	r.POST("/participation", basicAuth, handler.CreateWorkgroupMemberHandler())
	r.DELETE("/participation/:id", basicAuth, handler.DeleteWorkgroupMemberHandler())
	// TODO: BAS-29 r.POST("/workgroup/invite", handler.InviteToWorkgroupHandler())
	// full details of workgroup, including organization
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func basicAuth(c *gin.Context) {
	basicAuthUser, _ := viper.Get("API_UB_USER").(string)
	basicAuthPwd, _ := viper.Get("API_UB_PWD").(string)
	// Get the Basic Authentication credentials
	user, password, hasAuth := c.Request.BasicAuth()
	if hasAuth && user == basicAuthUser && password == basicAuthPwd {
		logger.Info("Basic auth successful")
	} else {
		logger.Error("Basic auth failed")
		c.Abort()
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		c.JSON(http.StatusForbidden, map[string]interface{}{"error": "auth failed"})
		return
	}
}

// discuss if we should use config struct or this is enough
func setupViper() {
	viper.AddConfigPath("../")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv() // Overwrite config with env variables if exist, important for debugging session

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(fmt.Printf("viper read config error %v\n", err))
	}
}

// migrate should be separate package, and we should have .sh script for running, see provide services
// leaving this for first version but we should separate definetely
func setupDb() {
	dbutil.InitDbIfNotExists()
	dbutil.PerformMigrations()
	// TODO: BAS-29 Add own org id to database with some dummy name
	dbutil.InitConnection()
}

func subscribeToWorkgroupMessages() {
	natsServerUrl, _ := viper.Get("NATS_URL").(string)
	natsToken := "testToken1" // TODO: Read from configuration
	logger.Infof("subscribeToWorkgroupMessages natsServerUrl %v", natsServerUrl)
	messagingClient := &messaging.NatsMessagingClient{}
	messagingClient.Subscribe(natsServerUrl, natsToken, "baseledger", receiveOffchainProcessMessage)
}

func receiveOffchainProcessMessage(sender string, natsMsg *nats.Msg) {
	// TODO: should we move this parsing to nats client and just get struct in this callback?
	var natsMessage types.NatsMessage
	err := json.Unmarshal(natsMsg.Data, &natsMessage)
	if err != nil {
		logger.Errorf("Error parsing nats message %v\n", err)
		return
	}

	logger.Infof("message received %v\n", natsMessage)

	natsMessage.ProcessMessage.Id = uuid.Nil // set to nil so that it can be created in the DB
	if !natsMessage.ProcessMessage.Create() {
		logger.Errorf("error when creating new offchain msg entry")
		return
	}

	entryType := common.SuggestionReceivedTrustmeshEntryType
	if natsMessage.ProcessMessage.EntryType == common.FeedbackSentTrustmeshEntryType {
		entryType = common.FeedbackReceivedTrustmeshEntryType
	}

	trustmeshEntry := &types.TrustmeshEntry{
		TendermintTransactionId:              natsMessage.ProcessMessage.BaseledgerTransactionIdOfStoredProof,
		OffchainProcessMessageId:             natsMessage.ProcessMessage.Id,
		SenderOrgId:                          natsMessage.ProcessMessage.SenderId,
		ReceiverOrgId:                        natsMessage.ProcessMessage.ReceiverId,
		WorkgroupId:                          uuid.FromStringOrNil(natsMessage.ProcessMessage.Topic),
		WorkstepType:                         natsMessage.ProcessMessage.WorkstepType,
		BaseledgerTransactionType:            natsMessage.ProcessMessage.BaseledgerTransactionType,
		BaseledgerTransactionId:              natsMessage.ProcessMessage.BaseledgerTransactionIdOfStoredProof,
		ReferencedBaseledgerTransactionId:    natsMessage.ProcessMessage.ReferencedBaseledgerTransactionId,
		BusinessObjectType:                   natsMessage.ProcessMessage.BusinessObjectType,
		BaseledgerBusinessObjectId:           natsMessage.ProcessMessage.BaseledgerBusinessObjectId,
		ReferencedBaseledgerBusinessObjectId: natsMessage.ProcessMessage.ReferencedBaseledgerBusinessObjectId,
		ReferencedProcessMessageId:           natsMessage.ProcessMessage.ReferencedOffchainProcessMessageId,
		SorBusinessObjectId:                  natsMessage.ProcessMessage.SorBusinessObjectId,
		TransactionHash:                      natsMessage.TxHash,
		EntryType:                            entryType,
	}

	if !trustmeshEntry.Create() {
		logger.Errorf("error when creating new trustmesh entry")
	}

}
