# API 约定

## 当前请求层

共享 helpers：

```text
src/request/request.ts
```

领域模块：

```text
src/request/login/index.tsx
src/request/userList/index.ts
```

公共 re-export：

```text
src/request/index.ts
```

## Response Envelope

当前前端期望后端返回类似下面的 envelope：

```ts
export interface fetchResponse<T> {
  code: number
  msg: string
  cost: number
  request_id: string
  data: T
}
```

后续重构请求层时，建议将类型名改为 PascalCase：`FetchResponse<T>`。

## Request Client 目标

共享 request client 后续应支持：

- `VITE_API_BASE_URL`
- 请求时读取 token
- 可选 auth header
- JSON body 处理
- URL query params
- typed response envelope
- 集中处理 401
- 集中处理网络错误和 JSON 解析错误
- 为 AI endpoint 提供可选 streaming request helper

## API 模块模式

新增 API 模块优先使用这种结构：

```ts
import { GET, POST } from '@src/request/request'
import type { FetchResponse } from '@src/request/request'

export type UserListItem = {
  id: number
  name: string
}

export type UserListResponse = {
  list: UserListItem[]
  total: number
}

export const getUserListAPI = (): Promise<FetchResponse<UserListResponse>> => {
  return GET<UserListResponse>('/admin-user/list')
}
```

## AI API 边界

前端 AI services 应调用后端拥有的接口，例如：

```text
POST /api/v1/ai/chat
POST /api/v1/ai/chat/stream
POST /api/v1/ai/log/explain
POST /api/v1/ai/table/filter-suggestion
```

浏览器绝不能直接调用模型服务商。

## Streaming

当后端 streaming 能力存在时，使用专用 helper，例如 `fetchStream`，并支持：

- abort 支持
- partial message callback
- final message callback
- error callback
- auth 处理

在后端 streaming 尚未存在前，保持 AI client 可 mock。
