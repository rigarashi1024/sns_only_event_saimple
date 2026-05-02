package firestore

import (
	"context"

	gofirestore "cloud.google.com/go/firestore"

	"github.com/rigarashi1024/sns_only_event_saimple/apps/backend/internal/config"
)

func NewClient(ctx context.Context, cfg config.Config) (*gofirestore.Client, error) {
	return gofirestore.NewClient(ctx, cfg.ProjectID)
}
