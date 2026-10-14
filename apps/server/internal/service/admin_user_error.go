package service

import "GinAdmin/internal/pkg/errors"

const (
	CodeUsernameExists = 20001 // 用户名已存在
	CodeUserInUse      = 20002 // 用户已绑定角色，不可删除
	CodeUserNotFound   = 20003 // 用户不存在
)

var (
	ErrUsernameExists = errors.NewBusinessError(CodeUsernameExists, "用户名已存在")
	ErrUserInUse      = errors.NewBusinessError(CodeUserInUse, "用户已绑定角色，不可删除")
	ErrUserNotFound   = errors.NewBusinessError(CodeUserNotFound, "用户不存在")
)
