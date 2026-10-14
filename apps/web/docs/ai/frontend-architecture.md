# 前端架构

## 启动流程

`src/main.tsx` 基于生成的 route tree 创建 TanStack Router 实例，并渲染 `RouterProvider`。

`src/pages/__root.tsx` 使用以下组件包裹应用：

- Ant Design `ConfigProvider`
- Ant Design `App`
- 根路由 `Outlet`

`src/pages/_layout.tsx` 挂载来自 `src/layout` 的后台共享布局。

## 布局

布局拆分为：

```text
src/layout/index.tsx
src/layout/components/Header.tsx
src/layout/components/Sider.tsx
src/layout/components/Content.tsx
src/layout/components/Footer.tsx
src/layout/constants.tsx
```

菜单配置位于 `src/layout/constants.tsx`。菜单 key 要与路由路径保持一致。

## 页面

页面存放在 `src/pages` 下。登录后的后台页面位于 `src/pages/_layout` 下。

现有页面示例：

```text
src/pages/_layout/index.tsx
src/pages/_layout/userList/index.tsx
src/pages/_layout/permission.tsx
src/pages/_layout/log/requestLog.tsx
src/pages/_layout/log/errorLog.tsx
src/pages/_layout/task.tsx
src/pages/_layout/setting.tsx
src/pages/login/index.tsx
```

## 路由本地辅助文件

用户列表页当前使用路由本地辅助文件：

```text
src/pages/_layout/userList/-hooks.tsx
src/pages/_layout/userList/-constant.tsx
src/pages/_layout/userList/-types.ts
```

这种约定适合小型、只归属于该路由的代码。对于可复用或更大的功能，优先使用 `src/features/<feature>`。

## 请求层

共享请求 helper 位于 `src/request/request.ts`。

接口模块位于 `src/request/<domain>` 下。对外公共 API 从 `src/request/index.ts` 统一导出。

## 状态

Zustand store 位于 `src/store` 下。store 应保持小而聚焦，主要管理客户端 UI 状态。

服务端状态通常应放在请求 hook 或未来的 server-cache library 中，不要默认放进全局 store。

## 样式

产品级视觉风格使用 `src/style/antdTheme.ts` 中的 Ant Design theme tokens。

Tailwind utilities 只用于简单布局、间距和尺寸。避免用 Tailwind 再搭一套设计系统。

## 推荐 Feature 结构

新增较大功能时，推荐使用：

```text
src/features/<feature-name>/
  components/
  hooks/
  services/
  types.ts
  constants.ts
```

Route files can then compose the feature.
