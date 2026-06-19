package model

import "time"

const (
	AdminUserStatusEnabled  uint8 = 1 // 启用
	AdminUserStatusDisabled uint8 = 0 // 禁用（数据库定义：1启用 0禁用）
)

// AdminUser 管理员用户表 对应数据库， admin_user 表
// 对应 gin-layout internal/model/admin_user.go AdminUser

type AdminUser struct {
	ContainsDeleteBaseModel
	IsSuperAdmin    uint8     `json:"is_super_admin" gorm:"column:is_super_admin;type:tinyint(1);not null;default:0"`     // 是否超级管理员
	Nickname        string    `json:"nickname" gorm:"column:nickname;type:varchar(30);not null;default:''"`               // 昵称
	Username        string    `json:"username" gorm:"column:username;type:varchar(30);not null;default:''"`               // 用户名
	Password        string    `json:"-" gorm:"column:password;type:varchar(255);not null;default:''"`                     // 密码（json:"−" 不序列化到响应体）
	PhoneNumber     string    `json:"phone_number" gorm:"column:phone_number;type:varchar(15);not null;default:''"`       // 手机号
	FullPhoneNumber string    `json:"full_phone_number" gorm:"column:full_phone_number;type:varchar(20);not null;default:''"` // 完整手机号
	CountryCode     string    `json:"country_code" gorm:"column:country_code;type:varchar(10);not null;default:''"`       // 国际区号
	Email           string    `json:"email" gorm:"column:email;type:varchar(120);not null;default:''"`                    // 邮箱
	Avatar          string    `json:"avatar" gorm:"column:avatar;type:varchar(255);not null;default:''"`                  // 头像
	Status          uint8     `json:"status" gorm:"column:status;type:tinyint(1);not null;default:1"`                     // 状态 1启用 0禁用
	LastLogin       time.Time `json:"last_login" gorm:"column:last_login;type:datetime"`                                  // 最后登录时间
	LastIp          string    `json:"last_ip" gorm:"column:last_ip;type:varchar(45);not null;default:''"`                 // 最后登录IP
}



// TableName 指定表名
func (AdminUser) TableName() string{
	return "admin_user"
}