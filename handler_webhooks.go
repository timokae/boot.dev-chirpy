package main

import (
	"encoding/json"
	"net/http"

	auth "github.com/timokae/boot.dev-chirpy-auth"
)

func (cfg *apiConfig) handlerWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId int `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetApiToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters")
		return
	}

	user, err := cfg.db.GetUserById(params.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Cannot find user")
		return
	}

	if params.Event == "user.upgraded" {
		user.IsChirpyRed = true
		_, err = cfg.db.UpdateUser(user)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Could not update user")
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
