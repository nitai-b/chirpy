package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	httpServer := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
