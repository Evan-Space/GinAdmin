# GinAdmin Web Claude 工作指南

本文档面向 Claude、Codex、Cursor 等 AI agent。进入 `web` 前端项目后，请优先按这里的规则理解、修改和扩展代码。

## 项目定位

`web` 是 GinAdmin 后台管理系统的前端项目。它不是营销官网，也不是展示型页面，而是面向用户管理、权限管理、日志、任务和系统设置等场景的管理控制台。

开发时应优先保证：

- 信息清晰
- 布局紧凑
- 操作稳定
- 状态明确
- 易于后续 AI agent 继续维护

## 技术栈

| 技术 | 用途 |
| --- | --- |
| Vite | 构建工具和开发服务器 |
| React | UI runtime |
| TypeScript | 类型系统 |
| Ant Design | 后台 UI 组件库 |
| TanStack Router | file-based routing |
| Zustand | 客户端 UI 状态 |
| ahooks | 请求和常用 hooks |
| Tailwind CSS | 少量布局、间距和尺寸 utilities |
| pnpm | 包管理 |

## 重要路径

```text
src/main.tsx                  App 入口和 RouterProvider
src/pages/__root.tsx          根路由，挂载 ConfigProvider 和 App
src/pages/_layout.tsx         后台布局路由
src/pages/_layout/            后台管理页面
src/layout/                   Header、Sider、Content、Footer
src/layout/constants.tsx      侧边栏菜单配置
src/request/                  API 请求层
src/store/                    Zustand stores
src/style/antdTheme.ts        Ant Design theme tokens
src/style/global.css          Tailwind import 和全局基础样式
src/routeTree.gen.ts          TanStack Router 生成文件，不要手动编辑
```

## 必读上下文

开始较大改动前，先阅读：

```text
.cursor/rules/project.mdc
docs/ai/project-overview.md
docs/ai/frontend-architecture.md
docs/ai/tasks.md
```

按任务类型继续阅读：

```text
.cursor/rules/react-frontend.mdc      React、TypeScript、hooks、store
.cursor/rules/router.mdc              TanStack Router 路由
.cursor/rules/api-request.mdc         API 请求和 DTO
.cursor/rules/antd-ui.mdc             Ant Design 和后台 UI
.cursor/rules/ai-feature.mdc          前端 AI 功能
docs/ai/security-and-privacy.md       AI 安全与隐私
docs/ai/ai-feature-blueprint.md       AI 功能蓝图
docs/ai/skills/use-antd.md            Ant Design 总规则
```

## 开发命令

```bash
pnpm dev
pnpm build
pnpm preview
```

当前 `package.json` 中没有真正可用的测试脚本。需要验证时，优先运行 `pnpm build`。

## 路由规则

- 使用 TanStack Router file-based routes。
- 后台页面放在 `src/pages/_layout` 下。
- 路由文件使用 `createFileRoute`。
- 不要手动编辑 `src/routeTree.gen.ts`。
- 如果新增页面需要出现在侧边栏，更新 `src/layout/constants.tsx`。
- 菜单 key 要和路由路径保持一致。

## 组件规则

- 使用函数组件。
- route component 保持轻量，复杂 UI 拆到 components 或 feature 目录。
- 数据请求和副作用放到 hooks 中。
- 新代码尽量避免 `any`。
- 面向用户的页面需要明确 loading、empty、error 和 disabled 状态。
- 后台 UI 优先使用 Ant Design 组件。
- Tailwind CSS 只用于少量 utilities，不要另起一套设计系统。

## API 规则

- 共享请求逻辑位于 `src/request/request.ts`。
- 领域 API 模块放在 `src/request/<domain>` 下。
- 公共 API 从 `src/request/index.ts` 导出。
- 新增 API 时定义 request params 和 response payload 类型。
- 不要在 feature module 中硬编码完整后端 origin。
- auth、401、错误处理应集中在共享 request helper 中。

## AI 化规则

前端 AI 化的目标是让项目更容易被 AI 理解和扩展，同时为后续 AI 辅助 UI 留好结构。

允许前端做：

- 收集用户意图
- 收集安全的页面 context
- 渲染 AI suggestions
- 渲染 streaming text
- 调用后端 AI endpoints 或 mock client

禁止前端做：

- 保存模型服务商 API Key
- 直接调用 OpenAI、Anthropic、Gemini 等模型服务商
- 绕过当前登录态或权限规则
- 自动执行删除、授权、导出、修改系统设置等高影响操作

AI 输出默认是建议。应用 AI 生成的 mutation 前，必须要求用户明确确认。

## 推荐 Feature 结构

较大功能优先使用：

```text
src/features/<feature-name>/
  components/
  hooks/
  services/
  types.ts
  constants.ts
```

小型路由专属代码可以沿用现有模式：

```text
src/pages/_layout/<page>/index.tsx
src/pages/_layout/<page>/-hooks.tsx
src/pages/_layout/<page>/-types.ts
src/pages/_layout/<page>/-constant.tsx
```

## 常用 Skill 文档

```text
docs/ai/skills/add-page.md           新增后台页面
docs/ai/skills/add-api-module.md     新增 API 模块
docs/ai/skills/add-table-page.md     新增表格管理页
docs/ai/skills/add-ai-feature.md     新增前端 AI 功能
docs/ai/skills/refactor-component.md 重构组件
docs/ai/skills/use-antd.md           Ant Design 总规则
docs/ai/skills/add-antd-table.md     新增或改造 Table
docs/ai/skills/add-antd-form.md      新增或改造 Form
docs/ai/skills/add-antd-drawer-modal.md Drawer / Modal 工作流
docs/ai/skills/add-filter-toolbar.md 筛选区和工具栏
```

## 修改原则

- 优先贴合现有代码风格。
- 不做无关重构。
- 不改动用户已有修改，除非任务明确要求。
- 不引入不必要的新依赖。
- 不把后端职责塞进前端。
- 修改架构性约定后，同步更新 `docs/ai` 或 `.cursor/rules`。
