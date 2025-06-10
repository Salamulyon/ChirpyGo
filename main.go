package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func isServerReady(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))

}

func (cfg *apiConfig) middlewareMetricsInc(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Add(1)

}

var updateMetrics http.Handler = http.HandlerFunc(middlewareMetricsInc)

func (cfg *apiConfig) writeNumberOfHits(w http.ResponseWriter, req *http.Request) {

	hits := cfg.fileserverHits.Load()
	data := fmt.Sprintf("Hits: %s", strconv.Itoa(int(hits)))
	w.Write([]byte(data))

}

func (cfg *apiConfig) resetNumberOfHits(w http.ResponseWriter, req *http.Request) {

	cfg.fileserverHits.Store(0)

}

func main() {

	filepathRoot := "."
	fileServer := http.FileServer(http.Dir(filepathRoot))
	apiCfg := apiConfig{}
	const port = "8080"

	mux := http.NewServeMux()

	mux.HandleFunc("/", isServerReady)
	mux.HandleFunc("/metrics", apiCfg.writeNumberOfHits)
	mux.HandleFunc("/reset", apiCfg.resetNumberOfHits)

	mux.Handle("/app/", (http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileServer))))

	server := &http.Server{
		Addr:        ":" + port,
		Handler:     mux,
		ReadTimeout: 10 * time.Second,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
