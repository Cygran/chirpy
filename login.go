package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cygran/chirpy/internal/auth"
	"github.com/cygran/chirpy/internal/database"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	const refreshTokenExpiry = 60 * 24 * time.Hour

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode paramaters", err)
		return
	}
	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to generate token", err)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to generate refresh token", err)
		return
	}
	_, err = cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(refreshTokenExpiry).UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to store refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			ChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refreshToken,
	})
}
