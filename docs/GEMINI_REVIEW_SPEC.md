# Gemini PR Review Spec

## 目的

このドキュメントは、GitHub Actions 上で動作する `Gemini PR Review` の仕様を人間向けに整理したものです。

主な目的は以下です。

- PR 作成時に AI レビューを自動実行する
- 重大な問題だけを早期に検知する
- スタイル論争ではなく、バグやセキュリティ事故の検出に集中する
- 人間レビューの補助として利用する

## 役割

このリポジトリにおける AI の役割分担は以下です。

- 実装: `Codex`
- 設計: 人間 + `Codex`
- レビュー: `Gemini`

Gemini は最終判断者ではなく、PR レビューの補助役です。最終的なマージ判断は人間が行います。

## 実行タイミング

Gemini レビューは以下のタイミングで実行されます。

- Pull Request の `opened`
- Pull Request の `synchronize`
- Pull Request の `reopened`
- PR に対して `/gemini-review` コメントが投稿されたとき

コメント経由の再実行を用意しているため、レビュー内容に違和感がある場合や、最新差分で再確認したい場合に再利用できます。

## 実行フロー

Gemini レビューは概ね以下の流れで動作します。

1. GitHub Actions が起動する
2. GitHub API から PR 情報を取得する
3. GitHub API から PR の changed files 一覧を取得する
4. レビュー対象外のファイルを除外する
5. 各ファイルごとに patch を Gemini に渡してレビューする
6. 結果を集約して PR コメントとして投稿する
7. 既存の Gemini コメントは編集せず、新しいコメントとして追加する

## Gemini に渡す情報

Gemini には、ファイル単位の patch とあわせて以下の情報を渡します。

- PR タイトル
- PR 本文
- ベースブランチ名
- ヘッドブランチ名
- changed files 一覧
- 対象ファイルの patch
- レビュールールの要約
- 未反映レビュー指摘 TODO

これにより、diff だけを単発で渡すよりも前提を理解しやすくし、見当違いなレビューを減らすことを狙っています。

未反映レビュー指摘 TODO は `docs/GEMINI_REVIEW_TODO.md` で管理します。Gemini は TODO に記載済みの項目を、人間が既に把握しているものとして扱い、同じ内容を繰り返し指摘しない方針です。

## レビュー対象

Gemini は、変更ファイルのうち patch を持つファイルを対象にレビューします。

ただし、以下のようなファイルは原則スキップします。

- lock file
- build 成果物
- 圧縮済みファイル
- `vendor/`
- `.next/`
- `.nuxt/`
- `.output/`
- `dist/`
- `build/`
- `go.sum`

また、1回の実行でレビューするファイル数には上限があります。

## レビュー方針

Gemini は以下の種類の問題のみを報告対象とします。

- `bug`
- `security`
- `logic`

反対に、以下は報告対象外です。

- スタイルの問題
- 命名規則
- フォーマット
- 微細なパフォーマンス改善
- 任意の改善提案
- patch に含まれていない範囲への推測
- 情報不足のままの断定

## 出力形式

重大な問題がある場合は、定型フォーマットで報告します。

```text
Issue-1:
    type: bug | security | logic
    file: <path>
    lines: <line or range>
    problem: <summary>
    reason: <why it is a problem>
    suggestion: <fix idea>
```

重大な問題がない場合は、以下を返します。

```text
重大な問題はありません。LGTM。
```

## コメント運用

Gemini のレビュー結果は PR にコメントとして投稿されます。

- 既存の Gemini コメントは編集しない
- レビュー実行ごとに新しいコメントを追加する
- これにより、レビューの時系列と判断の流れを追いやすくする
- 古いコメントは履歴として扱い、最新の状態は新しいコメントと後続の人間コメントで判断する

## TODO 管理

Gemini の指摘のうち、現時点では重大ではないが将来見直したいものは `docs/GEMINI_REVIEW_TODO.md` に追加します。

TODO に追加する例:

- 本番運用前には対応したいセキュリティ強化
- 今回の PR ではスコープ外の設計論点
- 将来の機能追加時に必要になる検証やガード
- 繰り返し指摘されやすいが、現状の脅威モデルでは即時対応しない項目

TODO に追加しない例:

- 明確な誤検知
- 既に実装または仕様書で対応済みの内容
- スタイルや任意改善の提案
- patch から断定できない推測

TODO の詳細な運用ルールは `docs/GEMINI_REVIEW_TODO_SPEC.md` に従います。

## モデル運用

Gemini は複数モデルを順番に試します。

- まず `GEMINI_API_KEY1` を利用する
- 429 が発生した場合は別モデルにフォールバックする
- すべて失敗した場合のみ `GEMINI_API_KEY2` にフォールバックする
- `GEMINI_API_KEY2` 利用時は Slack 通知する
- 運用上は `GEMINI_API_KEY1` を無料枠キー、`GEMINI_API_KEY2` を有料枠キーとして分けている

## 使用する環境変数

Gemini レビューでは、以下の 2 種類の環境変数を利用します。

- GitHub Secrets などに事前登録しておく固定値
- GitHub Actions 実行時に workflow が注入する実行時値

### GitHub Secrets などに事前登録する値

- `GEMINI_API_KEY1`
  - 主系で利用する Gemini API キー
  - 運用上は無料枠のキーを想定している
  - 未設定の場合は Gemini レビューをスキップする

- `GEMINI_API_KEY2`
  - `GEMINI_API_KEY1` 側の候補モデルがすべて 429 などで利用不能になった場合のフォールバック先
  - 運用上は有料枠のキーを想定している
  - 未設定でも動作はするが、主系が枯渇した場合はレビュー失敗となる

- `SLACK_WEBHOOK_URL`
  - `GEMINI_API_KEY2` にフォールバックした際の通知先
  - 未設定でもレビュー自体は継続する

### GitHub Actions が実行時に注入する値

- `GITHUB_TOKEN`
  - GitHub API 呼び出しに利用するトークン
  - PR 情報取得、changed files 取得、コメント投稿に使用する
  - 通常は GitHub Actions の標準トークンを利用する

- `PR_NUMBER`
  - レビュー対象 PR の番号
  - `pull_request` または `issue_comment` イベントから workflow 側で注入する
  - 事前にリポジトリ設定へ登録する値ではない

- `REPO_FULL_NAME`
  - GitHub の `owner/repo` 形式のリポジトリ名
  - GitHub API のエンドポイント組み立てに使用する
  - 通常は `github.repository` から実行時に注入する

## 現状の workflow における注入元

現状の workflow では、主に以下のように環境変数を設定しています。

- `GEMINI_API_KEY1`: GitHub Secrets
- `GEMINI_API_KEY2`: GitHub Secrets
- `SLACK_WEBHOOK_URL`: GitHub Secrets
- `GITHUB_TOKEN`: GitHub Actions 標準トークン
- `PR_NUMBER`: `github.event.pull_request.number || github.event.issue.number`
- `REPO_FULL_NAME`: `github.repository`

## 未設定時の挙動

- `GEMINI_API_KEY1` が無い
  - レビューはスキップされる
- `GITHUB_TOKEN` / `PR_NUMBER` / `REPO_FULL_NAME` が無い
  - スクリプトはエラー終了する
- `GEMINI_API_KEY2` が無い
  - 主系の全モデルが利用不可になった時点でフォールバックできない
- `SLACK_WEBHOOK_URL` が無い
  - Slack 通知だけスキップする

## 既知の制約

Gemini レビューには以下の制約があります。

- patch に現れない前提変更は理解しきれない場合がある
- 大きい PR ではレビュー対象ファイル数の上限により一部を見ない可能性がある
- 仕様変更とバグ修正の区別は PR 本文の質に影響される
- AI なので誤検知や見逃しはあり得る

## 運用上のおすすめ

- PR 本文に「何を変えたか」「なぜ変えたか」を短く書く
- PR を必要以上に大きくしない
- AI レビュー結果をそのまま採用せず、人間が妥当性を確認する
- `/gemini-review` で再実行できることを前提に運用する

## 関連ファイル

- workflow: `.github/workflows/gemini-pr-review.yaml`
- review script: `scripts/gemini-review.js`
- rules summary: `docs/REVIEW_RULES_SUMMARY.md`
- review TODO: `docs/GEMINI_REVIEW_TODO.md`
- review TODO spec: `docs/GEMINI_REVIEW_TODO_SPEC.md`
