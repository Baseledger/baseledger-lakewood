package handler

import (
	"encoding/json"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/unibrightio/proxy-api/logger"
	"github.com/unibrightio/proxy-api/restutil"
	"github.com/unibrightio/proxy-api/synctree"
	"github.com/unibrightio/proxy-api/types"
)

type syncTreeSunburst struct {
	Items  []sunburstItem
	Levels float64
}

type sunburstItem struct {
	Name     string         `json:"name"`
	Value    uint           `json:"value"`
	Children []sunburstItem `json:"children"`
}

func GetSunburstHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		txId := c.Param("txId")

		offchainMessage, err := types.GetOffchainMsgForSunburst(txId)

		if err != nil {
			restutil.RenderError(err.Error(), 422, c)
			return
		}

		syncTree := &synctree.BaseledgerSyncTree{}
		err = json.Unmarshal([]byte(offchainMessage.BaseledgerSyncTreeJson), &syncTree)
		if err != nil {
			logger.Errorf("Error unmarshalling sync tree", err.Error())
			return
		}
		logger.Infof("Sync tree unmarshalled", syncTree)

		sunburst := getSyncTreeSunburst(*syncTree)
		restutil.Render(sunburst, 200, c)
	}
}

func getSyncTreeSunburst(syncTree synctree.BaseledgerSyncTree) syncTreeSunburst {
	var rootNode synctree.SyncTreeNode
	for _, node := range syncTree.Nodes {
		if node.IsRoot {
			rootNode = node
			break
		}
	}

	children := getSunburstChildren(syncTree.Nodes, rootNode)
	return syncTreeSunburst{
		Items:  children,
		Levels: 1 + math.Floor(math.Log2(float64(len(syncTree.Nodes)))),
	}
}

func getSunburstChildren(nodes []synctree.SyncTreeNode, parentSyncTreeNode synctree.SyncTreeNode) []sunburstItem {
	var sunburstItems []sunburstItem

	for _, node := range nodes {
		if node.ParentNodeID == parentSyncTreeNode.SyncTreeNodeID {
			if len(node.Value) == 0 {
				continue
			}
			sunburstItems = append(sunburstItems, sunburstItem{
				Name:  node.Value,
				Value: 50,
				// stack overflow bug here
				Children: getSunburstChildren(nodes, node),
			})
		}
	}

	return sunburstItems
}
