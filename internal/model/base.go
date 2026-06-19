package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// BaseModel 提供模型通用字段 id + created_at + updated_at
// 对应 gin-layout internal/model/base.go BaseModel
type BaseModel struct {
	ID        uint      `json:"id" gorm:"column:id"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}


// ContainsDeleteBaseModel 在 BaseModel 基础上增加软删除字段
// 对应 gin-layout internal/model/base.go ContainsDeleteBaseModel
type ContainsDeleteBaseModel struct {
	BaseModel
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;type:int(11) unsigned;not null;default:0;index;" json:"-"`
}