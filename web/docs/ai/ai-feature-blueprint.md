# AI 功能蓝图

本文档说明如何新增 AI 辅助前端功能，同时不让浏览器负责模型访问。

## 前端职责

前端可以：

- 收集用户意图
- 收集安全的页面 context
- 对明显敏感值做脱敏
- 调用后端 AI endpoints 或 mock clients
- 渲染 streaming text
- 渲染 suggestions
- 要求用户确认建议操作

前端不能：

- 保存模型服务商 API Key
- 直接调用模型服务商
- 绕过当前 auth state
- 自动执行高影响操作
- 在可以避免时发送密钥或原始敏感数据

## 建议 Feature 目录

```text
src/features/ai/
  components/
    AiAssistantDrawer.tsx
    AiChatPanel.tsx
    AiMessageList.tsx
    AiPromptBox.tsx
  hooks/
    useAiChat.ts
    useAiPageContext.ts
  services/
    aiClient.ts
    mockAiClient.ts
    redaction.ts
  types.ts
  constants.ts
```

## 核心类型

推荐的领域概念：

```ts
type AiMessageRole = 'user' | 'assistant' | 'system'

type AiMessage = {
  id: string
  role: AiMessageRole
  content: string
  createdAt: string
  status?: 'streaming' | 'done' | 'error'
}

type AiPageContext = {
  route: string
  title?: string
  selectedText?: string
  filters?: Record<string, unknown>
  tableColumns?: string[]
}
```

## 优先实现的 AI 功能

### 后台 Assistant Drawer

在后台布局中添加右侧 drawer。它可以回答当前页面相关问题，并建议下一步操作。

### 自然语言筛选

让用户输入自然语言请求，并在应用前预览生成的表格筛选条件。

### 日志解释

在请求日志和错误日志页面，让用户对选中的日志行请求通俗解释。

### 权限草案

根据自然语言描述生成权限建议，但应用前必须要求用户明确确认。

## Mock First 策略

先使用 `mockAiClient.ts` 构建 UI 行为：

- streaming simulation
- error simulation
- static suggestions
- cancellation handling

当后端 endpoints 可用后，只替换 `aiClient.ts`。

## 验收清单

- AI UI 不需要前端 env 中存在模型服务商 key。
- Loading、streaming、error 和 cancel 状态可见。
- 用户可以在应用建议前先检查建议。
- 明显敏感字段在放入 context 前已脱敏。
- 功能可以基于 mock service 运行。
