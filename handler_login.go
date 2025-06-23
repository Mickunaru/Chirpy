package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Mickunaru/Chirpy/internal/auth"
	"github.com/Mickunaru/Chirpy/internal/database"
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
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil || auth.CheckPasswordHash(user.HashedPassword, params.Password) != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't make JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate refresh token", err)
		return
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(24*60)),
		RevokedAt: sql.NullTime{},
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't store refresh token to database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	})
}
