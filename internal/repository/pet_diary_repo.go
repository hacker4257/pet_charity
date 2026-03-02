package repository

import (
	"github.com/hacker4257/pet_charity/internal/database"
	"github.com/hacker4257/pet_charity/internal/model"
	"gorm.io/gorm"
)



type PetDiaryRepo struct {
	db *gorm.DB
}

func NewPetDiaryRepo() *PetDiaryRepo {
	return &PetDiaryRepo{
		db: database.DB,
	}
}


func (r *PetDiaryRepo) Create(diary *model.PetDiary) error {
	return r.db.Create(diary).Error
}

func (r *PetDiaryRepo) FindByID(id uint) (*model.PetDiary, error) {
	var petdiary model.PetDiary

	result := r.db.Preload("Pet").
		Preload("User").
		Preload("Images").
		First(&petdiary, id)
	if result.Error != nil {
		return &petdiary, result.Error
	}
	return &petdiary, nil
}

func (r *PetDiaryRepo) Update(id uint, fields map[string]interface{}) error {
	return r.db.Model(&model.PetDiary{}).Where("id = ?", id).Updates(fields).Error
}

func (r *PetDiaryRepo) Delete(id uint) error {
	return r.db.Delete(&model.PetDiary{}, id).Error
}

func (r *PetDiaryRepo) ListByPet(petID uint, page, pageSize int) ([]model.PetDiary, int64, error) {
	var items []model.PetDiary
	var total int64

	db := r.db.Model(&model.PetDiary{}).Where("pet_id = ?", petID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := db.Preload("User").Preload("Images").
		Offset(offset).Limit(pageSize).Order("id DESC").Find(&items).Error
	return items, total, err
}

func (r *PetDiaryRepo) ListByUser(userID uint, page, pageSize int) ([]model.PetDiary, int64, error) {
	var items []model.PetDiary
	var total int64

	db := r.db.Model(&model.PetDiary{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := db.Preload("User").Preload("Images").
		Offset(offset).Limit(pageSize).Order("id DESC").Find(&items).Error
	return items, total, err
}
func (r *PetDiaryRepo) ListPublic(page, pageSize int) ([]model.PetDiary, int64, error) {
	var items []model.PetDiary
	var total int64

	db := r.db.Model(&model.PetDiary{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := db.Preload("User").Preload("Pet").Preload("Images").
		Offset(offset).Limit(pageSize).Order("id DESC").Find(&items).Error
	return items, total, err
}

func (r *PetDiaryRepo) CreateImage(image *model.DiaryImage) error {
	return r.db.Create(image).Error
}

func (r *PetDiaryRepo) FindImageByID(imageID uint) (*model.DiaryImage, error) {
	var image model.DiaryImage
	err := r.db.First(&image, imageID).Error
	return &image, err
}

func (r *PetDiaryRepo) DeleteImage(imageID uint) error {
	return r.db.Unscoped().Delete(&model.DiaryImage{}, imageID).Error
}

func (r *PetDiaryRepo) ToggleLike(diaryID, userID uint) (bool, error) {
	var existing model.DiaryLike
	err := r.db.Where("diary_id = ? AND user_id = ?", diaryID, userID).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		//没点过---> 点赞
		like := model.DiaryLike{DiaryID: diaryID, UserID: userID}
		return true, r.db.Create(&like).Error
	}
	if err != nil {
		return false, err
	}

	//已点过赞----> 取消
	return false, r.db.Unscoped().Delete(&existing).Error
}

func (r *PetDiaryRepo) CountLikes(diaryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.DiaryLike{}).Where("diary_id = ?", diaryID).Count(&count).Error
	return count, err
}

func (r *PetDiaryRepo) IsLiked(diaryID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.DiaryLike{}).
		Where("diary_id = ? AND user_id = ?", diaryID, userID).Count(&count).Error

	return count > 0, err
}
