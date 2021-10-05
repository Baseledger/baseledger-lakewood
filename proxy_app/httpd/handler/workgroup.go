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

// @Security BasicAuth
// GetWorkgroups ... Get all workgroups
// @Summary Get all workgroups
// @Description get all workgroups
// @Tags Workgroups
// @Produce json
// @Success 200 {array} workgroupDetailsDto
// @Router /workgroup [get]
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

// @Security BasicAuth
// Create Workgroup ... Create Workgroup
// @Summary Create new workgroup based on parameters
// @Description Create new workgroup
// @Tags Workgroups
// @Accept json
// @Param user body createWorkgroupRequest true "Workgroup Request"
// @Success 200 {string} types.Workgroup
// @Failure 400,422,500 {string} errorMessage
// @Router /workgroup [post]
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

// @Security BasicAuth
// Delete Workgroup ... Delete Workgroup
// @Summary Delete workgroup
// @Description Delete workgroup
// @Tags Workgroups
// @Param id path string format "uuid" "id"
// @Success 204
// @Failure 404,500 {string} errorMessage
// @Router /workgroup/{id} [delete]
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

		restutil.Render(nil, 204, c)
	}
}

func newWorkgroup(req createWorkgroupRequest) *types.Workgroup {
	return &types.Workgroup{
		WorkgroupName: req.Name,
		PrivatizeKey:  req.PrivatizeKey,
	}
}
