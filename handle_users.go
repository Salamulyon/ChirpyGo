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
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	Password    string    `json:"hashed_password"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
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
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
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

	hashed_pw, err := auth.HashPassword(userUpdate.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}
	userUpdate.Password = hashed_pw

	id, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid user or password", err)
		return
	}

	newUser, _ := cfg.dbQueries.UpdateUserEmailAndPassword(r.Context(), database.UpdateUserEmailAndPasswordParams{
		ID:             id,
		Email:          userUpdate.Email,
		HashedPassword: hashed_pw,
	})

	respondWithJSON(w, http.StatusOK, User{
		ID:          newUser.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Email:       newUser.Email,
		IsChirpyRed: newUser.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlerUpgradeUserToChirpyRed(w http.ResponseWriter, r *http.Request) {

	type Datastruct struct {
		User_id string `json:"user_id"`
	}
	type Webhook struct {
		Event string     `json:"event"`
		Data  Datastruct `json:"data"`
	}

	requestApiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldnt get apikey", err)
		return
	}

	if requestApiKey != cfg.polkaAPIKey {
		respondWithError(w, http.StatusUnauthorized, "Wrong apikey", err)
		return
	}

	webhook := Webhook{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&webhook)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt decode the webhook", err)
		return
	}

	if webhook.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, "User should not be upgraded")
		return
	}

	parsedUserId, _ := uuid.Parse(webhook.Data.User_id)

	err = cfg.dbQueries.UpgradeUserToChirpyRed(r.Context(), parsedUserId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldnt find the user id", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")

}
