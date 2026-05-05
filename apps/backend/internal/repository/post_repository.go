package repository

import (
	"context"
	"time"

	gofirestore "cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// Post は Firestore の posts コレクションに保存する投稿情報です。
type Post struct {
	ID           string    `firestore:"id"`
	UserID       string    `firestore:"user_id"`
	Content      string    `firestore:"content"`
	LikeCount    int       `firestore:"like_count"`
	CommentCount int       `firestore:"comment_count"`
	CreatedAt    time.Time `firestore:"created_at"`
	UpdatedAt    time.Time `firestore:"updated_at"`
}

// PostRepository は posts コレクションへのアクセスをまとめます。
type PostRepository struct {
	client *gofirestore.Client
}

// NewPostRepository は Firestore client を使って PostRepository を生成します。
func NewPostRepository(client *gofirestore.Client) *PostRepository {
	return &PostRepository{client: client}
}

// ListByUserID は user_id に紐づく投稿を新しい順で返します。
func (r *PostRepository) ListByUserID(ctx context.Context, userID string) ([]Post, error) {
	iter := r.client.Collection("posts").
		Where("user_id", "==", userID).
		OrderBy("created_at", gofirestore.Desc).
		Documents(ctx)
	defer iter.Stop()

	posts := make([]Post, 0)
	for {
		doc, err := iter.Next()
		if err != nil {
			if err == iterator.Done {
				return posts, nil
			}
			return nil, err
		}

		var post Post
		if err := doc.DataTo(&post); err != nil {
			return nil, err
		}
		post.ID = doc.Ref.ID
		posts = append(posts, post)
	}
}
