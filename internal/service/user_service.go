package service

import (
	"errors"

	"github.com/hacker4257/pet_charity/internal/model"
	"github.com/hacker4257/pet_charity/internal/repository"
	"github.com/hacker4257/pet_charity/pkg/event"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo repository.UserRepository
	tokenRepo *repository.TokenRepo
}

func NewUserService(repo repository.UserRepository, tokenRepo *repository.TokenRepo) *UserService {
	return &UserService{
		repo:      repo,
		tokenRepo: tokenRepo,
	}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Nickname string `json:"nickname" binding:"max=50"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	User         model.User `json:"user"`
}

//更新资料请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
	Language string `json:"language" binding:"omitempty,oneof=zh-CN en-US"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=50"`
}

//注册
func (s *UserService) Register(req *RegisterRequest) (*model.User, error) {
	//1.查看用户名是否已存在
	_, err := s.repo.FindByUsername(req.Username)
	if err == nil {
		return nil, errors.New("Username already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	//2.检查邮箱是否存在
	_, err = s.repo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	//密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("password encryption failed")
	}

	//创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Role:     "user",
		Status:   "active",
		Language: "zh-CN",
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(req *LoginRequest) (*model.User, error) {
	//1.先按用户名查，查不到再按邮箱查
	user, err := s.repo.FindByUsername(req.Account)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		user, err = s.repo.FindByEmail(req.Account)
	}
	if err != nil {
		return nil, errors.New("account or password incorrect")
	}

	//2.检查账号状态
	if user.Status != "active" {
		return nil, errors.New("account is disabled")
	}

	//3. 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("account or password incorrect")
	}
	event.Publish("login", user.ID, 0)
	return user, nil
}

func (s *UserService) GetByID(id uint) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) FindByPhoneOrCreateByPhone(phone string) (*model.User, error) {
	//先查找
	user, err := s.repo.FindByPhone(phone)
	if err == nil {
		//找到了，检查状态
		if user.Status != "active" {
			return nil, errors.New("account is disabled")
		}
		return user, nil
	}

	//没有找到，自动注册
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user = &model.User{
		Phone:    phone,
		Username: "user_" + phone[len(phone)-4:],
		Role:     "user",
		Status:   "active",
		Language: "zh-CN",
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

//更新个人资料
func (s *UserService) UpdateProfile(userID uint, req *UpdateProfileRequest) (*model.User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	fields := map[string]interface{}{}

	//只更新非空字段
	if req.Nickname != "" {
		fields["nickname"] = req.Nickname
	}

	if req.Language != "" {
		fields["language"] = req.Language
	}
	//邮箱验证唯一性
	if req.Email != "" && req.Email != user.Email {
		existing, err := s.repo.FindByEmail(req.Email)
		if err == nil && existing.ID != userID {
			return nil, errors.New("email already taken")
		}
		fields["email"] = req.Email
	}
	//手机号唯一性
	if req.Phone != "" && req.Phone != user.Phone {
		existing, err := s.repo.FindByPhone(req.Phone)
		if err == nil && existing.ID != userID {
			return nil, errors.New("phone already taken")
		}
		fields["phone"] = req.Phone
	}

	if len(fields) == 0 {
		return user, nil
	}

	if err := s.repo.UpdateFields(userID, fields); err != nil {
		return nil, errors.New("update failed")
	}
	event.Publish("profile_update", userID, 0)
	return s.repo.FindByID(userID)
}

func (s *UserService) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("old password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("password encryption failed")
	}

	if err := s.repo.UpdateFields(userID, map[string]interface{}{
		"password": string(hashedPassword),
	}); err != nil {
		return err
	}
	s.tokenRepo.IncrVersion(userID)
	return nil
}

func (s *UserService) UpdateAvatar(userID uint, avatarURL string) error {
	return s.repo.UpdateFields(userID, map[string]interface{}{
		"avatar": avatarURL,
	})
}

// Count 统计用户总数
func (s *UserService) Count() (int64, error) {
	return s.repo.Count()
}

// List 用户列表（管理员用）
func (s *UserService) List(page, pageSize int, role string) ([]model.User, int64,
	error) {
	return s.repo.List(page, pageSize, role)
}

// UpdateStatus 修改用户状态
func (s *UserService) UpdateStatus(id uint, status string) error {
	validStatus := map[string]bool{"active": true, "disabled": true}
	if !validStatus[status] {
		return errors.New("invalid status")
	}
	return s.repo.UpdateStatus(id, status)
}

//UpdateRole 更新角色
func (s *UserService) UpdateRole(id uint, role string) error {
	roleMap := map[string]bool{"admin": true, "org": true, "user": true}
	if !roleMap[role] {
		return errors.New("role not postive")
	}
	return s.repo.UpdateRole(id, role)
}
