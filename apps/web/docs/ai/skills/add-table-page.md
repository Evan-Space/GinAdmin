# Skill: 新增表格管理页

当新增以列表、筛选和行操作为核心的后台页面时，使用本 skill。

## 推荐文件

小型 route-owned table 推荐：

```text
src/pages/_layout/<page>/index.tsx
src/pages/_layout/<page>/-hooks.tsx
src/pages/_layout/<page>/-types.ts
src/pages/_layout/<page>/-constant.tsx
```

较大的可复用 feature 推荐：

```text
src/features/<feature>/
  components/
  hooks/
  services/
  types.ts
  constants.ts
```

## 必需 UI

- Filter form。
- Table。
- Loading state。
- Empty state。
- Error feedback。
- 后端支持时添加 pagination。
- 根据需要添加 view、edit、delete、enable、disable 或 details 等 row actions。

## Ant Design Components

优先使用：

- `Form`
- `Input`
- `Select`
- `DatePicker`
- `Button`
- `Space`
- `Table`
- `Tag`
- `Drawer`
- `Modal`
- `Popconfirm`

## 数据流

1. Form state 生成 query params。
2. Hook 调用 API module。
3. Hook 返回 data、loading、error 和 action handlers。
4. Page 渲染 filters 和 table。
5. 破坏性 row actions 需要 confirmation。

## AI 增强选项

适合表格页面的 AI 增强：

- natural language filter suggestion
- column explanation
- selected-row summary
- anomaly explanation

AI suggestions 应在应用前先预览。

## 验收标准

- Table columns 有类型。
- Loading 和 empty states 不会让整个页面意外跳动。
- 破坏性操作需要确认。
- AI suggestions 不会静默应用。
