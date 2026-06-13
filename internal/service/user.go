package service

import (
	"GinAdmin/data"
	"GinAdmin/internal/model"
)



type UserService struct{}


func NewUserService() *UserService {
	return &UserService{}
}


func (s *UserService) ListUsers() ([]model.User, error) {
	var users []model.User

	err := data.GetDB().
		Select("id", "username", "nickname", "email", "status", "created_at", "updated_at").
		Order("id asc").
		Find(&users).Error

	return users, err
}

// GetUserByID
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := data.GetDB().
		Select("id", "username", "nickname", "email", "status", "created_at", "update_at").
		Where("id = ?", id).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}