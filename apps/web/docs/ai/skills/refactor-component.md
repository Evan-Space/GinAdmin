# Skill: 重构组件

当需要在不改变产品行为的前提下改进现有组件时，使用本 skill。

## 目标

- 提升可读性。
- 减少重复。
- 明确类型。
- 将副作用与渲染分离。
- 除非任务明确要求 UI 改动，否则保持现有视觉行为稳定。

## 步骤

1. 阅读相关 route、component、hook 和 constants。
2. 找到最小但有价值的抽取点。
3. 尽可能保持 public props 和 route behavior 不变。
4. 将 data fetching 移入 hooks。
5. 将重复 UI 移入小组件。
6. 复用的 magic values 移入 typed constants。
7. 不触碰 generated files。

## 规则

- 不重设计无关 UI。
- 除非明确要求，不重命名 routes。
- 除非明确要求，不修改 API contracts。
- 不因为某些状态当前未使用就随意移除。
- 不为小型清理引入新库。

## 验收标准

- TypeScript 仍能通过。
- 现有行为保持不变。
- 组件对其他 AI agent 或开发者来说更容易阅读。
- diff 聚焦。
