package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cygran/chirpy/internal/auth"
)

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password  string `json:"password"`
		Email     string `json:"email"`
		ExpiresIn *int   `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode paramaters", err)
		return
	}
	expiresIn := 3600
	if params.ExpiresIn != nil {
		userExpiresIn := *params.ExpiresIn
		if userExpiresIn > 0 && userExpiresIn < 3600 {
			expiresIn = userExpiresIn
		}
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
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(expiresIn)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to generate token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})
}
