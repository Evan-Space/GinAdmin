package model

// MenuApiMap 菜单-API关联表，对应数据库 menu_api_map 表
type MenuApiMap struct {
	BaseModel
	MenuId uint `json:"menu_id" gorm:"column:menu_id;type:int unsigned;not null;default:0"` // 菜单id
	ApiId  uint `json:"api_id" gorm:"column:api_id;type:int unsigned;not null;default:0"`   // 接口id
}

// TableName 指定表名
func (MenuApiMap) TableName() string {
	return "menu_api_map"
}