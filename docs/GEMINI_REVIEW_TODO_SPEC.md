# Gemini Review TODO Spec

## 目的

`docs/GEMINI_REVIEW_TODO.md` は、Gemini PR Review の指摘のうち、今すぐ修正しない項目を明示的に管理するためのファイルです。

この運用の目的は以下です。

- 重大でない指摘を毎回修正対象にしない
- 「認識済みだが後で対応する」項目を PR の流れから切り離して管理する
- Gemini が同じ指摘を繰り返さないよう、レビュー時の前提情報として渡す
- 人間が「今修正するもの」と「後で見るもの」を判断しやすくする

## 基本方針

Gemini の指摘は以下の 3 種類に分類します。

- 重大で今修正するもの: PR 内で修正する
- 重要だが今は修正しないもの: `docs/GEMINI_REVIEW_TODO.md` に追加する
- 誤検知または現時点で不要なもの: PR コメントで理由を説明し、必要に応じて TODO には追加しない

「セキュリティリスク」という表現が含まれていても、現在の実装スコープ、利用環境、脅威モデルに照らして重大でない場合は、即時修正ではなく TODO 管理またはコメント回答を選びます。

## TODO に追加する基準

以下に該当する場合は TODO に追加します。

- 将来の機能追加時には対応が必要になる
- 本番運用前には見直したい
- 現在の PR ではスコープ外だが、設計上の論点として残したい
- 同じ指摘が繰り返される可能性がある

以下は原則 TODO に追加しません。

- 明確な誤検知
- 既にコードまたは仕様書で対応済みの内容
- 好み、スタイル、任意のリファクタ
- patch から断定できない推測だけの指摘

## TODO の書き方

`docs/GEMINI_REVIEW_TODO.md` の `運用中の TODO` に、以下の形式で追記します。

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

## status の意味

- `open`: 未対応。将来見直す
- `resolved`: 対応済み。履歴として残す
- `wontfix`: 対応しない判断をした。理由を `notes` に残す

## Gemini への扱い

Gemini レビュー実行時には、`docs/GEMINI_REVIEW_TODO.md` の内容を prompt に含めます。

Gemini には以下の前提でレビューさせます。

- TODO に記載済みの内容は、人間が認識済みのため再指摘しない
- TODO に記載済みでも、今回の patch がリスクを明確に悪化させた場合のみ報告する
- TODO にない新しい重大問題だけを報告する

## PR コメント運用

Gemini レビューコメントは、既存コメントを編集せず、毎回新しいコメントとして追加します。

理由:

- レビューの流れを追えるようにするため
- どの push に対してどの指摘が出たかを残すため
- 人間の回答や TODO 追加判断との対応関係を見やすくするため

## 関連ファイル

- review script: `scripts/gemini-review.js`
- prompt rules: `docs/REVIEW_RULES_SUMMARY.md`
- TODO file: `docs/GEMINI_REVIEW_TODO.md`
- overall spec: `docs/GEMINI_REVIEW_SPEC.md`
