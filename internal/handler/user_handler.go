package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/hacker4257/pet_charity/internal/middleware"
	"github.com/hacker4257/pet_charity/internal/service"
	"github.com/hacker4257/pet_charity/pkg/response"
	"github.com/hacker4257/pet_charity/pkg/upload"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

//获取个人信息
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.userService.GetByID(userID)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	response.Success(c, user)
}

//更新个人资料
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, user)
}

//修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.ChangePassword(userID, &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *UserHandler) Updateavatar(c *gin.Context) {
	userID := middleware.GetUserID(c)

	//1.获取上传文件
	file, err := c.FormFile("avatar")
	if err != nil {
		response.BadRequest(c, "no file uploaded")
		return
	}

	//2.验证文件
	if err := upload.Validate(file); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	//3.保存文件
	url, err := upload.SaveFile(file, "avatar")
	if err != nil {
		response.ServerError(c, "upload failed")
		return
	}

	//4.更新数据库
	if err := h.userService.UpdateAvatar(userID, url); err != nil {
		response.BadRequest(c, "update avatar failed")
		return
	}

	response.Success(c, gin.H{"avatar": url})
}
