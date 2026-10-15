package service

import (
	"GinAdmin/data"
	"GinAdmin/internal/model"
	"GinAdmin/internal/pkg/errors"

	"GinAdmin/internal/validator/form"
	"fmt"

	"gorm.io/gorm"
)


type DeptService struct {}


func NewDeptService() *DeptService {
	return &DeptService{}
}


// List  返回部门树形列表
func (s *DeptService) List(params *form.DeptListQuery) ([]*model.Department, error) {
	var depts []model.Department
	query := data.GetDB().Where("deleted_at = 0")
	if params.Name != "" {
		query = query.Where("name LIKE ?", "%"+params.Name+"%")
	}

	if params.Pid != nil {
		query = query.Where("pid = ?", *params.Pid)
	}

	if err := query.Order("sort ASC, id ASC").Find(&depts).Error; err != nil {
		return nil, err
	}

	return s.buildTree(depts, 0), nil
}


// buildTree 递归构建部门树
func (s *DeptService) buildTree(all []model.Department, pid uint) []*model.Department {
	var tree []*model.Department
	for i := range all {
		if all[i].Pid == pid {
			all[i].ChildrenNum = uint(len(s.buildTree(all, all[i].ID)))
			tree = append(tree, &all[i])
		}
	}
	return tree
}


// Detail 部门详情
func (s *DeptService) Detail(id uint) (*model.Department, error) {
	var dept model.Department
	err := data.GetDB().Where("id = ? AND delete_at = 0", id).First(&dept).Error
	if err != nil {
		return nil, err
	}
	return &dept, nil
}


// Create 新增部门
func (s *DeptService) Create(params *form.CreateDeptForm) (*model.Department, error) {
	dept := model.Department {
		Name: params.Name,
		Pid: params.Pid,
		Description: params.Description,
		Sort: params.Sort,
		Level: 1,
	}
	// 如果有父部门，计算 level 和 pids
	if params.Pid > 0 {
		var parent model.Department
		if err := data.GetDB().Where("id = ? AND deleted_at = 0", params.Pid).First(&parent).Error; err != nil {
			return nil, errors.NewBusinessError(errors.NotFound, "父部门不存在")
		}
		dept.Level = parent.Level + 1
		if parent.Pids != "" {
			dept.Pids = parent.Pids + "," + fmt.Sprintf("%d", parent.ID)
		} else {
			dept.Pids = fmt.Sprintf("%d", parent.ID)
		}
	}

	err := data.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dept).Error; err != nil {
			return err
		}
		if params.Pid > 0 {
			tx.Model(&model.Department{}).Where("id = ?", params.Pid).
			Update("children_num", gorm.Expr("children_num + 1"))
		}
		return nil
	})

	return &dept, err
}


// Update 更新部门
func (s *DeptService) Update(params *form.UpdateDeptForm) error {
	updates := map[string]interface{} {}

	if params.Name != "" {
		updates["name"] = params.Name
	}

	if params.Description != "" {
		updates["description"] = params.Description
	}

	if params.Pid != nil {
		updates["pid"] = *params.Sort
	}

	if params.Sort != nil {
		updates["sort"] = *params.Sort
	}

	if len(updates) == 0 {
		return nil
	}
	return data.GetDB().
		Model(&model.Department{}).
		Where("id = ? AND deleted_at = 0", params.ID).
		Updates(updates).Error
}

// Delete 软删除部门
func (s *DeptService) Delete(id uint) error {
	// 检查子部门
	var count int64
	data.GetDB().Model(&model.Department{}).Where("pid = ? AND deleted_at = 0", id).Count(&count)
	if count > 0 {
		return errors.NewBusinessError(errors.AuthorizationErr, "该部门下有子部门，不能删除")
	}

	// 检查是否是系统部门
	var dept model.Department
	data.GetDB().Where("id = ? AND deleted_at = 0", id).First(&dept)
	if dept.IsSystem == 1 {
		return errors.NewBusinessError(errors.AuthorizationErr, "系统部门不能删除")
	}

	return data.GetDB().Transaction(func(tx *gorm.DB) error {
		tx.Model(&model.Department{}).Where("id = ?", id).
		Update("deleted_at", gorm.Expr("UNIX_TIMESTAMP()"))
		tx.Where("dept_id = ?", id).Delete(&model.Department{})
		return nil
	})
}






// BindRole 部门绑定角色（全量替换）
func (s *DeptService) BindRole(params *form.DeptBindRoleForm) error {
	return data.GetDB().Transaction(func(tx *gorm.DB) error {
		// 删除旧关联
		if err := tx.Where("dept_id = ?", params.DeptId).Delete(&model.DeptRoleMap{}).Error; err != nil {
			return err
		}
		// 插入新关联
		for _, roleID := range params.RoleIds {
			m := model.DeptRoleMap{DeptId: params.DeptId, RoleId: roleID}
			if err := tx.Create(&m).Error; err != nil {
				return err
			}
		}
		return nil
	})
}