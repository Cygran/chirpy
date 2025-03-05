package main

import (
	"net/http"
	"time"

	"github.com/cygran/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
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
	err = cfg.dbQueries.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
