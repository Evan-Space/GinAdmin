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