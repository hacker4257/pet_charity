package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hacker4257/pet_charity/internal/service"
	"github.com/hacker4257/pet_charity/pkg/response"
	"github.com/hacker4257/pet_charity/pkg/utils"
)

type AdminHandler struct {
	orgService      *service.OrgService
	userService     *service.UserService
	petService      *service.PetService
	donationService *service.DonationService
}

func NewAdminHandler(
	orgService *service.OrgService,
	userService *service.UserService,
	petService *service.PetService,
	donationService *service.DonationService,
) *AdminHandler {
	return &AdminHandler{
		orgService:      orgService,
		userService:     userService,
		petService:      petService,
		donationService: donationService,
	}
}

//待审核机构列表
func (h *AdminHandler) ListPendingOrgs(c *gin.Context) {
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)

	orgs, total, err := h.orgService.ListPending(page, pageSize)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}

	response.SuccessWithPage(c, orgs, total, page, pageSize)
}

//审核机构
func (h *AdminHandler) ReviewOrg(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req service.ReviewOrgRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.orgService.Review(uint(id), &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// Stats 后台统计数据
func (h *AdminHandler) Stats(c *gin.Context) {
	petStats, _ := h.petService.PublicStats()
	donationStats, _ := h.donationService.Stats()
	userCount, _ := h.userService.Count()
	orgCount, _ := h.orgService.CountApproved()

	response.Success(c, gin.H{
		"user_count":     userCount,
		"org_count":      orgCount,
		"pet_stats":      petStats,
		"donation_stats": donationStats,
	})
}

// ListUsers 用户列表
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)
	role := c.Query("role")

	users, total, err := h.userService.List(page, pageSize, role)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}

	response.SuccessWithPage(c, users, total, page, pageSize)
}

// UpdateUserRole 修改用户角色
func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.UpdateRole(uint(id), req.Role); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateUserStatus 禁用/启用用户
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.UpdateStatus(uint(id), req.Status); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}
