package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Salamulyon/ChirpyGo.git/internal/auth"
	"github.com/Salamulyon/ChirpyGo.git/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"` //might have to change the json to exp or remove completely because apparently it's always set to an hour
	}

	type UserResponse struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
	}

	const minutesInHour int = 60
	const secondsInMinutes int = 60

	login := LoginRequest{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&login)

	login.ExpiresInSeconds = minutesInHour * 1 * secondsInMinutes * int(time.Second)
	refreshTokenExpiryDate := time.Now().Add(time.Hour * 24 * 60) // adding 60 days to the current time for the resfresh token expiry date

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

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not make resfresh token", err)

	}

	cfg.dbQueries.CreateUserRefreshToken(r.Context(), database.CreateUserRefreshTokenParams{
		Token:     refreshToken,
		ExpiresAt: refreshTokenExpiryDate,
		UserID:    dbUser.ID,
	})

	respondWithJSON(w, http.StatusOK, UserResponse{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  dbUser.IsChirpyRed,
	})

}
