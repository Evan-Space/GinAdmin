package controller

import (
	"GinAdmin/internal/service"
	"GinAdmin/internal/validator"
	"GinAdmin/internal/validator/form"

	"github.com/gin-gonic/gin"
)




type DeptController struct {
	Api
	svc *service.DeptService
}

func NewDeptController() *DeptController {
	return &DeptController{svc: service.NewDeptService()}
}

// List 部门列表
func (ctl *DeptController) List(c *gin.Context) {
	params := form.NewDeptListQuery()
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


// Detail 部门详情
func (ctl *DeptController) Detail(c *gin.Context) {
	query := form.NewIdForm()
	if err := validator.CheckQueryParams(c, query); err != nil {
		return
	}

	dept, err := ctl.svc.Detail(query.ID)
	if err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, dept)
}

// Create 新增部门
func (ctl *DeptController) Create(c *gin.Context) {
	params := &form.CreateDeptForm{}
	if err := validator.CheckPostParams(c, params); err != nil {
		return
	}
	dept, err := ctl.svc.Create(params)
	if err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, dept)
}


// Update 更新部门
func (ctl *DeptController) Update(c *gin.Context) {
	params := &form.UpdateDeptForm{}
	if err := validator.CheckPostParams(c, params); err != nil {
		return
	}
	if err := ctl.svc.Update(params); err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, nil)
}

// Delete 删除部门
func (ctl *DeptController) Delete(c *gin.Context) {
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

// BindRole 部门绑定角色
func (ctl *DeptController) BindRole(c *gin.Context) {
	params := &form.DeptBindRoleForm{}
	if err := validator.CheckPostParams(c, params); err != nil {
		return
	}
	if err := ctl.svc.BindRole(params); err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, nil)
}