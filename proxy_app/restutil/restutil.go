package restutil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/kthomas/go.uuid"
	"github.com/spf13/viper"
	"github.com/unibrightio/proxy-api/common"
	"github.com/unibrightio/proxy-api/logger"
	"github.com/unibrightio/proxy-api/proxyutil"
	"github.com/unibrightio/proxy-api/types"
)

const defaultResponseContentType = "application/json; charset=UTF-8"

type SignAndBroadcastPayload struct {
	TransactionId string `json:"transaction_id"`
	Payload       string `json:"payload"`
	OpCode        uint32 `json:"op_code"`
}

func SignAndBroadcast(payload SignAndBroadcastPayload) *string {
	jsonValue, err := json.Marshal(payload)

	if err != nil {
		logger.Error("Error marshaling sign and broadcast json")
		return nil
	}

	resp, err := http.Post("http://"+viper.Get("BLOCKCHAIN_APP_URL").(string)+"/signAndBroadcast", "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		logger.Errorf("error while sending feedback request %v\n", err.Error())
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("error while reading sign and broadcast transaction response %v\n", err.Error())
		return nil
	}

	txHash := string(body)
	return &txHash
}

func HasEnoughBalance() bool {
	resp, err := http.Get("http://" + viper.Get("BLOCKCHAIN_APP_URL").(string) + "/balanceCheck")

	if err != nil {
		logger.Errorf("error while sending feedback request %v\n", err.Error())
		return false
	}

	return resp.StatusCode == 200
}

func SendRejectFeedback(offchainProcessMessage *types.OffchainProcessMessage, workgroupId string) {
	var feedback = &types.SynchronizationFeedback{
		WorkgroupId:        uuid.FromStringOrNil(workgroupId),
		BusinessObjectType: offchainProcessMessage.BusinessObjectType,
		Recipient:          offchainProcessMessage.ReceiverId.String(),
		Approved:           false,
		BaseledgerBusinessObjectIdOfApprovedObject: offchainProcessMessage.BaseledgerBusinessObjectId,
		HashOfObjectToApprove:                      offchainProcessMessage.BusinessObjectProof,
		OriginalBaseledgerTransactionId:            offchainProcessMessage.BaseledgerTransactionIdOfStoredProof.String(),
		OriginalOffchainProcessMessageId:           offchainProcessMessage.Id.String(),
		FeedbackMessage:                            "Rejected because Hashes do not match",
		BaseledgerProvenBusinessObjectJson:         offchainProcessMessage.BaseledgerSyncTreeJson,
	}

	transactionId := uuid.NewV4()

	offchainMsg := createFeedbackOffchainMessage(*feedback, transactionId, "Reject")

	if !offchainMsg.Create() {
		logger.Errorf("error when creating new offchain msg entry")
		return
	}

	payload := proxyutil.CreateBaseledgerTransactionFeedbackPayload(feedback, &offchainMsg)

	signAndBroadcastPayload := SignAndBroadcastPayload{
		TransactionId: transactionId.String(),
		Payload:       payload,
	}

	transactionHash := SignAndBroadcast(signAndBroadcastPayload)

	if transactionHash == nil {
		return
	}

	trustmeshEntry := createFeedbackSentTrustmeshEntry(*feedback, transactionId, offchainMsg, "Reject", *transactionHash)

	if !trustmeshEntry.Create() {
		logger.Errorf("error when creating new trustmesh entry")
		return
	}
}

// Render an object and status using the given gin context
func Render(obj interface{}, status int, c *gin.Context) {
	c.Header("content-type", defaultResponseContentType)
	c.Writer.WriteHeader(status)
	if &obj != nil && status != http.StatusNoContent {
		encoder := json.NewEncoder(c.Writer)
		encoder.SetIndent("", "    ")
		if err := encoder.Encode(obj); err != nil {
			panic(err)
		}
	} else {
		c.Header("content-length", "0")
	}
}

func RenderError(message string, status int, c *gin.Context) {
	err := map[string]*string{}
	err["message"] = &message
	Render(err, status, c)
}

func createFeedbackOffchainMessage(req types.SynchronizationFeedback, transactionId uuid.UUID, baseledgerTransactionType string) types.OffchainProcessMessage {
	offchainMessage := types.OffchainProcessMessage{
		SenderId:                             uuid.FromStringOrNil(viper.Get("ORGANIZATION_ID").(string)),
		ReceiverId:                           uuid.FromStringOrNil(req.Recipient),
		Topic:                                req.WorkgroupId.String(),
		WorkstepType:                         "Feedback",
		ReferencedOffchainProcessMessageId:   uuid.FromStringOrNil(req.OriginalOffchainProcessMessageId),
		BaseledgerSyncTreeJson:               req.BaseledgerProvenBusinessObjectJson,
		BusinessObjectProof:                  req.HashOfObjectToApprove,
		BaseledgerBusinessObjectId:           "",
		ReferencedBaseledgerBusinessObjectId: req.BaseledgerBusinessObjectIdOfApprovedObject,
		StatusTextMessage:                    req.FeedbackMessage,
		BaseledgerTransactionIdOfStoredProof: transactionId,
		TendermintTransactionIdOfStoredProof: transactionId,
		BusinessObjectType:                   req.BusinessObjectType,
		BaseledgerTransactionType:            baseledgerTransactionType,
		ReferencedBaseledgerTransactionId:    uuid.FromStringOrNil(req.OriginalBaseledgerTransactionId),
		EntryType:                            common.FeedbackSentTrustmeshEntryType,
	}

	return offchainMessage
}

func createFeedbackSentTrustmeshEntry(req types.SynchronizationFeedback, transactionId uuid.UUID, offchainMsg types.OffchainProcessMessage, feedbackMsg string, txHash string) *types.TrustmeshEntry {
	trustmeshEntry := &types.TrustmeshEntry{
		TendermintTransactionId:              transactionId,
		OffchainProcessMessageId:             offchainMsg.Id,
		SenderOrgId:                          uuid.FromStringOrNil(viper.Get("ORGANIZATION_ID").(string)),
		ReceiverOrgId:                        uuid.FromStringOrNil(req.Recipient),
		WorkgroupId:                          req.WorkgroupId,
		WorkstepType:                         offchainMsg.WorkstepType,
		BaseledgerTransactionType:            feedbackMsg,
		BaseledgerTransactionId:              transactionId,
		ReferencedBaseledgerTransactionId:    uuid.FromStringOrNil(req.OriginalBaseledgerTransactionId),
		BusinessObjectType:                   req.BusinessObjectType,
		BaseledgerBusinessObjectId:           offchainMsg.BaseledgerBusinessObjectId,
		ReferencedBaseledgerBusinessObjectId: offchainMsg.ReferencedBaseledgerBusinessObjectId,
		ReferencedProcessMessageId:           offchainMsg.ReferencedOffchainProcessMessageId,
		TransactionHash:                      txHash,
		EntryType:                            common.FeedbackSentTrustmeshEntryType,
	}

	return trustmeshEntry
}
