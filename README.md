# sns_only_event_saimple

イベント駆動アーキテクチャ、コレオグラフィー、Queue Service の練習用として作成する簡易 SNS プロジェクトです。

最小構成で開始しつつ、Cloud Run、Pub/Sub、Firestore を使って、イベントを主体的に発行しながらサービス間連携を行う構成を目指します。

## 目的

- イベント駆動設計の基本を、SNS の代表的なユースケースで試す
- コレオグラフィーパターンに寄せた非同期連携を体験する
- Pub/Sub を中心に、各処理が次のイベントを自発的に発行する流れを確認する
- Firestore を使って、低コスト寄りの NoSQL ベース構成を試す
- ローカルエミュレータを使い、開発環境でも疎結合な流れを再現する

## 想定サービス規模

現時点では、以下のイベントを中心に扱う想定です。

- `UserRegistered`
- `UserFollowed`
- `PostCreated`
- `PostLiked`
- `CommentCreated`
- `PostDeleted`
- `ProfileUpdated`

この段階では「SNS の主要導線を一通り触れる最小イベントセット」として扱います。

## 使用アーキテクチャ

利用する主要コンポーネントは以下です。

- フロントエンド: Cloud Run
- バックエンド: Cloud Run
- メッセージング: Pub/Sub
- データストア: Firestore

### バックエンド構成方針

- バックエンドは 1 つの `Cloud Run` サービスとして構築する
- ただしコード上の責務は機能別に分割する
- 今回は複数の `Cloud Run` サービスへは分割しない

### 構成イメージ

- ブラウザ -> フロント
- ブラウザ -> バック
- バック -> Pub/Sub
- Pub/Sub -> バック
- バック -> Pub/Sub -> ...

### 想定フロー

- 画面リソース取得: ブラウザ -> フロント
- データ取得: ブラウザ -> バック
- サービス間処理: ブラウザ -> バック -> Pub/Sub -> バック -> Pub/Sub -> ...

## 設計原則

- 各サービス内では、処理の完了に応じて次のイベント発行まで責務として持つ
- `cron` などの監視起点ではなく、処理主体が自ら `Pub/Sub` に push する
- 段階的な学習順として、まずは素朴な `Pub/Sub` ベースで始める
- 次の発展として `Outbox` パターンを検討する
- 最後に監視パターンや補助的な配信方式を比較対象として扱う
- コスト重視のため、RDB ではなく `Firestore` を採用する
- `Cloud Run` が URL を持つことを前提に、まずは最小構成で組む
- 開発環境ではローカルエミュレータを使用する

## 技術スタック

- フロントエンド: `Nuxt`
- バックエンド: `Go`
- API 仕様: `OpenAPI`
- データストア: `Firestore`
- メッセージング: `Pub/Sub`
- 実行基盤: `Cloud Run`

## Pub/Sub 方針

- `Pub/Sub` のトピックはイベントごとに個別作成する
- まずはイベントの流れを追いやすい構成を優先する
- 想定トピック例
  - `user-registered`
  - `user-followed`
  - `post-created`
  - `post-liked`
  - `comment-created`
  - `post-deleted`
  - `profile-updated`

## 目指すシステム像

このプロジェクトでは、画面表示と同期 API を最小限に保ちつつ、状態変化はイベントとして明示的に流す構成を目指します。

たとえば投稿作成であれば、API が `PostCreated` を発行し、その後の関連処理はイベントを購読したバックエンドが担当し、必要に応じてさらに別イベントを発行する、という流れを基本とします。

## API と開発ルール

- API 仕様書は `OpenAPI` の仕様に従って管理する
- API 管理方針は `schema-first` で進める
- 静的解析は必ず実施する
- 静的解析は `GitHub Actions` 側でも実行する

## AI の役割分担

- 実装: `Codex`
- 設計: 人間 + `Codex`
- レビュー: `Gemini`

## Firestore の初期コレクション案

初期案として、以下のコレクションをベースに進める。

- `users`
  - `id`
  - `name`
  - `email`
  - `nickname`
  - `created_at`
  - `updated_at`
- `sessions`
  - `id`: セッション ID。Cookie に保存する独自 `JWT` の `sid` と対応する
  - `user_id`: このセッションが属するユーザー ID
  - `provider_access_token_encrypted`: provider または local provider 相当の access token を AES-256-GCM で暗号化した値
  - `provider_access_token_expires_at`: provider access token の有効期限
  - `provider_refresh_token_encrypted`: provider または local provider 相当の refresh token を AES-256-GCM で暗号化した値
  - `provider_refresh_token_expires_at`: provider refresh token の有効期限
  - `internal_jwt_jti`: 発行した独自 `JWT` の `jti`
  - `internal_access_token_expires_at`: Cookie に保存する独自 `JWT` の有効期限
  - `provider_type`: `local`, `google` などの認証 provider 種別
  - `created_at`: セッション作成日時
  - `updated_at`: セッション更新日時
- `posts`
  - `id`
  - `user_id`
  - `content`
  - `like_count`
  - `comment_count`
  - `created_at`
  - `updated_at`
- `follows`
  - `id`
  - `user_id_from`
  - `user_id_to`
  - `created_at`
  - `updated_at`
- `comments`
  - `id`
  - `user_id`
  - `post_id`
  - `content`
  - `created_at`
  - `updated_at`
- `timelines`
  - `id`
  - `post_id`
  - `post_user_id`
  - `content`
  - `post_created_at`
  - `timeline_owner_user_id`
  - `created_at`
  - `updated_at`

### likes の扱い

- `posts.like_count` は表示高速化のために持つ
- ただし「誰がいいねしたか」の事実自体は別管理を前提とする
- そのため、`likes` コレクションを独立して持つ
- `likes` の初期案
  - `id`
  - `user_id`
  - `post_id`
  - `created_at`

### Firestore 設計メモ

- Firestore では `pk` と `id` を分けず、`document id` を `id` として扱う前提で進める
- タイムラインは非同期で事前構築する
- `timelines` は表示用の読み取りモデルとして、投稿内容の最小スナップショットを保持する
- `POST /posts` は保存完了時点で成功とし、タイムライン反映は非同期処理とする
- 削除は物理削除を前提とする
- イベントの再送や重複受信を考慮し、主要処理は冪等に実装する
- 冪等性は Firestore のドキュメント ID を複合主キー相当で設計して担保する
- 実行済みイベント保存用のコレクションを持ち、処理済みイベントを記録する
- `processed_events` には `event_id`, `event_type`, `topic_name`, `processed_at` を保存する
- `event_id` は重要イベントのみ付与して管理する

## 認証・トークン方針

- BFF を前提に、セッションと `JWT` を組み合わせた構成を採用する
- 初期実装では `id/password` ログイン画面を用意する
- 認証成功後、バックエンドは以下を払い出し、または保存する
  - `access_token`
  - `refresh_token`
  - フロント返却用の独自 `JWT`
- `access_token` の有効期限は 15 分とする
- `refresh_token` の有効期限は 1 日とする
- 独自 `JWT` の有効期限は `access_token` と同程度とする
- 独自 `JWT` は `Cookie` に保存する
- 独自 `JWT` には `session_id` と `jti` を含め、バックエンドが `sessions` の保存データと照合できるようにする
- `access_token` と `refresh_token` はフロントへ直接渡さず、バックエンド側で暗号化して `sessions` コレクションに保存する
- 保存先コレクションは `sessions` を使用し、独自 `JWT` の `session_id` と保存済み `access_token` / `refresh_token` を紐づけてセッション管理する
- リクエスト時は、フロントから送られた `Cookie` 内の独自 `JWT` をバックエンドで検証する
- 独自 `JWT` が期限切れ、または対応する `access_token` が期限切れの場合は、保存済み `refresh_token` を用いて `access_token` を更新する
- `refresh_token` も期限切れの場合は再ログインとする
- 将来的には `Google` などの `SSO` に切り替えやすい構成を前提とする

## トークン保存と暗号化方針

- `access_token` と `refresh_token` は DB 保存時に `AES-256-GCM` で暗号化する
- 暗号化鍵は `Secret Manager` に保存する
- `Cloud Run` 起動時に暗号化鍵を読み込み、メモリ上で利用する
- 費用を抑えるため、今回は `Cloud KMS` などの鍵管理機能ではなく `Secret Manager` を利用する
- 理想的には `KMS` などの専用鍵管理機能を用いるべきだが、今回は学習用の最小コスト構成を優先する

## ローカル開発方針

- ローカルでは各種エミュレータを利用して開発する
- 本番相当では `Cloud Run`、`Pub/Sub`、`Firestore` を利用する
- まずは最小構成で作り、必要に応じてサービス分割や補強を進める

## ディレクトリ構成

```text
apps/
  frontend/   Nuxt フロントエンド
  backend/    Go バックエンド
openapi/      OpenAPI 定義
packages/     共有定義や補助コード
infra/
  local/      ローカル開発用設定
  gcp/        GCP 向け設定
docs/         運用・設計ドキュメント
scripts/      補助スクリプト
```

## 今後の実装候補

- ユーザー登録
- フォロー
- 投稿作成
- いいね
- コメント
- 投稿削除
- プロフィール更新
- イベント購読処理
- Firestore への永続化
- OpenAPI ベースの API 定義
- GitHub Actions による lint / static analysis

## 現時点の前提

以下は現時点での仮置き前提です。

- まずは最小構成を優先する
- コレオグラフィーの学習が主目的である
- インフラの豪華さより、イベントの流れが分かることを優先する
- 将来的に `Outbox` や監視起点の方式と比較できる土台にする
- モノレポでフロントとバックを同居させる
- `likes` は独立コレクションとして保持し、`posts.like_count` を集約値として併用する
- `timelines` は `post_id`, `post_user_id`, `content`, `post_created_at`, `timeline_owner_user_id`, `created_at` を持つ
- 認証は初期実装として `id/password` ログイン + 独自 `JWT` 発行方式を採用する
- 冪等性は複合主キー相当のドキュメント ID 設計、処理済みイベント保存、重要イベントの `event_id` 管理で担保する
- OpenAPI からは Go の型定義とインタフェースのみ生成し、業務ロジックは手実装する
- OpenAPI 定義を更新したタイミングで生成コードも更新し、同じ差分として管理する
- OpenAPI 定義は `openapi/`、生成コードは `backend/gen/`、手実装コードは `backend/internal/` に配置する
- トークンは DB 保存時に `AES-256-GCM` で暗号化し、鍵は `Secret Manager` から起動時に読み込む
