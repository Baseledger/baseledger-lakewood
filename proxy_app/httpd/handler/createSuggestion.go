package handler

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	uuid "github.com/kthomas/go.uuid"
	"github.com/spf13/viper"
	"github.com/unibrightio/proxy-api/common"
	"github.com/unibrightio/proxy-api/logger"
	"github.com/unibrightio/proxy-api/proxyutil"
	"github.com/unibrightio/proxy-api/restutil"
	"github.com/unibrightio/proxy-api/synctree"
	"github.com/unibrightio/proxy-api/types"
)

type createInitialSuggestionRequest struct {
	WorkgroupId                          string   `json:"workgroup_id"`
	Recipient                            string   `json:"recipient"`
	WorkstepType                         string   `json:"workstep_type"`
	BusinessObjectType                   string   `json:"business_object_type"`
	BaseledgerBusinessObjectId           string   `json:"baseledger_business_object_id"`
	BusinessObjectJson                   string   `json:"business_object_json"`
	ReferencedBaseledgerBusinessObjectId string   `json:"referenced_baseledger_business_object_id"`
	ReferencedBaseledgerTransactionId    string   `json:"referenced_baseledger_transaction_id"`
	SorBusinessObjectId                  string   `json:"sor_message_id"`
	KnowledgeLimiters                    []string `json:"knowledge_limiters"`
}

func CreateInitialSuggestionRequestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !restutil.HasEnoughBalance() {
			restutil.RenderError("not enough token balance", 400, c)
			return
		}

		buf, err := c.GetRawData()
		if err != nil {
			restutil.RenderError(err.Error(), 400, c)
			return
		}

		req := &createInitialSuggestionRequest{}
		err = json.Unmarshal(buf, &req)
		if err != nil {
			restutil.RenderError(err.Error(), 422, c)
			return
		}

		syncReq := newSynchronizationRequest(*req)

		syncTree := synctree.CreateFromBusinessObjectJson(syncReq.BusinessObjectJson, syncReq.KnowledgeLimiters)
		logger.Infof("Sync tree %v", syncTree)

		syncTreeJson, err := json.Marshal(syncTree)
		if err != nil {
			logger.Errorf("error marshaling sync tree", err.Error())
			restutil.RenderError("error marshaling sync tree", 500, c)
			return
		}
		logger.Infof("Sync tree json %v", string(syncTreeJson))

		transactionId := uuid.NewV4()
		offchainMsg := createSuggestOffchainMessage(*req, transactionId, string(syncTreeJson), syncTree.RootProof)

		if !offchainMsg.Create() {
			logger.Errorf("error when creating new offchain msg entry")
			restutil.RenderError("error when creating new offchain msg entry", 500, c)
			return
		}

		payload := proxyutil.CreateBaseledgerTransactionPayload(syncReq, &offchainMsg)

		signAndBroadcastPayload := restutil.SignAndBroadcastPayload{
			TransactionId: transactionId.String(),
			Payload:       payload,
			OpCode:        uint32(getRandomSuggestionOpCode()),
		}

		transactionHash := restutil.SignAndBroadcast(signAndBroadcastPayload)

		if transactionHash == nil {
			restutil.RenderError("sign and broadcast transaction error", 500, c)
			return
		}

		trustmeshEntry := createSuggestionSentTrustmeshEntry(*req, transactionId, offchainMsg, *transactionHash)

		if !trustmeshEntry.Create() {
			logger.Errorf("error when creating new trustmesh entry")
			restutil.RenderError("error when creating new trustmesh entry", 500, c)
			return
		}

		restutil.Render(transactionHash, 200, c)
	}
}

func getRandomSuggestionOpCode() int {
	rand.Seed(time.Now().UnixNano())
	min := 7
	max := 8

	return rand.Intn(max-min+1) + min
}

func newSynchronizationRequest(req createInitialSuggestionRequest) *types.SynchronizationRequest {
	return &types.SynchronizationRequest{
		WorkgroupId:                          uuid.FromStringOrNil(req.WorkgroupId),
		Recipient:                            req.Recipient,
		WorkstepType:                         req.WorkstepType,
		BusinessObjectType:                   req.BusinessObjectType,
		BaseledgerBusinessObjectId:           req.BaseledgerBusinessObjectId,
		BusinessObjectJson:                   req.BusinessObjectJson,
		ReferencedBaseledgerBusinessObjectId: req.ReferencedBaseledgerBusinessObjectId,
		KnowledgeLimiters:                    req.KnowledgeLimiters,
	}
}

func createSuggestionSentTrustmeshEntry(req createInitialSuggestionRequest, transactionId uuid.UUID, offchainMsg types.OffchainProcessMessage, txHash string) *types.TrustmeshEntry {
	return &types.TrustmeshEntry{
		TendermintTransactionId:              transactionId,
		OffchainProcessMessageId:             offchainMsg.Id,
		SenderOrgId:                          uuid.FromStringOrNil(viper.Get("ORGANIZATION_ID").(string)),
		ReceiverOrgId:                        uuid.FromStringOrNil(req.Recipient),
		WorkgroupId:                          uuid.FromStringOrNil(req.WorkgroupId),
		WorkstepType:                         offchainMsg.WorkstepType,
		BaseledgerTransactionType:            "Suggest",
		BaseledgerTransactionId:              transactionId,
		ReferencedBaseledgerTransactionId:    uuid.FromStringOrNil(req.ReferencedBaseledgerTransactionId),
		BusinessObjectType:                   req.BusinessObjectType,
		BaseledgerBusinessObjectId:           offchainMsg.BaseledgerBusinessObjectId,
		ReferencedBaseledgerBusinessObjectId: offchainMsg.ReferencedBaseledgerBusinessObjectId,
		ReferencedProcessMessageId:           offchainMsg.ReferencedOffchainProcessMessageId,
		TransactionHash:                      txHash,
		EntryType:                            common.SuggestionSentTrustmeshEntryType,
		SorBusinessObjectId:                  req.SorBusinessObjectId,
	}
}

func createSuggestOffchainMessage(req createInitialSuggestionRequest, transactionId uuid.UUID, syncTreeJson string, rootProof string) types.OffchainProcessMessage {
	offchainMessage := types.OffchainProcessMessage{
		SenderId:                             uuid.FromStringOrNil(viper.Get("ORGANIZATION_ID").(string)),
		ReceiverId:                           uuid.FromStringOrNil(req.Recipient),
		Topic:                                req.WorkgroupId,
		WorkstepType:                         req.WorkstepType,
		ReferencedOffchainProcessMessageId:   uuid.FromStringOrNil(""),
		BaseledgerSyncTreeJson:               syncTreeJson,
		BusinessObjectProof:                  rootProof,
		BaseledgerBusinessObjectId:           req.BaseledgerBusinessObjectId,
		ReferencedBaseledgerBusinessObjectId: req.ReferencedBaseledgerBusinessObjectId,
		StatusTextMessage:                    req.WorkstepType + " suggested",
		BaseledgerTransactionIdOfStoredProof: transactionId,
		TendermintTransactionIdOfStoredProof: transactionId,
		BusinessObjectType:                   req.BusinessObjectType,
		BaseledgerTransactionType:            "Suggest",
		ReferencedBaseledgerTransactionId:    uuid.FromStringOrNil(req.ReferencedBaseledgerTransactionId),
		EntryType:                            common.SuggestionSentTrustmeshEntryType,
		SorBusinessObjectId:                  req.SorBusinessObjectId,
	}

	return offchainMessage
}
