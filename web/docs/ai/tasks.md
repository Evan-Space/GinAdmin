# AI 任务手册

当你要求 AI agent 修改前端时，使用本文档选择合适的 workflow。

## 新增后台页面

使用：

```text
web/docs/ai/skills/add-page.md
```

预期改动：

- `src/pages/_layout` 下的 route file
- 可选的 route-local helpers
- `src/layout/constants.tsx` 中的 menu item
- 如有需要，新增 API module

## 新增 API 模块

使用：

```text
web/docs/ai/skills/add-api-module.md
```

预期改动：

- typed DTOs
- endpoint functions
- 从 `src/request/index.ts` 导出
- 在 page 或 feature hook 中使用

## 新增表格管理页

使用：

```text
web/docs/ai/skills/add-table-page.md
web/docs/ai/skills/add-antd-table.md
web/docs/ai/skills/add-filter-toolbar.md
```

预期改动：

- filters
- table columns
- loading state
- empty state
- row actions
- 可选 create/edit drawer

## 新增或改造 Ant Design UI

使用：

```text
web/docs/ai/skills/use-antd.md
```

按具体场景继续使用：

```text
web/docs/ai/skills/add-antd-table.md
web/docs/ai/skills/add-antd-form.md
web/docs/ai/skills/add-antd-drawer-modal.md
web/docs/ai/skills/add-filter-toolbar.md
```

预期改动：

- 使用 Ant Design 组件作为主要 UI system
- 遵守 `src/style/antdTheme.ts`
- 补齐 loading、empty、error、disabled 状态
- 对危险操作增加二次确认

## 新增前端 AI 功能

使用：

```text
web/docs/ai/skills/add-ai-feature.md
web/docs/ai/skills/add-antd-drawer-modal.md
```

预期改动：

- AI feature folder
- mock AI client
- typed AI messages
- redaction helper
- assistant panel 或 inline suggestion UI

## 重构现有 UI

使用：

```text
web/docs/ai/skills/refactor-component.md
```

预期改动：

- 保持行为不变
- 提升可读性
- 保持 routing 和 API contracts 稳定
- 避免无关视觉重设计
