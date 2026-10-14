package controller

import (
	"GinAdmin/internal/pkg/errors"
	"GinAdmin/internal/pkg/utils"

	req "GinAdmin/internal/pkg/request"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	Api
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (ctl *AuthController) RefreshToken(c *gin.Context) {
	// 1. 提取当前  Token
	accessToken, err := req.GetAccessToken(c)
	if err != nil || accessToken == "" {
		ctl.FailCode(c, errors.NotLogin)
		return
	}

	// 2. 解析旧 token
	claims, err := utils.ParseToken(accessToken)
	if err != nil {
		ctl.FailCode(c, errors.NotLogin)
		return
	}

	// 3. 签发新的 token
	newToken, err := utils.RefreshToken(claims)
	if err != nil {
		ctl.Err(c, err)
		return
	}

	ctl.Success(c, gin.H{
		"access_token": newToken,
		"token_type":   "Bearer",
	})

}

// CheckToken 检查 token 是否有效
func (ctl *AuthController) CheckToken(c *gin.Context) {
	accessToken, err := req.GetAccessToken(c)
	if err != nil || accessToken == "" {
		ctl.Success(c, gin.H{"valid": false})
		return
	}

	_, err = utils.ParseToken(accessToken)
	ctl.Success(c, gin.H{"valid": err == nil})
}
