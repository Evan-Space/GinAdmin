# AI-Ready Frontend 设计文档

## 目标

为 GinAdmin 前端项目创建 AI 辅助配置，让 Cursor/Claude Code/Copilot 等 AI 工具生成符合项目约定的代码。

## 技术栈

- React 19 + TypeScript 7 + Vite 8
- Ant Design 6 + Tailwind CSS 4
- TanStack Router（文件路由，自动生成 routeTree）
- Zustand（状态管理）
- ahooks（工具 hooks）
- pnpm（包管理）

## 文件结构

```
web/
├── CLAUDE.md                          # 入口：项目概述 + 技术栈 + 指向各规则
├── .cursor/
│   └── rules/
│       ├── architecture.mdc           # 路由机制、目录结构、状态管理
│       ├── components-and-styling.mdc # Ant Design 6 + Tailwind CSS 4 规范
│       ├── api-and-request.mdc        # request 层封装、API 模块约定
│       └── conventions.mdc            # 命名、TypeScript、import 顺序
```

## 各文件内容设计

### 1. CLAUDE.md（入口文件）

- 项目一句话简介
- 技术栈列表（每项说明用途）
- 指向 `.cursor/rules/` 下各规则文件
- 关键约定速览：`@src/` 路径别名、pnpm 命令、开发服务器端口 3002

### 2. architecture.mdc（架构规约）

**触发场景：** 创建新页面、修改路由、调整布局、使用状态管理

**内容要点：**
- **路由机制：** TanStack Router 文件路由，`src/pages/` 下文件自动生成 `routeTree.gen.ts`；文件前缀 `-` 不参与路由（如 `-constant.tsx`, `-hooks.tsx`, `-types.ts`）
- **目录结构：**
  - `layout/` — 布局组件（Header, Sider, Content, Footer）
  - `pages/` — 页面，按功能分目录
  - `store/` — Zustand store
  - `request/` — API 请求层的封装
  - `style/` — 全局样式、Ant Design 主题配置
- **Zustand 约定：** store 文件放在 `src/store/`，使用 Zustand 5.x API
- **新增页面步骤：** 在 `src/pages/` 下创建目录 → 创建 `index.tsx` → 路由自动生成

### 3. components-and-styling.mdc（组件与样式规范）

**触发场景：** 编写 UI 组件、样式

**内容要点：**
- **Ant Design 6 规范：** 明确标注使用 Ant Design 6.x API（非 5.x），常用组件写法示例
- **Tailwind CSS 4：** 优先使用 Tailwind utility 类处理间距、颜色、布局；复杂组件内部样式使用 Ant Design 的 `style`/`className` prop
- **组件拆分约定：**
  - `index.tsx` — 页面视图
  - `-constant.tsx` — 常量（表格列定义等）
  - `-hooks.tsx` — 业务逻辑 hook
  - `-types.ts` — 类型定义
- **主题定制：** `src/style/antdTheme.ts` 中配置 Ant Design 主题 token

### 4. api-and-request.mdc（请求层规范）

**触发场景：** 新增 API 调用、修改请求逻辑

**内容要点：**
- **封装实例位置：** `src/request/request.ts`
- **模块拆分约定：** 按业务域在 `src/request/` 下创建子目录（如 `login/`, `userList/`），每个目录一个 `index.ts` 导出该模块的 API 方法
- **登录 Token：** token 获取和传递方式（从 request.ts 中推断）
- **统一错误处理：** 错误处理方式（从 request.ts 中推断）
- **调用方式示例：** 在页面中如何引入和调用 API

### 5. conventions.mdc（编码约定）

**触发场景：** 任何代码编写

**内容要点：**
- **TypeScript：** strict 模式，必须定义类型；`noUnusedLocals`/`noUnusedParameters` 已开启
- **命名规范：**
  - 组件 PascalCase，文件 camelCase
  - Hook 以 `use` 开头
  - 辅助文件前缀 `-`
- **Import 顺序：** 第三方库 → `@src/` 别名 → 相对路径
- **路径别名：** `@src/*` → `./src/*`
- **Prettier：** 已配置，无需修改

## 实施注意事项

- 创建文件前需阅读现有源码，确保规则与代码实际行为一致
- 规则要具体，用真实文件路径和代码片段作为示例
- 每个 `.mdc` 文件配置 `globs` 和 `alwaysApply` front-matter，精确控制触发范围
