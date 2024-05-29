package main

import (
	"log"
	"net/http"
	"sort"
	"strconv"

	database "github.com/timokae/boot.dev-chirpy-database"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps()
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirps")
	}

	chirps := []database.Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, database.Chirp{
			Id:   dbChirp.Id,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Provide a valid ID")
		return
	}

	chirp, err := cfg.db.GetChirp(int(id))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find chirp with id")
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
