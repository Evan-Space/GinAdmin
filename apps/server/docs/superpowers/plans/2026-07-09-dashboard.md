# Dashboard 仪表盘 — 实现文档

> 参照 gin-layout 仪表盘实现，简化版：统计系统核心实体数量。
> gin-layout 完整版依赖 `admin_login_logs`、`request_logs` 等审计日志表，当前阶段先做基础统计。

---

## 接口

**路由：** `GET /api/v1/dashboard/overview`

**认证：** `AuthLogin`（只需登录，无需权限）

**返回示例：**

```json
{
  "code": 0,
  "msg": "OK",
  "data": {
    "user_count": 3,
    "role_count": 2,
    "menu_count": 5,
    "dept_count": 0
  }
}
```

---

## 步骤 1：新建 `internal/controller/dashboard.go`

```go
package controller

import (
	"GinAdmin/data"
	"GinAdmin/internal/model"

	"github.com/gin-gonic/gin"
)

// DashboardController 仪表盘控制器
type DashboardController struct {
	Api
}

// NewDashboardController 创建仪表盘控制器实例
func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

// Overview 仪表盘概览
// 返回系统的核心统计数据：用户数、角色数、菜单数、部门数
func (ctl *DashboardController) Overview(c *gin.Context) {
	var userCount, roleCount, menuCount, deptCount int64

	// Count 是 GORM 的链式方法，等价于 SELECT COUNT(*) FROM table WHERE ...
	data.GetDB().Model(&model.AdminUser{}).Where("deleted_at = 0").Count(&userCount)
	data.GetDB().Model(&model.Role{}).Where("deleted_at = 0").Count(&roleCount)
	data.GetDB().Model(&model.Menu{}).Where("deleted_at = 0").Count(&menuCount)
	data.GetDB().Model(&model.Department{}).Where("deleted_at = 0").Count(&deptCount)

	ctl.Success(c, gin.H{
		"user_count": userCount,
		"role_count": roleCount,
		"menu_count": menuCount,
		"dept_count": deptCount,
	})
}
```

**代码说明：**

| 行 | 说明 |
|----|------|
| `&model.AdminUser{}` | 告诉 GORM 操作 `admin_user` 表 |
| `Where("deleted_at = 0")` | 过滤掉软删除的记录 |
| `Count(&userCount)` | 执行 `SELECT COUNT(*)`，结果写入 `userCount` 变量 |
| `gin.H{...}` | `gin.H` 是 `map[string]interface{}` 的别名，`gin.H{"user_count": userCount}` 等价于 `{"user_count": 3}` |

---

## 步骤 2：修改 `internal/routers/router.go`

在 `AdminRouteTree()` 的 `Children` 中追加仪表盘路由。放在登录接口下面、管理员路由上面：

```go
func AdminRouteTree() RouteGroupDef {
	loginCtrl     := controller.NewLoginController()
	dashboardCtrl := controller.NewDashboardController()   // ← 新加
	adminUserCtrl := controller.NewAdminUserController()
	// ... 其他 controller ...

	return RouteGroupDef{
		Prefix: "api/v1",
		Children: []RouteGroupDef{
			// 公开接口
			{
				Routes: []RouteDef{
					POST("login", "登录", AuthNone, loginCtrl.Login),
				},
			},
			// 只需登录（仪表盘） ← 新加整个分组
			{
				Routes: []RouteDef{
					GET("dashboard/overview", "仪表盘概览", AuthLogin, dashboardCtrl.Overview),
				},
			},
			// 需要登录（个人资料）
			{
				Prefix: "admin-user",
				Routes: []RouteDef{
					// ... 保持不变 ...
				},
			},
			// 需要权限
			{
				Middleware: []gin.HandlerFunc{middleware.AdminAuthHandler()},
				Children: []RouteGroupDef{
					// ... 保持不变 ...
				},
			},
		},
	}
}
```

---

## 新建/修改文件汇总

| 文件 | 操作 | 内容 |
|------|------|------|
| `internal/controller/dashboard.go` | **新建** | 仪表盘控制器 + Overview 方法 |
| `internal/routers/router.go` | 修改 | 追加仪表盘路由 |

---

## 测试验证

```bash
# 1. 编译
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go build ./...

# 2. 登录获取 token
TOKEN=$(curl -s -X POST http://127.0.0.1:8080/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"super_admin","password":"123456"}' | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)

# 3. 访问仪表盘
curl http://127.0.0.1:8080/api/v1/dashboard/overview -H "Authorization: Bearer $TOKEN"
```
