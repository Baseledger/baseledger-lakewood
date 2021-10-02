package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/unibrightio/proxy-api/logger"
	"github.com/unibrightio/proxy-api/proxyutil"
	"github.com/unibrightio/proxy-api/restutil"
)

type sendOffchainMessageRequest struct {
	WorkgroupId string `json:"workgroup_id"`
	RecipientId string `json:"recipient_id"`
	Payload     string `json:"payload"`
}

func SendOffchainMessageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Infof("sendOffchainMessageHandler initiated\n")
		buf, err := c.GetRawData()
		if err != nil {
			restutil.RenderError(err.Error(), 400, c)
			return
		}

		req := &sendOffchainMessageRequest{}
		err = json.Unmarshal(buf, &req)
		if err != nil {
			restutil.RenderError(err.Error(), 422, c)
			return
		}

		err = proxyutil.SendOffchainMessage([]byte(req.Payload), req.WorkgroupId, req.RecipientId)

		if err != nil {
			restutil.RenderError(err.Error(), 404, c)
			return
		}

		restutil.Render(nil, 200, c)
	}
}
