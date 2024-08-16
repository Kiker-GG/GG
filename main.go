package main

import (
	"flag"
	"http_server/http"
	"http_server/storage"
	"log"
)

// @title My API
// @version 1.0
// @description This is a sample server.
// @host 127.0.0.1:8000
// @BasePath /
func main() {
	addr := flag.String("addr", "0.0.0.0:8000", "address for http server")

	s := storage.NewDatabase()

	log.Printf("Starting server on %s", *addr)
	if err := http.CreateNewServer(s, *addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
