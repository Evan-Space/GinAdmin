package service

import (
	"GinAdmin/data"
	"GinAdmin/internal/model"
	"GinAdmin/internal/validator/form"

	"gorm.io/gorm"

	"GinAdmin/internal/pkg/errors"
)

type AdminUserService struct{}

func NewAdminUserService() *AdminUserService {
	return &AdminUserService{}
}

/**
* GetUserInfo 根据 ID 获取用户详情
* @description: 获取用户信息
* @param: id uint
* @return: *model.AdminUser, error
 */
func (s *AdminUserService) GetUserInfo(id uint) (*model.AdminUser, error) {
	var user model.AdminUser

	err := data.GetDB().
		Where("id = ? AND deleted_at = 0", id).
		First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

/**
* List 分页查询用户列表
 */
func (s *AdminUserService) List(params *form.AdminUserList) (map[string]interface{}, error) {
	var users []model.AdminUser
	var total int64

	query := data.GetDB().Model(&model.AdminUser{}).Where("deleted_at = 0")

	if params.Username != "" {
		// query = query.Where("username = LIKE", "%"+params.Username+"%")
		query = query.Where("username = ?", params.Username)
	}

	if params.Nickname != "" {
		query = query.Where("nickname = ?", params.Nickname)
	}

	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	if params.Email != "" {
		query = query.Where("email = ?", params.Email)
	}
	if params.Age != "" {
		query = query.Where("age = ?", params.Age)
	}

	// 总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页
	offset := (params.CurrentPage - 1) * params.PageSize
	if err := query.Order("id ASC").Offset(offset).Limit(params.PageSize).Find(&users).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"list":        users,
		"total":       total,
		"currentPage": params.CurrentPage,
		"pageSize":    params.PageSize,
	}, nil
}

/**
* 新增用户
 */
func (s *AdminUserService) Create(params *form.CreateAdminUser) (*model.AdminUser, error) {
	user := model.AdminUser{
		Username:    params.Username,
		Password:    params.Password,
		Nickname:    params.Nickname,
		Email:       params.Email,
		PhoneNumber: params.Phone,
		Status:      1,
	}

	if params.Status != nil {
		user.Status = *params.Status
	}
	if err := data.GetDB().Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

/**
* Update 更新用户
 */
func (s *AdminUserService) Update(params *form.UpdateAdminUser) error {
	updates := map[string]interface{}{}
	if params.Username != "" {
		updates["username"] = params.Username
	}
	if params.Password != "" {
		updates["password"] = params.Password
	}
	if params.Nickname != "" {
		updates["nickname"] = params.Nickname
	}
	if params.Email != "" {
		updates["email"] = params.Email
	}
	if params.Phone != "" {
		updates["phone_number"] = params.Phone
	}
	if params.Status != nil {
		updates["status"] = *params.Status
	}

	if len(updates) == 0 {
		return nil
	}

	return data.GetDB().Model(&model.AdminUser{}).
		Where("id = ? AND deleted_at = 0", params.ID).
		Updates(updates).Error
}

/**
* UpdateProfile 更新个人用户
 */
func (s *AdminUserService) UpdateProfile(uid uint, params *form.UpdateProfile) error {
	updates := map[string]interface{}{}
	if params.Nickname != "" {
		updates["nickname"] = params.Nickname
	}
	if params.Password != "" {
		updates["password"] = params.Password
	}
	if params.Phone != "" {
		updates["phone_number"] = params.Phone
	}
	if params.Email != "" {
		updates["email"] = params.Email
	}
	if len(updates) == 0 {
		return nil
	}

	return data.GetDB().
		Model(&model.AdminUser{}).
		Where("id = ? AND deleted_at = 0", uid).
		Updates(updates).Error
}

/**
* Delete 软删除用户
 */
func (s *AdminUserService) Delete(id uint) error {
	return data.GetDB().
		Model(&model.AdminUser{}).
		Where("id = ? AND deleted_at = 0", id).
		Update("deleted_at", gorm.Expr("UNIX_TIMESTAMP()")).Error
}

// BindRole 为用户绑定角色（全量替换）
func (s *AdminUserService) BindRole(params *form.BindRoleForm) error {
	return data.GetDB().Transaction(func(tx *gorm.DB) error {
		// 删除旧的关联
		if err := tx.Where("uid = ?", params.UserID).Delete(&model.AdminUserRoleMap{}).Error; err != nil {
			return err
		}
		// 插入新的关联
		for _, roleID := range params.RoleIDs {
			m := model.AdminUserRoleMap{Uid: params.UserID, RoleId: roleID}
			if err := tx.Create(&m).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (ctl *AdminUserService) UserNameOptions() ([]map[string]any, error) {
	var users []model.AdminUser

	if err := data.GetDB().
		Model(&model.AdminUser{}).
		Select("id, nickname").
		Where("deleted_at = 0").
		Find(&users).Error; err != nil {
		return nil, errors.NewBusinessError(errors.ServerErr, "error fetching user name options: %s")
	}

	result := make([]map[string]any, 0, len(users))

	for _, u := range users {
		result = append(result, map[string]any{
			"label": u.Nickname,
			"value": u.ID,
		})
	}

	return result, nil
}
