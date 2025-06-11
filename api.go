package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type respBody struct {
	body string `json: "body"`
}

type errorReply struct {
	error string `json:"error"`
}

type validResponse struct {
	valid bool `json:"valid"`
}

func reqHandler(w http.ResponseWriter, r *http.Request) {
	response := respBody{
		body: "",
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&response)
	if err != nil {
		errorResponse := errorReply{
			error: "Something went wrong",
		}
		w.WriteHeader(500)
		err := respHandler(errorResponse, w)
		if err != nil {
			log.Printf("error : %s", err)
		}
		return
	}

	respwithoutspaces := strings.ReplaceAll(response.body, " ", "")
	if len(respwithoutspaces) > 140 {
		errorResponse := errorReply{
			error: "Chirp is too long",
		}
		w.WriteHeader(400)
		err := respHandler(errorResponse, w)
		if err != nil {
			log.Printf("error : %s", err)
		}
		return
	} else {
		valid := validResponse{
			valid: true,
		}
		w.WriteHeader(200)
		err := respHandler(valid, w)
		if err != nil {
			log.Printf("error : %s", err)
		}
		return
	}

}

func respHandler(resp any, w http.ResponseWriter) error {
	data, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	return nil
}
