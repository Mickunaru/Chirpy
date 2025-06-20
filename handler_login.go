package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Mickunaru/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil || auth.CheckPasswordHash(user.HashedPassword, params.Password) != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	var expiresIn int
	if params.ExpiresInSeconds == 0 {
		expiresIn = 60 * 60
	} else {
		expiresIn = params.ExpiresInSeconds
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Second*time.Duration(expiresIn))

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't make JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})
}
