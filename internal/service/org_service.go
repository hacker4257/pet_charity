package service

import (
	"errors"

	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/hacker4257/pet_charity/internal/repository"
	"github.com/hacker4257/pet_charity/pkg/geo"
)

type OrgService struct {
	orgRepo  repository.OrgRepository
	userRepo repository.UserRepository
}

func NewOrgService(orgRepo repository.OrgRepository, userRepo repository.UserRepository) *OrgService {
	return &OrgService{
		orgRepo:  orgRepo,
		userRepo: userRepo,
	}
}

// 申请入驻请求
type CreateOrgRequest struct {
	Name          string  `json:"name" binding:"required,max=100"`
	LiscenseNo    string  `json:"license_no" binding:"required,max=100"`
	Description   string  `json:"description" binding:"required"`
	Address       string  `json:"address" binding:"required,max=500"`
	ContactPhone  string  `json:"contact_phone" binding:"required"`
	OpeningHours  string  `json:"opening_hours" binding:"max=100"`
	AcceptSpecies string  `json:"accept_species" binding:"max=200"`
	Capacity      int     `json:"capacity"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
}

// 更新机构信息请求
type UpdateOrgRequest struct {
	Description   string  `json:"description"`
	Address       string  `json:"address" binding:"max=500"`
	ContactPhone  string  `json:"contact_phone"`
	OpeningHours  string  `json:"opening_hours" binding:"max=100"`
	AcceptSpecies string  `json:"accept_species" binding:"max=200"`
	Capacity      int     `json:"capacity"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
}

// 审核请求
type ReviewOrgRequest struct {
	Status       string `json:"status" binding:"required,oneof=approved rejected"`
	RejectReason string `json:"reject_reason"`
}

// 申请入驻
func (s *OrgService) Create(userID uint, req *CreateOrgRequest) (*model.Organization, error) {
	//1. 检查是否已申请过
	existing, err := s.orgRepo.FindByUserID(userID)
	if err == nil {
		if existing.Status == "pending" {
			return nil, errors.New("you already have a pending application")
		}
		if existing.Status == "approved" {
			return nil, errors.New("you already have an approved organization")
		}
	}

	//如果是被拒绝过的，更新而不是新建
	if existing != nil && existing.Status == "rejected" {
		fields := map[string]interface{}{
			"name":           req.Name,
			"license_no":     req.LiscenseNo,
			"description":    req.Description,
			"address":        req.Address,
			"contact_phone":  req.ContactPhone,
			"opening_hours":  req.OpeningHours,
			"accept_species": req.AcceptSpecies,
			"capacity":       req.Capacity,
			"longitude":      req.Longitude,
			"latitude":       req.Latitude,
			"status":         "pending",
		}

		if err := s.orgRepo.UpdateFields(existing.ID, fields); err != nil {
			return nil, err
		}
		return s.orgRepo.FindByID(existing.ID)
	}

	//新建申请
	org := &model.Organization{
		UserID:        userID,
		Name:          req.Name,
		LicenseNo:     req.LiscenseNo,
		Description:   req.Description,
		Address:       req.Address,
		ContactPhone:  req.ContactPhone,
		OpeningHours:  req.OpeningHours,
		AcceptSpecies: req.AcceptSpecies,
		Capacity:      req.Capacity,
		Longitude:     req.Longitude,
		Latitude:      req.Latitude,
		Status:        "pending",
	}

	if err := s.orgRepo.Create(org); err != nil {
		return nil, err
	}
	return org, nil
}

// 更新机构信息
func (s *OrgService) Update(userID uint, req *UpdateOrgRequest) (*model.Organization, error) {
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("organization not found")
	}
	if org.Status != "approved" {
		return nil, errors.New("only approved organization can update info")
	}

	fields := map[string]interface{}{}

	if req.Description != "" {
		fields["description"] = req.Description
	}
	if req.Address != "" {
		fields["address"] = req.Address
	}
	if req.ContactPhone != "" {
		fields["contact_phone"] = req.ContactPhone
	}
	if req.OpeningHours != "" {
		fields["opening_hours"] = req.OpeningHours
	}
	if req.AcceptSpecies != "" {
		fields["accept_species"] = req.AcceptSpecies
	}
	if req.Capacity > 0 {
		fields["capacity"] = req.Capacity
	}
	if req.Longitude != 0 {
		fields["longitude"] = req.Longitude
	}
	if req.Latitude != 0 {
		fields["latitude"] = req.Latitude
	}

	if len(fields) == 0 {
		return org, nil
	}

	if err := s.orgRepo.UpdateFields(org.ID, fields); err != nil {
		return nil, err
	}

	if req.Longitude != 0 || req.Latitude != 0 {
		updateOrg, _ := s.orgRepo.FindByID(org.ID)
		if updateOrg != nil {
			geo.Add(geo.KeyOrgs, org.ID, updateOrg.Longitude, updateOrg.Latitude)
		}
	}

	return s.orgRepo.FindByID(org.ID)
}

// 获取机构详情
func (s *OrgService) GetByID(id uint) (*model.Organization, error) {
	return s.orgRepo.FindByID(id)
}

// 已审核机构列表
func (s *OrgService) ListApproved(page, pageSize int) ([]model.Organization, int64, error) {
	return s.orgRepo.ListApproved(page, pageSize)
}

// 审核机构（管理员）
func (s *OrgService) Review(orgID uint, req *ReviewOrgRequest) error {
	org, err := s.orgRepo.FindByID(orgID)
	if err != nil {
		return errors.New("organization not found")
	}
	if org.Status != "pending" {
		return errors.New("only pending organization can be reviewed")
	}
	if req.Status == "approved" {
		return s.orgRepo.ApproveWithTX(orgID, org.UserID)
	}
	//拒绝
	fields := map[string]interface{}{
		"status": "rejected",
	}
	return s.orgRepo.UpdateFields(orgID, fields)
}

// 待审核列表
func (s *OrgService) ListPending(page, pageSize int) ([]model.Organization, int64, error) {
	return s.orgRepo.ListPending(page, pageSize)
}

// 附近救助站
func (s *OrgService) FindNearby(lng, lat float64, radiusKm float64) ([]model.Organization, error) {
	if radiusKm <= 0 {
		radiusKm = 10 //默认10公里
	}
	if radiusKm > 50 {
		radiusKm = 50 //最大50公里
	}

	return s.orgRepo.FindNearby(lng, lat, radiusKm, 20)
}

func (s *OrgService) CountApproved() (int64, error) {
	return s.orgRepo.CountByStatus("approved")
}

func (s *OrgService) FindByUserID(userID uint) (*model.Organization, error) {
	return s.orgRepo.FindByUserID(userID)
}
