package main

import (
	"net/http"
	"strconv"

	auth "github.com/timokae/boot.dev-chirpy-auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	userId, err := auth.VerifyJWTToken(authToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Provide a valid ID")
		return
	}

	chirp, err := cfg.db.GetChirp(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	if chirp.AuthorId != userId {
		respondWithError(w, http.StatusForbidden, "cannot delete foreign chirps")
		return
	}

	cfg.db.DeleteChirp(id)

	w.WriteHeader(http.StatusNoContent)
}
