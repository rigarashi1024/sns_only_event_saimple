package main

import (
	"log"
	"net/http"

	api "github.com/rigarashi1024/sns_only_event_saimple/apps/backend/gen"
	httpHandler "github.com/rigarashi1024/sns_only_event_saimple/apps/backend/internal/http"
)

func main() {
	handler := httpHandler.NewHandler()
	server := api.Handler(handler)
	server = httpHandler.WithCORS(server)

	log.Println("listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", server))
}
