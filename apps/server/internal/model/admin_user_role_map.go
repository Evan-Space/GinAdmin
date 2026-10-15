package model

// AdminUserRoleMap 用户-角色关联表，对应数据库 admin_user_role_map 表
type AdminUserRoleMap struct {
	BaseModel
	Uid    uint `json:"uid" gorm:"column:uid;type:int unsigned;not null;default:0"`         // admin_user表id
	RoleId uint `json:"role_id" gorm:"column:role_id;type:int unsigned;not null;default:0"` // 角色id
}

// TableName 指定表名
func (AdminUserRoleMap) TableName() string {
	return "admin_user_role_map"
}