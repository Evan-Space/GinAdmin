# Skill: 新增后台页面

当需要在登录后的后台区域新增页面时，使用本 skill。

## 需要明确的输入

- 页面中文名和英文名。
- Route path。
- 是否需要出现在侧边栏菜单。
- 是否需要后端数据。
- 是否需要 table、form、drawer、modal 或 tabs。

## 步骤

1. 在 `src/pages/_layout` 下创建 route file。
2. 使用 `createFileRoute` 定义路由。
3. route component 保持轻量。
4. 只有页面专属逻辑才添加 route-local helpers。
5. 较大页面应在 `src/features/<feature>` 下添加 feature-level files。
6. 如果页面需要导航，在 `src/layout/constants.tsx` 添加菜单项。
7. 使用 Ant Design primitives 构建布局和控件。
8. 如果页面请求数据，需要包含 loading、empty 和 error 状态。

## Route 模板

```tsx
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_layout/example/')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Example</div>
}
```

## 验收标准

- Route 能被 TanStack Router 发现。
- 没有手动编辑 `src/routeTree.gen.ts`。
- 如果添加菜单项，menu key 与 route path 匹配。
- 页面符合现有后台视觉风格。
- 遵守 TypeScript strict mode。
