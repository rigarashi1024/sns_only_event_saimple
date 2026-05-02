package main

import (
	"context"
	"log"
	"net/http"

	api "github.com/rigarashi1024/sns_only_event_saimple/apps/backend/gen"
	"github.com/rigarashi1024/sns_only_event_saimple/apps/backend/internal/auth"
	"github.com/rigarashi1024/sns_only_event_saimple/apps/backend/internal/config"
	firestoreclient "github.com/rigarashi1024/sns_only_event_saimple/apps/backend/internal/firestore"
	httpHandler "github.com/rigarashi1024/sns_only_event_saimple/apps/backend/internal/http"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	// API 起動時に Firestore client を作成し、各ハンドラから共有して利用する。
	firestoreClient, err := firestoreclient.NewClient(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to create firestore client: %v", err)
	}
	defer firestoreClient.Close()

	// local では環境変数、dev/prd では将来 Secret Manager から token 用の秘密値を取得する。
	tokenSecret, err := config.GetDBEncryptionKey(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to load token encryption key: %v", err)
	}
	// tokenService は internal JWT 発行、provider token 発行、DB 保存用暗号化をまとめて扱う。
	tokenService, err := auth.NewTokenService(tokenSecret)
	if err != nil {
		log.Fatalf("failed to create token service: %v", err)
	}

	// OpenAPI から生成したルーティングに、実装したハンドラと CORS 設定を接続する。
	// local は HTTP で動かすため Secure=false、dev/prd は HTTPS 前提で Secure=true にする。
	cookieSecure := cfg.Env != config.EnvLocal
	handler := httpHandler.NewHandler(firestoreClient, tokenService, cookieSecure)
	server := api.Handler(handler)
	server = httpHandler.WithCORS(server)

	log.Println("listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", server))
}
