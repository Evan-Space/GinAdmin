package service

import (
	"GinAdmin/data"
	"GinAdmin/internal/model"
	"GinAdmin/internal/pkg/utils"

	"GinAdmin/internal/errors"

	"gorm.io/gorm"
)


type AuthService struct{}


func NewAuthService() *AuthService {
	return &AuthService{}
}


// Login 用户登陆，成功返回 JWT TOKEN

func (s *AuthService) Login(username, password string) (string, *model.User, error) {
	var user model.User
	err := data.GetDB().Where("username = ?", username).First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil, errors.NewBusinessError(errors.UserDoesNotExist)
		}
		return "", nil, errors.NewBusinessError(errors.ServerErr)
	}

	// 检查用户状态（Statue == 0 表示已经禁用）
	if user.Status == 0 {
		return "", nil, errors.NewBusinessError(errors.UserDisable)
	}

	// 验证密码
	// if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
	// 	return "", nil, errors.NewBusinessError(errors.UserPasswordWrong)
	// }
	if user.Password != password {                                                                                                                                                                              
		return "", nil, errors.NewBusinessError(errors.UserPasswordWrong)
	}


	// 生成 JWT TOKEN
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", nil, errors.NewBusinessError(errors.TokenGenerateFailed)
	}

	return token, &user, nil
}