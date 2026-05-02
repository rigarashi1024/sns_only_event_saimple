package config

import (
	"context"
	"fmt"
	"os"
)

const tokenEncryptionKeyEnvName = "TOKEN_ENCRYPTION_KEY"

// GetDBEncryptionKey は DB 保存用トークン暗号化キーを実行環境に応じて取得します。
// local では環境変数から TOKEN_ENCRYPTION_KEY を読み込みます。
// dev/prd では将来的に Secret Manager から読み込む想定です。
func GetDBEncryptionKey(ctx context.Context, cfg Config) (string, error) {
	switch cfg.Env {
	case EnvLocal:
		return getDBEncryptionKeyFromEnv()
	case EnvDev, EnvPrd:
		return getDBEncryptionKeyFromSecretManager(ctx, cfg)
	default:
		return "", fmt.Errorf("unsupported APP_ENV %q: expected %q, %q, or %q", cfg.Env, EnvLocal, EnvDev, EnvPrd)
	}
}

func getDBEncryptionKeyFromEnv() (string, error) {
	key := os.Getenv(tokenEncryptionKeyEnvName)
	if key == "" {
		return "", fmt.Errorf("%s is required when APP_ENV=%s", tokenEncryptionKeyEnvName, EnvLocal)
	}

	return key, nil
}

func getDBEncryptionKeyFromSecretManager(ctx context.Context, cfg Config) (string, error) {
	// TODO: dev/prd では Secret Manager から TOKEN_ENCRYPTION_KEY を取得する。
	// その際は project id、secret id、version の指定と、実行サービスアカウントの
	// secretAccessor 権限をセットで整備する。
	_ = ctx
	_ = cfg
	return "", fmt.Errorf("Secret Manager key loading is not implemented yet")
}
