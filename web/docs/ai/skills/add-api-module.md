# Skill: 新增 API 模块

当前端需要新增后端 endpoint 调用时，使用本 skill。

## 步骤

1. 在 `src/request/<domain>` 下新增或更新模块。
2. 定义 request parameter 和 response payload 类型。
3. 使用 `src/request/request.ts` 中的共享 helpers。
4. 如果该模块需要对外使用，从 `src/request/index.ts` 导出。
5. 在 hook 或 feature service 中使用 API，不要散落在多个组件中直接调用。

## 模块模板

```ts
import { GET, POST } from '@src/request/request'
import type { fetchResponse } from '@src/request/request'

export type ExampleItem = {
  id: number
  name: string
}

export type ExampleListResponse = {
  list: ExampleItem[]
  total: number
}

export const getExampleListAPI = (): Promise<fetchResponse<ExampleListResponse>> => {
  return GET<ExampleListResponse>('/example/list')
}
```

## 规则

- 对已知 payload 不要使用 `any`。
- 不要在 feature module 中硬编码完整后端 origin。
- 不要在每个 API module 中重复 auth 处理。
- API function 名称要贴合领域，并表达调用意图。

## 验收标准

- request 和 response 类型明确。
- 调用方不需要了解底层 fetch 细节。
- error 和 auth 行为保持集中处理。
