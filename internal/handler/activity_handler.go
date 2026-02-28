package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hacker4257/pet_charity/internal/middleware"
	"github.com/hacker4257/pet_charity/internal/service"
	"github.com/hacker4257/pet_charity/pkg/response"
	"github.com/hacker4257/pet_charity/pkg/utils"
)

type ActivityHandler struct {
	activityService *service.ActivityService
}

func NewActivityHandler(actService *service.ActivityService) *ActivityHandler {
	return &ActivityHandler{activityService: actService}
}

//排行榜
func (h *ActivityHandler) Leaderboard(c *gin.Context) {
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)

	items, err := h.activityService.GetLeaderboard(page, pageSize)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}
	response.Success(c, items)
}

//我的排名
func (h *ActivityHandler) MyRank(c *gin.Context) {
	userID := middleware.GetUserID(c)

	info, err := h.activityService.GetMyRank(userID)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}
	response.Success(c, info)
}
