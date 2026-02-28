package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hacker4257/pet_charity/internal/middleware"
	"github.com/hacker4257/pet_charity/internal/service"
	"github.com/hacker4257/pet_charity/pkg/response"
)

type FavoriteHandler struct {
	favoService *service.FavoriteService
}

func NewFavoriteHandler(favo *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{
		favoService: favo,
	}
}

// Toggle 收藏/取消收藏
func (h *FavoriteHandler) Toggle(c *gin.Context) {
	userID := middleware.GetUserID(c)
	petID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid pet id")
		return
	}

	favorited, err := h.favoService.Toggle(userID, uint(petID))
	if err != nil {
		response.ServerError(c, "operation failed")
		return
	}

	response.Success(c, gin.H{"favorited": favorited})
}

// GetStatus 收藏状态
func (h *FavoriteHandler) GetStatus(c *gin.Context) {
	userID := middleware.GetUserID(c)
	petID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid pet id")
		return
	}

	status, err := h.favoService.GetStatus(userID, uint(petID))
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}

	response.Success(c, status)
}

// ListMine 我的收藏列表
func (h *FavoriteHandler) ListMine(c *gin.Context) {
	userID := middleware.GetUserID(c)

	pets, err := h.favoService.ListMyFavorites(userID)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}

	response.Success(c, pets)
}
