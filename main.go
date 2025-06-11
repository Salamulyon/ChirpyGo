package main

import (
	"log"
	"net/http"
	"time"
)

func main() {

	filepathRoot := "."
	fileServer := http.FileServer(http.Dir(filepathRoot))
	apiCfg := apiConfig{}
	const port = "8080"

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", isServerReady)

	mux.HandleFunc("POST /admin/reset", apiCfg.middlewareMetricsReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.middlewareMetricsWrite)

	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileServer)))

	server := &http.Server{
		Addr:        ":" + port,
		Handler:     mux,
		ReadTimeout: 10 * time.Second,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
