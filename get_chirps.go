package main

import (
	"net/http"
	"sort"

	"github.com/cygran/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		dbChirps []database.Chirp
		err      error
	)
	sortQuery := r.URL.Query().Get("sort")
	if sortQuery == "" {
		sortQuery = "asc"
	}
	if sortQuery != "asc" && sortQuery != "desc" {
		respondWithError(w, http.StatusBadRequest, "Invalid sort query", err)
		return
	}
	idString := r.URL.Query().Get("author_id")
	if idString != "" {
		authorId, err := uuid.Parse(idString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id", err)
			return
		}
		dbChirps, err = cfg.dbQueries.GetChirpsByUserId(r.Context(), authorId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrive chirps by author", err)
			return
		}
	} else {
		dbChirps, err = cfg.dbQueries.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
			return
		}
		if sortQuery == "desc" {
			sort.Slice(dbChirps, func(i, j int) bool {
				return dbChirps[i].CreatedAt.After(dbChirps[j].CreatedAt)
			})
		} else if sortQuery == "asc" {
			sort.Slice(dbChirps, func(i, j int) bool {
				return dbChirps[i].CreatedAt.Before(dbChirps[j].CreatedAt)
			})
		}
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
