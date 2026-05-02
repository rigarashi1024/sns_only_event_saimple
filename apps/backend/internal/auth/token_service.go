package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 24 * time.Hour
)

// TokenService は internal JWT の発行と DB 保存用トークン暗号化を担当します。
type TokenService struct {
	jwtSigningKey []byte
	dbEncryptKey  []byte
}

// IssuedTokens はログイン成功時に発行されるアプリ内部用トークン一式です。
type IssuedTokens struct {
	SessionID                       string
	JTI                             string
	InternalAccessToken             string
	ProviderAccessToken             string
	ProviderRefreshToken            string
	ProviderAccessTokenEncrypted    string
	ProviderRefreshTokenEncrypted   string
	InternalAccessTokenExpiresAt    time.Time
	ProviderAccessTokenExpiresAt    time.Time
	ProviderRefreshTokenExpiresAt   time.Time
	InternalAccessTokenExpiresInSec int32
}

type internalJWTClaims struct {
	Subject   string `json:"sub"` // 認証済みユーザー ID
	SessionID string `json:"sid"` // sessions ドキュメント ID と対応するセッション ID
	JTI       string `json:"jti"` // 発行した JWT 自体を識別する ID
	IssuedAt  int64  `json:"iat"` // JWT 発行日時の Unix time
	ExpiresAt int64  `json:"exp"` // JWT 失効日時の Unix time
}

// NewTokenService は環境ごとの秘密値から JWT 署名鍵と DB 暗号化鍵を導出します。
func NewTokenService(secret string) (*TokenService, error) {
	if secret == "" {
		return nil, errors.New("token secret is required")
	}

	return &TokenService{
		jwtSigningKey: deriveKey(secret, "internal-jwt-signing"),
		dbEncryptKey:  deriveKey(secret, "db-token-encryption"),
	}, nil
}

// IssueSessionTokens は session_id / jti / internal JWT / local provider token を発行します。
func (s *TokenService) IssueSessionTokens(userID string, now time.Time) (*IssuedTokens, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	// session_id はサーバー側でログインセッションを識別するための ID。
	sessionID, err := randomID("session")
	if err != nil {
		return nil, fmt.Errorf("failed to generate session id: %w", err)
	}
	// jti は JWT そのものを識別するための ID。初期実装では session_id とは別に発行する。
	jti, err := randomID("jti")
	if err != nil {
		return nil, fmt.Errorf("failed to generate jti: %w", err)
	}
	// local ログインでは、将来の OAuth provider から受け取る token と同じ意味の値を自前発行する。
	providerAccessToken, err := randomID("provider_access")
	if err != nil {
		return nil, fmt.Errorf("failed to generate provider access token: %w", err)
	}
	providerRefreshToken, err := randomID("provider_refresh")
	if err != nil {
		return nil, fmt.Errorf("failed to generate provider refresh token: %w", err)
	}

	internalAccessTokenExpiresAt := now.Add(accessTokenTTL)
	providerAccessTokenExpiresAt := now.Add(accessTokenTTL)
	providerRefreshTokenExpiresAt := now.Add(refreshTokenTTL)
	// フロント用 access token は、session_id / user_id / jti / exp を含む internal JWT として発行する。
	internalAccessToken, err := s.signInternalJWT(internalJWTClaims{
		Subject:   userID,
		SessionID: sessionID,
		JTI:       jti,
		IssuedAt:  now.Unix(),
		ExpiresAt: internalAccessTokenExpiresAt.Unix(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to sign internal jwt: %w", err)
	}

	// provider token は将来の外部 OAuth token と同じ扱いで、Firestore には暗号化済み値だけを保存する。
	providerAccessTokenEncrypted, err := s.encryptForDB(providerAccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt provider access token: %w", err)
	}
	providerRefreshTokenEncrypted, err := s.encryptForDB(providerRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt provider refresh token: %w", err)
	}

	return &IssuedTokens{
		SessionID:                       sessionID,
		JTI:                             jti,
		InternalAccessToken:             internalAccessToken,
		ProviderAccessToken:             providerAccessToken,
		ProviderRefreshToken:            providerRefreshToken,
		ProviderAccessTokenEncrypted:    providerAccessTokenEncrypted,
		ProviderRefreshTokenEncrypted:   providerRefreshTokenEncrypted,
		InternalAccessTokenExpiresAt:    internalAccessTokenExpiresAt,
		ProviderAccessTokenExpiresAt:    providerAccessTokenExpiresAt,
		ProviderRefreshTokenExpiresAt:   providerRefreshTokenExpiresAt,
		InternalAccessTokenExpiresInSec: int32(accessTokenTTL / time.Second),
	}, nil
}

func (s *TokenService) signInternalJWT(claims internalJWTClaims) (string, error) {
	// JWT header は現時点では HMAC-SHA256 のみを使う。
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	signingInput := base64.RawURLEncoding.EncodeToString(headerJSON) + "." +
		base64.RawURLEncoding.EncodeToString(claimsJSON)
	// header.payload を HMAC-SHA256 で署名し、JWT の第 3 要素として付与する。
	mac := hmac.New(sha256.New, s.jwtSigningKey)
	_, _ = mac.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + signature, nil
}

func (s *TokenService) encryptForDB(plaintext string) (string, error) {
	// AES-GCM は暗号化と改ざん検知を同時に行えるため、DB 保存用 token 暗号化に使う。
	block, err := aes.NewCipher(s.dbEncryptKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// GCM の nonce は同じ鍵で再利用してはいけないため、暗号化ごとにランダム生成する。
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	// 復号時に nonce が必要になるため、nonce + ciphertext をまとめて保存する。
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	payload := append(nonce, ciphertext...)
	return base64.RawURLEncoding.EncodeToString(payload), nil
}

func deriveKey(secret string, purpose string) []byte {
	// 同じ TOKEN_ENCRYPTION_KEY から用途別の 32 bytes 鍵を作り、署名用と暗号化用を分離する。
	hash := sha256.Sum256([]byte(purpose + ":" + secret))
	return hash[:]
}

func randomID(prefix string) (string, error) {
	// 256bit のランダム値を URL-safe base64 にして、推測困難な ID/token として使う。
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return prefix + "_" + base64.RawURLEncoding.EncodeToString(token), nil
}
