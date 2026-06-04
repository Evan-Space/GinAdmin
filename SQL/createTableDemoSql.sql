-- ========================================
-- GinAdmin Demo 测试数据
-- 数据库：MySQL 8.0+
-- ========================================

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username    VARCHAR(64)  NOT NULL UNIQUE COMMENT '登录名',
    password    VARCHAR(128) NOT NULL COMMENT 'bcrypt 加密密码',
    nickname    VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '昵称',
    email       VARCHAR(128) NOT NULL DEFAULT '' COMMENT '邮箱',
    status      TINYINT      NOT NULL DEFAULT 1 COMMENT '1=启用 0=禁用',
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(64)  NOT NULL UNIQUE COMMENT '角色唯一标识，如 admin',
    label       VARCHAR(64)  NOT NULL DEFAULT '' COMMENT '角色显示名',
    status      TINYINT      NOT NULL DEFAULT 1,
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 用户-角色关联
CREATE TABLE IF NOT EXISTS user_roles (
    user_id     BIGINT UNSIGNED NOT NULL,
    role_id     BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (user_id, role_id)
);

-- 菜单表（树形结构）
CREATE TABLE IF NOT EXISTS menus (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    parent_id   BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '0 = 顶级菜单',
    name        VARCHAR(64)  NOT NULL COMMENT '菜单标识',
    title       VARCHAR(64)  NOT NULL COMMENT '菜单显示名',
    path        VARCHAR(128) NOT NULL DEFAULT '' COMMENT '前端路由路径',
    sort        INT          NOT NULL DEFAULT 0
);

-- API 资源表
CREATE TABLE IF NOT EXISTS apis (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    method      VARCHAR(16)  NOT NULL COMMENT 'GET/POST/PUT/DELETE',
    path        VARCHAR(256) NOT NULL COMMENT '接口路径，如 /api/v1/users',
    title       VARCHAR(128) NOT NULL COMMENT '接口名称',
    auth_mode   TINYINT      NOT NULL DEFAULT 0 COMMENT '0=无需登录 1=需登录 2=需权限',
    created_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_method_path (method, path)
);

-- 角色-API 权限关联
CREATE TABLE IF NOT EXISTS role_apis (
    role_id     BIGINT UNSIGNED NOT NULL,
    api_id      BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (role_id, api_id)
);


-- ========================================
-- 测试数据
-- ========================================

-- 角色（admin / editor / viewer）
INSERT INTO roles (id, name, label) VALUES
(1, 'admin',  '超级管理员'),
(2, 'editor', '内容编辑'),
(3, 'viewer', '只读用户');

-- 用户（密码均为 123456 的 bcrypt 哈希）
INSERT INTO users (id, username, password, nickname, email, status) VALUES
(1, 'admin',   '$2a$10$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPGga31lW', '管理员',  'admin@example.com',  1),
(2, 'editor1', '$2a$10$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPGga31lW', '编辑小王', 'editor1@example.com', 1),
(3, 'viewer1', '$2a$10$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPGga31lW', '游客小李', 'viewer1@example.com', 1),
(4, 'disabled','$2a$10$EixZaYVK1fsbw1ZfbX3OXePaWxn96p36WQoeG6Lruj3vjPGga31lW', '已禁用账号','disabled@example.com', 0);

-- 用户-角色绑定
INSERT INTO user_roles (user_id, role_id) VALUES
(1, 1),  -- admin → 超级管理员
(2, 2),  -- editor1 → 内容编辑
(3, 3),  -- viewer1 → 只读用户
(1, 2);  -- admin 同时拥有 editor 角色（测试多角色）

-- 菜单树
INSERT INTO menus (id, parent_id, name, title, path, sort) VALUES
(1, 0, 'dashboard',   '控制台',     '/dashboard',        1),
(2, 0, 'system',      '系统管理',   '/system',           2),
(3, 2, 'sys-user',    '用户管理',   '/system/users',     1),
(4, 2, 'sys-role',    '角色管理',   '/system/roles',     2),
(5, 2, 'sys-menu',    '菜单管理',   '/system/menus',     3),
(6, 2, 'sys-api',     'API管理',    '/system/apis',      4),
(7, 0, 'content',     '内容管理',   '/content',          3),
(8, 7, 'article',     '文章管理',   '/content/articles', 1);

-- API 资源
INSERT INTO apis (id, method, path, title, auth_mode) VALUES
(1,  'POST', '/api/v1/login',          '用户登录',     0),
(2,  'POST', '/api/v1/logout',         '用户登出',     1),
(3,  'GET',  '/api/v1/users',          '用户列表',     2),
(4,  'POST', '/api/v1/users',          '新建用户',     2),
(5,  'PUT',  '/api/v1/users/:id',      '更新用户',     2),
(6,  'DELETE','/api/v1/users/:id',     '删除用户',     2),
(7,  'GET',  '/api/v1/roles',          '角色列表',     2),
(8,  'POST', '/api/v1/roles',          '新建角色',     2),
(9,  'GET',  '/api/v1/menus',          '菜单列表',     1),
(10, 'GET',  '/api/v1/apis',           'API列表',      2),
(11, 'GET',  '/api/v1/ping',           '健康检查',     0);

-- 角色-API 权限（admin 拥有全部，editor 只有内容相关，viewer 只读）
-- admin 拥有所有 API
INSERT INTO role_apis (role_id, api_id)
SELECT 1, id FROM apis;

-- editor 可以登录/登出 + 菜单查看
INSERT INTO role_apis (role_id, api_id) VALUES
(2, 1), (2, 2), (2, 9), (2, 11);

-- viewer 只能登录/登出/健康检查
INSERT INTO role_apis (role_id, api_id) VALUES
(3, 1), (3, 2), (3, 11);