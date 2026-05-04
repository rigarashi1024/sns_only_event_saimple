# Local Development Setup

このドキュメントは、Nuxt フロントエンド、Go バックエンド、Firestore Emulator、Pub/Sub Emulator を `docker compose` で動かすための初期手順です。

## 前提

- Docker Desktop または Docker Engine + Docker Compose Plugin が利用できる
- ホスト側で `npm` と `go` を使って初期化コマンドを実行できる
- このリポジトリのルートにいる

## 推奨ランタイム

- Node.js は `24 LTS` を推奨する
- `npm` は Node.js に同梱されるものを基本利用でよい
- Go は `1.25` 系を推奨する
- まずは以下で確認する

```bash
node -v
npm -v
go version
```

補足:

- Nuxt 初期化時は `npm` の細かい版より Node.js の互換性のほうが重要
- 新規構築では `Current` より `LTS` を優先する
- バックエンド開発コンテナでは `air@latest` を利用するため、Go も `1.25` 系を前提にする
- frontend 開発コンテナも `Node 24` 系を前提にする

## 全体方針

- Nuxt プロジェクト自体はホスト側で初期化する
- Go バックエンドのモジュール作成もホスト側で行う
- 開発時の実行は `docker compose` でまとめる
- `apps/frontend` と `apps/backend` は bind mount する
- Firestore / Pub/Sub はローカルエミュレータで起動する
- Firestore / Pub/Sub のコンテナは emulator 同梱の `google-cloud-cli:emulators` イメージを利用する

## 1. Nuxt フロントエンドの初期化

`apps/frontend` で Nuxt を作成します。

```bash
cd apps/frontend
npm create nuxt@latest .
```

補足:

- 既存ファイルの上書き確認が出たら、README などを退避してから実行してください
- 初期化後は `.env.example` を参考に `.env` を作成してください

## 2. Go バックエンドの初期化

`apps/backend` で Go モジュールを初期化します。

例:

```bash
cd apps/backend
go mod init github.com/rigarashi1024/sns_only_event_saimple/apps/backend
```

最初の疎通確認用として、最低限 `cmd/api/main.go` を作成してください。

例:

```go
package main

import (
  "log"
  "net/http"
)

func main() {
  http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
    _, _ = w.Write([]byte("ok"))
  })

  log.Println("listening on :8081")
  log.Fatal(http.ListenAndServe(":8081", nil))
}
```

## 3. 環境変数の準備

必要に応じて以下をコピーして利用してください。

- `apps/frontend/.env.example`
- `apps/backend/.env.example`

例:

```bash
cp apps/frontend/.env.example apps/frontend/.env
cp apps/backend/.env.example apps/backend/.env
```

## 4. docker compose 起動

ルートディレクトリで以下を実行します。

```bash
docker compose -f infra/local/docker-compose.yml up --build
```

## 5. 想定ポート

- Nuxt: `3001`
- Go API: `8081`
- Firestore Emulator: `8080`
- Pub/Sub Emulator: `8085`

## 6. 起動確認

### フロント

ブラウザで以下にアクセス:

```text
http://localhost:3001
```

### バックエンド

ヘルスチェック例:

```bash
curl http://localhost:8081/healthz
```

### Firestore Emulator

ホスト側から以下のような環境変数で接続します。

```bash
export FIRESTORE_EMULATOR_HOST=localhost:8080
```

### Pub/Sub Emulator

ホスト側から以下のような環境変数で接続します。

```bash
export PUBSUB_EMULATOR_HOST=localhost:8085
export PUBSUB_PROJECT_ID=sns-only-event-local
```

## 7. Seed データ投入

`users` と `sessions` のダミーデータを Firestore Emulator に投入できます。

```bash
cd apps/backend
FIRESTORE_EMULATOR_HOST=localhost:8080 PUBSUB_PROJECT_ID=sns-only-event-local go run ./cmd/seed
```

投入される主なデータ:

- `users/test-user`
- `users/user-002`
- `sessions/session-test-user-001`

## 8. エミュレータ内データの確認

### Firestore

現在の `docker compose` 構成では、`gcloud` ベースの Firestore Emulator を直接起動しているため、Firestore データを確認する専用 UI は同梱していません。

確認方法の候補:

- Go / Node の確認スクリプトを作る
- Firestore Emulator に接続する簡易 API を作る
- 将来的に Firebase Local Emulator Suite へ切り替えて UI を使う

補足:

- Firebase Local Emulator Suite には Firestore をブラウザで確認できる UI があります
- ただし今の構成は Firebase Emulator Suite ではなく、`gcloud` の emulator を直接使う構成です

### Pub/Sub

Pub/Sub Emulator については、Google Cloud の公式ドキュメントでもコンソール UI や `gcloud pubsub` コマンドはサポート対象外です。確認はアプリケーションコードや補助スクリプト経由で行う前提になります。

## 9. よくある注意点

- Nuxt 初期化前に `docker compose up` すると、`package.json` が無いため frontend コンテナは待機メッセージを出して停止せず待機します
- Go 初期化前に `docker compose up` すると、`go.mod` が無いため backend コンテナは待機メッセージを出して停止せず待機します
- Nuxt の `node_modules` は named volume で持つため、ホストの `node_modules` と混ざりません
- バックエンドは `air` でホットリロードする前提です

## 10. 今後追加したいもの

- Firestore 初期データ投入スクリプト
- Pub/Sub の topic / subscription 初期化スクリプト
- Nuxt と Go の初期テンプレート
- `make` または `task` ベースの補助コマンド
