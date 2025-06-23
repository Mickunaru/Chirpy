package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Mickunaru/Chirpy/internal/auth"
	"github.com/Mickunaru/Chirpy/internal/database"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get bearer token", err)
		return
	}

	err = cfg.db.UpdateRefreshTokenRevokedAt(r.Context(), database.UpdateRefreshTokenRevokedAtParams{
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Token: token,
	})
	if err != nil {
		fmt.Printf("DEBUG: Error revoking token: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
