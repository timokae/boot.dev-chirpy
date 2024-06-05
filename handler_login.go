package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

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

	jwtToken, err := auth.NewJWTToken(cfg.jwtSecret, user.Id, auth.JwtTokenDefaultExpiration)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	refreshToken, err := auth.NewRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate refresh token")
		return
	}

	user.RefreshTokenExpiration = time.Now().UTC().Add(auth.RefreshTokenDefaultExpiration)
	user.RefreshToken = refreshToken
	user, err = cfg.db.UpdateUser(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate refresh token")
	}

	respondWithJSON(w, http.StatusOK, struct {
		Id           int    `json:"id"`
		Email        string `json:"email"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		Id:           user.Id,
		Email:        user.Email,
		Token:        jwtToken,
		RefreshToken: refreshToken,
	})
}
