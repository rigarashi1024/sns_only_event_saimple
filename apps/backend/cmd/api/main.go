package main

import (
	"context"
	"log"
	"net/http"

	api "github.com/rigarashi1024/sns_only_event_saimple/apps/backend/gen"
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

	// OpenAPI から生成したルーティングに、実装したハンドラと CORS 設定を接続する。
	handler := httpHandler.NewHandler(firestoreClient)
	server := api.Handler(handler)
	server = httpHandler.WithCORS(server)

	log.Println("listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", server))
}
