# 路由指南

## Router

项目使用 TanStack Router 的 file-based routing。

重要文件：

```text
src/pages/__root.tsx
src/pages/_layout.tsx
src/routeTree.gen.ts
vite.config.ts
```

Vite router plugin 在 `vite.config.ts` 中配置：

```ts
tanstackRouter({
  target: 'react',
  autoCodeSplitting: true,
  routesDirectory: './src/pages',
  generatedRouteTree: './src/routeTree.gen.ts',
})
```

## 不要编辑生成的 Route Tree

不要手动编辑 `src/routeTree.gen.ts`。它由 `src/pages` 下的文件自动生成。

## 当前路由模型

- `__root.tsx` 是根级 wrapper。
- `_layout.tsx` 是后台管理布局路由。
- `_layout/*` 包含登录后的后台页面。
- `login/index.tsx` 位于后台布局之外。

## 新增后台页面

示例：新增角色管理。

1. 创建路由文件：

```text
src/pages/_layout/role/index.tsx
```

2. 定义路由：

```tsx
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_layout/role/')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Role management</div>
}
```

3. 在 `src/layout/constants.tsx` 添加菜单项：

```tsx
{ key: '/role', label: '角色管理', icon: <UserOutlined /> }
```

4. 运行前端 build 或 dev server，让 route tree 重新生成。

## 路由文件建议

- 路由文件保持轻量。
- 表格列定义、service 调用和大块 UI 移到同级辅助文件或 `src/features` 中。
- 路由按产品领域命名，不按实现细节命名。
