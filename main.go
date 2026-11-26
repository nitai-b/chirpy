package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
)

func main() {
	mux := http.NewServeMux()

	httpServer := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	var cfg apiConfig

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/healthz", handlerHealth)
	mux.HandleFunc("/metrics", handlerMetrics)
	mux.HandleFunc("/reset", handlerResetMetrics)

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

func handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	var cfg apiConfig

	_, err := w.Write([]byte("Hits: " + strconv.Itoa(int(cfg.fileServerHits.Load()))))
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}

func handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	var cfg apiConfig
	cfg.fileServerHits.Store(0)

	_, err := w.Write([]byte("Hits: " + strconv.Itoa(int(cfg.fileServerHits.Load()))))
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
