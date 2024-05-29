package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameter")
		return
	}

	user, err := cfg.db.CreateUser(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create user")
	}

	respondWithJSON(w, http.StatusCreated, struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}{
		Id:    user.Id,
		Email: user.Email,
	})
}
