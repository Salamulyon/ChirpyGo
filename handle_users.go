package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Salamulyon/ChirpyGo.git/internal/auth"
	"github.com/Salamulyon/ChirpyGo.git/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"hashed_password"`
}

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't deecode parameters", err)
		return
	}

	hashed_pw, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
		return
	}
	//params.HashedPassword = hashed_pw

	user, err := apiCfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{Email: params.Email, HashedPassword: hashed_pw})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}

func (cfg *apiConfig) handlerUpdateUserEmailAndPassword(w http.ResponseWriter, r *http.Request) {
	type UserUpdate struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithJSON(w, http.StatusUnauthorized, "")
		return
	}

	userUpdate := UserUpdate{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&userUpdate)

	hashed_pw, _ := auth.HashPassword(userUpdate.Password)
	userUpdate.Password = hashed_pw

	id, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "invalid user or password", err)
		return
	}

	newUser, _ := cfg.dbQueries.UpdateUserEmailAndPassword(r.Context(), database.UpdateUserEmailAndPasswordParams{
		ID:             id,
		Email:          userUpdate.Email,
		HashedPassword: hashed_pw,
	})

	respondWithJSON(w, http.StatusOK, User{
		ID:        newUser.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     newUser.Email,
	})
}
