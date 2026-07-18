# Skill: 新增前端 AI 功能

当前端需要新增 AI 辅助 UI 时，使用本 skill。

## 默认做法

先基于可 mock 的 service interface 构建前端。不要等后端 AI endpoints 完成后才开始打磨 UI 行为。

## 推荐文件

```text
src/features/ai/
  components/
  hooks/
  services/
  types.ts
  constants.ts
```

## 步骤

1. 定义 `AiMessage`、`AiRequest` 和 `AiResponse` 类型。
2. 在 `services/aiClient.ts` 中创建 service interface。
3. 创建 `services/mockAiClient.ts`。
4. 添加 `services/redaction.ts`，用于前端侧尽力脱敏。
5. 创建类似 `useAiChat` 的 hook。
6. 创建类似 `AiAssistantDrawer` 的 UI surface。
7. 添加 loading、streaming、error、retry 和 cancel 状态。
8. AI suggestions 与已应用 actions 必须分离。

## UX 要求

- assistant thinking 或 streaming 时要可见。
- 允许用户取消长响应。
- 被复制或应用的 suggestions 要有明确反馈。
- 不要隐藏错误。
- 有帮助时展示所使用的 page context。

## 安全要求

- 前端代码中没有 provider API key。
- 浏览器不直接调用 provider APIs。
- 发送 context 前脱敏明显敏感字段。
- 应用 AI 生成的 mutation 前必须确认。

## 初始功能示例

- 后台布局中的 assistant drawer。
- 解释选中的 request log row。
- 根据自然语言生成 table filter draft。
- 为非敏感字段建议 form values。

## 验收标准

- 可以基于 mock AI client 工作。
- 后续可以切换到 backend client，而无需重写 UI。
- loading、error 和 cancel 行为可见。
- 不会自动应用 suggestions。
