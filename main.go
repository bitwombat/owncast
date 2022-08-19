package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bitwombat/listen-later/service"
)

func main() {
	log.Print("INFO: starting server...")
	http.HandleFunc("/", service.UpdateHandler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("INFO: defaulting to port %s", port)
	}

	log.Printf("INFO: listening on port %s", port)

	// Start HTTP server.
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Fatal error serving: ", err)
	}
}
