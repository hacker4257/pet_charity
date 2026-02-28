package repository

import (
	"github.com/hacker4257/pet_charity/internal/database"
	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/hacker4257/pet_charity/pkg/geo"
	"gorm.io/gorm"
)

type RescueRepo struct {
	db *gorm.DB
}

func NewRescueRepo() *RescueRepo {
	return &RescueRepo{
		db: database.DB,
	}
}

func (r *RescueRepo) Create(rescue *model.Rescue) error {
	return r.db.Create(rescue).Error
}

// 根据id查找
func (r *RescueRepo) FindByID(id uint) (*model.Rescue, error) {
	var rescue model.Rescue
	result := r.db.
		Preload("Reporter").
		Preload("Images").
		Preload("Follows").
		Preload("Follows.User").
		First(&rescue, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &rescue, nil
}

// 更新指定字段
func (r *RescueRepo) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&model.Rescue{}).Where("id = ?",
		id).Updates(fields).Error
}

func (r *RescueRepo) List(filter RescueFilter, page, pageSize int) ([]model.Rescue,
	int64, error) {
	var rescues []model.Rescue
	var total int64

	db := r.db.Model(&model.Rescue{})

	if filter.Species != "" {
		db = db.Where("species = ?", filter.Species)
	}
	if filter.Urgency != "" {
		db = db.Where("urgency = ?", filter.Urgency)
	}
	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	} else {
		// 默认不显示已关闭的
		db = db.Where("status != ?", "closed")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("Reporter").Preload("Images").
		Offset(offset).Limit(pageSize).
		Order("FIELD(urgency, 'critical', 'high', 'medium', 'low'), id DESC").
		Find(&rescues).Error; err != nil {
		return nil, 0, err
	}

	return rescues, total, nil
}

type RescueFilter struct {
	Species string
	Urgency string
	Status  string
}

func (r *RescueRepo) ListForMap() ([]model.Rescue, error) {
	var rescues []model.Rescue
	result := r.db.
		Select("id, title, species, urgency, status, longitude, latitude, address,created_at").
		Where("status != ? AND longitude != 0 AND latitude != 0", "closed").
		Order("FIELD(urgency, 'critical', 'high', 'medium', 'low')").
		Limit(500).
		Find(&rescues)
	if result.Error != nil {
		return nil, result.Error
	}
	return rescues, nil
}

func (r *RescueRepo) FindNearby(lng, lat float64, radiusKm float64, limit int) ([]model.Rescue, error) {
	items, err := geo.Search(geo.KeyRescues, lng, lat, radiusKm, limit)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return []model.Rescue{}, nil
	}

	//2.ids
	ids := make([]uint, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}

	//3.mysql
	var rescues []model.Rescue
	if err := r.db.Where("id IN ? AND status != ?", ids, "closed").
		Preload("Reporter").
		Preload("Images").
		Find(&rescues).Error; err != nil {
		return nil, err
	}
	rescueMap := make(map[uint]model.Rescue, len(rescues))
	for _, r := range rescues {
		rescueMap[r.ID] = r
	}
	//排序
	sorted := make([]model.Rescue, 0, len(items))
	for _, item := range items {
		if r, ok := rescueMap[item.ID]; ok {
			sorted = append(sorted, r)
		}
	}
	return sorted, nil
}

// 添加图片
func (r *RescueRepo) CreateImage(image *model.RescueImage) error {
	return r.db.Create(image).Error
}

// 添加跟进记录
func (r *RescueRepo) CreateFollow(follow *model.RescueFollow) error {
	return r.db.Create(follow).Error
}

// 创建认领记录
func (r *RescueRepo) CreateClaim(claim *model.RescueClaim) error {
	return r.db.Create(claim).Error
}

// 查找某救助的认领记录
func (r *RescueRepo) FindClaimByRescueAndOrg(rescueID, orgID uint) (*model.RescueClaim, error) {
	var claim model.RescueClaim
	result := r.db.
		Where("rescue_id = ? AND org_id = ?", rescueID, orgID).
		First(&claim)
	if result.Error != nil {
		return nil, result.Error
	}
	return &claim, nil
}

// 查找某救助的活跃认领（非完成状态的）
func (r *RescueRepo) FindActiveClaimByRescue(rescueID uint) (*model.RescueClaim, error) {
	var claim model.RescueClaim
	result := r.db.
		Where("rescue_id = ? AND status IN ?", rescueID, []string{"claimed", "in_progress"}).
		Preload("Organization").
		First(&claim)
	if result.Error != nil {
		return nil, result.Error
	}
	return &claim, nil
}

// 更新认领记录
func (r *RescueRepo) UpdateClaimFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&model.RescueClaim{}).Where("id = ?", id).Updates(fields).Error
}
