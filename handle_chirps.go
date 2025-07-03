package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Salamulyon/ChirpyGo.git/internal/auth"
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

	bearerToken, bearerErr := auth.GetBearerToken(r.Header)
	if bearerErr != nil {
		respondWithError(w, http.StatusUnauthorized, "couldnt retrieve bearer", bearerErr)
		return
	}

	id, validateErr := auth.ValidateJWT(bearerToken, cfg.secretKey)
	if validateErr != nil {
		respondWithError(w, http.StatusUnauthorized, "couldnt parse id from jwt", validateErr)
		return
	}

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
		UserID: id})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldnt create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		Id:         chirp.ID,
		Body:       params.Body,
		Created_at: chirp.CreatedAt,
		Updated_at: chirp.UpdatedAt,
		User_id:    id,
	})

}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	var totalChirps []database.Chirp

	userId := r.URL.Query().Get("author_id")
	if userId != "" {
		parsedUserId, err := uuid.Parse(userId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "couldnt parse user", err)
			return
		}
		totalChirps, err := cfg.dbQueries.GetChirpsFromUser(r.Context(), parsedUserId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "couldnt encode", err)
			return
		}
		responseChirps := make([]Chirp, len(totalChirps))
		//converting database.Chirp to chirp struct
		for i, dbChirp := range totalChirps {
			responseChirps[i] = Chirp{
				Id:         dbChirp.ID,
				Body:       dbChirp.Body,
				Created_at: dbChirp.CreatedAt,
				Updated_at: dbChirp.UpdatedAt,
				User_id:    dbChirp.UserID,
			}
		}
		respondWithJSON(w, http.StatusOK, responseChirps)
		return

	}

	totalChirps, err := cfg.dbQueries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt encode", err)
		return
	}
	responseChirps := make([]Chirp, len(totalChirps))
	//converting database.Chirp to chirp struct
	for i, dbChirp := range totalChirps {
		responseChirps[i] = Chirp{
			Id:         dbChirp.ID,
			Body:       dbChirp.Body,
			Created_at: dbChirp.CreatedAt,
			Updated_at: dbChirp.UpdatedAt,
			User_id:    dbChirp.UserID,
		}
	}
	respondWithJSON(w, http.StatusOK, responseChirps)

}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {

	rawID := r.PathValue("chirpID")
	parsedID, err := uuid.Parse(rawID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't parse ID", err)
		return
	}

	foundChirp, err := cfg.dbQueries.GetChirp(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldnt find chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		Id:         foundChirp.ID,
		Body:       foundChirp.Body,
		Created_at: foundChirp.CreatedAt,
		Updated_at: foundChirp.UpdatedAt,
		User_id:    foundChirp.UserID,
	})

}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't retrieve token", err)
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "invalid user or password", err)
		return
	}

	rawID := r.PathValue("chirpID")
	parsedID, err := uuid.Parse(rawID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't parse ID", err)
		return
	}

	foundChirp, err := cfg.dbQueries.GetChirp(r.Context(), parsedID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldnt find chirp", err)
		return
	}

	if userId != foundChirp.UserID {
		respondWithError(w, http.StatusForbidden, "wrong user", err)
		return
	}

	err = cfg.dbQueries.DeleteChirp(r.Context(), foundChirp.ID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp could not be deleted", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
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
