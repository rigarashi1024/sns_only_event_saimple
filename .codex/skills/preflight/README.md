# Preflight Skill - Nuxt 向けガイド

## 概要

Preflight skill は、Nuxt プロジェクトでコミットや PR の前に行いたい品質チェックをまとめたものです。

## 実行する検査

1. Format
2. Lint
3. Type Check
4. Test

## 前提

- このリポジトリでは Nuxt フロントエンドを想定する
- script 名はプロジェクトごとに異なる可能性がある
- 実行前に `package.json` を確認する

## 想定する script 名

- `format`
- `format:check`
- `lint`
- `lint:fix`
- `typecheck` または `type-check`
- `test`

Nuxt では `nuxi typecheck` を直接使う構成もあるため、script が無い場合はそちらも候補にする。

## 使い方

```text
/preflight
```

## 期待する動き

- format と lint は自動修正可能なら先に直す
- type と test は失敗原因を整理して修正案を返す
- 途中で失敗しても、可能な範囲で最後まで確認する

## 出力イメージ

```text
Preflight Check Results

1. [OK] Format
2. [OK] Lint
3. [NG] Type
4. [SKIP] Test

必要な修正:
- components/PostForm.vue:24
  問題: ref の型推論が崩れている
  修正案: ref<string>('') のように明示する
```

## この skill でやらないこと

- PR 作成
- コミット
- 破壊的な一括変更
- script が存在しない検査の無理な実行

## 関連

- skill 本体: `.codex/skills/preflight/SKILL.md`
- プロジェクト概要: `README.md`
