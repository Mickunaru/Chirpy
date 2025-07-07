package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/Mickunaru/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	queryParams := r.URL.Query()

	authorID := queryParams.Get("author_id")

	var chirps []database.Chirp
	var err error
	if authorID == "" {
		chirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
			return
		}
	} else {
		parsedID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't parse author ID", err)
			return
		}

		chirps, err = cfg.db.GetChirpsByUserId(r.Context(), parsedID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
			return
		}
	}

	res := []Chirp{}
	for _, c := range chirps {
		res = append(res, Chirp(c))
	}

	sortType := queryParams.Get("sort")

	if sortType == "desc" {
		sort.Slice(res, func(i, j int) bool {
			return res[i].CreatedAt.UnixMilli() > res[j].CreatedAt.UnixMilli()
		})
	} else {
		sort.Slice(res, func(i, j int) bool {
			return res[i].CreatedAt.UnixMilli() < res[j].CreatedAt.UnixMilli()
		})
	}

	respondWithJSON(w, http.StatusOK, res)
}
