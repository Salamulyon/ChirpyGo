package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Salamulyon/ChirpyGo.git/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		pw        []byte    `json:"hashed_password"`
	}
	user := User{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&user)

	dbUser, err := cfg.dbQueries.GetUserUsingEmail(r.Context(), user.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User not Found", err)
		return
	}

	pw_err := auth.ComparePasswords(string(user.pw), dbUser.HashedPassword)
	if pw_err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", pw_err)
		return
	}

	respondWithJSON(w, http.StatusOK, dbUser)

}
