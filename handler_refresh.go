package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Mickunaru/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type Token struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get bearer token", err)
		return
	}

	refreshTokenObj, err := cfg.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token", err)
		return
	}

	if refreshTokenObj.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired", err)
		return
	}

	if refreshTokenObj.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token was revoked", errors.New("refresh token was revoked"))
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshTokenObj.Token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user from refresh token", err)
		return
	}

	newToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't create JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Token{
		Token: newToken,
	})
}
