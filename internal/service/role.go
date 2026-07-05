package service

import (
	"GinAdmin/data"
	"GinAdmin/internal/errors"
	"GinAdmin/internal/model"
	"GinAdmin/internal/validator/form"
	"fmt"

	"gorm.io/gorm"
)

type RoleService struct{}

func NewRoleService() *RoleService {
	return &RoleService{}
}

// List 角色列表
func (s *RoleService) List(params *form.RoleListQuery) (map[string]interface{}, error) {
	var roles []model.Role
	var total int64

	query := data.GetDB().Model(&model.Role{}).Where("deleted_at = 0")

	if params.Name != "" {
		query = query.Where("name LIKE ?", "%"+params.Name+"%")
	}
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}

	if params.Pid != nil {
		query = query.Where("pid = ?", *params.Pid)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PerPage
	query = query.Order("sort ASC, id ASC").Offset(offset).Limit(params.PerPage)
	if err := query.Find(&roles).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"list":     roles,
		"total":    total,
		"page":     params.Page,
		"per_page": params.PerPage,
	}, nil
}

// 角色详情
func (s *RoleService) Detail(id uint) (*model.Role, error) {
	var role model.Role
	err := data.GetDB().Where("id = ? AND deleted_at = 0", id).First(&role).Error
	if err != nil {
		return nil, err
	}

	return &role, nil
}

// Create 新增角色
func (s *RoleService) Create(params *form.CreateRoleForm) (*model.Role, error) {
	role := model.Role{
		Code:        params.Code,
		Name:        params.Name,
		Description: params.Description,
		Pid:         params.Pid,
		Sort:        params.Sort,
		Status:      1,
		Level:       1,
	}

	// 如果有父角色，计算 level 和 pids
	if params.Pid > 0 {
		var parent model.Role
		if err := data.GetDB().Where("id = ? AND deleted_at = 0", params.Pid).First(&parent).Error; err != nil {
			return nil, err
		}

		role.Level = parent.Level + 1

		if parent.Pids != "" {
			role.Pids = parent.Pids + "," + fmt.Sprintf("%d", parent.ID)
		} else {
			role.Pids = fmt.Sprintf("%d", parent.ID)
		}
	}

	err := data.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&role).Error; err != nil {
			return err
		}
		// 更新父角色的 children_num
		if params.Pid > 0 {
			tx.Model(&model.Role{}).Where("id = ?", params.Pid).Update("children_num", gorm.Expr("children_num + 1"))
		}
		return nil
	})

	return &role, err
}

// Update 更新角色
func (s *RoleService) Update(params *form.UpdateRoleForm) error {
	updates := map[string]interface{}{}
	if params.Code != "" {
		updates["code"] = params.Code
	}
	if params.Name != "" {
		updates["name"] = params.Name
	}
	if params.Description != "" {
		updates["description"] = params.Description
	}
	if params.Pid != nil {
		updates["pid"] = *params.Pid
	}
	if params.Sort != nil {
		updates["sort"] = *params.Sort
	}
	if params.Status != nil {
		updates["status"] = *params.Status
	}
	if len(updates) == 0 {
		return nil
	}
	return data.GetDB().
		Model(&model.Role{}).
		Where("id = ? AND deleted_at = 0", params.ID).
		Updates(updates).Error
}


// Delete 软删除角色（同时清理关联）
func (s *RoleService) Delete(id uint) error {
	// is_system 角色不允许删除
	var role model.Role
	if err := data.GetDB().Where("id = ? AND deleted_at = 0", id).First(&role).Error; err != nil {
		return err
	}
	if role.IsSystem == 1 {
		return errors.NewBusinessError(errors.AuthorizationErr)
	}

	return data.GetDB().Transaction(func(tx *gorm.DB) error {
		// 软删除角色
		tx.Model(&model.Role{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("UNIX_TIMESTAMP()"))
		// 删除角色-用户关联
		tx.Where("role_id = ?", id).Delete(&model.AdminUserRoleMap{})
		// 删除角色-菜单关联
		tx.Where("role_id = ?", id).Delete(&model.RoleMenuMap{})
		return nil
	})
}