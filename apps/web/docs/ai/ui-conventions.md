# UI 约定

## 产品气质

GinAdmin 是后台管理系统。界面应当：

- 冷静
- 紧凑
- 可读
- 可预测
- 适合操作

避免营销式布局、过大的 hero text、纯装饰视觉和大面积空 card。

## Design System

使用 Ant Design 作为主要 UI system。

Use:

- `Layout` for shell structure
- `Menu` for navigation
- `Table` for data grids
- `Form`, `Input`, `Select`, `DatePicker` for filters and editing
- `Drawer` for secondary workflows
- `Modal` for confirmation or focused editing
- `App` message and notification APIs for feedback
- `Tag` for status and categories

Tailwind CSS 只在能让代码更清晰的小型布局场景中使用。

## Theme

Theme tokens 位于：

```text
src/style/antdTheme.ts
```

尊重现有的中性、克制主题。只有重复 UI 确实需要时，才添加新的 theme token。

## Tables

表格页面应当：

- 将筛选器放在表格上方
- 主操作靠近筛选区或表格 header
- 暴露 loading state
- 暴露 empty state
- 保持行操作视觉稳定
- 对破坏性操作二次确认

## Forms

表单应当：

- 使用 validation rules
- submitting 时禁用提交
- 清楚展示请求错误
- 避免静默丢弃用户输入

## AI UI

推荐的 AI 交互位置：

- 右侧 assistant drawer
- 表单或表格附近的 inline suggestions
- 日志页面的可展开 explanation panel
- 不阻塞主流程的 recommendation area

AI UI 应展示：

- 使用了什么 context
- 回答是否正在 streaming
- 某个 action 是否只是建议
- 应用变更前的明确 confirmation step

不要让 AI 输出看起来像系统事实。它是辅助建议。
