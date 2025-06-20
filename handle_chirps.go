package main

import (
	"encoding/json"
	"net/http"
	"strings"
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

type parameters struct {
	Body    string    `json:"body"`
	User_ID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Message could not be decoded", err)
		return
	}

	cleaned := handlerChirpsValidate(w, params.Body)
	params.Body = cleaned
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{Body: params.Body,
		UserID: params.User_ID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt create chirp", err)
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		Id:         chirp.ID,
		Body:       params.Body,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		User_id:    chirp.UserID,
	})

}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	allChirps := []Chirp{}
	// type allChirps struct {
	// 	Chirps []Chirp
	// }
	//totalChirps := allChirps{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&allChirps); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt decode chirps", err)
	}

	err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt get chirps", err)
	}

	for _, chirp := range allChirps {
		respondWithJSON(w, http.StatusCreated, chirp)
	}

	//espondWithJSON(w, http.StatusCreated, totalChirps)
}

func handlerChirpsValidate(w http.ResponseWriter, Body string) string {

	const maxChirpLength = 140
	if len(Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return ""
	}

	cleaned := getCleanedBody(Body)

	return cleaned

}

func getCleanedBody(body string) string {

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
