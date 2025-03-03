package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirpsByIdHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Chirp
	}

	idString := r.PathValue("chirpID")
	id, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID format", err)
		return
	}
	chirp, err := cfg.dbQueries.ChirpsById(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Chirp: Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		},
	})
}
