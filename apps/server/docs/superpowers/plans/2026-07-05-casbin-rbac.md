# Casbin RBAC 权限控制 — 完整实现文档

> 参照 gin-layout 实现 Casbin RBAC 权限控制：用户 → 角色 → 菜单 → API 权限链。
> 简化掉 Redis 缓存 + 权限变更传播，先做最核心的 enforcer + 中间件 + 策略同步。

---

## 架构概览

```
请求进入
  │
  ├─ ParseToken (全局中间件，静默解析 JWT，写入 uid 到 context)
  │
  ├─ AdminAuthHandler (路由组中间件，只对 AuthPerm 路由生效)
  │     │
  │     ├─ uid=0 → 401 请先登录
  │     ├─ uid=1 (超级管理员) → 直接放行
  │     └─ 其他用户：
  │           enforcer.Enforce("adminUser:{uid}", "/api/v1/admin-user/list", "GET")
  │             ├─ 匹配成功 → 放行
  │             └─ 匹配失败 → 403 暂无权限
  │
  └─ Controller 处理
```

---

## 前置：添加依赖

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin
go get github.com/casbin/casbin/v3
go get github.com/casbin/gorm-adapter/v3
```

---

## 步骤 1：创建 Casbin 模型文件

### 新建 `rbac_model.conf`（项目根目录）

**对应 gin-layout：** `rbac_model.conf`

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = (g(r.sub, p.sub) && keyMatch3(r.obj, p.obj) && regexMatch(r.act, p.act))
```

**各部分作用：**
- `request_definition` — 定义每次鉴权请求的结构：sub(谁)、obj(哪个资源路径)、act(什么操作)
- `policy_definition` — 定义策略规则的结构，和请求结构相同
- `role_definition` — `g = _, _` 表示"谁属于哪个角色/组"，Casbin 用这个函数做用户分组匹配
- `policy_effect` — 只要匹配到一条 allow 策略就通过
- `matchers` — 匹配规则：
  - `g(r.sub, p.sub)` — 用户和策略主体通过 g 函数关联
  - `keyMatch3(r.obj, p.obj)` — 路径通配符匹配（如 `/api/v1/*` 匹配 `/api/v1/login`）
  - `regexMatch(r.act, p.act)` — HTTP 方法正则匹配

---

## 步骤 2：添加全局常量

### 修改 `global/context_keys.go`

追加 Casbin 相关常量：

```go
// Casbin 前缀常量
const (
    SuperAdminId        uint = 1
    CasbinAdminUserPrefix   = "adminUser"
    CasbinSeparator         = ":"
)
```

### 修改 `global/api_auth_mode.go`

追加 `RequiresAPIPermission` 方法：

```go
// RequiresAPIPermission 返回该模式是否要求 API 权限校验
func (m ApiAuthMode) RequiresAPIPermission() bool {
    return m == ApiAuthModeAuth
}
```

---

## 步骤 3：Casbin Enforcer 封装

### 新建 `internal/access/casbin/enforcer.go`

**对应 gin-layout：** `internal/access/casbin/casbin.go` + `adapter.go` + `enforcer_init.go` 合并简化

```go
package 

import (
    "fmt"
    "os"
    "path/filepath"
    "sync"

    "GinAdmin/data"

    "github.com/casbin/casbin/v3"
    "github.com/casbin/casbin/v3/model"
    gormadapter "github.com/casbin/gorm-adapter/v3"
)

// CasbinEnforcer 封装 Casbin Enforcer
type CasbinEnforcer struct {
    *casbin.Enforcer
}

var (
    enforcer *CasbinEnforcer
    once     sync.Once
    initErr  error
)

// InitEnforcer 初始化 Casbin Enforcer（只执行一次）
// 启动时在 data.Initialize() 之后调用
func InitEnforcer() error {
    once.Do(func() {
        initErr = doInit()
    })
    return initErr
}

func doInit() error {
    // 1. 获取模型文件路径
    cwd, _ := os.Getwd()
    modelPath := filepath.Join(cwd, "rbac_model.conf")
    if _, err := os.Stat(modelPath); os.IsNotExist(err) {
        return fmt.Errorf("Casbin 模型文件不存在: %s", modelPath)
    }

    // 2. 加载模型
    m, err := model.NewModelFromFile(modelPath)
    if err != nil {
        return fmt.Errorf("加载 Casbin 模型失败: %w", err)
    }

    // 3. 获取数据库连接
    db := data.GetDB()
    if db == nil {
        return fmt.Errorf("数据库未初始化")
    }

    // 4. 创建 GORM 适配器（策略存储在 casbin_rule 表）
    gormadapter.TurnOffAutoMigrate(db)        // 不使用自动迁移
    adapter, err := gormadapter.NewAdapterByDB(db)
    if err != nil {
        return fmt.Errorf("创建 Casbin 适配器失败: %w", err)
    }

    // 5. 创建 Enforcer（自动从 casbin_rule 表加载策略）
    e, err := casbin.NewEnforcer(m, adapter)
    if err != nil {
        return fmt.Errorf("创建 Enforcer 失败: %w", err)
    }
    e.EnableAutoSave(true) // 调用 AddPolicy 等方法时自动写入数据库

    enforcer = &CasbinEnforcer{Enforcer: e}
    return nil
}

// GetEnforcer 返回已初始化的 Enforcer 实例
func GetEnforcer() (*CasbinEnforcer, error) {
    if enforcer == nil {
        if err := InitEnforcer(); err != nil {
            return nil, err
        }
    }
    return enforcer, nil
}

// ReloadPolicy 重新从数据库加载策略（策略变更后调用）
func ReloadPolicy() error {
    e, err := GetEnforcer()
    if err != nil {
        return err
    }
    return e.LoadPolicy()
}
```

---

## 步骤 4：策略同步服务

### 新建 `internal/service/casbin_sync.go`

**作用：** 从数据库读出"用户→角色→菜单→API"的完整权限链，生成 Casbin 策略写入 `casbin_rule` 表。

```go
package service

import (
    "GinAdmin/data"
    "GinAdmin/internal/access/casbinx"
    "GinAdmin/internal/model"
    "fmt"
    "log"
)

// SyncAllPolicies 同步所有用户的 Casbin 策略
// 策略格式：p, adminUser:{uid}, /api/v1/xxx, GET
// 调用时机：启动时 + 角色/菜单关联变更时
func SyncAllPolicies() error {
    enforcer, err := casbinx.GetEnforcer()
    if err != nil {
        return err
    }

    db := data.GetDB()

    // 1. 查询所有角色-菜单-API 关联
    // role_menu_map  JOIN menu_api_map  JOIN api
    type RoleApiRow struct {
        RoleID uint
        Method string
        Route  string
    }
    var roleApis []RoleApiRow
    err = db.Table("role_menu_map").
        Select("role_menu_map.role_id, api.method, api.route").
        Joins("JOIN menu_api_map ON menu_api_map.menu_id = role_menu_map.menu_id").
        Joins("JOIN api ON api.id = menu_api_map.api_id").
        Where("api.is_effective = 1").
        Scan(&roleApis).Error
    if err != nil {
        return fmt.Errorf("查询角色API关联失败: %w", err)
    }

    // 2. 查询所有用户-角色关联
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

    // 3. 构建策略：user → policies
    // 先建立 roleID → policies 的映射
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
```

---

## 步骤 5：鉴权中间件

### 新建 `internal/middleware/admin_auth.go`

**对应 gin-layout：** `internal/middleware/admin_auth.go`

```go
package middleware

import (
    "GinAdmin/global"
    "GinAdmin/internal/access/casbinx"
    "GinAdmin/internal/pkg/errors"
    "GinAdmin/internal/pkg/response"
    "fmt"

    "github.com/gin-gonic/gin"
)

// AdminAuthHandler Casbin 权限校验中间件
// 依赖 ParseToken 中间件预先写入 uid 到 context
// 只对路由树中声明 AuthPerm 的路由生效
func AdminAuthHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 获取当前用户 ID（ParseToken 中间件写入）
        uid := c.GetUint(global.ContextKeyUID)
        if uid == 0 {
            response.FailCode(c, errors.NotLogin)
            c.Abort()
            return
        }

        // 2. 超级管理员直接放行
        if uid == global.SuperAdminId {
            c.Next()
            return
        }

        // 3. 获取 Casbin Enforcer
        enforcer, err := casbinx.GetEnforcer()
        if err != nil {
            response.FailCode(c, errors.ServerErr)
            c.Abort()
            return
        }

        // 4. 执行权限检查
        userKey := fmt.Sprintf("%s%s%d", global.CasbinAdminUserPrefix, global.CasbinSeparator, uid)
        path := c.Request.URL.Path
        method := c.Request.Method

        ok, err := enforcer.Enforce(userKey, path, method)
        if err != nil {
            response.FailCode(c, errors.ServerErr)
            c.Abort()
            return
        }

        if !ok {
            response.FailCode(c, errors.AuthorizationErr)
            c.Abort()
            return
        }

        c.Next()
    }
}
```

**工作流程：**
1. 从 context 取出 uid（`ParseToken` 中间件已解析 JWT 并写入）
2. `uid == 1` → 超级管理员，跳过所有权限检查
3. 构建 subject：`adminUser:{uid}`（如 `adminUser:3`）
4. 调用 `enforcer.Enforce("adminUser:3", "/api/v1/admin-user/list", "GET")`
5. Casbin 去 `casbin_rule` 表查找有没有 `p, adminUser:3, /api/v1/admin-user/list, GET` 这条策略
6. 有 → 放行，没有 → 403

---

## 步骤 6：路由改造 — 区分公开/登录/鉴权路由

### 修改 `internal/routers/router.go`

**核心改动：** 将需要权限校验的路由放到一个独立的路由组里，给这个组挂载 `AdminAuthHandler` 中间件。

```go
func AdminRouteTree() RouteGroupDef {
    loginCtrl := controller.NewLoginController()
    adminUserCtrl := controller.NewAdminUserController()
    roleCtrl := controller.NewRoleController()
    menuCtrl := controller.NewMenuController()

    return RouteGroupDef{
        Prefix: "api/v1",
        Children: []RouteGroupDef{
            // ---- 公开接口（AuthNone）----
            {
                Routes: []RouteDef{
                    POST("login", "登录", AuthNone, loginCtrl.Login),
                },
            },
            // ---- 只需登录（AuthLogin）----
            {
                Routes: []RouteDef{
                    // 用户 - 个人操作
                    GET("admin-user/get", "获取当前用户信息", AuthLogin, adminUserCtrl.GetUserInfo),
                    POST("admin-user/update-profile", "更新个人资料", AuthLogin, adminUserCtrl.UpdateProfile),
                },
            },
            // ---- 需要权限（AuthPerm）+ 挂载 AdminAuthHandler ----
            {
                Middleware: []gin.HandlerFunc{middleware.AdminAuthHandler()},
                Children: []RouteGroupDef{
                    {
                        Prefix: "admin-user",
                        Routes: []RouteDef{
                            GET("list", "用户列表", AuthPerm, adminUserCtrl.List),
                            GET("detail", "用户详情", AuthPerm, adminUserCtrl.Detail),
                            POST("create", "新增用户", AuthPerm, adminUserCtrl.Create),
                            POST("update", "更新用户", AuthPerm, adminUserCtrl.Update),
                            POST("delete", "删除用户", AuthPerm, adminUserCtrl.Delete),
                            POST("bind-role", "绑定角色", AuthPerm, adminUserCtrl.BindRole),
                        },
                    },
                    {
                        Prefix: "role",
                        Routes: []RouteDef{
                            GET("list", "角色列表", AuthPerm, roleCtrl.List),
                            GET("detail", "角色详情", AuthPerm, roleCtrl.Detail),
                            POST("create", "新增角色", AuthPerm, roleCtrl.Create),
                            POST("update", "更新角色", AuthPerm, roleCtrl.Update),
                            POST("delete", "删除角色", AuthPerm, roleCtrl.Delete),
                        },
                    },
                    {
                        Prefix: "menu",
                        Routes: []RouteDef{
                            GET("list", "菜单列表", AuthPerm, menuCtrl.List),
                            GET("detail", "菜单详情", AuthPerm, menuCtrl.Detail),
                            POST("create", "新增菜单", AuthPerm, menuCtrl.Create),
                            POST("update", "更新菜单", AuthPerm, menuCtrl.Update),
                            POST("delete", "删除菜单", AuthPerm, menuCtrl.Delete),
                        },
                    },
                },
            },
        },
    }
}
```

**路由分组说明：**

| 组 | 中间件 | 路由标记 | 效果 |
|----|--------|---------|------|
| 公开 | 无 | AuthNone | 不需要登录，任何人都能访问（login） |
| 登录 | ParseToken | AuthLogin | 需要登录但不需要权限（个人信息） |
| 鉴权 | AdminAuthHandler | AuthPerm | 需要登录 + 需要 Casbin 权限 |

---

## 步骤 7：启动初始化

### 修改 `cmd/root.go`

在 `PersistentPreRunE` 的数据库初始化之后，加入 Casbin 初始化和策略同步：

```go
// 初始化数据库
err := data.Initialize()
if err != nil {
    return err
}

// 初始化 Casbin（新加）
if err := casbinx.InitEnforcer(); err != nil {
    return fmt.Errorf("初始化 Casbin 失败: %w", err)
}

// 同步权限策略（新加）
if err := service.SyncAllPolicies(); err != nil {
    return fmt.Errorf("同步权限策略失败: %w", err)
}
```

注意：`service.SyncAllPolicies()` 需要 import `"GinAdmin/internal/service"`。

---

## 步骤 8：角色/菜单变更后刷新策略

### 修改 `internal/service/role.go` — 在 Create/Update/Delete 末尾

```go
// 角色变更后刷新 Casbin 策略
func (s *RoleService) refreshPolicy() {
    if err := SyncAllPolicies(); err != nil {
        // 日志记录，不影响正常返回
        log.Printf("刷新 Casbin 策略失败: %v", err)
    }
}
```

在 `Create`、`Update`、`Delete` 方法的 return 之前调用 `s.refreshPolicy()`。

---

## 步骤 9：填充测试数据

### 向 `menu_api_map` 和 `api` 表插入测试数据

```sql
-- API 接口记录（当前已实现的所有接口）
INSERT INTO `api` (`id`, `code`, `group_code`, `name`, `method`, `route`, `is_auth`) VALUES
(1,  MD5('POST/api/v1/login'),               'login',     '用户登录',        'POST', '/api/v1/login',                 0),
(2,  MD5('GET/api/v1/admin-user/get'),        'adminUser', '获取当前用户信息', 'GET',  '/api/v1/admin-user/get',         1),
(3,  MD5('POST/api/v1/admin-user/update-profile'), 'adminUser', '更新个人资料', 'POST', '/api/v1/admin-user/update-profile', 1),
(4,  MD5('GET/api/v1/admin-user/list'),       'adminUser', '用户列表',        'GET',  '/api/v1/admin-user/list',        2),
(5,  MD5('GET/api/v1/admin-user/detail'),     'adminUser', '用户详情',        'GET',  '/api/v1/admin-user/detail',      2),
(6,  MD5('POST/api/v1/admin-user/create'),    'adminUser', '新增用户',        'POST', '/api/v1/admin-user/create',     2),
(7,  MD5('POST/api/v1/admin-user/update'),    'adminUser', '更新用户',        'POST', '/api/v1/admin-user/update',     2),
(8,  MD5('POST/api/v1/admin-user/delete'),    'adminUser', '删除用户',        'POST', '/api/v1/admin-user/delete',     2),
(9,  MD5('POST/api/v1/admin-user/bind-role'), 'adminUser', '绑定角色',        'POST', '/api/v1/admin-user/bind-role',  2),
(10, MD5('GET/api/v1/role/list'),            'role',      '角色列表',        'GET',  '/api/v1/role/list',              2),
(11, MD5('GET/api/v1/role/detail'),          'role',      '角色详情',        'GET',  '/api/v1/role/detail',            2),
(12, MD5('POST/api/v1/role/create'),         'role',      '新增角色',        'POST', '/api/v1/role/create',           2),
(13, MD5('POST/api/v1/role/update'),         'role',      '更新角色',        'POST', '/api/v1/role/update',           2),
(14, MD5('POST/api/v1/role/delete'),         'role',      '删除角色',        'POST', '/api/v1/role/delete',           2),
(15, MD5('GET/api/v1/menu/list'),            'menu',      '菜单列表',        'GET',  '/api/v1/menu/list',              2),
(16, MD5('GET/api/v1/menu/detail'),          'menu',      '菜单详情',        'GET',  '/api/v1/menu/detail',            2),
(17, MD5('POST/api/v1/menu/create'),         'menu',      '新增菜单',        'POST', '/api/v1/menu/create',           2),
(18, MD5('POST/api/v1/menu/update'),         'menu',      '更新菜单',        'POST', '/api/v1/menu/update',           2),
(19, MD5('POST/api/v1/menu/delete'),         'menu',      '删除菜单',        'POST', '/api/v1/menu/delete',           2);

-- 超级管理员角色(role_id=1) 拥有所有菜单(menu_id 1-5)
INSERT IGNORE INTO `role_menu_map` (`role_id`, `menu_id`) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5);

-- 菜单-API 绑定：权限管理的三个菜单各自绑定对应的 API
INSERT IGNORE INTO `menu_api_map` (`menu_id`, `api_id`) VALUES
-- 用户管理菜单(menu_id=3) → 用户相关 API
(3, 4), (3, 5), (3, 6), (3, 7), (3, 8), (3, 9),
-- 角色管理菜单(menu_id=4) → 角色相关 API
(4, 10), (4, 11), (4, 12), (4, 13), (4, 14),
-- 菜单管理菜单(menu_id=5) → 菜单相关 API
(5, 15), (5, 16), (5, 17), (5, 18), (5, 19);
```

---

## 修改文件汇总

| 文件 | 操作 | 内容 |
|------|------|------|
| `rbac_model.conf` | **新建** | Casbin 模型定义 |
| `global/context_keys.go` | 修改 | 追加 Casbin 常量 |
| `global/api_auth_mode.go` | 修改 | 追加 `RequiresAPIPermission()` |
| `internal/access/casbin/enforcer.go` | **新建** | Casbin 单例 + GORM 适配器 |
| `internal/service/casbin_sync.go` | **新建** | 权限策略同步 |
| `internal/middleware/admin_auth.go` | **新建** | 鉴权中间件 |
| `internal/routers/router.go` | 修改 | 路由分组 + 挂载 AdminAuthHandler |
| `cmd/root.go` | 修改 | 启动时初始化 Casbin + 同步策略 |
| `internal/service/role.go` | 修改 | Create/Update/Delete 后刷新策略 |
| SQL 测试数据 | 手动执行 | api + role_menu_map + menu_api_map |

---

## 测试验证

```bash
# 1. 编译
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go build ./...

# 2. 导入测试数据
docker exec -i $(docker ps -qf "name=mysql") mysql -u root -p123456 < 上面的SQL

# 3. 启动服务
go run main.go server

# 4. super_admin 登录（uid=1，超级管理员，所有接口都可访问）
TOKEN=$(curl -s -X POST http://127.0.0.1:8080/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"super_admin","password":"123456"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['data']['access_token'])")

# 5. 访问需要权限的接口
curl http://127.0.0.1:8080/api/v1/admin-user/list -H "Authorization: Bearer $TOKEN"
# 预期: 200，返回用户列表

# 6. 创建一个普通用户，测试权限拦截
# （未完，数据库里再插入一个用户，不给他分配角色，用他的 token 访问 admin-user/list）
# 预期: {"code":403,"msg":"暂无权限",...}
```
