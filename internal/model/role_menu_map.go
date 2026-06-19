package model

// RoleMenuMap 角色-菜单关联表，对应数据库 role_menu_map 表
type RoleMenuMap struct {
	BaseModel
	RoleId uint `json:"role_id" gorm:"column:role_id;type:int unsigned;not null;default:0"` // 角色id
	MenuId uint `json:"menu_id" gorm:"column:menu_id;type:int unsigned;not null;default:0"` // 菜单id
}

// TableName 指定表名
func (RoleMenuMap) TableName() string {
	return "role_menu_map"
}