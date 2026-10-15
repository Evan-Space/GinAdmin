# Skill: 新增筛选区和工具栏

当为表格页新增查询筛选、批量操作、刷新、重置、AI 筛选建议时，使用本 skill。

## 适用场景

- 表格顶部筛选区。
- 列表页主操作工具栏。
- 批量操作区。
- 自然语言筛选建议入口。

## 推荐布局

基础结构：

```tsx
<Space direction="vertical" size="middle" style={{ display: 'flex' }}>
  <Form layout="inline">{/* filters */}</Form>
  <Table />
</Space>
```

复杂页面可以将筛选区抽成组件：

```text
components/<Feature>FilterToolbar.tsx
```

## 筛选区规则

- 查询条件使用 Ant Design `Form`。
- 常用字段默认展开。
- 低频字段可以放到高级筛选中。
- 查询按钮使用 primary。
- 重置按钮必须清空表单并刷新查询。
- 字段 label 使用业务语言，不用接口字段名。

## 工具栏规则

- 主操作靠左或靠近筛选区。
- 刷新、导出、批量操作等次级操作放在视觉次级位置。
- 批量操作需要依赖 selected rows。
- 导出、批量删除、批量禁用等高影响操作必须二次确认。

## AI 筛选建议

可以增加自然语言筛选入口，例如：

- 用户输入“查最近 7 天失败的登录请求”。
- AI 返回 filter draft。
- UI 展示将要应用的字段和值。
- 用户确认后写入 Form 并触发查询。

AI filter draft 不允许静默应用。

## 状态规则

- 查询 pending 时禁用查询按钮。
- 重置 pending 时避免重复触发。
- error 时保留用户筛选条件。
- loading 不应导致筛选区跳动。

## 常见错误

- 筛选表单和请求参数没有清晰映射。
- 重置只清空 UI，没有重新查询。
- 高级筛选挤占表格空间。
- AI 返回筛选条件后直接查询，没有 preview。

## 验收标准

- 筛选字段和后端 query params 映射明确。
- 查询、重置、刷新行为清晰。
- loading 和 error 状态可见。
- 批量操作依赖 selected rows。
- AI 筛选建议必须先 preview，再应用。
