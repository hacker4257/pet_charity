package repository

import (
	"github.com/hacker4257/pet_charity/internal/database"
	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/hacker4257/pet_charity/pkg/geo"
	"gorm.io/gorm"
)

type OrgRepo struct {
	db *gorm.DB
}

func NewOrgRepo() *OrgRepo {
	return &OrgRepo{
		db: database.DB,
	}
}

// 根据id查找
func (r *OrgRepo) FindByID(id uint) (*model.Organization, error) {
	var org model.Organization

	result := r.db.Preload("User").First(&org, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &org, nil
}

// 根据用户id
func (r *OrgRepo) FindByUserID(userID uint) (*model.Organization, error) {
	var org model.Organization

	result := r.db.Where("user_id = ?", userID).First(&org)
	if result.Error != nil {
		return nil, result.Error
	}

	return &org, nil
}

// 创建机构
func (r *OrgRepo) Create(org *model.Organization) error {
	return r.db.Create(org).Error
}

// 更新机构
func (r *OrgRepo) UpdateFields(id uint, fields map[string]interface{}) error {
	return r.db.Model(&model.Organization{}).Where("id = ?", id).Updates(fields).Error
}

// 已审核通过的机构列表 （分页）
func (r *OrgRepo) ListApproved(page, pageSize int) ([]model.Organization, int64, error) {
	var orgs []model.Organization
	var total int64

	db := r.db.Model(&model.Organization{}).Where("status = ?", "approved")

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("User").Offset(offset).Limit(pageSize).Order("id DESC").Find(&orgs).Error; err != nil {
		return nil, 0, err
	}

	return orgs, total, nil
}

// 待审核机构列表
func (r *OrgRepo) ListPending(page, pageSize int) ([]model.Organization, int64, error) {
	var orgs []model.Organization
	var total int64

	db := r.db.Model(&model.Organization{}).Where("status = ?", "pending")
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := db.Preload("User").Offset(offset).Limit(pageSize).Order("id DESC").Find(&orgs).Error; err != nil {
		return nil, 0, err
	}

	return orgs, total, nil
}

// 附近的救助站
func (r *OrgRepo) FindNearby(lng, lat float64, radiusKm float64, limit int) ([]model.Organization, error) {
	//1.redis 拿到id+距离
	items, err := geo.Search(geo.KeyOrgs, lng, lat, radiusKm, limit)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return []model.Organization{}, nil
	}

	//2，提取id
	ids := make([]uint, len(items))
	for i, itme := range items {
		ids[i] = itme.ID
	}

	//3.MySQL拿出数据
	var orgs []model.Organization
	if err := r.db.Where("id IN ? AND status = ?", ids, "approved").Preload("User").Find(&orgs).Error; err != nil {
		return nil, err
	}

	//4.按redis返回距离
	orgMap := make(map[uint]model.Organization, len(orgs))
	for _, org := range orgs {
		orgMap[org.ID] = org
	}

	sorted := make([]model.Organization, 0, len(items))
	for _, item := range items {
		if org, ok := orgMap[item.ID]; ok {
			sorted = append(sorted, org)
		}
	}

	return sorted, nil
}

// 用事务更新机构状态和用户角色
func (r *OrgRepo) ApproveWithTX(orgID uint, userID uint) error {
	//先查出机构坐标
	org, err := r.FindByID(orgID)
	if err != nil {
		return err
	}
	err = r.db.Transaction(func(tx *gorm.DB) error {
		//更新机构状态
		if err := tx.Model(&model.Organization{}).Where("id = ?", orgID).Updates(map[string]interface{}{
			"status":      "approved",
			"verified_at": gorm.Expr("NOW()"),
		}).Error; err != nil {
			return err
		}

		// 更新用户角色
		if err := tx.Model(&model.User{}).Where("id = ?", userID).Update("role", "org").Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	geo.Add(geo.KeyOrgs, orgID, org.Longitude, org.Latitude)
	return nil
}

func (r *OrgRepo) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Organization{}).Where("status = ?", status).Count(&count).Error
	return count, err
}
