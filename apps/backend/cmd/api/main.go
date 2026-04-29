package main

import (
  "log"
  "net/http"
)

func main() {
  http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
    _, _ = w.Write([]byte("ok"))
  })

  log.Println("listening on :8081")
  log.Fatal(http.ListenAndServe(":8081", nil))
}
