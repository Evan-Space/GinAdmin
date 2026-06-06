package model

import "time"


type User struct {
	ID        uint64      `json:"id" gorm:"column:id"`
	Username  string    `json:"username" gorm:"column:username"`
	Nickname  string    `json:"nickname" gorm:"column:nickname"`
	Email     string    `json:"email" gorm:"column:email"`
	Status    int       `json:"status" gorm:"column:status"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "users"
}