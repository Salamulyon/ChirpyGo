package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Salamulyon/ChirpyGo.git/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type UserResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}

	const minutesInHour int = 60
	const secondsInMinutes int = 60

	login := LoginRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&login)

	if login.ExpiresInSeconds == 0 || login.ExpiresInSeconds*int(time.Second) > minutesInHour*secondsInMinutes*1*int(time.Second) {
		login.ExpiresInSeconds = minutesInHour * 1 * secondsInMinutes * int(time.Second)
	}

	dbUser, err := cfg.dbQueries.GetUserUsingEmail(r.Context(), login.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User not Found", err)
		return

	}

	pw_err := auth.ComparePasswords(dbUser.HashedPassword, login.Password)
	if pw_err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", pw_err)
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.secretKey, time.Duration(login.ExpiresInSeconds))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not create token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, UserResponse{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Token:     token,
	})

}
