# Repo Skills

このディレクトリは、このリポジトリ向けに管理する skill の置き場所です。

## 例

- `.codex/skills/my-skill/SKILL.md`

## 標準フロー

このリポジトリでは、AI 支援の PR 作業は原則として以下の順で進める。

1. `preflight`: PR 前の品質チェック
2. `pr-creator`: コミット、push、PR 作成
3. `fix-loop`: Gemini レビュー指摘の分類、必要な修正、TODO 管理
4. `sync-main`: マージ後に local `main` を最新化

## 基本ルール

- 1 skill につき 1 ディレクトリを作る
- skill 本体の説明は `SKILL.md` に書く
- 追加の資料が必要なら skill 配下に `references/` や `assets/` を置く
- 複数 skill で使う実行スクリプトは repository root の `scripts/codex-*.sh` に置く
- `SKILL.md` は長い手順書にせず、可能な限り `scripts/codex-*.sh` の呼び出しに寄せる
- 移植元ツール向けの README や slash-command 前提の説明は置かず、Codex が読む `SKILL.md` に集約する

## 推奨構成

```text
.codex/skills/
  my-skill/
    SKILL.md
    references/
    scripts/
    assets/
```
