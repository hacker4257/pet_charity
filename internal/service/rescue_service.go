package service

import (
	"errors"

	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/hacker4257/pet_charity/internal/repository"
	"github.com/hacker4257/pet_charity/pkg/event"
	"github.com/hacker4257/pet_charity/pkg/geo"
	"gorm.io/gorm"
)

type RescueService struct {
	rescueRepo repository.RescueRepository
	orgRepo    repository.OrgRepository
	petRepo    repository.PetRepository
}

func NewRescueService(rescueRepo repository.RescueRepository, orgRepo repository.OrgRepository, petRepo repository.PetRepository) *RescueService {
	return &RescueService{
		rescueRepo: rescueRepo,
		orgRepo:    orgRepo,
		petRepo:    petRepo,
	}
}

// 上报请求
type CreateRescueRequest struct {
	Title        string  `json:"title" binding:"required,max=100"`
	Description  string  `json:"description" binding:"required"`
	Species      string  `json:"species" binding:"required,oneof=cat dog other"`
	Urgency      string  `json:"urgency" binding:"required,oneof=low medium high critical"`
	Longitude    float64 `json:"longitude" binding:"required"`
	Latitude     float64 `json:"latitude" binding:"required"`
	Address      string  `json:"address" binding:"required,max=500"`
	ContactPhone string  `json:"contact_phone"`
}

// 添加跟进记录请求
type CreateFollowRequest struct {
	Content string `json:"content" binding:"required"`
}

// 认领请求
type CreateClaimRequest struct {
	Note string `json:"note"`
}

// 更新认领状态请求
type UpdateClaimRequest struct {
	Status string `json:"status" binding:"required,oneof=in_progress completed"`
	Note   string `json:"note"`
}

// 救助转领养请求
type ConvertToPetRequest struct {
	Name         string `json:"name" binding:"required,max=50"`
	Breed        string `json:"breed" binding:"max=50"`
	Age          int    `json:"age"`
	Gender       string `json:"gender" binding:"required,oneof=male female unknown"`
	HealthStatus string `json:"health_status"`
	Description  string `json:"description" binding:"required"`
	CoverImage   string `json:"cover_image"`
}

// 上报流浪动物
func (s *RescueService) Create(userID uint, req *CreateRescueRequest) (*model.Rescue, error) {
	rescue := &model.Rescue{
		ReporterID:   userID,
		Title:        req.Title,
		Description:  req.Description,
		Species:      req.Species,
		Urgency:      req.Urgency,
		Status:       "reported",
		Longitude:    req.Longitude,
		Latitude:     req.Latitude,
		Address:      req.Address,
		ContactPhone: req.ContactPhone,
	}

	if err := s.rescueRepo.Create(rescue); err != nil {
		return nil, err
	}
	geo.Add(geo.KeyRescues, rescue.ID, rescue.Longitude, rescue.Latitude)
	event.Publish("rescue", userID, 0)

	return s.rescueRepo.FindByID(rescue.ID)
}

// 救助详情
func (s *RescueService) GetByID(id uint) (*model.Rescue, error) {
	return s.rescueRepo.FindByID(id)
}

// 救助列表
func (s *RescueService) List(filter repository.RescueFilter, page, pageSize int) ([]model.Rescue, int64, error) {
	return s.rescueRepo.List(filter, page, pageSize)
}

// 地图数据
func (s *RescueService) ListForMap() ([]model.Rescue, error) {
	return s.rescueRepo.ListForMap()
}

// 更新救助信息（只有上报人可以）
func (s *RescueService) Update(userID uint, rescueID uint, req *CreateRescueRequest) (*model.Rescue, error) {
	rescue, err := s.rescueRepo.FindByID(rescueID)
	if err != nil {
		return nil, errors.New("rescue not found")
	}
	if rescue.ReporterID != userID {
		return nil, errors.New("you can only edit your own rescue reports")
	}

	fields := map[string]interface{}{
		"title":         req.Title,
		"description":   req.Description,
		"species":       req.Species,
		"urgency":       req.Urgency,
		"longitude":     req.Longitude,
		"latitude":      req.Latitude,
		"address":       req.Address,
		"contact_phone": req.ContactPhone,
	}

	if err := s.rescueRepo.UpdateFields(rescueID, fields); err != nil {
		return nil, err
	}
	geo.Add(geo.KeyRescues, rescueID, req.Longitude, req.Latitude)
	return s.rescueRepo.FindByID(rescueID)
}

func (s *RescueService) AddImage(userID uint, rescueID uint, imageURL string, sortOrder int) (*model.RescueImage, error) {
	rescue, err := s.rescueRepo.FindByID(rescueID)
	if err != nil {
		return nil, errors.New("rescue not found")
	}
	if rescue.ReporterID != userID {
		return nil, errors.New("permission denied")
	}

	image := &model.RescueImage{
		RescueID:  rescueID,
		ImageURL:  imageURL,
		SortOrder: sortOrder,
	}

	if err := s.rescueRepo.CreateImage(image); err != nil {
		return nil, err
	}

	return image, nil
}

// 添加跟进记录（任何登录用户都可以）
func (s *RescueService) AddFollow(userID uint, rescueID uint, req *CreateFollowRequest) (*model.RescueFollow, error) {
	_, err := s.rescueRepo.FindByID(rescueID)
	if err != nil {
		return nil, errors.New("rescue not found")
	}

	follow := &model.RescueFollow{
		RescueID: rescueID,
		UserID:   userID,
		Content:  req.Content,
	}

	if err := s.rescueRepo.CreateFollow(follow); err != nil {
		return nil, err
	}

	return follow, nil
}

// 机构认领救助
func (s *RescueService) Claim(userID uint, rescueID uint, req *CreateClaimRequest) (*model.RescueClaim, error) {
	// 1. 验证是机构用户
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("you are not an approved organization")
	}
	if org.Status != "approved" {
		return nil, errors.New("your organization is not approved")
	}

	// 2. 验证救助信息存在
	rescue, err := s.rescueRepo.FindByID(rescueID)
	if err != nil {
		return nil, errors.New("rescue not found")
	}
	if rescue.Status == "closed" || rescue.Status == "rescued" {
		return nil, errors.New("this rescue has already been handled")
	}

	// 3. 检查是否已认领
	_, err = s.rescueRepo.FindClaimByRescueAndOrg(rescueID, org.ID)
	if err == nil {
		return nil, errors.New("you have already claimed this rescue")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 4. 创建认领
	claim := &model.RescueClaim{
		RescueID: rescueID,
		OrgID:    org.ID,
		Status:   "claimed",
		Note:     req.Note,
	}

	if err := s.rescueRepo.CreateClaim(claim); err != nil {
		return nil, err
	}

	// 5. 更新救助状态为"响应中"
	if err := s.rescueRepo.UpdateFields(rescueID, map[string]interface{}{
		"status": "responding",
	}); err != nil {
		return nil, err
	}

	return claim, nil
}

// 更新认领进度
func (s *RescueService) UpdateClaim(userID uint, rescueID uint, req *UpdateClaimRequest) error {
	// 1. 验证机构身份
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("you are not an organization")
	}

	// 2. 查找认领记录
	claim, err := s.rescueRepo.FindClaimByRescueAndOrg(rescueID, org.ID)
	if err != nil {
		return errors.New("claim not found")
	}

	// 3. 更新认领状态
	fields := map[string]interface{}{
		"status": req.Status,
	}
	if req.Note != "" {
		fields["note"] = req.Note
	}

	if err := s.rescueRepo.UpdateClaimFields(claim.ID, fields); err != nil {
		return err
	}

	// 4. 如果认领完成，更新救助状态
	if req.Status == "completed" {
		if err := s.rescueRepo.UpdateFields(rescueID, map[string]interface{}{
			"status": "rescued",
		}); err != nil {
			return err
		}
	}

	return nil
}

//救助完成 -> 转为待领养
func (s *RescueService) ConverToPet(userID uint, rescueID uint, req *ConvertToPetRequest) (*model.Pet, error) {
	//1. 验证机构身份
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("you are not an organization")
	}

	//2.验证救助信息
	rescue, err := s.rescueRepo.FindByID(rescueID)
	if err != nil {
		return nil, errors.New("rescue not found")
	}
	if rescue.Status != "rescued" {
		return nil, errors.New("only rescued animals can be converted to pets")
	}

	//3. 验证是否认领救助的机构
	_, err = s.rescueRepo.FindClaimByRescueAndOrg(rescueID, org.ID)
	if err != nil {
		return nil, errors.New("you did not claim this rescue")
	}

	//4.创建记录
	pet := &model.Pet{
		OrgID:        org.ID,
		Name:         req.Name,
		Species:      rescue.Species,
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

	//5.更新救助状态为已关闭
	if err := s.rescueRepo.UpdateFields(rescueID, map[string]interface{}{
		"status": "closed",
	}); err != nil {
		return nil, err
	}
	geo.Remove(geo.KeyRescues, rescueID)
	return s.petRepo.FindByID(pet.ID)
}
