package controller

import (
	"GinAdmin/data"
	"GinAdmin/internal/model"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	Api
}

// NewDashboardController 创建Dashboard控制器实例
func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

// Overview 仪表盘概览
func (ctl *DashboardController) Overview(c *gin.Context) {
	var userCount, roleCount, menuCount, deptCount int64

	data.GetDB().Model(&model.AdminUser{}).Where("delete_at = 0").Count(&userCount)
	data.GetDB().Model(&model.Role{}).Where("delete_at = 0").Count(&roleCount)
	data.GetDB().Model(&model.Menu{}).Where("delete_at = 0").Count(&menuCount)
	data.GetDB().Model(&model.Department{}).Where("delete_at = 0").Count(&deptCount)

	ctl.Success(c, gin.H{
		"user_count": userCount,
		"role_count": roleCount,
		"menu_count": menuCount,
		"dept_count": deptCount,
	})

}
