package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hacker4257/pet_charity/internal/middleware"
	"github.com/hacker4257/pet_charity/internal/service"
	"github.com/hacker4257/pet_charity/pkg/response"
	"github.com/hacker4257/pet_charity/pkg/utils"
)

type NotificationHandler struct {
	notifyService *service.NotificationService
}

func NewNotificationHandler(notifyService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifyService: notifyService}
}

// 通知列表
func (h *NotificationHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)

	list, total, err := h.notifyService.List(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}
	response.SuccessWithPage(c, list, total, page, pageSize)
}

// 未读数
func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID := middleware.GetUserID(c)
	count, err := h.notifyService.UnreadCount(userID)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}
	response.Success(c, gin.H{"count": count})
}

// 标记单条已读
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.notifyService.MarkRead(userID, uint(id)); err != nil {
		response.ServerError(c, "operation failed")
		return
	}
	response.Success(c, nil)
}

// 全部已读
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.notifyService.MarkAllRead(userID); err != nil {
		response.ServerError(c, "operation failed")
		return
	}
	response.Success(c, nil)
}
