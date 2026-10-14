package model

const (
	MenuTypeCatalogue uint8 = 1 // 目录
	MenuTypeMenu      uint8 = 2 // 菜单
	MenuTypeButton    uint8 = 3 // 按钮
)

// Menu 菜单权限表，对应数据库 menu 表
// 对应 gin-layout internal/model/menu.go Menu
type Menu struct {
	ContainsDeleteBaseModel
	Icon            string  `json:"icon" gorm:"column:icon;type:varchar(255);not null;default:''"`                              // 图标
	Code            string  `json:"code" gorm:"column:code;type:varchar(120);not null;default:''"`                              // 前端权限标识
	Path            string  `json:"path" gorm:"column:path;type:varchar(255);not null;default:''"`                              // 前端路由路径
	FullPath        string  `json:"full_path" gorm:"column:full_path;type:varchar(255);not null;default:''"`                    // 完整路由路径
	Redirect        string  `json:"redirect" gorm:"column:redirect;type:varchar(255);not null;default:''"`                      // 重定向路由名称
	Name            string  `json:"name" gorm:"column:name;type:varchar(255);not null;default:''"`                              // 路由名称
	Component       string  `json:"component" gorm:"column:component;type:varchar(255);not null;default:''"`                    // 组件路径
	AnimateEnter    string  `json:"animate_enter" gorm:"column:animate_enter;type:varchar(60);not null;default:''"`             // 进入动画
	AnimateLeave    string  `json:"animate_leave" gorm:"column:animate_leave;type:varchar(60);not null;default:''"`             // 离开动画
	AnimateDuration float32 `json:"animate_duration" gorm:"column:animate_duration;type:float(2,2);not null;default:0.00"`      // 动画时长
	IsShow          uint8   `json:"is_show" gorm:"column:is_show;type:tinyint unsigned;not null;default:1"`                     // 是否显示 1是 0否
	Status          uint8   `json:"status" gorm:"column:status;type:tinyint unsigned;not null;default:1"`                       // 状态 1启用 0禁用
	IsAuth          uint8   `json:"is_auth" gorm:"column:is_auth;type:tinyint unsigned;not null;default:0"`                     // 是否需要授权 1是 0否
	IsExternalLinks uint8   `json:"is_external_links" gorm:"column:is_external_links;type:tinyint unsigned;not null;default:0"` // 是否外链
	IsNewWindow     uint8   `json:"is_new_window" gorm:"column:is_new_window;type:tinyint unsigned;not null;default:0"`         // 是否新窗口打开
	Sort            uint    `json:"sort" gorm:"column:sort;type:int unsigned;not null;default:0"`                               // 排序
	Type            uint8   `json:"type" gorm:"column:type;type:tinyint unsigned;not null;default:1"`                           // 菜单类型 1目录 2菜单 3按钮
	Pid             uint    `json:"pid" gorm:"column:pid;type:int unsigned;not null;default:0"`                                 // 上级菜单id
	Level           uint8   `json:"level" gorm:"column:level;type:tinyint unsigned;not null;default:0"`                         // 层级
	Pids            string  `json:"pids" gorm:"column:pids;type:varchar(255);not null;default:''"`                              // 层级序列逗号分隔
	ChildrenNum     uint    `json:"children_num" gorm:"column:children_num;type:int unsigned;not null;default:0"`               // 子集数量
	Description     string  `json:"description" gorm:"column:description;type:varchar(255);not null;default:''"`                // 描述
}


// TableName 指定表名
func (Menu) TableName() string {
	return "menu"
}