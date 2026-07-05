package controller

import (
	"GinAdmin/internal/service"
	"GinAdmin/internal/validator"
	"GinAdmin/internal/validator/form"

	"github.com/gin-gonic/gin"
)

// 定义 MenuController 的结构
type MenuController struct {
	Api
	svc *service.MenuService
}



func NewMenuController() *MenuController {
	return &MenuController{ // 返回一个 MenuController 的实例
		svc: service.NewMenuService(),
	}
}


// List 菜单列表（树形）
func (ctl *MenuController) List(c *gin.Context) {
	params := form.NewMenuListQuery()
	if err := validator.CheckQueryParams(c, params); err != nil {
		return
	}
	result, err := ctl.svc.List(params)
	if err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, result)
}


// Detail 菜单详情
func (ctl *MenuController) Detail(c *gin.Context) {
	query := form.NewIdForm()
	if err := validator.CheckQueryParams(c, query); err != nil {
		return
	}
	menu, err := ctl.svc.Detail(query.ID)
	if err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, menu)
}

// Create 新增菜单
func (ctl *MenuController) Create(c *gin.Context) {
	params := &form.CreateMenuForm{}
	if err := validator.CheckPostParams(c, params); err != nil {
		return
	}
	menu, err := ctl.svc.Create(params)
	if err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, menu)
}

// Update 更新菜单
func (ctl *MenuController) Update(c *gin.Context) {
	params := &form.UpdateMenuForm{}
	if err := validator.CheckPostParams(c, params); err != nil {
		return
	}
	if err := ctl.svc.Update(params); err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, nil)
}


// Delete 删除菜单
func (ctl *MenuController) Delete(c *gin.Context) {
	params := form.NewIdForm()
	if err := validator.CheckPostParams(c, params); err != nil {
		return
	}
	if err := ctl.svc.Delete(params.ID); err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, nil)
}