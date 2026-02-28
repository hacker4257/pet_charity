package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hacker4257/pet_charity/internal/repository"
	"github.com/hacker4257/pet_charity/pkg/response"
	"github.com/hacker4257/pet_charity/pkg/utils"
)

type FeedHandler struct {
	feedRepo *repository.FeedRepo
}

func NewFeedHandler(feedRepo *repository.FeedRepo) *FeedHandler {
	return &FeedHandler{feedRepo: feedRepo}
}

//list 全站动态
func (h *FeedHandler) List(c *gin.Context) {
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)

	items, err := h.feedRepo.List(page, pageSize)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}
	response.Success(c, items)
}
