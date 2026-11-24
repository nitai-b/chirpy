package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	httpServer := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", handlerHealth)

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte("OK"))
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}
