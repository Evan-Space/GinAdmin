package service

import (
	"GinAdmin/data"
	casbinx "GinAdmin/internal/access/casbin"
	"fmt"
	"log"
)

// Syn AllPolicies 同步所有用户的 Casbin 策略
// 策略格式： p, adminUser:{uid}, /api/v1/xxx, GET
// 调用时机，启动时 + 角色/菜单关联变更时
// 在这里把业务表摊平写入 casbin_rule 表
func SyncAllPolicies() error {
	enforcer, err := casbinx.GetEnforcer()
	if err != nil {
		return err
	}

	db := data.GetDB()

	// 1. 查询所有角色-菜单- API 关联
	// role_menu_map JSON menu_api_map
	type RoleApiRow struct {
		RoleID uint
		Method string
		Route  string
	}

	var roleApis []RoleApiRow
	err = db.Table("role_menu_map").
		Select("role_menu_map.role_id, api.method, api.route").
		Joins("JOIN menu_api_map ON menu_api_map.menu_id = role_menu_map.menu_id").
		Joins("JOIN api On api.id = menu_api_map.api_id").
		Where("api.is_effective = 1").
		Scan(&roleApis).Error

	if err != nil {
		return fmt.Errorf("查询角色API 关联失败：%w", err)
	}

	// 2. 查询所有用户- 角色关联
	type UserRoleRow struct {
		Uid    uint
		RoleID uint
	}

	var userRoles []UserRoleRow
	err = db.Table("admin_user_role_map").
		Select("uid, role_id").
		Scan(&userRoles).Error

	if err != nil {
		return fmt.Errorf("查询用户角色关联失败: %w", err)
	}

	// 3. 构建策略 user -> policies
	// 先建立 roleID -> policies 的映射
	rolePolicyMap := make(map[uint][][]string)
	for _, ra := range roleApis {
		rolePolicyMap[ra.RoleID] = append(rolePolicyMap[ra.RoleID], []string{ra.Route, ra.Method})
	}

	// 再建立 userID → policies 的映射
	userPolicyMap := make(map[uint][][]string)
	for _, ur := range userRoles {
		if policies, ok := rolePolicyMap[ur.RoleID]; ok {
			userPolicyMap[ur.Uid] = append(userPolicyMap[ur.Uid], policies...)
		}
	}

	// 4. 写入 Casbin
	for uid, policies := range userPolicyMap {
		subject := fmt.Sprintf("adminUser:%d", uid)

		// 删除旧策略
		enforcer.DeletePermissionsForUser(subject)

		// 添加新策略
		rules := make([][]string, 0, len(policies))
		for _, p := range policies {
			rules = append(rules, []string{subject, p[0], p[1]})
		}
		if len(rules) > 0 {
			if _, err := enforcer.AddPolicies(rules); err != nil {
				log.Printf("为用户 %s 添加策略失败: %v", subject, err)
			}
		}
	}

	// 5. 重新加载策略到内存
	enforcer.LoadPolicy()
	return nil

}
