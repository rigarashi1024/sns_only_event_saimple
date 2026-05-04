package firestore

import (
	"context"

	gofirestore "cloud.google.com/go/firestore"

	"github.com/rigarashi1024/sns_only_event_saimple/apps/backend/internal/config"
)

// NewClient は設定された project id で Firestore client を生成します。
// FIRESTORE_EMULATOR_HOST が設定されている場合は、SDK が自動的に Emulator へ接続します。
func NewClient(ctx context.Context, cfg config.Config) (*gofirestore.Client, error) {
	return gofirestore.NewClient(ctx, cfg.ProjectID)
}
