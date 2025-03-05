package main

import (
	"net/http"
	"time"

	"github.com/cygran/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or missing Authorization header", err)
		return
	}
	tokenInfo, err := cfg.dbQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil || tokenInfo.ExpiresAt.Before(time.Now().UTC()) || tokenInfo.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token", err)
		return
	}
	userRecord, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to retrieve user for refresh token", err)
		return
	}
	newToken, err := auth.MakeJWT(userRecord.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate new token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		Token: newToken,
	})
}
