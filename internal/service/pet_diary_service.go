package service

import (
	"errors"

	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/hacker4257/pet_charity/internal/repository"
	"github.com/hacker4257/pet_charity/pkg/upload"
	"mime/multipart"
)

type PetDiaryService struct {
	diaryRepo    repository.PetDiaryRepository
	adoptionRepo repository.AdoptionRepository
}

func NewPetDiaryService(diaryRepo repository.PetDiaryRepository,
	adoptionRepo repository.AdoptionRepository) *PetDiaryService {
	return &PetDiaryService{
		diaryRepo:    diaryRepo,
		adoptionRepo: adoptionRepo,
	}
}

// ---- 请求 DTO ----

type CreateDiaryRequest struct {
	PetID   uint   `json:"pet_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateDiaryRequest struct {
	Content string `json:"content" binding:"required"`
}

// ---- 方法 ----

// Create 创建日记
func (s *PetDiaryService) Create(userID uint, req *CreateDiaryRequest) (*model.PetDiary, error) {
	// 验证领养关系
	_, err := s.adoptionRepo.FindApprovedByUserAndPet(userID, req.PetID)
	if err != nil {
		return nil, errors.New("you have not adopted this pet")
	}

	diary := &model.PetDiary{
		UserID:  userID,
		PetID:   req.PetID,
		Content: req.Content,
	}

	if err := s.diaryRepo.Create(diary); err != nil {
		return nil, errors.New("create diary failed")
	}

	return diary, nil
}

// GetByID 查看日记详情（公开，不需要权限）
func (s *PetDiaryService) FindByID(id uint) (*model.PetDiary, error) {
	diary, err := s.diaryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("diary not found")
	}
	return diary, nil
}

// Update 编辑日记（只有作者可以）
func (s *PetDiaryService) Update(userID uint, diaryID uint, req *UpdateDiaryRequest) (*model.PetDiary, error) {
	// 查日记
	diary, err := s.diaryRepo.FindByID(diaryID)
	if err != nil {
		return nil, errors.New("diary not found")
	}

	// 验证是不是你写的
	if diary.UserID != userID {
		return nil, errors.New("you can only edit your own diary")
	}

	// 更新
	fields := map[string]interface{}{
		"content": req.Content,
	}
	if err := s.diaryRepo.Update(diaryID, fields); err != nil {
		return nil, errors.New("update diary failed")
	}

	// 重新查一次返回最新数据
	diary, _ = s.diaryRepo.FindByID(diaryID)
	return diary, nil
}

// Delete 删除日记（只有作者可以）
func (s *PetDiaryService) Delete(userID uint, diaryID uint) error {
	diary, err := s.diaryRepo.FindByID(diaryID)
	if err != nil {
		return errors.New("diary not found")
	}

	if diary.UserID != userID {
		return errors.New("you can only delete your own diary")
	}

	return s.diaryRepo.Delete(diaryID)
}

// AddImage 上传日记图片（只有作者可以）
func (s *PetDiaryService) AddImage(userID uint, diaryID uint, file *multipart.FileHeader, sortOrder int) (*model.DiaryImage, error) {
	diary, err := s.diaryRepo.FindByID(diaryID)
	if err != nil {
		return nil, errors.New("diary not found")
	}

	if diary.UserID != userID {
		return nil, errors.New("permission denied")
	}

	// 校验文件
	if err := upload.Validate(file); err != nil {
		return nil, err
	}

	// 保存文件
	url, err := upload.SaveFile(file, "diary")
	if err != nil {
		return nil, errors.New("upload failed")
	}

	image := &model.DiaryImage{
		DiaryID:   diaryID,
		ImageURL:  url,
		SortOrder: sortOrder,
	}
	if err := s.diaryRepo.CreateImage(image); err != nil {
		return nil, errors.New("save image failed")
	}

	return image, nil
}

// DeleteImage 删除日记图片（只有日记作者可以）
func (s *PetDiaryService) DeleteImage(userID uint, diaryID uint, imageID uint) error {
	// 验证日记归属
	diary, err := s.diaryRepo.FindByID(diaryID)
	if err != nil {
		return errors.New("diary not found")
	}
	if diary.UserID != userID {
		return errors.New("permission denied")
	}

	// 验证图片属于这篇日记
	image, err := s.diaryRepo.FindImageByID(imageID)
	if err != nil {
		return errors.New("image not found")
	}
	if image.DiaryID != diaryID {
		return errors.New("image does not belong to this diary")
	}

	return s.diaryRepo.DeleteImage(imageID)
}

// ToggleLike 点赞/取消点赞（任何登录用户都可以）
func (s *PetDiaryService) ToggleLike(userID uint, diaryID uint) (bool, error) {
	// 只需验证日记存在
	_, err := s.diaryRepo.FindByID(diaryID)
	if err != nil {
		return false, errors.New("diary not found")
	}

	return s.diaryRepo.ToggleLike(diaryID, userID)
}

// ListByPet 按宠物查日记（公开）
func (s *PetDiaryService) ListByPet(petID uint, page, pageSize int) ([]model.PetDiary, int64, error) {
	return s.diaryRepo.ListByPet(petID, page, pageSize)
}

// ListByUser 查某用户的日记
func (s *PetDiaryService) ListByUser(userID uint, page, pageSize int) ([]model.PetDiary, int64, error) {
	return s.diaryRepo.ListByUser(userID, page, pageSize)
}

// ListPublic 公开时间线
func (s *PetDiaryService) ListPublic(page, pageSize int) ([]model.PetDiary, int64, error) {
	return s.diaryRepo.ListPublic(page, pageSize)
}
