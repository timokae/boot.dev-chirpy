package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	database "github.com/timokae/boot.dev-chirpy-database"
)

func main() {
	const port = "8080"
	const fileRootPath = "."
	const dbPath = "./database.json"

	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Could not read .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalln("JWT_SECRET not present")
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *dbg {
		err := database.DropDB(dbPath)
		if err != nil {
			log.Panicln(err)
		} else {
			log.Println("Database dropped")
		}
	}

	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Println(err)
	}

	apiCfg := apiConfig{
		fileServerHits: 0,
		db:             db,
		jwtSecret:      jwtSecret,
	}

	mux := http.NewServeMux()

	// mux.Handle("GET /app/admin/", )
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(fileRootPath)))))

	// ADMIN
	mux.HandleFunc("GET /admin/metrics/", apiCfg.handlerMetrics)

	// API
	mux.HandleFunc("GET /api/healthz", apiCfg.handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)

	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerGetChirp)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("DELETE /api/chirps/{id}", apiCfg.handlerDeleteChirp)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	server.ListenAndServe()
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
}
