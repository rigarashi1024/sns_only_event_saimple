package repository

import (
	"context"
	"errors"
	"strings"

	gofirestore "cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrUserNotFound はログイン ID または email に該当するユーザーが存在しない場合に返します。
var ErrUserNotFound = errors.New("user not found")

// User は Firestore の users コレクションに保存するユーザー情報です。
type User struct {
	ID           string `firestore:"id"`
	Name         string `firestore:"name"`
	Email        string `firestore:"email"`
	Nickname     string `firestore:"nickname"`
	PasswordHash string `firestore:"password_hash"`
}

// UserRepository は users コレクションへのアクセスをまとめます。
type UserRepository struct {
	client *gofirestore.Client
}

// NewUserRepository は Firestore client を使って UserRepository を生成します。
func NewUserRepository(client *gofirestore.Client) *UserRepository {
	return &UserRepository{client: client}
}

// FindByID は users/{userID} を取得します。
func (r *UserRepository) FindByID(ctx context.Context, userID string) (*User, error) {
	doc, err := r.client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	var user User
	if decodeErr := doc.DataTo(&user); decodeErr != nil {
		return nil, decodeErr
	}
	user.ID = doc.Ref.ID
	return &user, nil
}

// FindByLoginID は document id と email の両方でログイン対象ユーザーを検索します。
func (r *UserRepository) FindByLoginID(ctx context.Context, loginID string) (*User, error) {
	// email 形式の入力は email 検索を優先し、document ID と email が衝突した場合の誤認証を避ける。
	if strings.Contains(loginID, "@") {
		return r.findByEmail(ctx, loginID)
	}

	// まず users/{loginID} を直接取得し、通常のログイン ID として扱う。
	doc, err := r.client.Collection("users").Doc(loginID).Get(ctx)
	if err == nil {
		var user User
		if decodeErr := doc.DataTo(&user); decodeErr != nil {
			return nil, decodeErr
		}
		user.ID = doc.Ref.ID
		return &user, nil
	}
	// NotFound 以外は接続失敗などの可能性があるため、そのまま上位に返す。
	if status.Code(err) != codes.NotFound {
		return nil, err
	}

	return r.findByEmail(ctx, loginID)
}

func (r *UserRepository) findByEmail(ctx context.Context, email string) (*User, error) {
	// document id で見つからない場合は、email として一致するユーザーを探す。
	iter := r.client.Collection("users").Where("email", "==", email).Limit(1).Documents(ctx)
	defer iter.Stop()

	snap, iterErr := iter.Next()
	if iterErr != nil {
		if errors.Is(iterErr, iterator.Done) {
			return nil, ErrUserNotFound
		}
		return nil, iterErr
	}

	var user User
	if decodeErr := snap.DataTo(&user); decodeErr != nil {
		return nil, decodeErr
	}
	// Firestore の document id をアプリケーション上の正規ユーザー ID として扱う。
	user.ID = snap.Ref.ID
	return &user, nil
}
