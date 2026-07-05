package controller

import (
	"GinAdmin/internal/errors"
	"GinAdmin/internal/service"
	"GinAdmin/internal/validator"
	"GinAdmin/internal/validator/form"
	"fmt"

	"github.com/gin-gonic/gin"
)

type AdminUserController struct {
	Api
	svc *service.AdminUserService
}

func NewAdminUserController() *AdminUserController {
	return &AdminUserController{
		svc: service.NewAdminUserService(),
	}
}

/**
* 获取当前登录用户的信息
 */
func (ctl *AdminUserController) GetUserInfo(c *gin.Context) {
	uid := ctl.GetCurrentUserID(c)
	if uid == 0 {
		ctl.FailCode(c, errors.NotLogin)
		return
	}

	fmt.Println("uid", uid)

	user, err := ctl.svc.GetUserInfo(uid)
	fmt.Println("user", user)
	if err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, user)
}

/**
* UpdateProfile 更新个人资料
 */
func (ctl *AdminUserController) UpdateProfile(c *gin.Context) {
	uid := ctl.GetCurrentUserID(c)
	if uid == 0 {
		ctl.FailCode(c, errors.NotLogin)
		return
	}
	params := &form.UpdateProfile{}
	if err := validator.CheckPostParams(c, params); err != nil {
		return
	}

	if err := ctl.svc.UpdateProfile(uid, params); err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, nil)
}

// List 用户列表
func (ctl *AdminUserController) List(c *gin.Context) {
	params := &form.AdminUserList{ 
		Paginate: form.Paginate{ Page: 1, PerPage: 10 },
	 }
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

// Detail 用户详情
func (ctl *AdminUserController) Detail(c *gin.Context) {
	query := form.NewIdForm()
	if err := validator.CheckQueryParams(c, query); err != nil {
		return
	}
	user, err := ctl.svc.GetUserInfo(query.ID)
	if err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, user)
}

// Create 新增用户
func (ctl *AdminUserController) Create(c *gin.Context) {
	params := &form.CreateAdminUser{}
	if err := validator.CheckPostParams(c, params); err != nil {
		return
	}
	user, err := ctl.svc.Create(params)
	if err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, user)

}

// Update 更新用户
func (ctl *AdminUserController) Update(c *gin.Context) {
	params := &form.UpdateAdminUser{}
	if err := validator.CheckPostParams(c, params); err != nil {
		return
	}
	if err := ctl.svc.Update(params); err != nil {
		ctl.Err(c, err)
		return
	}

	ctl.Success(c, nil)
}

// Delete 删除用户（软删除）
func (ctl *AdminUserController) Delete(c *gin.Context) {
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


// BindRole 为用户绑定角色
func (ctl *AdminUserController) BindRole(c *gin.Context) {
	params := &form.BindRoleForm{}
	if err := validator.CheckPostParams(c, params); err != nil {
		return
	}
	if err := ctl.svc.BindRole(params); err != nil {
		ctl.Err(c, err)
		return
	}
	ctl.Success(c, nil)
}
