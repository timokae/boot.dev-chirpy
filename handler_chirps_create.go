package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	auth "github.com/timokae/boot.dev-chirpy-auth"
	database "github.com/timokae/boot.dev-chirpy-database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decorder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decorder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters")
		return
	}

	authToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
	}

	id, err := auth.VerifyJWTToken(authToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirp, err := cfg.db.CreateChirp(cleaned, id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp")
	}

	respondWithJSON(w, http.StatusCreated, database.Chirp{
		Id:       chirp.Id,
		Body:     chirp.Body,
		AuthorId: chirp.AuthorId,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

func getCleanedBody(message string, badWords map[string]struct{}) string {
	words := strings.Split(message, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		_, ok := badWords[loweredWord]
		if ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
