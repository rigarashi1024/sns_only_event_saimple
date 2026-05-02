package config

import "os"

const defaultProjectID = "sns-only-event-local"
const defaultEnv = "local"
const defaultTokenEncryptionKeySecretID = "TOKEN_ENCRYPTION_KEY"
const defaultTokenEncryptionKeySecretVersion = "latest"

const (
	// EnvLocal はローカル PC や Docker Compose での開発環境です。
	EnvLocal = "local"
	// EnvDev はクラウド上の開発環境です。
	EnvDev = "dev"
	// EnvPrd は本番環境向けの実行環境です。
	EnvPrd = "prd"
)

// Config はバックエンド起動時に参照する環境設定です。
type Config struct {
	Env                             string
	ProjectID                       string
	TokenEncryptionKeySecretID      string
	TokenEncryptionKeySecretVersion string
}

// Load は環境変数から設定を読み込みます。
// ローカル開発では PUBSUB_PROJECT_ID と Firestore の project id を共用します。
func Load() Config {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = defaultEnv
	}

	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		projectID = os.Getenv("PUBSUB_PROJECT_ID")
	}
	if projectID == "" {
		projectID = defaultProjectID
	}

	tokenEncryptionKeySecretID := os.Getenv("TOKEN_ENCRYPTION_KEY_SECRET_ID")
	if tokenEncryptionKeySecretID == "" {
		tokenEncryptionKeySecretID = defaultTokenEncryptionKeySecretID
	}

	tokenEncryptionKeySecretVersion := os.Getenv("TOKEN_ENCRYPTION_KEY_SECRET_VERSION")
	if tokenEncryptionKeySecretVersion == "" {
		tokenEncryptionKeySecretVersion = defaultTokenEncryptionKeySecretVersion
	}

	return Config{
		Env:                             env,
		ProjectID:                       projectID,
		TokenEncryptionKeySecretID:      tokenEncryptionKeySecretID,
		TokenEncryptionKeySecretVersion: tokenEncryptionKeySecretVersion,
	}
}
