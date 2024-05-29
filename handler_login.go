package main

import (
	"encoding/json"
	"log"
	"net/http"

	auth "github.com/timokae/boot.dev-chirpy-auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
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

	user, err := cfg.db.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}{
		Id:    user.Id,
		Email: user.Email,
	})
}
