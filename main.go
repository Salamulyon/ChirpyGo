package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func serverReady(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))

}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	cfg.fileserverHits.Add(1)
	return next
}

func (cfg *apiConfig) getNumberOfHits(next http.Handler) atomic.Int32 {

	return cfg.fileserverHits

}

func main() {

	filepathRoot := "."
	fileServer := http.FileServer(http.Dir(filepathRoot))
	apiCfg := apiConfig{}
	const port = "8080"

	mux := http.NewServeMux()
	mux.HandleFunc("/", serverReady)
	mux.Handle("/app/", (http.StripPrefix("/app", fileServer)))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServer))

	//mux.Handle("assets/logo.png", http.FileServer(http.Dir(imagePath)))
	server := &http.Server{
		Addr:        ":" + port,
		Handler:     mux,
		ReadTimeout: 10 * time.Second,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
