package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type respBody struct {
	Body string `json: "body"`
}

type errorReply struct {
	Bodyerror string `json:"error"`
}

type validResponse struct {
	Cleaned_body string `json:"cleaned_body"`
}

var errorResponse = errorReply{
	Bodyerror: "Something went wrong",
}

func reqHandler(w http.ResponseWriter, r *http.Request) {
	response := respBody{
		Body: "",
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&response)
	if err != nil {
		w.WriteHeader(500)
		if err := respHandler(errorResponse, w); err != nil {
			w.Write([]byte(`{"error": "Something went wrong"}`))
		}
	}

	respwithoutspaces := strings.ReplaceAll(response.Body, " ", "")
	if len(respwithoutspaces) > 140 {
		errorResponse := errorReply{
			Bodyerror: "Chirp is too long",
		}
		w.WriteHeader(400)
		if err := respHandler(errorResponse, w); err != nil {
			w.Write([]byte(`{"error": "Something went wrong"}`))
		}

		return
	} else {
		newText := cleanWords(response.Body)
		valid := validResponse{
			Cleaned_body: newText,
		}
		w.WriteHeader(200)
		if err := respHandler(valid, w); err != nil {
			w.Write([]byte(`{"error": "Something went wrong"}`))
		}
		return
	}

}

func respHandler(resp any, w http.ResponseWriter) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	return nil
}

func cleanWords(text string) string {

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(text, " ")
	for i, str := range words {
		if _, ok := badWords[strings.ToLower(str)]; ok {
			words[i] = "****"
		}
	}
	cleaned_body := strings.Join(words, " ")
	return cleaned_body
}
