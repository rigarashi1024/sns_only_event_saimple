package repository

import (
	"context"
	"time"

	gofirestore "cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// Timeline は Firestore の timelines コレクションに保存する表示用エントリです。
type Timeline struct {
	ID                  string    `firestore:"id"`
	PostID              string    `firestore:"post_id"`
	PostUserID          string    `firestore:"post_user_id"`
	Content             string    `firestore:"content"`
	PostCreatedAt       time.Time `firestore:"post_created_at"`
	TimelineOwnerUserID string    `firestore:"timeline_owner_user_id"`
	CreatedAt           time.Time `firestore:"created_at"`
	UpdatedAt           time.Time `firestore:"updated_at"`
}

// TimelineRepository は timelines コレクションへのアクセスをまとめます。
type TimelineRepository struct {
	client *gofirestore.Client
}

// NewTimelineRepository は Firestore client を使って TimelineRepository を生成します。
func NewTimelineRepository(client *gofirestore.Client) *TimelineRepository {
	return &TimelineRepository{client: client}
}

// ListByOwnerUserID は timeline_owner_user_id に紐づく表示済みタイムラインを新しい投稿順で返します。
func (r *TimelineRepository) ListByOwnerUserID(ctx context.Context, ownerUserID string) ([]Timeline, error) {
	iter := r.client.Collection("timelines").
		Where("timeline_owner_user_id", "==", ownerUserID).
		OrderBy("post_created_at", gofirestore.Desc).
		Documents(ctx)
	defer iter.Stop()

	timelines := make([]Timeline, 0)
	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				return timelines, nil
			}
			return nil, err
		}

		var timeline Timeline
		if err := doc.DataTo(&timeline); err != nil {
			return nil, err
		}
		timeline.ID = doc.Ref.ID
		timelines = append(timelines, timeline)
	}
}
