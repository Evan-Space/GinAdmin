# 切换到 gin-layout 表结构 — 完整代码参考

> 每层代码按依赖顺序排列，你可以逐文件手写。

---

## 前置：添加依赖

```bash
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin
go get gorm.io/plugin/soft_delete@latest
```

---

## 1. SQL — 替换 `SQL/createTableDemoSql.sql`

**对应 gin-layout：** `data/migrations/20260425000001_init_table.up.sql` + `20260425000002_init_data.up.sql`

```sql
-- ========================================
-- GinAdmin 初始化 SQL（对齐 gin-layout 表结构）
-- ========================================

CREATE DATABASE IF NOT EXISTS `GinAdminMysql` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE `GinAdminMysql`;

-- ----------------------------
-- 1. admin_user 用户表
-- 对应 gin-layout internal/model/admin_users.go AdminUser
-- 包含软删除字段 deleted_at (int unix时间戳, 0=未删除)
-- AUTO_INCREMENT 从 10000 开始（gin-layout 惯例）
-- ----------------------------
DROP TABLE IF EXISTS `admin_user`;
CREATE TABLE `admin_user` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `nickname` varchar(30) NOT NULL DEFAULT '' COMMENT '昵称',
  `username` varchar(30) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
  `phone_number` varchar(15) NOT NULL DEFAULT '' COMMENT '手机号',
  `full_phone_number` varchar(20) NOT NULL DEFAULT '' COMMENT '带区号的手机号',
  `country_code` varchar(10) NOT NULL DEFAULT '' COMMENT '国际区号',
  `email` varchar(120) NOT NULL DEFAULT '' COMMENT '邮箱',
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态 1启用 0禁用',
  `is_super_admin` tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否超级管理员 1是 0不是',
  `last_login` datetime DEFAULT NULL COMMENT '最后登录时间',
  `last_ip` varchar(45) NOT NULL DEFAULT '' COMMENT '最后登录IP',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` int unsigned NOT NULL DEFAULT 0 COMMENT '软删除时间戳(0=未删除)',
  PRIMARY KEY (`id`),
  UNIQUE KEY `adu_u_d` (`username`, `deleted_at`),
  KEY `idx_status_deleted_at` (`status`, `deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- 2. role 角色表（树形层级结构）
-- 对应 gin-layout internal/model/role.go Role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(60) NOT NULL DEFAULT '' COMMENT '角色业务编码',
  `is_system` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '是否系统保留对象 1是 0否',
  `pid` int unsigned NOT NULL DEFAULT 0 COMMENT '上级角色id',
  `pids` varchar(255) NOT NULL DEFAULT '' COMMENT '所有上级id路径 逗号分隔',
  `name` varchar(60) NOT NULL DEFAULT '' COMMENT '角色名称',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
  `level` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '层级',
  `sort` mediumint unsigned NOT NULL DEFAULT 0 COMMENT '排序 数字越大越靠前',
  `children_num` int unsigned NOT NULL DEFAULT 0 COMMENT '子集数量',
  `status` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '是否启用 1启用 0禁用',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` int unsigned NOT NULL DEFAULT 0 COMMENT '软删除时间戳',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_code_deleted_at` (`code`, `deleted_at`),
  KEY `idx_pid_deleted_at_sort_id` (`pid`, `deleted_at`, `sort`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

-- ----------------------------
-- 3. menu 菜单表（树形层级结构，25个字段）
-- 对应 gin-layout internal/model/menu.go Menu
-- ----------------------------
DROP TABLE IF EXISTS `menu`;
CREATE TABLE `menu` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `icon` varchar(255) NOT NULL DEFAULT '' COMMENT '图标（Ant Design 图标名）',
  `code` varchar(120) NOT NULL DEFAULT '' COMMENT '前端权限标识',
  `path` varchar(255) NOT NULL DEFAULT '' COMMENT '前端路由路径',
  `full_path` varchar(255) NOT NULL DEFAULT '' COMMENT '完整路由路径',
  `redirect` varchar(255) NOT NULL DEFAULT '' COMMENT '重定向路由名称',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '前端路由名称',
  `component` varchar(255) NOT NULL DEFAULT '' COMMENT '前端组件路径',
  `animate_enter` varchar(60) NOT NULL DEFAULT '' COMMENT '进入动画',
  `animate_leave` varchar(60) NOT NULL DEFAULT '' COMMENT '离开动画',
  `animate_duration` float(2,2) NOT NULL DEFAULT 0.00 COMMENT '动画时长',
  `is_show` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '是否显示 1是 0否',
  `status` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '状态 1启用 0禁用',
  `is_auth` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '是否需要授权 1是 0否',
  `is_external_links` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '是否外链 1是 0否',
  `is_new_window` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '是否新窗口打开 1是 0否',
  `sort` int unsigned NOT NULL DEFAULT 0 COMMENT '排序 数字越大越靠前',
  `type` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '菜单类型 1目录 2菜单 3按钮',
  `pid` int unsigned NOT NULL DEFAULT 0 COMMENT '上级菜单id',
  `level` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '层级',
  `pids` varchar(255) NOT NULL DEFAULT '' COMMENT '层级序列 逗号分隔',
  `children_num` int unsigned NOT NULL DEFAULT 0 COMMENT '子集数量',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` int unsigned NOT NULL DEFAULT 0 COMMENT '软删除时间戳',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_code_del` (`code`, `deleted_at`),
  KEY `idx_pid_deleted_at_sort_id` (`pid`, `deleted_at`, `sort`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

-- ----------------------------
-- 4. admin_user_role_map 用户-角色关联表
-- 用户通过此表绑定一个或多个角色
-- 注意：关联表无软删除，使用 BaseModel（无 deleted_at）
-- ----------------------------
DROP TABLE IF EXISTS `admin_user_role_map`;
CREATE TABLE `admin_user_role_map` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `uid` int unsigned NOT NULL DEFAULT 0 COMMENT 'admin_user表id',
  `role_id` int unsigned NOT NULL DEFAULT 0 COMMENT '角色id',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uid_role_id` (`uid`, `role_id`),
  KEY `idx_role_id_uid` (`role_id`, `uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

-- ----------------------------
-- 5. role_menu_map 角色-菜单关联表
-- 角色通过此表获得菜单访问权限
-- ----------------------------
DROP TABLE IF EXISTS `role_menu_map`;
CREATE TABLE `role_menu_map` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `role_id` int unsigned NOT NULL DEFAULT 0 COMMENT '角色id',
  `menu_id` int unsigned NOT NULL DEFAULT 0 COMMENT '菜单id',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_role_id_menu_id` (`role_id`, `menu_id`),
  KEY `idx_menu_id_role_id` (`menu_id`, `role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

-- ----------------------------
-- 6. menu_api_map 菜单-API关联表
-- 菜单下的按钮/操作的 API 权限通过此表关联
-- ----------------------------
DROP TABLE IF EXISTS `menu_api_map`;
CREATE TABLE `menu_api_map` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `menu_id` int unsigned NOT NULL DEFAULT 0 COMMENT '菜单id',
  `api_id` int unsigned NOT NULL DEFAULT 0 COMMENT '接口id',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_menu_id_api_id` (`menu_id`, `api_id`),
  KEY `idx_api_id_menu_id` (`api_id`, `menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

-- ----------------------------
-- 7. api 接口权限表
-- 系统内所有 API 接口的元数据（由路由声明自动同步）
-- ----------------------------
DROP TABLE IF EXISTS `api`;
CREATE TABLE `api` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `code` char(32) NOT NULL DEFAULT '' COMMENT 'MD5(method_path)唯一标识',
  `group_code` varchar(60) NOT NULL DEFAULT '' COMMENT '分组编码',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT '接口名称',
  `method` varchar(16) NOT NULL DEFAULT '' COMMENT '请求方式 GET/POST/PUT/DELETE',
  `route` varchar(256) NOT NULL DEFAULT '' COMMENT '接口路径 如 /api/v1/login',
  `is_auth` tinyint NOT NULL DEFAULT 0 COMMENT '鉴权模式 0无需登录 1需登录 2需权限',
  `is_effective` tinyint NOT NULL DEFAULT 1 COMMENT '是否有效 1有效 0无效（过期路由标记）',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_is_auth` (`is_auth`, `is_effective`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

-- ----------------------------
-- 8. department 部门表（树形层级，后续使用）
-- ----------------------------
DROP TABLE IF EXISTS `department`;
CREATE TABLE `department` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(60) NOT NULL DEFAULT '' COMMENT '部门业务编码',
  `is_system` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '是否系统保留',
  `pid` int unsigned NOT NULL DEFAULT 0 COMMENT '上级部门id',
  `pids` varchar(255) NOT NULL DEFAULT '' COMMENT '上级id路径',
  `name` varchar(60) NOT NULL DEFAULT '' COMMENT '部门名称',
  `description` varchar(255) NOT NULL DEFAULT '' COMMENT '描述',
  `level` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '层级',
  `sort` mediumint unsigned NOT NULL DEFAULT 0 COMMENT '排序',
  `children_num` int unsigned NOT NULL DEFAULT 0 COMMENT '子集数量',
  `user_number` int unsigned NOT NULL DEFAULT 0 COMMENT '部门用户数量',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` int unsigned NOT NULL DEFAULT 0 COMMENT '软删除时间戳',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_code_deleted_at` (`code`, `deleted_at`),
  KEY `idx_pid_deleted_at_sort_id` (`pid`, `deleted_at`, `sort`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

-- ========================================
-- 种子数据
-- ========================================

-- 默认超级管理员角色（is_system=1 不允许删除）
INSERT INTO `role` (`id`, `code`, `is_system`, `pid`, `pids`, `name`, `description`, `level`, `sort`, `status`) VALUES
(1, 'super_admin', 1, 0, '', '超级管理员', '系统内置超级管理员角色', 1, 0, 1);

-- 超级管理员用户（id=1, is_super_admin=1, 密码明文 123456 测试用）
INSERT INTO `admin_user` (`id`, `username`, `password`, `nickname`, `email`, `is_super_admin`, `status`) VALUES
(1, 'super_admin', '123456', '超级管理员', 'admin@example.com', 1, 1);

-- 用户-角色绑定：超级管理员拥有超级管理员角色
INSERT INTO `admin_user_role_map` (`uid`, `role_id`) VALUES
(1, 1);

-- 默认菜单树（pid/level 形成树形层级关系）
INSERT INTO `menu` (`id`, `code`, `name`, `path`, `icon`, `type`, `pid`, `level`, `is_show`, `is_auth`, `status`, `sort`) VALUES
(1, 'dashboard',  '控制台',     '/dashboard',     'DashboardOutlined',  1, 0, 0, 1, 0, 1, 1),
(2, 'system',     '系统管理',   '/system',        'SettingOutlined',    1, 0, 0, 1, 0, 1, 2),
(3, 'user',       '用户管理',   '/system/user',   'UserOutlined',       2, 2, 1, 1, 1, 1, 1),
(4, 'role',       '角色管理',   '/system/role',   'TeamOutlined',       2, 2, 1, 1, 1, 1, 2),
(5, 'menu-mgr',   '菜单管理',   '/system/menu',   'MenuOutlined',       2, 2, 1, 1, 1, 1, 3);
```

**关键设计说明：**
- `admin_user`、`role`、`menu`、`department` 四张主表都有 `deleted_at` 软删除字段（int unix 时间戳，0=未删除）
- 关联表（`admin_user_role_map`、`role_menu_map`、`menu_api_map`）无软删除，记录直接物理删除
- `admin_user` 的 AUTO_INCREMENT 从 10000 开始，系统内置用户 id=1 手动指定
- `role.is_system=1` 的角色不允许删除
- `menu.type` 区分目录(1)/菜单(2)/按钮(3)，对应前端不同渲染方式
- `api` 表用于存储路由元数据，后续 Casbin RBAC 会用到

重新导入数据库：

```bash
docker exec -i $(docker ps -qf "name=ginAdminContainer") mysql -u root -p123456 < SQL/createTableDemoSql.sql
```

---

## 2. config/config.go — 删除旧的 basePath 字段

**文件：** `config/config.go`

gin-layout 的 config 中有 `BasePath`，但我们目前未用到。先保持现有 config 结构不变，只确认 JWT 和 MySQL 配置正确即可。

---

## 3. Model 层

### 3.1 `internal/model/base.go` — 新建

**对应 gin-layout：** `internal/model/base.go`

```go
package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

// BaseModel 提供模型通用字段 id + created_at + updated_at
// 对应 gin-layout internal/model/base.go BaseModel
type BaseModel struct {
	ID        uint      `gorm:"column:id;type:int unsigned AUTO_INCREMENT;not null;primarykey" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;<-:create" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

// ContainsDeleteBaseModel 在 BaseModel 基础上增加软删除字段
// 对应 gin-layout internal/model/base.go ContainsDeleteBaseModel
type ContainsDeleteBaseModel struct {
	BaseModel
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;type:int unsigned;not null;default:0;index" json:"-"`
}
```

**作用：**
- `BaseModel` — 所有表公用的 id、created_at、updated_at，供关联表（map 表）嵌入
- `ContainsDeleteBaseModel` — 在 BaseModel 上叠加软删除 deleted_at 字段，供 admin_user、role、menu 嵌入
- `soft_delete.DeletedAt` 是 GORM 软删除插件提供的时间戳类型，值为 0 表示未删除，非 0 表示已删除

---

### 3.2 `internal/model/admin_user.go` — 重写（替换原 user.go）

**对应 gin-layout：** `internal/model/admin_users.go`

```go
package model

import "time"

// 管理员状态常量
const (
	AdminUserStatusEnabled  uint8 = 1 // 启用
	AdminUserStatusDisabled uint8 = 0 // 禁用
)

// AdminUser 管理员用户表，对应数据库 admin_user 表
// 对应 gin-layout internal/model/admin_users.go AdminUser
type AdminUser struct {
	ContainsDeleteBaseModel
	IsSuperAdmin    uint8     `json:"is_super_admin" gorm:"column:is_super_admin;type:tinyint(1);not null;default:0"`     // 是否超级管理员
	Nickname        string    `json:"nickname" gorm:"column:nickname;type:varchar(30);not null;default:''"`               // 昵称
	Username        string    `json:"username" gorm:"column:username;type:varchar(30);not null;default:''"`               // 用户名
	Password        string    `json:"-" gorm:"column:password;type:varchar(255);not null;default:''"`                     // 密码（json:"−" 不序列化到响应体）
	PhoneNumber     string    `json:"phone_number" gorm:"column:phone_number;type:varchar(15);not null;default:''"`       // 手机号
	FullPhoneNumber string    `json:"full_phone_number" gorm:"column:full_phone_number;type:varchar(20);not null;default:''"` // 完整手机号
	CountryCode     string    `json:"country_code" gorm:"column:country_code;type:varchar(10);not null;default:''"`       // 国际区号
	Email           string    `json:"email" gorm:"column:email;type:varchar(120);not null;default:''"`                    // 邮箱
	Avatar          string    `json:"avatar" gorm:"column:avatar;type:varchar(255);not null;default:''"`                  // 头像
	Status          uint8     `json:"status" gorm:"column:status;type:tinyint(1);not null;default:1"`                     // 状态 1启用 0禁用
	LastLogin       time.Time `json:"last_login" gorm:"column:last_login;type:datetime"`                                  // 最后登录时间
	LastIp          string    `json:"last_ip" gorm:"column:last_ip;type:varchar(45);not null;default:''"`                 // 最后登录IP
}

// TableName 指定表名
func (AdminUser) TableName() string {
	return "admin_user"
}
```

**与原 GinAdmin user.go 的差异：**

- 表名 `users` → `admin_user`
- 结构体名 `User` → `AdminUser`
- 嵌入了 `ContainsDeleteBaseModel`（继承 id、created_at、updated_at、deleted_at）
- 新增字段：is_super_admin、phone_number、full_phone_number、country_code、avatar、last_login、last_ip
- Status 类型 `int` → `uint8`，与 gin-layout 一致

---

### 3.3 `internal/model/role.go` — 新建

**对应 gin-layout：** `internal/model/role.go`

```go
package model

// Role 角色表，对应数据库 role 表
// 对应 gin-layout internal/model/role.go Role
type Role struct {
	ContainsDeleteBaseModel
	Code        string `json:"code" gorm:"column:code;type:varchar(60);not null;default:''"`                // 角色业务编码
	IsSystem    uint8  `json:"is_system" gorm:"column:is_system;type:tinyint unsigned;not null;default:0"`   // 是否系统保留对象
	Pid         uint   `json:"pid" gorm:"column:pid;type:int unsigned;not null;default:0"`                   // 上级id
	Pids        string `json:"pids" gorm:"column:pids;type:varchar(255);not null;default:''"`                // 所有上级id路径
	Name        string `json:"name" gorm:"column:name;type:varchar(60);not null;default:''"`                 // 角色名称
	Description string `json:"description" gorm:"column:description;type:varchar(255);not null;default:''"`  // 描述
	Level       uint8  `json:"level" gorm:"column:level;type:tinyint unsigned;not null;default:1"`           // 层级
	Sort        uint   `json:"sort" gorm:"column:sort;type:mediumint unsigned;not null;default:0"`           // 排序
	ChildrenNum uint   `json:"children_num" gorm:"column:children_num;type:int unsigned;not null;default:0"` // 子集数量
	Status      uint8  `json:"status" gorm:"column:status;type:tinyint unsigned;not null;default:1"`         // 状态 1启用 0禁用
}

// TableName 指定表名
func (Role) TableName() string {
	return "role"
}
```

**作用：**
- 角色支持树形层级（pid + pids + level），后续实现角色继承
- is_system=1 的角色不允许删除

---

### 3.4 `internal/model/menu.go` — 新建

**对应 gin-layout：** `internal/model/menu.go`

```go
package model

// 菜单类型常量
const (
	MenuTypeCatalogue uint8 = 1 // 目录
	MenuTypeMenu      uint8 = 2 // 菜单
	MenuTypeButton    uint8 = 3 // 按钮
)

// Menu 菜单权限表，对应数据库 menu 表
// 对应 gin-layout internal/model/menu.go Menu
type Menu struct {
	ContainsDeleteBaseModel
	Icon            string  `json:"icon" gorm:"column:icon;type:varchar(255);not null;default:''"`                    // 图标
	Code            string  `json:"code" gorm:"column:code;type:varchar(120);not null;default:''"`                    // 前端权限标识
	Path            string  `json:"path" gorm:"column:path;type:varchar(255);not null;default:''"`                    // 前端路由路径
	FullPath        string  `json:"full_path" gorm:"column:full_path;type:varchar(255);not null;default:''"`          // 完整路由路径
	Redirect        string  `json:"redirect" gorm:"column:redirect;type:varchar(255);not null;default:''"`            // 重定向路由名称
	Name            string  `json:"name" gorm:"column:name;type:varchar(255);not null;default:''"`                    // 路由名称
	Component       string  `json:"component" gorm:"column:component;type:varchar(255);not null;default:''"`          // 组件路径
	AnimateEnter    string  `json:"animate_enter" gorm:"column:animate_enter;type:varchar(60);not null;default:''"`   // 进入动画
	AnimateLeave    string  `json:"animate_leave" gorm:"column:animate_leave;type:varchar(60);not null;default:''"`   // 离开动画
	AnimateDuration float32 `json:"animate_duration" gorm:"column:animate_duration;type:float(2,2);not null;default:0.00"` // 动画时长
	IsShow          uint8   `json:"is_show" gorm:"column:is_show;type:tinyint unsigned;not null;default:1"`           // 是否显示 1是 0否
	Status          uint8   `json:"status" gorm:"column:status;type:tinyint unsigned;not null;default:1"`             // 状态 1启用 0禁用
	IsAuth          uint8   `json:"is_auth" gorm:"column:is_auth;type:tinyint unsigned;not null;default:0"`           // 是否需要授权 1是 0否
	IsExternalLinks uint8   `json:"is_external_links" gorm:"column:is_external_links;type:tinyint unsigned;not null;default:0"` // 是否外链
	IsNewWindow     uint8   `json:"is_new_window" gorm:"column:is_new_window;type:tinyint unsigned;not null;default:0"` // 是否新窗口打开
	Sort            uint    `json:"sort" gorm:"column:sort;type:int unsigned;not null;default:0"`                     // 排序
	Type            uint8   `json:"type" gorm:"column:type;type:tinyint unsigned;not null;default:1"`                 // 菜单类型 1目录 2菜单 3按钮
	Pid             uint    `json:"pid" gorm:"column:pid;type:int unsigned;not null;default:0"`                       // 上级菜单id
	Level           uint8   `json:"level" gorm:"column:level;type:tinyint unsigned;not null;default:0"`               // 层级
	Pids            string  `json:"pids" gorm:"column:pids;type:varchar(255);not null;default:''"`                    // 层级序列逗号分隔
	ChildrenNum     uint    `json:"children_num" gorm:"column:children_num;type:int unsigned;not null;default:0"`     // 子集数量
	Description     string  `json:"description" gorm:"column:description;type:varchar(255);not null;default:''"`      // 描述
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menu"
}
```

**作用：**
- 菜单同样是树形结构（pid + pids + level）
- type 区分目录/菜单/按钮三种类型，前端据此渲染不同组件
- is_auth 标记该菜单是否需要权限校验

---

### 3.5 `internal/model/admin_user_role_map.go` — 新建

**对应 gin-layout：** `internal/model/admin_user_role_map.go`

```go
package model

// AdminUserRoleMap 用户-角色关联表，对应数据库 admin_user_role_map 表
type AdminUserRoleMap struct {
	BaseModel
	Uid    uint `json:"uid" gorm:"column:uid;type:int unsigned;not null;default:0"`         // admin_user表id
	RoleId uint `json:"role_id" gorm:"column:role_id;type:int unsigned;not null;default:0"` // 角色id
}

// TableName 指定表名
func (AdminUserRoleMap) TableName() string {
	return "admin_user_role_map"
}
```

---

### 3.6 `internal/model/role_menu_map.go` — 新建

**对应 gin-layout：** `internal/model/role_menu_map.go`

```go
package model

// RoleMenuMap 角色-菜单关联表，对应数据库 role_menu_map 表
type RoleMenuMap struct {
	BaseModel
	RoleId uint `json:"role_id" gorm:"column:role_id;type:int unsigned;not null;default:0"` // 角色id
	MenuId uint `json:"menu_id" gorm:"column:menu_id;type:int unsigned;not null;default:0"` // 菜单id
}

// TableName 指定表名
func (RoleMenuMap) TableName() string {
	return "role_menu_map"
}
```

---

### 3.7 `internal/model/menu_api_map.go` — 新建

**对应 gin-layout：** `internal/model/menu_api_map.go`

```go
package model

// MenuApiMap 菜单-API关联表，对应数据库 menu_api_map 表
type MenuApiMap struct {
	BaseModel
	MenuId uint `json:"menu_id" gorm:"column:menu_id;type:int unsigned;not null;default:0"` // 菜单id
	ApiId  uint `json:"api_id" gorm:"column:api_id;type:int unsigned;not null;default:0"`   // 接口id
}

// TableName 指定表名
func (MenuApiMap) TableName() string {
	return "menu_api_map"
}
```

---

## 4. data/mysql.go — 改造

**文件：** `data/mysql.go`

```go
package data

import (
	"GinAdmin/config"
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db    *gorm.DB   // 全局数据库实例
	dbMu  sync.Mutex // 初始化锁
	dbErr error      // 初始化错误
)

// Initialize 初始化 MySQL 连接
func Initialize() error {
	cfg := config.GetConfig()
	if !cfg.MySQl.Enable {
		return nil
	}

	dbMu.Lock()
	defer dbMu.Unlock()

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQl.Username,
		cfg.MySQl.Password,
		cfg.MySQl.Host,
		cfg.MySQl.Port,
		cfg.MySQl.Database,
	)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		dbErr = err
		return err
	}

	db = conn
	return nil
}

// GetDB 返回全局数据库实例
func GetDB() *gorm.DB {
	return db
}

// SetDB 设置数据库实例（用于测试或事务场景注入）
func SetDB(conn *gorm.DB) {
	db = conn
}

// Shutdown 关闭数据库连接
func Shutdown() error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
```

---

## 5. `internal/pkg/request/request.go` — 确认无需改动

保持当前代码即可，`GetAccessToken` 返回 `(string, error)` 已经对齐 gin-layout。

---

## 6. Service 层

### 6.1 `internal/service/auth.go` — 改造

**文件：** `internal/service/auth.go`

改动点：
- `model.User` → `model.AdminUser`
- 表查询不变（GORM 通过 TableName() 自动路由到 `admin_user`）
- 密码比对逻辑不变

```go
package service

import (
	"GinAdmin/data"
	"GinAdmin/internal/model"
	"GinAdmin/internal/pkg/errors"
	"GinAdmin/internal/pkg/utils"

	"gorm.io/gorm"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// Login 用户登录，成功返回 JWT Token
// 对应 gin-layout internal/service/auth/login.go validateUser + Login
func (s *AuthService) Login(username, password string) (string, *model.AdminUser, error) {
	var user model.AdminUser
	err := data.GetDB().Where("username = ?", username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", nil, errors.NewBusinessError(errors.UserDoesNotExist)
		}
		return "", nil, errors.NewBusinessError(errors.ServerErr)
	}

	// 检查用户状态（Status == 0 已禁用）
	if user.Status == model.AdminUserStatusDisabled {
		return "", nil, errors.NewBusinessError(errors.UserDisable)
	}

	// 明文比对密码
	if user.Password != password {
		return "", nil, errors.NewBusinessError(errors.UserPasswordWrong)
	}

	// 生成 JWT
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", nil, errors.NewBusinessError(errors.TokenGenerateFailed)
	}

	return token, &user, nil
}
```

---

### 6.2 `internal/service/user.go` — 改造

**文件：** `internal/service/user.go`

改动点：`model.User` → `model.AdminUser`

```go
package service

import (
	"GinAdmin/data"
	"GinAdmin/internal/model"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

// ListUsers 查询所有未删除的用户
func (s *UserService) ListUsers() ([]model.AdminUser, error) {
	var users []model.AdminUser
	err := data.GetDB().
		Where("deleted_at = 0").
		Select("id", "username", "nickname", "email", "phone_number", "avatar", "status", "is_super_admin", "last_login", "created_at", "updated_at").
		Order("id asc").
		Find(&users).Error
	return users, err
}

// GetUserByID 按 ID 查询用户（不含密码，不含已删除记录）
func (s *UserService) GetUserByID(id uint) (*model.AdminUser, error) {
	var user model.AdminUser
	err := data.GetDB().
		Where("id = ? AND deleted_at = 0", id).
		Select("id", "username", "nickname", "email", "phone_number", "avatar", "status", "is_super_admin", "last_login", "created_at", "updated_at").
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
```

---

## 7. Controller 层

### 7.1 `internal/controller/login.go` — 只改 import，逻辑不变

适配 `AdminUser` 类型：

```go
package controller

import (
	"GinAdmin/internal/service"
	"GinAdmin/internal/validator"
	"GinAdmin/internal/validator/form"

	"github.com/gin-gonic/gin"
)

type LoginController struct {
	Api
	authService *service.AuthService
}

func NewLoginController() *LoginController {
	return &LoginController{
		authService: service.NewAuthService(),
	}
}

// Login 用户登录
func (ctl *LoginController) Login(c *gin.Context) {
	params := &form.LoginForm{}
	if err := validator.CheckPostParams(c, params); err != nil {
		return
	}

	token, user, err := ctl.authService.Login(params.Username, params.Password)
	if err != nil {
		ctl.Err(c, err)
		return
	}

	ctl.Success(c, gin.H{
		"access_token": token,
		"token_type":   "Bearer",
		"user": gin.H{
			"id":             user.ID,
			"username":       user.Username,
			"nickname":       user.Nickname,
			"email":          user.Email,
			"is_super_admin": user.IsSuperAdmin,
		},
	})
}
```

---

### 7.2 `internal/controller/user.go` — 重写

```go
package controller

import (
	"GinAdmin/internal/pkg/errors"
	"GinAdmin/internal/service"
	"GinAdmin/internal/validator"
	"GinAdmin/internal/validator/form"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Api
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

// List 用户列表
func (ctl *UserController) List(c *gin.Context) {
	params := form.NewPaginate()
	if err := validator.CheckQueryParams(c, params); err != nil {
		return
	}

	users, err := ctl.userService.ListUsers()
	if err != nil {
		ctl.Fail(c, 500, "查询用户列表失败")
		return
	}

	ctl.Success(c, users)
}

// Info 获取当前登录用户信息
func (ctl *UserController) Info(c *gin.Context) {
	uid := ctl.GetCurrentUserID(c)
	if uid == 0 {
		ctl.FailCode(c, errors.NotLogin)
		return
	}

	user, err := ctl.userService.GetUserByID(uid)
	if err != nil {
		ctl.FailCode(c, errors.NotFound)
		return
	}

	ctl.Success(c, gin.H{
		"id":             user.ID,
		"username":       user.Username,
		"nickname":       user.Nickname,
		"email":          user.Email,
		"phone_number":   user.PhoneNumber,
		"avatar":         user.Avatar,
		"status":         user.Status,
		"is_super_admin": user.IsSuperAdmin,
		"last_login":     user.LastLogin,
		"created_at":     user.CreatedAt,
	})
}
```

---

## 8. `internal/validator/form/common.go` — 无需改动

保持当前 Paginate 定义即可。

---

## 9. 其他无需改动的文件

| 文件 | 说明 |
|------|------|
| `internal/pkg/utils/jwt.go` | JWT Claims 中的 Username 字段不变 |
| `middleware/parse_token.go` | 从 context 取 uid，不依赖具体 model |
| `middleware/request_meta.go` | 纯请求元数据，不涉及表结构 |
| `middleware/cors.go` | 跨域处理，无关 |
| `internal/routers/*` | 路由定义，handler 签名不变 |
| `cmd/*` | 启动流程不变 |
| `internal/pkg/response/*` | 响应封装，无关 |
| `internal/pkg/errors/*` | 错误码体系，无关 |
| `internal/pkg/request/*` | 请求提取，无关 |
| `internal/validator/runtime.go` | 参数校验引擎，无关 |

---

## 10. 改动汇总

| 文件 | 操作 | 关键变化 |
|------|------|---------|
| `SQL/createTableDemoSql.sql` | 重写 | 8 张 gin-layout 风格表 + 种子数据 |
| `internal/model/base.go` | 新建 | BaseModel + ContainsDeleteBaseModel |
| `internal/model/admin_user.go` | 重写 | `User` → `AdminUser`，嵌入 ContainsDeleteBaseModel，新增 is_super_admin/phone/avatar/last_login/last_ip 等字段 |
| `internal/model/role.go` | 新建 | 角色表模型，树形结构（pid/pids/level） |
| `internal/model/menu.go` | 新建 | 菜单表模型，25 个字段，type 区分目录/菜单/按钮 |
| `internal/model/admin_user_role_map.go` | 新建 | 用户-角色关联表 |
| `internal/model/role_menu_map.go` | 新建 | 角色-菜单关联表 |
| `internal/model/menu_api_map.go` | 新建 | 菜单-API 关联表 |
| `data/mysql.go` | 改造 | 加锁保护 + dbErr 记录 |
| `internal/service/auth.go` | 修改 | `User` → `AdminUser`，status 用常量 |
| `internal/service/user.go` | 修改 | `User` → `AdminUser`，查询加 `deleted_at = 0` 条件 |
| `internal/controller/login.go` | 修改 | 响应字段适配 AdminUser |
| `internal/controller/user.go` | 修改 | 响应字段适配 AdminUser |
| `go.mod` | 修改 | 新增 `gorm.io/plugin/soft_delete` |

---

## 11. 验证步骤

```bash
# 1. 重新导入 SQL（清空旧数据 + 建新表 + 种子数据）
docker exec -i $(docker ps -qf "name=mysql") mysql -u root -p123456 < SQL/createTableDemoSql.sql

# 2. 编译检查
cd /Users/wangchuang/GoLang/studyAdmin/GinAdmin && go build ./...

# 3. 启动服务
go run main.go server

# 4. 测试登录（种子用户 super_admin / 123456）
curl -X POST http://127.0.0.1:8080/api/v1/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"super_admin","password":"123456"}'

# 5. 测试用户信息
curl http://127.0.0.1:8080/api/v1/user/info \
  -H 'Authorization: Bearer <上面返回的token>'

# 6. 测试用户列表
curl http://127.0.0.1:8080/api/v1/user/list \
  -H 'Authorization: Bearer <token>'
```
