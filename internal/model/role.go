package model

type Role struct {
	ContainsDeleteBaseModel
	Code        string `json:"code" gorm:"column:code;type:varchar(60);not null;default:''"`                // 角色业务编码
	IsSystem    uint8  `json:"is_system" gorm:"column:is_system;type:tinyint unsigned;not null;default:0"`   // 是否系统保留对象
	Pid         uint   `json:"pid" gorm:"column:pid;type:int unsigned;not null;default:0"`                   // 上级id
	Pids        string `json:"pids" gorm:"column:pids;type:varchar(255);not null;default:''"`                // 所有上级id路径
	Name        string `json:"name" gorm:"column:name;type:varchar(60);not null;default:''"`                 // 角色名称
	Description string `json:"description" gorm:"column:description;type:varchar(255);not null;default:''"`  // 描述
	Level       uint8  `json:"level" gorm:"column:level;type:tinyint unsigned;not null;default:1"`           // 层级
	Sort        uint   `json:"sort" gorm:"column:sort;type:mediumint unsigned;not null;default:0"`           // 排序
	ChildrenNum uint   `json:"children_num" gorm:"column:children_num;type:int unsigned;not null;default:0"` // 子集数量
	Status      uint8  `json:"status" gorm:"column:status;type:tinyint unsigned;not null;default:1"`         // 状态 1启用 0禁用
}

// TableName 指定表名
func (Role) TableName() string {
	return "role"
}