package controller

import (
	"GinAdmin/internal/pkg/response"

	"github.com/gin-gonic/gin"
)



type LoginController struct {}


func NewLoginController() *LoginController {
	return &LoginController{}
}


func (c *LoginController) Login(ctx *gin.Context) {
	// 登录逻辑
	response.Ok(ctx, gin.H{"msg": "ok"})
}


