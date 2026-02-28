package repository

import (
	"github.com/hacker4257/pet_charity/internal/database"
	"github.com/hacker4257/pet_charity/internal/model"
	"gorm.io/gorm"
)

type AdoptionRepo struct {
	db *gorm.DB
}

func NewAdoptionRepo() *AdoptionRepo {
	return &AdoptionRepo{
		db: database.DB,
	}
}

// 创建领养申请
func (r *AdoptionRepo) Create(adoption *model.Adoption) error {
	return r.db.Create(&adoption).Error
}

func (r *AdoptionRepo) FindByID(id uint) (*model.Adoption, error) {
	var adoption model.Adoption

	result := r.db.
		Preload("User").
		Preload("Pet").
		Preload("Pet.Organization").
		First(&adoption, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &adoption, nil
}

func (r *AdoptionRepo) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&model.Adoption{}).Where("id = ?", id).Updates(fields).Error
}

func (r *AdoptionRepo) FindPendingByUserAndPet(userID, petID uint) (*model.Adoption, error) {
	var adoption model.Adoption
	result := r.db.
		Where("user_id = ? AND pet_id = ? AND status = ?", userID, petID, "pending").
		First(&adoption)
	if result.Error != nil {
		return nil, result.Error
	}
	return &adoption, nil
}

// 用户申请表
func (r *AdoptionRepo) ListByUser(userID uint, page, pageSize int) ([]model.Adoption, int64, error) {
	var adoptions []model.Adoption
	var total int64

	db := r.db.Model(&model.Adoption{}).Where("user_id = ?", userID)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("Pet").Preload("Pet.Images").
		Offset(offset).Limit(pageSize).Order("id DESC").
		Find(&adoptions).Error; err != nil {
		return nil, 0, err
	}

	return adoptions, total, nil
}

// 机构收到的申请列表
func (r *AdoptionRepo) ListByOrg(orgID uint, status string, page, pageSize int) ([]model.Adoption, int64, error) {
	var adoptions []model.Adoption
	var total int64

	db := r.db.Model(&model.Adoption{}).
		Joins("JOIN pets ON pets.id = adoptions.pet_id").
		Where("pets.org_id = ?", orgID)
	if status != "" {
		db = db.Where("adoptions.status = ?", status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("User").Preload("Pet").Offset(offset).Limit(pageSize).Order("adoptions.id DESC").Find(&adoptions).Error; err != nil {
		return nil, 0, err
	}

	return adoptions, total, nil
}

func (r *AdoptionRepo) WithTx(tx *gorm.DB) AdoptionRepository {
	return &AdoptionRepo{db: tx}
}
