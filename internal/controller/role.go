package controller

import (
	"GinAdmin/internal/service"
	"GinAdmin/internal/validator"
	"GinAdmin/internal/validator/form"

	"github.com/gin-gonic/gin"
)




type RoleController struct {
	Api
	svc *service.RoleService
}

func NewRoleController() *RoleController {
	return &RoleController{
		svc: service.NewRoleService(),
	}
}


// List 角色列表
func (ctl *RoleController) List(c *gin.Context) {
	params := form.NewRoleListQuery()
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


// Detail 角色详情
func (ctl *RoleController) Detail(c *gin.Context) {
	query := form.NewIdForm()
	if err := validator.CheckQueryParams(c, query); err != nil {
		return
	}
	role, err := ctl.svc.Detail(query.ID)
	if err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, role)
}

// Create 新增角色
func (ctl *RoleController) Create(c *gin.Context) {
	params := &form.CreateRoleForm{}
	if err := validator.CheckPostParams(c, params); err != nil {
		return
	}
	role, err := ctl.svc.Create(params)
	if err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, role)
}


// Update 更新角色
func (ctl *RoleController) Update(c *gin.Context) {
	params := &form.UpdateRoleForm{}
	if err := validator.CheckPostParams(c, params); err != nil {
		return
	}
	if err := ctl.svc.Update(params); err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, nil)
}

// Delete 删除角色
func (ctl *RoleController) Delete(c *gin.Context) {
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