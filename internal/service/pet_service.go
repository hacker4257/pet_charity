package service

import (
	"errors"

	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/hacker4257/pet_charity/internal/repository"
)

type PetService struct {
	petRepo repository.PetRepository
	orgRepo repository.OrgRepository
}

func NewPetService(petRepo repository.PetRepository, orgRepo repository.OrgRepository) *PetService {
	return &PetService{
		petRepo: petRepo,
		orgRepo: orgRepo,
	}
}

type CreatePetRequest struct {
	Name         string `json:"name" binding:"required,max=50"`
	Species      string `json:"species" binding:"required,oneof=cat dog other"`
	Breed        string `json:"breed" binding:"max=50"`
	Age          int    `json:"age" binding:"min=0"`
	Gender       string `json:"gender" binding:"required,oneof=male female unknown"`
	HealthStatus string `json:"health_status" binding:"max=50"`
	Description  string `json:"description" binding:"required"`
	CoverImage   string `json:"cover_image"`
}

//更新宠物请求
type UpdatePetRequest struct {
	Name         string `json:"name" binding:"required,max=50"`
	Breed        string `json:"breed" binding:"max=50"`
	Age          int    `json:"age" binding:"min=0"`
	Gender       string `json:"gender" binding:"required,oneof=male female unknown"`
	HealthStatus string `json:"health_status" binding:"max=50"`
	Description  string `json:"description" binding:"required"`
	CoverImage   string `json:"cover_image"`
	Status       string `json:"status" binding:"omitempty,oneof=available reserved adopted"`
}

//发布宠物，只有机构可以
func (s *PetService) Create(userID uint, req *CreatePetRequest) (*model.Pet, error) {
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("you don't have an approved organization")
	}
	if org.Status != "approved" {
		return nil, errors.New("your organization is not approved yet")
	}

	pet := &model.Pet{
		OrgID:        org.ID,
		Name:         req.Name,
		Species:      req.Species,
		Breed:        req.Breed,
		Age:          req.Age,
		Gender:       req.Gender,
		HealthStatus: req.HealthStatus,
		Description:  req.Description,
		CoverImage:   req.CoverImage,
		Status:       "available",
	}
	if err := s.petRepo.Create(pet); err != nil {
		return nil, err
	}

	return s.petRepo.FindByID(pet.ID)
}

// 更新宠物信息
func (s *PetService) Update(userID uint, petID uint, req *UpdatePetRequest) (*model.Pet, error) {
	// 验证权限：必须是该宠物所属机构的用户
	pet, err := s.petRepo.FindByID(petID)
	if err != nil {
		return nil, errors.New("pet not found")
	}

	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil || org.ID != pet.OrgID {
		return nil, errors.New("you can only edit your own organization's pets")
	}

	fields := map[string]interface{}{}
	if req.Name != "" {
		fields["name"] = req.Name
	}
	if req.Breed != "" {
		fields["breed"] = req.Breed
	}
	if req.Age > 0 {
		fields["age"] = req.Age
	}
	if req.Gender != "" {
		fields["gender"] = req.Gender
	}
	if req.HealthStatus != "" {
		fields["health_status"] = req.HealthStatus
	}
	if req.Description != "" {
		fields["description"] = req.Description
	}
	if req.CoverImage != "" {
		fields["cover_image"] = req.CoverImage
	}
	if req.Status != "" {
		fields["status"] = req.Status
	}

	if len(fields) == 0 {
		return pet, nil
	}

	if err := s.petRepo.UpdateFields(petID, fields); err != nil {
		return nil, err
	}

	return s.petRepo.FindByID(petID)
}

func (s *PetService) Delete(userID uint, petID uint) error {
	pet, err := s.petRepo.FindByID(petID)
	if err != nil {
		return errors.New("pet not found")
	}
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil || org.ID != pet.OrgID {
		return errors.New("you can only delete your own organization's pets")
	}

	return s.petRepo.Delete(petID)
}

func (s *PetService) GetByID(id uint) (*model.Pet, error) {
	return s.petRepo.FindByID(id)
}

func (s *PetService) List(filter repository.PetFilter, page, pageSize int) ([]model.Pet, int64, error) {
	return s.petRepo.List(filter, page, pageSize)
}

//上传宠物图片
func (s *PetService) AddImage(userID uint, petID uint, imageURL string, sortOrder int) (*model.PetImage, error) {
	pet, err := s.petRepo.FindByID(petID)
	if err != nil {
		return nil, errors.New("pet not found")
	}
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil || org.ID != pet.OrgID {
		return nil, errors.New("permission denied")
	}

	image := &model.PetImage{
		PetID:     petID,
		ImageURL:  imageURL,
		SortOrder: sortOrder,
	}
	if err := s.petRepo.CreateImage(image); err != nil {
		return nil, err
	}
	return image, nil
}

func (s *PetService) DeleteImage(userID uint, petID uint, imageID uint) error {
	pet, err := s.petRepo.FindByID(petID)
	if err != nil {
		return errors.New("pet not found")
	}
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil || org.ID != pet.OrgID {
		return errors.New("permission denied")
	}
	image, err := s.petRepo.FindImageByID(imageID)
	if err != nil {
		return errors.New("image not found")
	}
	if image.PetID != petID {
		return errors.New("image does not belong to this pet")
	}

	return s.petRepo.DeleteImage(imageID)
}

func (s *PetService) PublicStats() (map[string]int64, error) {
	return s.petRepo.PublicStats(), nil
}
