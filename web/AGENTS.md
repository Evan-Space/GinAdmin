# GinAdmin Web AI Agent 指南

本文件是 `web` 前端项目的通用 AI agent 入口，适用于 Codex、Claude、Cursor 和其他能读取仓库上下文的编程助手。

## 项目定位

`web` 是 GinAdmin 后台管理系统的前端项目，面向用户管理、权限管理、请求日志、错误日志、任务中心和系统设置等管理场景。

开发时请把它当作现代后台控制台，而不是营销官网或展示页。优先保证信息清晰、操作稳定、状态明确、布局紧凑。

## 技术栈

- Vite
- React
- TypeScript
- Ant Design
- TanStack Router
- Zustand
- ahooks
- Tailwind CSS
- pnpm

## 必读文件

开始较大修改前，先阅读：

```text
web/.cursor/rules/project.mdc
web/docs/ai/project-overview.md
web/docs/ai/frontend-architecture.md
web/docs/ai/tasks.md
```

按任务类型继续阅读：

```text
web/.cursor/rules/react-frontend.mdc
web/.cursor/rules/router.mdc
web/.cursor/rules/api-request.mdc
web/.cursor/rules/antd-ui.mdc
web/.cursor/rules/ai-feature.mdc
web/docs/ai/security-and-privacy.md
web/docs/ai/ai-feature-blueprint.md
```

## 常用 Skill

```text
web/docs/ai/skills/add-page.md
web/docs/ai/skills/add-api-module.md
web/docs/ai/skills/add-table-page.md
web/docs/ai/skills/add-ai-feature.md
web/docs/ai/skills/refactor-component.md
web/docs/ai/skills/use-antd.md
web/docs/ai/skills/add-antd-table.md
web/docs/ai/skills/add-antd-form.md
web/docs/ai/skills/add-antd-drawer-modal.md
web/docs/ai/skills/add-filter-toolbar.md
```

## 关键路径

```text
web/src/main.tsx
web/src/pages/__root.tsx
web/src/pages/_layout.tsx
web/src/pages/_layout/
web/src/layout/
web/src/layout/constants.tsx
web/src/request/
web/src/store/
web/src/style/antdTheme.ts
web/src/style/global.css
web/src/routeTree.gen.ts
```

不要手动编辑 `web/src/routeTree.gen.ts`。

## 开发命令

在 `web` 目录下运行：

```bash
pnpm dev
pnpm build
pnpm preview
```

当前 `package.json` 中没有真正可用的测试脚本。需要验证前端可构建时，优先运行 `pnpm build`。

## 工作原则

- 只修改与当前任务相关的文件。
- 不做无关重构。
- 不覆盖用户已有改动。
- 新代码遵守 TypeScript strict 配置。
- 新增页面使用 TanStack Router file-based routes。
- 新增后台 UI 优先使用 Ant Design。
- Tailwind CSS 只用于简单布局、间距和尺寸 utilities。
- 新增 API 调用要定义 request 和 response 类型。
- 前端不能保存或暴露模型服务商 API Key。
- 前端 AI 功能必须调用后端 API 或 mock client，不能直接调用模型服务商。

## Ant Design 规则

Ant Design 是本项目的主要 UI system。修改表格、表单、筛选区、Drawer、Modal、布局和反馈状态时，优先阅读：

```text
web/docs/ai/skills/use-antd.md
web/docs/ai/skills/add-antd-table.md
web/docs/ai/skills/add-antd-form.md
web/docs/ai/skills/add-antd-drawer-modal.md
web/docs/ai/skills/add-filter-toolbar.md
```

通用要求：

- 使用 `src/style/antdTheme.ts` 中已有主题。
- 数据页面必须有 loading、empty 和 error 状态。
- 表单提交 pending 时禁用提交按钮。
- 删除、授权、导出、修改系统设置等高影响操作必须二次确认。
- AI 生成的 UI 建议不能自动应用，必须由用户确认。

## AI 化边界

允许前端做：

- 收集用户意图
- 收集安全的页面 context
- 渲染 AI suggestions
- 渲染 streaming text
- 调用后端 AI endpoints 或 mock client

禁止前端做：

- 保存 OpenAI、Anthropic、Gemini 等模型服务商 API Key
- 直接从浏览器调用模型服务商 API
- 绕过当前登录态或权限规则
- 自动执行删除、授权、导出、修改系统设置等高影响操作

AI 输出默认是建议。应用任何 AI 生成的 mutation 前，必须要求用户明确确认。
