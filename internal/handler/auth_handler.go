package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hacker4257/pet_charity/internal/config"
	"github.com/hacker4257/pet_charity/internal/middleware"
	"github.com/hacker4257/pet_charity/internal/repository"
	"github.com/hacker4257/pet_charity/internal/service"
	"github.com/hacker4257/pet_charity/pkg/response"
	"github.com/hacker4257/pet_charity/pkg/utils"
)

type AuthHandler struct {
	userService *service.UserService
	smsService  *service.SmsService
	tokenRepo   *repository.TokenRepo
}

func NewAuthHandler(userService *service.UserService, smsService *service.SmsService, tokenRepo *repository.TokenRepo) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		smsService:  smsService,
		tokenRepo: tokenRepo,
	}
}

//注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, user)
}

// 登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.userService.Login(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	jwtCfg := config.Global.JWT
	ver := h.tokenRepo.GetVersion(user.ID)
	accessToken, err := utils.GenerateToken(user.ID, user.Role, jwtCfg.Secret, jwtCfg.AccessExpire, ver)
	if err != nil {
		response.ServerError(c, "generate token failed")
		return
	}
	refreshToken, err := utils.GenerateToken(user.ID, user.Role, jwtCfg.Secret, jwtCfg.RefreshExpire, ver)
	if err != nil {
		response.ServerError(c, "generate token failed")
		return
	}

	response.Success(c, service.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	})
}

//获取当前登录用户信息
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.userService.GetByID(userID)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}
	response.Success(c, user)
}

//发送验证码
func (h *AuthHandler) SendSmsCode(c *gin.Context) {
	var req service.SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.smsService.SendCode(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, nil)
}

//验证码登录
func (h *AuthHandler) SmsLogin(c *gin.Context) {
	var req service.SmsLoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	//1.验证验证码
	if err := h.smsService.VerifyCode(req.Phone, req.Code, "login"); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	//2.查找或创建用户(手机号不存在则自动注册)
	user, err := h.userService.FindByPhoneOrCreateByPhone(req.Phone)
	if err != nil {
		response.BadRequest(c, "login failed")
		return
	}

	//生成token
	jwtCfg := config.Global.JWT
	ver := h.tokenRepo.GetVersion(user.ID)
	accessToken, err := utils.GenerateToken(user.ID, user.Role, jwtCfg.Secret, jwtCfg.AccessExpire, ver)
	if err != nil {
		response.ServerError(c, "generate token failed")
		return
	}
	refreshToken, err := utils.GenerateToken(user.ID, user.Role, jwtCfg.Secret, jwtCfg.RefreshExpire, ver)
	if err != nil {
		response.ServerError(c, "generate token failed")
		return
	}
	response.Success(c, service.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	})
}
