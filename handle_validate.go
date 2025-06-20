package main

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type parameters struct {
	Body   string
	UserID uuid.NullUUID
}

func handlerChirpsValidate(w http.ResponseWriter, Body string) string {

	// decoder := json.NewDecoder(r.Body)
	// params := parameters{}
	// err := decoder.Decode(&params)
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
	// 	return
	// }

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
