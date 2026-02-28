package service

import (
	"errors"
	"time"

	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/hacker4257/pet_charity/internal/repository"
	"github.com/hacker4257/pet_charity/pkg/event"
	"gorm.io/gorm"
)

type AdoptionService struct {
	adoptionRepo repository.AdoptionRepository
	petRepo      repository.PetRepository
	orgRepo      repository.OrgRepository
	txManager    repository.TransactionManager // 新增
}

func NewAdoptionService(
	adRepo repository.AdoptionRepository,
	petRepo repository.PetRepository,
	orgRepo repository.OrgRepository,
	txManager repository.TransactionManager,
) *AdoptionService {
	return &AdoptionService{
		adoptionRepo: adRepo,
		petRepo:      petRepo,
		orgRepo:      orgRepo,
		txManager:    txManager,
	}
}

//提交领养申请
type CreateAdoptionRequest struct {
	PetID           uint   `json:"pet_id" binding:"required"`
	Reason          string `json:"reason" binding:"required"`
	LivingCondition string `json:"living_condition" binding:"required,max=200"`
	Experience      string `json:"experience"`
}

//审核领养申请
type ReviewAdoptionRequest struct {
	Status       string `json:"status" binding:"required,oneof=approved rejected"`
	RejectReason string `json:"reject_reason"`
}

func (s *AdoptionService) Create(userID uint, req *CreateAdoptionRequest) (*model.Adoption, error) {
	//1.检查宠物是否存在且可领养
	pet, err := s.petRepo.FindByID(req.PetID)
	if err != nil {
		return nil, errors.New("pet not found")
	}
	if pet.Status != "available" {
		return nil, errors.New("this pet is not available for adoption")
	}

	//2. 不能领养自己机构的宠物
	org, err := s.orgRepo.FindByUserID(userID)
	if err == nil && org.ID == pet.OrgID {
		return nil, errors.New("you cannot adopt your own organization's pet")
	}

	//3. 检查是否已有pending 申请
	_, err = s.adoptionRepo.FindPendingByUserAndPet(userID, req.PetID)
	if err == nil {
		return nil, errors.New("you already have pending application for this pet")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	//4.创建申请
	adoption := &model.Adoption{
		UserID:          userID,
		PetID:           req.PetID,
		Reason:          req.Reason,
		LivingCondition: req.LivingCondition,
		Experience:      req.Experience,
		Status:          "pending",
	}

	if err := s.adoptionRepo.Create(adoption); err != nil {
		return nil, err
	}
	event.Publish("adoption", userID, 0)
	return s.adoptionRepo.FindByID(adoption.ID)
}

//审核领养申请
func (s *AdoptionService) Review(userID uint, adoptionID uint, req *ReviewAdoptionRequest) error {
	// 1.查找申请
	adoption, err := s.adoptionRepo.FindByID(adoptionID)
	if err != nil {
		return errors.New("adoption not found")
	}
	if adoption.Status != "pending" {
		return errors.New("this application has already been reviewed")
	}

	//2.验证是否是该宠物所属机构的用户
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("you are not an organization")
	}
	if org.ID != adoption.Pet.OrgID {
		return errors.New("this pet does not belong to your organization")
	}

	now := time.Now()
	fields := map[string]interface{}{
		"status":      req.Status,
		"reviewed_by": userID,
		"reviewed_at": now,
	}

	if req.Status == "rejected" {
		fields["reject_reason"] = req.RejectReason
	}

	//3.事务：更新申请状态 + 更新宠物状态
	err = s.txManager.Transaction(func(tx repository.TransactionContext) error {
		if err := tx.AdoptionRepo().UpdateFields(adoptionID, fields); err != nil {
			return err
		}
		if req.Status == "approved" {
			if err := tx.PetRepo().UpdateFields(adoption.PetID, map[string]interface{}{
				"status": "reserved",
			}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	// 事务成功后才发事件
	if req.Status == "approved" {
		event.Publish("adoption_approved", adoption.UserID, adoptionID)
	} else {
		event.Publish("adoption_rejected", adoption.UserID, adoptionID)
	}
	return nil
}

//确认领养完成
func (s *AdoptionService) Complete(userID uint, adoptionID uint) error {
	adoption, err := s.adoptionRepo.FindByID(adoptionID)
	if err != nil {
		return errors.New("adoption not found")
	}
	if adoption.Status != "approved" {
		return errors.New("only approved applications can be comleted")
	}
	org, err := s.orgRepo.FindByUserID(userID)
	if err != nil || org.ID != adoption.Pet.OrgID {
		return errors.New("permission denied")
	}

	//事务

	return s.txManager.Transaction(func(tx repository.TransactionContext) error {

		if err := tx.AdoptionRepo().UpdateFields(adoptionID, map[string]interface{}{
			"status": "completed",
		}); err != nil {
			return err
		}
		if err := tx.PetRepo().UpdateFields(adoption.PetID, map[string]interface{}{
			"status": "adopted",
		}); err != nil {
			return err
		}
		return nil
	})
}

//申请详情
func (s *AdoptionService) GetByID(userID uint, adoptionID uint) (*model.Adoption, error) {
	adoption, err := s.adoptionRepo.FindByID(adoptionID)
	if err != nil {
		return nil, errors.New("adoption not found")
	}

	//只有申请人或者对于机构可以查看
	if adoption.UserID != userID {
		org, err := s.orgRepo.FindByUserID(userID)
		if err != nil || org.ID != adoption.Pet.OrgID {
			return nil, errors.New("permission denied")
		}
	}

	return adoption, nil
}

func (s *AdoptionService) ListByUser(userID uint, page, pageSize int) ([]model.Adoption, int64, error) {
	return s.adoptionRepo.ListByUser(userID, page, pageSize)
}

func (s *AdoptionService) ListByOrg(usrID uint, status string, page, pageSize int) ([]model.Adoption, int64, error) {
	org, err := s.orgRepo.FindByUserID(usrID)
	if err != nil {
		return nil, 0, errors.New("you are not an organization")
	}
	return s.adoptionRepo.ListByOrg(org.ID, status, page, pageSize)
}
