# Skill: 新增 Drawer 或 Modal 工作流

当新增二级操作、详情、编辑、确认或 AI 辅助面板时，使用本 skill。

## 选择 Drawer 还是 Modal

优先使用 `Drawer`：

- 查看详情。
- 新增或编辑较多字段。
- AI assistant panel。
- 日志解释。
- 不希望打断主页面上下文的流程。

优先使用 `Modal`：

- 少量字段编辑。
- 删除、禁用、授权等高影响确认。
- 需要用户明确选择确认或取消。
- 需要强打断的警告。

## Drawer 规则

- 标题要说明当前动作，例如 `编辑用户`、`请求日志详情`。
- 宽度要适合内容，不要铺满整个屏幕。
- 表单型 Drawer 底部放取消和确认按钮。
- 打开时初始化数据，关闭时清理临时状态。
- 提交 pending 时禁用确认按钮。

## Modal 规则

- 内容要简短清晰。
- 危险操作标题和按钮文案要明确。
- 使用 `okButtonProps={{ danger: true }}` 标识危险确认。
- 不要把大表单塞进很小的 Modal。

## AI Assistant Drawer

AI assistant 建议使用右侧 `Drawer`：

- 保持紧凑。
- 显示当前 page context。
- 支持 loading、streaming、error、cancel。
- AI suggestions 与 apply action 分离。
- 应用任何 mutation 前要求用户确认。

## 状态管理

推荐状态：

```ts
type DrawerMode = 'create' | 'edit' | 'detail'
```

需要保存：

- open
- mode
- current record
- submitting
- error

## 常见错误

- Drawer 关闭后没有清理 record。
- 编辑成功后列表没有刷新。
- Modal 文案模糊。
- 危险操作没有 danger 样式。
- AI Drawer 直接执行修改。

## 验收标准

- Drawer/Modal 用途选择合理。
- open、close、submit、cancel 状态清晰。
- pending 状态下不可重复提交。
- 高影响操作有明确确认。
- AI assistant 只提供建议，不自动执行 mutation。
