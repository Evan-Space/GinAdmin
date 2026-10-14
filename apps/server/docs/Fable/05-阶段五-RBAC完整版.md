# 阶段五：RBAC 权限体系完整版

> 目标：把阶段一跑通的"最小 Casbin"升级成 gin-layout 那套完整权限体系：
> API 管理 → 菜单绑 API → 角色绑菜单 → 用户绑角色 → Casbin 策略自动同步。
> 这是整个项目最有含金量的部分，预计 7~10 天，不要赶。

## 权限链路全景（先看懂再动手）

```text
路由树(代码里声明的 method+path+auth_mode)
      │ 启动时同步
      ▼
apis 表(API 资产清单)
      ▼
menu_api_map(菜单 → 需要哪些 API)
      ▼
role_menu_map(角色 → 拥有哪些菜单)
      ▼
admin_user_role_map(用户 → 属于哪些角色)
      │ 任一环节变化时
      ▼
Casbin 策略(p: role:x, path, method / g: admin_user:y, role:x)
      ▼
AdminAuthHandler 中间件 Enforce
```

前端拿到的「菜单树 + 按钮权限」和后端的「接口能不能调」都由这条链驱动 —— 这正是你做前端时对接的权限系统的服务端全貌。

## 任务 5.1 API 资产管理

1. 建 `apis` 表 + `internal/model/api.go`：method、path、title、auth_mode、group 等。
2. **启动时自动同步**：遍历路由树（阶段一思考题埋的伏笔），把所有路由 upsert 进 apis 表；代码里删掉的路由标记为失效。
3. 接口：`GET /api/v1/api/list`（给前端权限配置页用）。

对照阅读：

```text
gin-layout/internal/routers/meta.go
gin-layout/internal/service/api_permission/api.go
gin-layout/internal/model/api.go
```

思考题：

- 「API 清单以代码为准自动同步」对比「后台手动录入」，各有什么优缺点？

## 任务 5.2 菜单绑 API

1. 你已有 `menu_api_map` 表结构，把它用起来：
   - 菜单编辑接口支持传 `api_ids`
   - `GET /api/v1/menu/detail` 返回已绑定的 API
2. 菜单类型区分：目录 / 菜单页 / 按钮。按钮也是菜单树节点（前端按钮级权限就是这么来的）。

对照阅读：

```text
gin-layout/internal/model/menu.go
gin-layout/internal/model/menu_api_map.go
gin-layout/internal/service/menu/menu_edit.go
```

## 任务 5.3 Casbin 策略自动同步（本阶段核心）

把你阶段一的 `casbin_sync.go` 升级成完整版。触发点：

| 事件 | 同步动作 |
|---|---|
| 角色绑定/解绑菜单 | 重算该角色的 p 规则 |
| 菜单绑定的 API 变化 | 重算所有含该菜单的角色的 p 规则 |
| 用户绑定/解绑角色 | 更新该用户的 g 规则 |
| 删除角色/菜单/用户 | 清理相关规则 |

要求：

1. 同步逻辑收敛到一个 service（如 `internal/service/access/`），其他 service 只调它，不直接摸 enforcer。
2. **数据库写入和策略更新要在同一事务语义下**：菜单绑定写库成功但策略同步失败时怎么办？（看 gin-layout 的 `transaction.go` 和 `coordinator.go` 的思路：失败回滚或重载兜底）
3. 同步后调 `LoadPolicy()` 或增量更新，保证运行中的服务立即生效。

对照阅读：

```text
gin-layout/internal/service/access/user_permission_sync.go
gin-layout/internal/service/access/graph_loader.go
gin-layout/internal/service/access/coordinator.go
gin-layout/internal/service/access/transaction.go
```

思考题：

- gin-layout 为什么把「读权限图谱」（graph_loader）和「写策略」分开？
- 全量重算一个角色的策略 vs 增量 diff，gin-layout 选了哪种？为什么？

## 任务 5.4 用户菜单接口（打通前端）

1. `GET /api/v1/get-menus`（AuthLogin）：返回当前用户可见的菜单树 + 按钮权限标识。
2. 超管返回全部；普通用户按角色合并去重。
3. 用你熟悉的任意前端框架（或 Postman）验证：不同角色的用户拿到不同菜单。

对照阅读：

```text
gin-layout/internal/service/admin/admin_user.go（GetUserMenus 相关）
gin-layout/internal/service/menu/menu_query.go
```

## 任务 5.5 resources 层（响应视图）

你会发现 model 直接序列化给前端有问题：密码字段、内部字段泄露，格式不友好。gin-layout 用 `internal/resources/` 解决：

1. 新建 `internal/resources/base.go` + `admin_user.go` + `role.go` + `menu.go`。
2. 每个 resource 定义「这个模型对外长什么样」：隐藏 password、格式化时间、组装关联数据（用户的角色名列表等）。
3. controller 返回统一改为 `resources.XxxResource(model)`。

对照阅读：

```text
gin-layout/internal/resources/base.go
gin-layout/internal/resources/admin_user.go
```

思考题：

- resource 层和前端的 DTO/ViewModel 是不是一回事？没有它时序列化 model 会踩什么坑（提示：`json:"-"` 忘写）？

## 任务 5.6 部门与数据权限（选做，可延后）

如果精力允许：建 dept 表 + 用户关联部门，实现「本部门数据」「本部门及子部门」的数据范围过滤。对照 `gin-layout/internal/service/dept/` 和 `internal/service/access/scope_resolver.go`。做不动可以跳过，不阻塞晋级。

## 晋级标准

```text
[ ] 启动时 apis 表与路由树自动同步
[ ] 菜单可绑 API，角色可绑菜单，用户可绑角色，任一变更后 Casbin 策略自动更新且立即生效
[ ] 完整演示：新建角色 → 绑菜单 → 新建用户绑角色 → 该用户只能调被授权的接口(403/200 对比)
[ ] get-menus 按角色返回正确的菜单树
[ ] 所有对外响应经过 resources 层，任何接口都不返回 password 等敏感字段
[ ] casbin 同步 service 有单元测试（可用 sqlite 内存库或 mock）
[ ] 阶段五笔记完成（重点：画出你自己理解的权限链路图）
```
