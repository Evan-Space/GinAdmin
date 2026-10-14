# Skill: 新增或改造 Ant Design Table

当新增或改造数据表格时，使用本 skill。

## 适用场景

- 用户列表、角色列表、权限列表。
- 请求日志、错误日志。
- 任务列表、系统配置列表。
- 任意以查询、列表、分页和行操作为核心的页面。

## 推荐文件组织

小型 route-owned table：

```text
src/pages/_layout/<page>/index.tsx
src/pages/_layout/<page>/-hooks.tsx
src/pages/_layout/<page>/-types.ts
src/pages/_layout/<page>/-constant.tsx
```

较大 feature：

```text
src/features/<feature>/
  components/
  hooks/
  services/
  types.ts
  constants.ts
```

## 类型规则

为 row data 定义明确类型：

```ts
export type UserListItem = {
  id: number
  username: string
  status: 'enabled' | 'disabled'
  createdAt: string
}
```

Table columns 使用 Ant Design 类型：

```ts
import type { TableColumnsType } from 'antd'

export const columns: TableColumnsType<UserListItem> = []
```

## 必需状态

- `loading`
- `empty`
- `error`
- pagination state，如果后端支持分页
- selected rows state，如果页面有批量操作

## 行操作规则

- 查看详情可以使用 `Drawer`。
- 编辑可以使用 `Drawer` 或独立页面。
- 删除、禁用、重置密码、授权等操作必须使用 `Popconfirm` 或 `Modal.confirm`。
- 操作列建议固定宽度，避免不同状态下抖动。

## AI 增强

适合 Table 的 AI 能力：

- 自然语言生成筛选条件。
- 解释选中行。
- 总结当前页数据。
- 解释异常日志。
- 推荐下一步操作。

AI 生成的筛选条件必须先展示 preview，再由用户点击应用。

## 常见错误

- columns 使用 `any`。
- dataSource 没有稳定 `rowKey`。
- loading 时整个布局跳动。
- 空数据没有说明。
- 错误只打印到 console。
- 删除操作没有二次确认。

## 验收标准

- `Table` 有明确泛型类型。
- `rowKey` 稳定。
- loading、empty、error 状态完整。
- row actions 稳定且危险操作有确认。
- 与页面筛选区和分页状态联动清晰。
