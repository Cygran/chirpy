package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cygran/chirpy/internal/auth"
	"github.com/cygran/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode paramaters", err)
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}
	createUserParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}
	user, err := cfg.dbQueries.CreateUser(r.Context(), createUserParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			ChirpyRed: user.IsChirpyRed,
		},
	})

}

func (cfg *apiConfig) usersHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or missing Authorization header", err)
		return
	}
	uuid, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to validate token", err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode paramaters", err)
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}
	updateUserParams := database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             uuid,
	}
	updatedUser, err := cfg.dbQueries.UpdateUser(r.Context(), updateUserParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        updatedUser.ID,
			CreatedAt: updatedUser.CreatedAt,
			UpdatedAt: updatedUser.UpdatedAt,
			Email:     updatedUser.Email,
			ChirpyRed: updatedUser.IsChirpyRed,
		},
	})

}

func (cfg *apiConfig) userUpgradeHandler(w http.ResponseWriter, r *http.Request) {
	type webhookRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	headerKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or missing Authorization header", err)
		return
	}
	if headerKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API Key", err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := webhookRequest{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode paramaters", err)
		return
	}
	eventType := params.Event
	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Unable to parse UUID", err)
	}
	if eventType != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	err = cfg.dbQueries.UpgradeUser(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Database query failed", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
