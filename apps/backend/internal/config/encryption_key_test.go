package config

import (
	"context"
	"testing"
)

func TestGetDBEncryptionKeyLocal(t *testing.T) {
	t.Setenv("TOKEN_ENCRYPTION_KEY", "local-dev-key")

	key, err := GetDBEncryptionKey(context.Background(), Config{Env: EnvLocal})
	if err != nil {
		t.Fatalf("GetDBEncryptionKey returned error: %v", err)
	}
	if key != "local-dev-key" {
		t.Fatalf("key = %q, want %q", key, "local-dev-key")
	}
}

func TestGetDBEncryptionKeyLocalRequiresEnv(t *testing.T) {
	t.Setenv("TOKEN_ENCRYPTION_KEY", "")

	_, err := GetDBEncryptionKey(context.Background(), Config{Env: EnvLocal})
	if err == nil {
		t.Fatal("GetDBEncryptionKey returned nil error")
	}
}

func TestGetDBEncryptionKeyRejectsUnsupportedEnv(t *testing.T) {
	t.Setenv("TOKEN_ENCRYPTION_KEY", "local-dev-key")

	_, err := GetDBEncryptionKey(context.Background(), Config{Env: "stg"})
	if err == nil {
		t.Fatal("GetDBEncryptionKey returned nil error")
	}
}
