package controller

import (
	"GinAdmin/internal/service"
	"GinAdmin/internal/validator"
	"GinAdmin/internal/validator/form"
	"fmt"

	"github.com/gin-gonic/gin"
)



type LoginController struct {
	Api
	authService *service.AuthService
}


func NewLoginController() *LoginController {
	return &LoginController{
		authService: service.NewAuthService(),
	}
}


// Login 用户登录
func (ctl *LoginController) Login(c *gin.Context) {
	// 登录逻辑
	params := &form.LoginForm{}
	if err := validator.CheckPostParams(c, params); err != nil {
		return
	}

	fmt.Println("params", params)

	token, user, err := ctl.authService.Login(params.Username, params.Password)

	if err != nil {
		ctl.Err(c, err)
		return
	}

	ctl.Success(c, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"email":    user.Email,
		},
	})
}


