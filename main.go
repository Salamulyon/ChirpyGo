package main

import (
	"log"
	"net/http"
	"time"
)

func main() {

	const port = "8080"

	serveMux := &http.ServeMux{}

	server := &http.Server{
		Addr:        ":" + port,
		Handler:     serveMux,
		ReadTimeout: 10 * time.Second,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}
