package repository

import (
	"github.com/hacker4257/pet_charity/internal/database"
	"github.com/hacker4257/pet_charity/internal/model"
	"gorm.io/gorm"
)

type PetRepo struct {
	db *gorm.DB
}

func NewPetRepo() *PetRepo {
	return &PetRepo{
		db: database.DB,
	}
}

// 创建宠物
func (r *PetRepo) Create(pet *model.Pet) error {
	return r.db.Create(pet).Error
}

// 根据ID查找 预加载机构和图片
func (r *PetRepo) FindByID(id uint) (*model.Pet, error) {
	var pet model.Pet
	result := r.db.
		Preload("Organization").
		First(&pet, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pet, nil
}

// 更新指定字段
func (r *PetRepo) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&model.Pet{}).Where("id = ?", id).Updates(fields).Error
}

// 软删除
func (r *PetRepo) Delete(id uint) error {
	return r.db.Delete(&model.Pet{}, id).Error
}

// 宠物列表
func (r *PetRepo) List(filter PetFilter, page, pageSize int) ([]model.Pet, int64, error) {
	var pets []model.Pet
	var total int64

	db := r.db.Model(&model.Pet{})

	//动态查询
	if filter.Species != "" {
		db = db.Where("species = ?", filter.Species)
	}
	if filter.Breed != "" {
		db = db.Where("breed LIKE ?", "%"+filter.Breed+"%")
	}
	if filter.Gender != "" {
		db = db.Where("gender = ?", filter.Gender)
	}
	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	} else {
		db = db.Where("status = ?", "adoptable")
	}
	if filter.OrgID > 0 {
		db = db.Where("org_id = ?", filter.OrgID)
	}
	if filter.AgeMin > 0 {
		db = db.Where("age >= ?", filter.AgeMin)
	}
	if filter.AgeMax > 0 {
		db = db.Where("age <= ?", filter.AgeMax)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("Organization").Preload("Images").Offset(offset).Limit(pageSize).Order("id DESC").Find(&pets).Error; err != nil {
		return nil, 0, err
	}

	return pets, total, nil
}

type PetFilter struct {
	Species string
	Breed   string
	Gender  string
	Status  string
	OrgID   uint
	AgeMin  int
	AgeMax  int
}

func (r *PetRepo) CreateImage(image *model.PetImage) error {
	return r.db.Create(image).Error
}

func (r *PetRepo) DeleteImage(imageID uint) error {
	return r.db.Delete(&model.PetImage{}, imageID).Error
}

// 查找图片
func (r *PetRepo) FindImageByID(imageID uint) (*model.PetImage, error) {
	var image model.PetImage
	result := r.db.First(&image, imageID)
	if result.Error != nil {
		return nil, result.Error
	}

	return &image, nil
}

func (r *PetRepo) PublicStats() map[string]int64 {
	stats := map[string]int64{}

	var petCount int64
	r.db.Model(&model.Pet{}).Where("status = ?", "adoptable").Count(&petCount)

	var adoptedCount int64
	r.db.Model(&model.Pet{}).Where("status = ?", "adopted").Count(&adoptedCount)

	var rescueCount int64
	r.db.Model(&model.Rescue{}).Count(&rescueCount)

	stats["adoptable_pet"] = petCount
	stats["adopted_pets"] = adoptedCount
	stats["rescue_count"] = rescueCount
	return stats

}

func (r *PetRepo) WithTx(tx *gorm.DB) PetRepository {
	return &PetRepo{db: tx}
}
