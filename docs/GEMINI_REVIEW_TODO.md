# Gemini Review TODO

このファイルは、Gemini PR Review が指摘したもののうち、今すぐ修正しない項目を管理するための TODO です。

Gemini にはこのファイルの内容をレビュー時に渡します。ここに記載済みの項目は、人間が既に把握して管理しているものとして扱い、同じ内容を繰り返し指摘しない方針です。

## 運用中の TODO

### TODO-20260505-001: タイムライン取得 API にページングを追加する

- status: open
- source: PR #12 / Gemini comment https://github.com/rigarashi1024/sns_only_event_saimple/pull/12#issuecomment-4376426968
- category: logic
- priority: medium
- target: `openapi/api.yaml`, `apps/backend/internal/repository/timeline_repository.go`, `apps/frontend/app/pages/timeline.vue`
- reason_to_defer: 現在は seed データ確認用の最小実装を優先しており、一括取得で十分に検証できるため
- revisit_timing: タイムライン件数が増える前、または投稿機能と無限スクロールを実装するタイミング
- notes: `limit` と `cursor` を API に追加し、`next_cursor` や `has_more` を返す形を想定する

## 記載テンプレート

```md
### TODO-YYYYMMDD-001: タイトル

- status: open | resolved | wontfix
- source: PR #<number> / Gemini comment URL
- category: bug | security | logic | docs | operational
- priority: high | medium | low
- target: 対象ファイルや領域
- reason_to_defer: 今すぐ修正しない理由
- revisit_timing: 見直すタイミング
- notes: 補足
```
