package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	uuid "github.com/kthomas/go.uuid"
	"github.com/unibrightio/proxy-api/dbutil"
	"github.com/unibrightio/proxy-api/logger"
	"github.com/unibrightio/proxy-api/restutil"
	"github.com/unibrightio/proxy-api/types"
)

type workgroupDetailsDto struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Key  string    `json:"key"`
}

type createWorkgroupRequest struct {
	Name         string `json:"name"`
	PrivatizeKey string `json:"privatize_key"`
}

type deleteWorkgroupRequest struct {
	Id string `json:"workgroup_id"`
}

func GetWorkgroupsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var workgroups []types.Workgroup
		dbutil.Db.GetConn().Find(&workgroups)

		var workgroupsDtos []workgroupDetailsDto

		for i := 0; i < len(workgroups); i++ {
			workgroupsDto := &workgroupDetailsDto{}
			workgroupsDto.Id = workgroups[i].Id
			workgroupsDto.Name = workgroups[i].WorkgroupName
			workgroupsDto.Key = workgroups[i].PrivatizeKey
			workgroupsDtos = append(workgroupsDtos, *workgroupsDto)
		}

		restutil.Render(workgroupsDtos, 200, c)
	}
}

func CreateWorkgroupHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, err := c.GetRawData()
		if err != nil {
			restutil.RenderError(err.Error(), 400, c)
			return
		}

		req := &createWorkgroupRequest{}
		err = json.Unmarshal(buf, &req)
		if err != nil {
			restutil.RenderError(err.Error(), 422, c)
			return
		}

		newWorkgroup := newWorkgroup(*req)

		if !newWorkgroup.Create() {
			logger.Errorf("error when creating new workgroup")
			restutil.RenderError("error when creating new workgroup", 500, c)
			return
		}

		restutil.Render(newWorkgroup, 200, c)
	}
}

func DeleteWorkgroupHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		workgroupId := c.Param("id")

		var existingWorkgroup types.Workgroup
		dbError := dbutil.Db.GetConn().First(&existingWorkgroup, "id = ?", workgroupId).Error

		if dbError != nil {
			logger.Errorf("error trying to fetch workgroup with id %s\n", workgroupId)
			restutil.RenderError("workgroup not found", 404, c)
			return
		}

		if !existingWorkgroup.Delete() {
			logger.Errorf("error when deleting workgroup")
			restutil.RenderError("error when deleting workgroup", 500, c)
			return
		}

		restutil.Render(nil, 200, c)
	}
}

func newWorkgroup(req createWorkgroupRequest) *types.Workgroup {
	return &types.Workgroup{
		WorkgroupName: req.Name,
		PrivatizeKey:  req.PrivatizeKey,
	}
}
