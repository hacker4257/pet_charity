package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hacker4257/pet_charity/internal/middleware"
	"github.com/hacker4257/pet_charity/internal/service"
	"github.com/hacker4257/pet_charity/pkg/response"
	"github.com/hacker4257/pet_charity/pkg/utils"
)

type PetDiaryHandler struct {
	diaryService *service.PetDiaryService
}

func NewPetDiaryHandler(diaryService *service.PetDiaryService) *PetDiaryHandler {
	return &PetDiaryHandler{
		diaryService: diaryService,
	}
}

func (h *PetDiaryHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.CreateDiaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	diary, err := h.diaryService.Create(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, diary)
}

func (h *PetDiaryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	diary, err := h.diaryService.FindByID(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, diary)
}

func (h *PetDiaryHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)

	diaryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req service.UpdateDiaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	diary, err := h.diaryService.Update(userID, uint(diaryID), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, diary)
}

func (h *PetDiaryHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)

	diaryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.diaryService.Delete(userID, uint(diaryID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *PetDiaryHandler) UploadImage(c *gin.Context) {
	userID := middleware.GetUserID(c)

	diaryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		response.BadRequest(c, "no file uploaded")
		return
	}

	sortOrder, _ := strconv.Atoi(c.DefaultPostForm("sort_order", "0"))

	image, err := h.diaryService.AddImage(userID, uint(diaryID), file, sortOrder)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, image)
}

func (h *PetDiaryHandler) DeleteImage(c *gin.Context) {
	userID := middleware.GetUserID(c)

	diaryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	imageID, err := strconv.ParseUint(c.Param("imageId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid image id")
		return
	}

	if err := h.diaryService.DeleteImage(userID, uint(diaryID), uint(imageID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *PetDiaryHandler) ToggleLike(c *gin.Context) {
	userID := middleware.GetUserID(c)

	diaryID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	liked, err := h.diaryService.ToggleLike(userID, uint(diaryID))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"liked": liked})
}

func (h *PetDiaryHandler) ListByPet(c *gin.Context) {
	petID, err := strconv.ParseUint(c.Param("petId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid pet id")
		return
	}

	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)

	diaries, total, err := h.diaryService.ListByPet(uint(petID), page, pageSize)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}
	response.SuccessWithPage(c, diaries, total, page, pageSize)
}

func (h *PetDiaryHandler) ListPublic(c *gin.Context) {
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)

	diaries, total, err := h.diaryService.ListPublic(page, pageSize)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}
	response.SuccessWithPage(c, diaries, total, page, pageSize)
}
