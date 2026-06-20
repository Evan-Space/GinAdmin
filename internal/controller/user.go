package controller

import (
	"GinAdmin/internal/errors"
	"GinAdmin/internal/service"
	"GinAdmin/internal/validator"
	"GinAdmin/internal/validator/form"
	"fmt"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Api
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

func (ctl *UserController) List(c *gin.Context) {
	params := form.NewPaginate()
	if err := validator.CheckQueryParams(c, params); err != nil {
		return
	}
	users, err := ctl.userService.ListUsers()
	if err != nil {
		ctl.Fail(c, 500, "查询用户列表失败")
		return
	}

	ctl.Success(c, users)
}

// Info 获取当前登录用户信息
func (ctl *UserController) Info(c *gin.Context) {
	fmt.Println("Info")
	fmt.Println("ctl.GetCurrentUserID(c)", ctl.GetCurrentUserID(c))
	uid := ctl.GetCurrentUserID(c)
	fmt.Println("uid", uid)
	if uid == 0 {
		ctl.FailCode(c, errors.NotLogin)
		return
	}

	user, err := ctl.userService.GetUserByID(uid)
	fmt.Println("user", user)
	fmt.Println("err", err)
	if err != nil {
		ctl.FailCode(c, errors.NotFound)
		return
	}

	ctl.Success(c, gin.H{
		"id":             user.ID,
		"username":       user.Username,
		"nickname":       user.Nickname,
		"email":          user.Email,
		"phone_number":   user.PhoneNumber,
		"avatar":         user.Avatar,
		"status":         user.Status,
		"is_super_admin": user.IsSuperAdmin,
		"last_login":     user.LastLogin,
		"created_at":     user.CreatedAt,
	})
}
