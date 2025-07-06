package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/Mickunaru/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != os.Getenv("POLKA_KEY") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type eventData struct {
		UserID string `json:"user_id"`
	}
	type parameters struct {
		Event string    `json:"event"`
		Data  eventData `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if decoder.Decode(&params) != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID", err)
		return
	}

	_, err = cfg.db.UpdateUserToRed(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't upgrade user to red", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
