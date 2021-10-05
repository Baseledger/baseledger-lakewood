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
	"github.com/unibrightio/proxy-api/types"
)

type createSynchronizationFeedbackRequest struct {
	WorkgroupId                                string `json:"workgroup_id"`
	BusinessObjectType                         string `json:"business_object_type"`
	Recipient                                  string `json:"recipient"`
	Approved                                   bool   `json:"approved"`
	BaseledgerBusinessObjectIdOfApprovedObject string `json:"baseledger_business_object_id_of_approved_object"`
	OriginalBaseledgerTransactionId            string `json:"original_baseledger_transaction_id"`
	OriginalOffchainProcessMessageId           string `json:"original_offchain_process_message_id"`
	FeedbackMessage                            string `json:"feedback_message"`
}

// @Security BasicAuth
// Create Feedback ... Create Feedback
// @Summary Create new feedback based on parameters
// @Description Create new feedback
// @Tags Feedbacks
// @Accept json
// @Param user body createSynchronizationFeedbackRequest true "Feedback Request"
// @Success 200 {string} txHash
// @Failure 400,422,500 {string} errorMessage
// @Router /feedback [post]
func CreateSynchronizationFeedbackHandler() gin.HandlerFunc {
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

		req := &createSynchronizationFeedbackRequest{}
		err = json.Unmarshal(buf, &req)
		if err != nil {
			restutil.RenderError(err.Error(), 422, c)
			return
		}

		feedbackOffchainMessageId, err := uuid.FromString(req.OriginalOffchainProcessMessageId)

		if err != nil {
			restutil.RenderError(err.Error(), 400, c)
			return
		}

		feedbackOffchainMessage, err := types.GetOffchainMsgById(feedbackOffchainMessageId)

		if err != nil {
			restutil.RenderError(err.Error(), 400, c)
			return
		}

		createFeedbackReq := newFeedbackRequest(*req, feedbackOffchainMessage)

		transactionId := uuid.NewV4()
		feedbackMsg := "Approve"
		if !req.Approved {
			feedbackMsg = "Reject"
		}

		offchainMsg := createFeedbackOffchainMessage(*req, feedbackOffchainMessage, transactionId, feedbackMsg)

		if !offchainMsg.Create() {
			logger.Errorf("error when creating new offchain msg entry")
			restutil.RenderError("error when creating new offchain msg entry", 500, c)
			return
		}

		payload := proxyutil.CreateBaseledgerTransactionFeedbackPayload(createFeedbackReq, &offchainMsg)

		signAndBroadcastPayload := restutil.SignAndBroadcastPayload{
			TransactionId: transactionId.String(),
			Payload:       payload,
			OpCode:        uint32(getRandomFeedbackOpCode()),
		}

		transactionHash := restutil.SignAndBroadcast(signAndBroadcastPayload)

		if transactionHash == nil {
			restutil.RenderError("sign and broadcast transaction error", 500, c)
			return
		}

		trustmeshEntry := createFeedbackSentTrustmeshEntry(*req, transactionId, offchainMsg, feedbackMsg, *transactionHash)

		if !trustmeshEntry.Create() {
			logger.Errorf("error when creating new trustmesh entry")
			restutil.RenderError("error when creating new trustmesh entry", 500, c)
			return
		}

		restutil.Render(transactionHash, 200, c)
	}
}

func getRandomFeedbackOpCode() int {
	rand.Seed(time.Now().UnixNano())
	min := 7
	max := 8

	return rand.Intn(max-min+1) + min
}

func createFeedbackOffchainMessage(
	req createSynchronizationFeedbackRequest,
	feedbackOffchainMessage *types.OffchainProcessMessage,
	transactionId uuid.UUID,
	baseledgerTransactionType string,
) types.OffchainProcessMessage {
	offchainMessage := types.OffchainProcessMessage{
		SenderId:                             uuid.FromStringOrNil(viper.Get("ORGANIZATION_ID").(string)),
		ReceiverId:                           uuid.FromStringOrNil(req.Recipient),
		Topic:                                req.WorkgroupId,
		WorkstepType:                         "Feedback",
		ReferencedOffchainProcessMessageId:   uuid.FromStringOrNil(req.OriginalOffchainProcessMessageId),
		BaseledgerSyncTreeJson:               feedbackOffchainMessage.BaseledgerSyncTreeJson,
		BusinessObjectProof:                  feedbackOffchainMessage.BusinessObjectProof,
		BaseledgerBusinessObjectId:           "",
		ReferencedBaseledgerBusinessObjectId: req.BaseledgerBusinessObjectIdOfApprovedObject,
		StatusTextMessage:                    req.FeedbackMessage,
		BaseledgerTransactionIdOfStoredProof: transactionId,
		TendermintTransactionIdOfStoredProof: transactionId,
		BusinessObjectType:                   req.BusinessObjectType,
		BaseledgerTransactionType:            baseledgerTransactionType,
		ReferencedBaseledgerTransactionId:    uuid.FromStringOrNil(req.OriginalBaseledgerTransactionId),
		EntryType:                            common.FeedbackSentTrustmeshEntryType,
		SorBusinessObjectId:                  feedbackOffchainMessage.SorBusinessObjectId,
	}

	return offchainMessage
}

func createFeedbackSentTrustmeshEntry(req createSynchronizationFeedbackRequest, transactionId uuid.UUID, offchainMsg types.OffchainProcessMessage, feedbackMsg string, txHash string) *types.TrustmeshEntry {
	trustmeshEntry := &types.TrustmeshEntry{
		TendermintTransactionId:              transactionId,
		OffchainProcessMessageId:             offchainMsg.Id,
		SenderOrgId:                          uuid.FromStringOrNil(viper.Get("ORGANIZATION_ID").(string)),
		ReceiverOrgId:                        uuid.FromStringOrNil(req.Recipient),
		WorkgroupId:                          uuid.FromStringOrNil(req.WorkgroupId),
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
		SorBusinessObjectId:                  offchainMsg.SorBusinessObjectId,
	}

	return trustmeshEntry
}

func newFeedbackRequest(req createSynchronizationFeedbackRequest, feedbackOffchainMessage *types.OffchainProcessMessage) *types.SynchronizationFeedback {
	return &types.SynchronizationFeedback{
		WorkgroupId:                        uuid.FromStringOrNil(req.WorkgroupId),
		BaseledgerProvenBusinessObjectJson: feedbackOffchainMessage.BaseledgerSyncTreeJson,
		Recipient:                          req.Recipient,
		Approved:                           req.Approved,
		BaseledgerBusinessObjectIdOfApprovedObject: req.BaseledgerBusinessObjectIdOfApprovedObject,
		HashOfObjectToApprove:                      feedbackOffchainMessage.BusinessObjectProof,
		OriginalBaseledgerTransactionId:            req.OriginalBaseledgerTransactionId,
		OriginalOffchainProcessMessageId:           req.OriginalOffchainProcessMessageId,
		FeedbackMessage:                            req.FeedbackMessage,
		BusinessObjectType:                         req.BusinessObjectType,
	}
}
