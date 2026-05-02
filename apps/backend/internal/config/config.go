package config

import "os"

const defaultProjectID = "sns-only-event-local"

// Config はバックエンド起動時に参照する環境設定です。
type Config struct {
	ProjectID string
}

// Load は環境変数から設定を読み込みます。
// ローカル開発では PUBSUB_PROJECT_ID と Firestore の project id を共用します。
func Load() Config {
	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		projectID = os.Getenv("PUBSUB_PROJECT_ID")
	}
	if projectID == "" {
		projectID = defaultProjectID
	}

	return Config{
		ProjectID: projectID,
	}
}
