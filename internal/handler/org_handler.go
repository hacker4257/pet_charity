package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hacker4257/pet_charity/internal/middleware"
	"github.com/hacker4257/pet_charity/internal/service"
	"github.com/hacker4257/pet_charity/pkg/response"
	"github.com/hacker4257/pet_charity/pkg/utils"
)

type OrgHandler struct {
	orgService *service.OrgService
}

func NewOrgHandler(orgService *service.OrgService) *OrgHandler {
	return &OrgHandler{
		orgService: orgService,
	}
}

//申请入驻
func (h *OrgHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.CreateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	org, err := h.orgService.Create(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, org)
}

// 更新机构信息
func (h *OrgHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.UpdateOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	org, err := h.orgService.Update(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, org)
}

//机构列表（公开）
func (h *OrgHandler) List(c *gin.Context) {
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)

	orgs, total, err := h.orgService.ListApproved(page, pageSize)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}

	response.SuccessWithPage(c, orgs, total, page, pageSize)
}

//机构详情（公开）
func (h *OrgHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	org, err := h.orgService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "organization not found")
		return
	}
	response.Success(c, org)
}

//附件救助站（公开）
func (h *OrgHandler) Nearby(c *gin.Context) {
	lng, err := strconv.ParseFloat(c.Query("lng"), 64)
	if err != nil {
		response.BadRequest(c, "invalid longitude")
		return
	}
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		response.BadRequest(c, "invalid latitude")
		return
	}

	radius, _ := strconv.ParseFloat(c.DefaultQuery("radius", "10"), 64)

	orgs, err := h.orgService.FindNearby(lng, lat, radius)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}
	response.Success(c, orgs)
}

// GetMine 获取当前用户的机构
func (h *OrgHandler) GetMine(c *gin.Context) {
	userID := middleware.GetUserID(c)
	org, err := h.orgService.FindByUserID(userID)
	if err != nil {
		response.NotFound(c, "未注册机构")
		return
	}
	response.Success(c, org)
}
