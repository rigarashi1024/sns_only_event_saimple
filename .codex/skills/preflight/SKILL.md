---
name: preflight
description: Nuxt プロジェクト向けの事前品質チェック。format/lint/type/test を順に確認し、可能なものは自動修正、難しいものは修正提案まで行う。
disable-model-invocation: false
project: true
---

# Preflight - Nuxt 向けコード品質事前検査

この skill は、Nuxt フロントエンドを含むリポジトリで、コミットや PR 作成前に最低限の品質確認を行うためのものです。

## 目的

- フォーマット崩れを事前に直す
- Lint エラーを事前に減らす
- TypeScript / Vue / Nuxt まわりの型エラーを検出する
- テストがあれば実行し、壊れていないことを確認する

## 基本方針

- 検査は `format -> lint -> type -> test` の順で行う
- 自動修正できるものは先に直す
- 型エラーやテスト失敗は、原因分析と修正提案を返す
- すべての検査を順番に実行し、途中失敗しても可能な範囲で最後まで確認する
- PR 作成やコミットはこの skill の責務に含めない

## 実行前の確認

まず、リポジトリに定義されている script を確認する。

```bash
cat package.json
```

以下のような script が存在するか確認する。

- `format`
- `format:check`
- `lint`
- `lint:fix`
- `typecheck`
- `type-check`
- `test`

## 実行順序

### STEP 1: Format Check

優先順:

1. `npm run format:check`
2. `pnpm format:check`
3. `yarn format:check`

失敗時:

- `format` script があれば自動修正を試す
- 例:

```bash
npm run format
```

自動修正後は「何ファイル程度変わったか」を報告する。

### STEP 2: Lint Check

優先順:

1. `npm run lint`
2. `pnpm lint`
3. `yarn lint`

失敗時:

- `lint:fix` があれば自動修正を試す
- 例:

```bash
npm run lint:fix
```

自動修正後もエラーが残る場合:

1. 対象ファイルと行番号を整理する
2. Vue / Nuxt の文脈で原因を説明する
3. 必要なら修正案を提示する

### STEP 3: Type Check

Nuxt ではプロジェクトによって script 名が揺れるため、以下の順で試す。

1. `npm run typecheck`
2. `npm run type-check`
3. `npx nuxi typecheck`

失敗時:

1. エラーメッセージからファイルと行番号を特定する
2. `Vue`, `TypeScript`, `Nuxt composables`, `runtimeConfig`, `defineProps`, `ref`, `computed` などの観点で原因を整理する
3. 修正案を提示する

### STEP 4: Test

`test` script が存在する場合のみ実行する。

優先順:

1. `npm test`
2. `pnpm test`
3. `yarn test`

失敗時:

1. 失敗したテストケースを特定する
2. テストコードと対象コードのどちらに問題がありそうか切り分ける
3. 修正案を提示する

## Nuxt で特に見るポイント

- `pages/`, `components/`, `composables/`, `plugins/`, `server/` 配下の責務が壊れていないか
- `useFetch`, `useAsyncData`, `navigateTo`, `useRuntimeConfig` の使い方に不整合がないか
- `ref` / `computed` / `reactive` の型が崩れていないか
- `defineProps`, `defineEmits`, `defineModel` の型が安全か
- SSR 前提の処理で `window`, `document`, `localStorage` を不適切に触っていないか
- `server/api` とフロント側の呼び出し責務が混ざっていないか

## 最終レポート形式

結果は以下のようにまとめる。

```text
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Preflight Check Results
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

1. [OK/NG/SKIP] Format
2. [OK/NG/SKIP] Lint
3. [OK/NG/SKIP] Type
4. [OK/NG/SKIP] Test

総合結果:
- OK: コミット前チェックを通過
- NG: 修正が必要

必要な修正:
- file:line
  問題:
  修正案:
```

## 注意事項

1. この skill は PR 作成を行わない
2. script 名はリポジトリごとに異なるため、決め打ちせず `package.json` を確認する
3. 自動修正は format / lint に限定する
4. type / test は安易に無理やり直さず、原因を説明してから修正する
