package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hacker4257/pet_charity/internal/middleware"
	"github.com/hacker4257/pet_charity/internal/service"
	"github.com/hacker4257/pet_charity/pkg/response"
	"github.com/hacker4257/pet_charity/pkg/utils"
)

type AdoptionHandler struct {
	adoptionService *service.AdoptionService
}

func NewAdoptionHandler(adService *service.AdoptionService) *AdoptionHandler {
	return &AdoptionHandler{
		adoptionService: adService,
	}
}

func (h *AdoptionHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.CreateAdoptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	adoption, err := h.adoptionService.Create(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, adoption)
}

func (h *AdoptionHandler) ListMine(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)

	adoptions, total, err := h.adoptionService.ListByUser(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}
	response.SuccessWithPage(c, adoptions, total, page, pageSize)
}

func (h *AdoptionHandler) GetByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	adoption, err := h.adoptionService.GetByID(userID, uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, adoption)
}

// 机构：收到的申请列表
func (h *AdoptionHandler) ListByOrg(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)
	status := c.Query("status")

	adoptions, total, err := h.adoptionService.ListByOrg(userID, status, page,
		pageSize)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithPage(c, adoptions, total, page, pageSize)
}

// 机构：审核申请
func (h *AdoptionHandler) Review(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req service.ReviewAdoptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.adoptionService.Review(userID, uint(id), &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// 机构：确认领养完成
func (h *AdoptionHandler) Complete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.adoptionService.Complete(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}
