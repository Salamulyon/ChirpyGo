package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Salamulyon/ChirpyGo.git/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Body       string    `json:"body"`
	Id         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	User_id    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type CreateChirpParams struct {
		Body    string    `json:"body"`
		User_ID uuid.UUID `json:"user_id"`
	}

	params := CreateChirpParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Message could not be decoded", err)
		return
	}

	cleaned := handlerChirpsValidate(w, params.Body)
	params.Body = cleaned
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{Body: params.Body,
		User_ID: params.User_ID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt create chirp", err)
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		Id:         chirp.ID,
		Body:       params.Body,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		User_id:    chirp.User_ID,
	})

}
