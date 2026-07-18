# AI 安全与隐私

## 浏览器边界

浏览器不是存放模型服务商密钥的安全位置。

不要在前端代码或前端环境变量中存放：

- OpenAI API key
- Anthropic API key
- Gemini API key
- 暗含私有路由的 provider base URL
- 后端服务凭证
- 特权 token

前端 AI 功能应调用 GinAdmin 后端接口。后端负责模型服务商凭证、权限检查、日志、限流和成本控制。

## 数据最小化

发送给 AI 功能的 context 应保持最小但有用。

优先发送：

- route name，而不是完整页面 dump
- selected rows，而不是整张表
- field labels，而不是隐藏内部值
- summarized errors，而不是完整 request payload

## Redaction

将前端 context 发送给 AI endpoints 前，脱敏明显敏感字段：

- `password`
- `token`
- `authorization`
- `cookie`
- `secret`
- `access_token`
- `refresh_token`
- phone numbers
- emails

前端脱敏是有帮助的一层，但不是最终安全边界。后端必须继续校验和脱敏。

## 用户确认

对于以下操作，AI 输出必须被视为建议：

- 删除记录
- 修改角色
- 修改权限
- 导出数据
- 创建任务
- 修改系统设置

应用这些操作前，必须要求用户明确确认。

## Prompt Injection

用户内容、日志、表格单元格和错误信息都应视为不可信数据。

不要让应用内容变成 system instructions。页面数据应放在 context 字段中，实际 instruction hierarchy 由后端 prompt 定义。

## 可审计性

AI 辅助工作流应保留：

- 用户意图
- 发送的 context
- 生成的 suggestion
- 最终用户 action
- 后端返回的 request id，如果存在

这对权限、导出和系统设置尤其重要。
