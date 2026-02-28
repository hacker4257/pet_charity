package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/hacker4257/pet_charity/internal/config"
	"github.com/hacker4257/pet_charity/internal/middleware"
	"github.com/hacker4257/pet_charity/internal/service"
	"github.com/hacker4257/pet_charity/pkg/response"
	"github.com/hacker4257/pet_charity/pkg/utils"
)

type DonationHandler struct {
	donationService *service.DonationService
}

func NewDonationHandler(doService *service.DonationService) *DonationHandler {
	return &DonationHandler{
		donationService: doService,
	}
}

//创建捐赠
func (h *DonationHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req service.CreateDonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.donationService.Create(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

//查看订单状态
func (h *DonationHandler) GetStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	donation, err := h.donationService.GetStatus(userID, uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, donation)
}

//我的捐赠记录
func (h *DonationHandler) ListMine(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)

	donations, total, err := h.donationService.ListByUser(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}

	response.SuccessWithPage(c, donations, total, page, pageSize)
}

func (h *DonationHandler) ListPublic(c *gin.Context) {
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)
	targetType := c.Query("target_type")
	targetID, _ := strconv.ParseUint(c.Query("target_id"), 10, 64)

	donations, total, err := h.donationService.ListPublic(targetType, uint(targetID), page, pageSize)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}
	response.SuccessWithPage(c, donations, total, page, pageSize)
}

//微信支付回调
func (h *DonationHandler) WechatNotify(c *gin.Context) {
	if err := h.donationService.HandleWechatNotify(c.Request); err != nil {
		c.JSON(500, gin.H{"code": "FAIL", "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"code": "SUCCESS", "message": "OK"})
}

func (h *DonationHandler) AlipayNotify(c *gin.Context) {
	if err := h.donationService.HandleAlipayNotify(c.Request); err != nil {
		c.String(200, "fail")
		return
	}
	c.String(200, "success")
}

//统计
func (h *DonationHandler) Stats(c *gin.Context) {
	_ = config.Global
	stats, err := h.donationService.Stats()
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}

	response.Success(c, stats)
}
