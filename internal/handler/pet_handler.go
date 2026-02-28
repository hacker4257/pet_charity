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

type PetHandler struct {
	petService *service.PetService
}

func NewPetHandler(petService *service.PetService) *PetHandler {
	return &PetHandler{
		petService: petService,
	}
}

func (h *PetHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.CreatePetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	pet, err := h.petService.Create(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, pet)
}

func (h *PetHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	petID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req service.UpdatePetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	pet, err := h.petService.Update(userID, uint(petID), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, pet)
}

func (h *PetHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	petID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.petService.Delete(userID, uint(petID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *PetHandler) List(c *gin.Context) {
	page := utils.GetPage(c)
	pageSize := utils.GetPageSize(c)

	orgID, _ := strconv.ParseUint(c.Query("org_id"), 10, 64)
	ageMin, _ := strconv.Atoi(c.Query("age_min"))
	ageMax, _ := strconv.Atoi((c.Query("age_max")))

	filter := repository.PetFilter{
		Species: c.Query("species"),
		Breed:   c.Query("breed"),
		Gender:  c.Query("gender"),
		Status:  c.Query("status"),
		OrgID:   uint(orgID),
		AgeMin:  ageMin,
		AgeMax:  ageMax,
	}

	pets, total, err := h.petService.List(filter, page, pageSize)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}
	response.SuccessWithPage(c, pets, total, page, pageSize)
}

func (h *PetHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	pet, err := h.petService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "pet not found")
		return
	}

	response.Success(c, pet)
}

func (h *PetHandler) UploadImage(c *gin.Context) {
	userID := middleware.GetUserID(c)
	petID, err := strconv.ParseUint(c.Param("id"), 10, 64)
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
	url, err := upload.SaveFile(file, "pet")
	if err != nil {
		response.ServerError(c, "upload failed")
		return
	}

	sortOrder, _ := strconv.Atoi(c.DefaultPostForm("sort_order", "0"))

	image, err := h.petService.AddImage(userID, uint(petID), url, sortOrder)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, image)
}

func (h *PetHandler) DeleteImage(c *gin.Context) {
	userID := middleware.GetUserID(c)
	petID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	imageID, err := strconv.ParseUint(c.Param("imageId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid image id")
		return
	}
	if err := h.petService.DeleteImage(userID, uint(petID), uint(imageID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *PetHandler) PublicStats(c *gin.Context) {
	stats, err := h.petService.PublicStats()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)

}
