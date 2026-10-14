# Skill: 使用 Ant Design

当新增或修改后台 UI 时，使用本 skill。Ant Design 是本项目的主要 UI system，Tailwind CSS 只作为少量 utilities 辅助。

## 适用场景

- 新增页面布局。
- 新增或修改 Table。
- 新增或修改 Form。
- 新增 Drawer、Modal、Popconfirm。
- 新增筛选区、工具栏、状态标签、操作按钮。
- 调整 Ant Design theme token。

## 组件选择

- 页面骨架使用 `Layout`、`Space`、`Flex`。
- 数据列表使用 `Table`。
- 筛选区使用 `Form`、`Input`、`Select`、`DatePicker`、`Button`。
- 编辑流程优先使用 `Drawer`。
- 确认、警告、少量字段编辑可以使用 `Modal`。
- 危险操作使用 `Popconfirm` 或 `Modal.confirm`。
- 状态展示使用 `Tag`。
- 页面反馈使用 Ant Design `App` 提供的 `message` 和 `notification`。
- 图标优先使用 `@ant-design/icons`。

## 视觉规则

- 遵守 `src/style/antdTheme.ts` 中的主题 token。
- 保持后台系统的紧凑、克制、清晰。
- 不要把后台页面做成 landing page 或营销页。
- 不要用 Tailwind 重写 Ant Design 的完整视觉体系。
- 不要嵌套无意义的 Card。
- 工具栏、筛选区、表格、分页和操作区在 loading 前后要保持稳定。

## 状态规则

- 数据页面必须有 loading state。
- 无数据时必须有 empty state。
- 请求失败时必须有 error feedback。
- 表单提交中必须禁用提交按钮。
- 破坏性操作必须二次确认。

## AI 功能规则

- AI suggestions 应以辅助建议展示，不应伪装成系统事实。
- AI 生成的筛选条件、表单值、权限配置等必须先预览，再由用户确认应用。
- AI UI 推荐使用紧凑 `Drawer`、inline suggestion 或日志解释 panel。

## 常见错误

- 用一堆 `div` 和 Tailwind 手写已有 Ant Design 组件。
- 表格没有 loading、empty、error 状态。
- 删除按钮没有确认。
- 表单提交时按钮仍可重复点击。
- 把 AI 建议直接写入业务数据。

## 验收标准

- UI 使用 Ant Design 作为主要组件体系。
- 风格与 `src/style/antdTheme.ts` 一致。
- 状态反馈完整。
- TypeScript 类型明确。
- 代码结构便于后续 AI agent 继续维护。
