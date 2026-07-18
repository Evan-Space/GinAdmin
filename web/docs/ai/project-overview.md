# GinAdmin Web AI 上下文

本文档为 AI agent 提供在 `web` 前端项目中安全工作的最小必要上下文。

## 项目目的

`web` 是 GinAdmin 的浏览器前端。它是一个后台管理系统，用于登录后的管理流程，例如用户管理、权限管理、日志、任务和系统设置。

前端的 AI 化分为两个方向：

- AI agent 能基于一致规则理解和修改代码。
- UI 后续可以承载 AI 辅助管理流程，同时不在浏览器暴露密钥。

## 当前技术栈

- 构建工具：Vite
- UI runtime：React
- 语言：TypeScript
- 路由：TanStack Router file-based routes
- UI 组件库：Ant Design
- 图标：`@ant-design/icons`
- 请求生命周期辅助：ahooks
- 客户端状态：Zustand
- 样式：Ant Design theme tokens + Tailwind CSS utilities

## 重要路径

```text
web/src/main.tsx                  App bootstrap and router setup
web/src/pages/__root.tsx          Root route, Ant Design provider, outlet
web/src/pages/_layout.tsx         Authenticated layout route
web/src/pages/_layout/            Admin pages
web/src/layout/                   Header, sider, content, footer shell
web/src/layout/constants.tsx      Menu configuration
web/src/request/                  API request helpers and endpoint modules
web/src/style/antdTheme.ts        Ant Design theme tokens
web/src/style/global.css          Tailwind import and global base styles
web/src/store/                    Zustand stores
web/src/routeTree.gen.ts          Generated route tree, do not edit manually
```

## 产品方向

这不是内容站点。产品应保持现代后台管理系统的气质：

- 紧凑
- 结构化
- 安静
- 可靠
- 易扫描
- 对破坏性操作谨慎

## AI 方向

适合此前端的 AI 能力：

- 后台布局内的 assistant drawer
- 表格自然语言筛选建议
- 错误日志解释 UI
- 请求日志总结 UI
- 权限草案预览
- 表单填写建议

不适合的 AI 做法：

- 浏览器侧保存模型 API Key
- 隐藏式自动修改数据
- 忽略当前页面上下文的泛泛 chatbot UI
- 绕过当前登录态或权限规则的 AI 功能

## 不可妥协项

- 前端代码中不能出现模型服务商密钥。
- 浏览器不能直接调用模型服务商。
- AI 生成的变更在用户确认前必须只是建议。
- 不要手动编辑生成的路由文件。
- 保持 TypeScript 严格配置。
