package main

import (
	"log"
	"net/http"
	"time"

	auth "github.com/timokae/boot.dev-chirpy-auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := cfg.db.GetUserByRefreshToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if user.RefreshTokenExpiration.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	jwtToken, err := auth.NewJWTToken(cfg.jwtSecret, user.Id, auth.RefreshTokenDefaultExpiration)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: jwtToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := cfg.db.GetUserByRefreshToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user.RefreshToken = ""
	user.RefreshTokenExpiration = time.Time{}
	cfg.db.UpdateUser(user)

	w.WriteHeader(http.StatusNoContent)
}
