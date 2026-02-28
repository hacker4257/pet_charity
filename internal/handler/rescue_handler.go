package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hacker4257/pet_charity/internal/middleware"
	"github.com/hacker4257/pet_charity/internal/repository"
	"github.com/hacker4257/pet_charity/internal/service"
	"github.com/hacker4257/pet_charity/pkg/response"
	"github.com/hacker4257/pet_charity/pkg/upload"
	"github.com/hacker4257/pet_charity/pkg/utils"
)

type RescueHandler struct {
	rescueService *service.RescueService
}

func NewRescueHandler(reService *service.RescueService) *RescueHandler {
	return &RescueHandler{
		rescueService: reService,
	}
}

// 上报流浪动物
func (h *RescueHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.CreateRescueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rescue, err := h.rescueService.Create(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, rescue)
}

// 救助详情（公开）
func (h *RescueHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	rescue, err := h.rescueService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "rescue not found")
		return
	}

	response.Success(c, rescue)
}

// 救助列表（公开）
func (h *RescueHandler) List(c *gin.Context) {
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)

	filter := repository.RescueFilter{
		Species: c.Query("species"),
		Urgency: c.Query("urgency"),
		Status:  c.Query("status"),
	}

	rescues, total, err := h.rescueService.List(filter, page, pageSize)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}

	response.SuccessWithPage(c, rescues, total, page, pageSize)
}

// 地图数据（公开）
func (h *RescueHandler) MapData(c *gin.Context) {
	rescues, err := h.rescueService.ListForMap()
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}

	response.Success(c, rescues)
}

// 更新救助信息
func (h *RescueHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req service.CreateRescueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rescue, err := h.rescueService.Update(userID, uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, rescue)
}

// 上传救助现场图片
func (h *RescueHandler) UploadImage(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		response.BadRequest(c, "no file uploaded")
		return
	}

	if err := upload.Validate(file); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	url, err := upload.SaveFile(file, "rescue")
	if err != nil {
		response.ServerError(c, "upload failed")
		return
	}

	sortOrder, _ := strconv.Atoi(c.DefaultPostForm("sort_order", "0"))

	image, err := h.rescueService.AddImage(userID, uint(id), url, sortOrder)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, image)
}

// 添加跟进记录
func (h *RescueHandler) AddFollow(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req service.CreateFollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	follow, err := h.rescueService.AddFollow(userID, uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, follow)
}

// 机构认领救助
func (h *RescueHandler) Claim(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req service.CreateClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	claim, err := h.rescueService.Claim(userID, uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, claim)
}

// 更新认领进度
func (h *RescueHandler) UpdateClaim(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req service.UpdateClaimRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.rescueService.UpdateClaim(userID, uint(id), &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// 救助转领养宠物
func (h *RescueHandler) ConvertToPet(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req service.ConvertToPetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	pet, err := h.rescueService.ConverToPet(userID, uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, pet)
}
