package config

import (
	"os"
	"strings"
)

const defaultProjectID = "sns-only-event-local"
const defaultEnv = "local"
const defaultTokenEncryptionKeySecretID = "TOKEN_ENCRYPTION_KEY"
const defaultTokenEncryptionKeySecretVersion = "latest"
const defaultCORSAllowedOrigins = "http://localhost:3001"

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
	CORSAllowedOrigins              []string
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

	corsAllowedOrigins := parseCSVEnv("CORS_ALLOWED_ORIGINS", defaultCORSAllowedOrigins)

	return Config{
		Env:                             env,
		ProjectID:                       projectID,
		TokenEncryptionKeySecretID:      tokenEncryptionKeySecretID,
		TokenEncryptionKeySecretVersion: tokenEncryptionKeySecretVersion,
		CORSAllowedOrigins:              corsAllowedOrigins,
	}
}

func parseCSVEnv(name string, fallback string) []string {
	value := os.Getenv(name)
	if value == "" {
		value = fallback
	}

	rawItems := strings.Split(value, ",")
	items := make([]string, 0, len(rawItems))
	for _, item := range rawItems {
		item = strings.TrimSpace(item)
		if item != "" {
			items = append(items, item)
		}
	}
	return items
}
