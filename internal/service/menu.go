package service

import (
	"GinAdmin/data"
	"GinAdmin/internal/errors"
	"GinAdmin/internal/model"
	"GinAdmin/internal/validator/form"
	"fmt"

	"gorm.io/gorm"
)

type MenuService struct{}

func NewMenuService() *MenuService {
	return &MenuService{}
}

// List 返回菜单树
// 先查全部，再内存中构建树形结构
func (s *MenuService) List(params *form.MenuListQuery) ([]*model.Menu, error) {
	var menus []model.Menu
	query := data.GetDB().Where("delete_at = 0")

	if params.Keyword != "" {
		query = query.Where("name LIKE ? OR LIKE ?", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}

	if err := query.Order("sort ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, err
	}

	// 构建树： pid = 0 为根结点
	tree := s.buildMenuTree(menus, 0)
	return tree, nil
}

// buildMenuTree 递归构建菜单树
func (s *MenuService) buildMenuTree(all []model.Menu, pid uint) []*model.Menu {
	var tree []*model.Menu
	for i := range all {
		if all[i].Pid == pid {
			children := s.buildMenuTree(all, all[i].ID)
			// 为每个菜单项填充子节点数
			all[i].ChildrenNum = uint(len(children))
			tree = append(tree, &all[i])
			tree = append(tree, children...)
		}
	}
	return tree
}

// Detail 菜单详情
func (s *MenuService) Detail(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := data.GetDB().Where("id = ? AND deleted_at = 0", id).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// 新增菜单
func (s *MenuService) Create(params *form.CreateMenuForm) (*model.Menu, error) {
	menu := model.Menu{
		Code:        params.Code,
		Name:        params.Name,
		Icon:        params.Icon,
		Path:        params.Path,
		Component:   params.Component,
		Type:        params.Type,
		Pid:         params.Pid,
		Sort:        params.Sort,
		Description: params.Description,
		Level:       1,
		IsShow:      1,
		Status:      1,
		IsAuth:      0,
	}

	if params.IsShow != nil {
		menu.IsShow = *params.IsShow
	}

	if params.Status != nil {
		menu.Status = *params.Status
	}

	if params.IsAuth != nil {
		menu.IsAuth = *params.IsAuth
	}


	// 计算层级和 pids
	if params.Pid > 0 {
		var parent model.Menu
		if err := data.GetDB().Where("id = ? AND delete_at = 0", params.Pid).First(&parent).Error; err != nil {
			return nil, err
		}

		menu.Level = parent.Level + 1
		if parent.Pids != "" {
			menu.Pids = parent.Pids + "," + fmt.Sprintf("%d", parent.ID)
		} else {
			menu.Pids = fmt.Sprintf("%d", parent.ID)
		}
	}

	err := data.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&menu).Error; err != nil {
			return err
		}
		// 更新父菜单 children_num
		if params.Pid > 0 {
			tx.Model(&model.Menu{}).Where("id = ?", params.Pid).Update("children_num", gorm.Expr("children_num + 1"))
		}
		return nil
	})

	return &menu, err
}

// update 更新菜单
func (s *MenuService) Update(params *form.UpdateMenuForm) error {
	updates := map[string]interface{}{}
	if params.Code != "" {
		updates["code"] = params.Code
	}
	if params.Name != "" {
		updates["name"] = params.Name
	}
	if params.Icon != "" {
		updates["icon"] = params.Icon
	}
	if params.Path != "" {
		updates["path"] = params.Path
	}
	if params.Component != "" {
		updates["component"] = params.Component
	}
	if params.Type != nil {
		updates["type"] = *params.Type
	}
	if params.Pid != nil {
		updates["pid"] = *params.Pid
	}
	if params.Sort != nil {
		updates["sort"] = *params.Sort
	}
	if params.IsShow != nil {
		updates["is_show"] = *params.IsShow
	}
	if params.Status != nil {
		updates["status"] = *params.Status
	}
	if params.IsAuth != nil {
		updates["is_auth"] = *params.IsAuth
	}
	if params.Description != "" {
		updates["description"] = params.Description
	}
	if len(updates) == 0 {
		return nil
	}

	return data.GetDB().
		Model(&model.Menu{}).
		Where("id = ? AND deleted_at = 0", params.ID).
		Updates(updates).Error
}


// Delete 软删除菜单（检查是否有子节点，同时清理关联）
func (s *MenuService) Delete(id uint) error {
	// 检查子节点
	var count int64
	data.GetDB().Model(&model.Menu{}).Where("pid = ? AND deleted_at = 0", id).Count(&count)
	if count > 0 {
		return errors.NewBusinessError(errors.AuthorizationErr, "菜单下有子菜单，无法删除")
	}

	return data.GetDB().Transaction(func(tx *gorm.DB) error {
		// 软删除
		tx.Model(&model.Menu{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("UNIX_TIMESTAMP()"))
		// 删除菜单-API 关联
		tx.Where("menu_id = ?", id).Delete(&model.MenuApiMap{})
		return nil
	})
}