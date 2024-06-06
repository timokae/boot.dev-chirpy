package main

import (
	"log"
	"net/http"

	database "github.com/timokae/boot.dev-chirpy-database"
)

type apiConfig struct {
	fileServerHits int
	db             *database.DB
	jwtSecret      string
	polkaKey       string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits += 1
		log.Printf("%v: %v", r.URL.Path, cfg.fileServerHits)
		next.ServeHTTP(w, r)
	})
}
