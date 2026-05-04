package config

import (
	"context"
	"fmt"
	"os"
	"sync"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

const tokenEncryptionKeyEnvName = "TOKEN_ENCRYPTION_KEY"

var (
	secretManagerClientOnce sync.Once
	secretManagerClient     *secretmanager.Client
	secretManagerClientErr  error
)

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
	client, err := getSecretManagerClient(ctx)
	if err != nil {
		return "", err
	}

	name := fmt.Sprintf(
		"projects/%s/secrets/%s/versions/%s",
		cfg.ProjectID,
		cfg.TokenEncryptionKeySecretID,
		cfg.TokenEncryptionKeySecretVersion,
	)
	result, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	})
	if err != nil {
		return "", fmt.Errorf("failed to access token encryption key secret: %w", err)
	}

	key := string(result.Payload.Data)
	if key == "" {
		return "", fmt.Errorf("token encryption key secret %q is empty", name)
	}
	return key, nil
}

func getSecretManagerClient(ctx context.Context) (*secretmanager.Client, error) {
	secretManagerClientOnce.Do(func() {
		secretManagerClient, secretManagerClientErr = secretmanager.NewClient(ctx)
	})
	if secretManagerClientErr != nil {
		return nil, fmt.Errorf("failed to create Secret Manager client: %w", secretManagerClientErr)
	}
	return secretManagerClient, nil
}
