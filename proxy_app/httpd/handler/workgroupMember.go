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

type workgroupMemberDetailsDto struct {
	Id                   uuid.UUID `json:"id"`
	WorkgroupId          string    `json:"workgroup_id"`
	OrganizationId       string    `json:"organization_id"`
	OrganizationEndpoint string    `json:"organization_endpoint"`
	OrganizationToken    string    `json:"organization_token"`
}

type getWorkgroupMemberRequest struct {
	WorkgroupId string `json:"workgroup_id"`
}

type createWorkgroupMemberRequest struct {
	WorkgroupId          string `json:"workgroup_id"`
	OrganizationId       string `json:"organization_id"`
	OrganizationEndpoint string `json:"organization_endpoint"`
	OrganizationToken    string `json:"organization_token"`
}

type deleteWorkgroupMemberRequest struct {
	Id string `json:"workgroup_member_id"`
}

func GetWorkgroupMemberHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, err := c.GetRawData()
		if err != nil {
			restutil.RenderError(err.Error(), 400, c)
			return
		}

		req := &getWorkgroupMemberRequest{}
		err = json.Unmarshal(buf, &req)
		if err != nil {
			restutil.RenderError(err.Error(), 422, c)
			return
		}

		var workgroupMembers []types.WorkgroupMember

		dbutil.Db.GetConn().Where("workgroup_id=?", req.WorkgroupId).Find(&workgroupMembers)

		var workgroupMembersDtos []workgroupMemberDetailsDto

		for i := 0; i < len(workgroupMembers); i++ {
			workgroupMemberDto := &workgroupMemberDetailsDto{}
			workgroupMemberDto.Id = workgroupMembers[i].Id
			workgroupMemberDto.WorkgroupId = workgroupMembers[i].WorkgroupId
			workgroupMemberDto.OrganizationId = workgroupMembers[i].OrganizationId
			workgroupMemberDto.OrganizationEndpoint = workgroupMembers[i].OrganizationEndpoint
			workgroupMemberDto.OrganizationToken = workgroupMembers[i].OrganizationToken
			workgroupMembersDtos = append(workgroupMembersDtos, *workgroupMemberDto)
		}

		restutil.Render(workgroupMembersDtos, 200, c)
	}
}

func CreateWorkgroupMemberHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, err := c.GetRawData()
		if err != nil {
			restutil.RenderError(err.Error(), 400, c)
			return
		}

		req := &createWorkgroupMemberRequest{}
		err = json.Unmarshal(buf, &req)
		if err != nil {
			restutil.RenderError(err.Error(), 422, c)
			return
		}

		newWorkgroupMember := newWorkgroupMember(*req)

		if !newWorkgroupMember.Create() {
			logger.Errorf("error when creating new workgroup member")
			restutil.RenderError("error when creating new workgroup member", 500, c)
			return
		}

		restutil.Render(newWorkgroupMember.Id, 200, c)
	}
}

func DeleteWorkgroupMemberHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		membershipId := c.Param("id")

		var existingWorkgroupMember types.WorkgroupMember
		dbError := dbutil.Db.GetConn().First(&existingWorkgroupMember, "id = ?", membershipId).Error

		if dbError != nil {
			logger.Errorf("error trying to fetch workgroup member with id %s\n", membershipId)
			restutil.RenderError("workgroup member not found", 404, c)
			return
		}

		if !existingWorkgroupMember.Delete() {
			logger.Errorf("error when deleting workgroup member")
			restutil.RenderError("error when deleting workgroup member", 500, c)
			return
		}

		restutil.Render(nil, 200, c)
	}
}

func newWorkgroupMember(req createWorkgroupMemberRequest) *types.WorkgroupMember {
	return &types.WorkgroupMember{
		WorkgroupId:          req.WorkgroupId,
		OrganizationId:       req.OrganizationId,
		OrganizationEndpoint: req.OrganizationEndpoint,
		OrganizationToken:    req.OrganizationToken,
	}
}
