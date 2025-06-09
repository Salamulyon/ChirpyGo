package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

func serverReady(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))

}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetrics(next http.Handler) http.Handler {

}

func main() {

	filepathRoot := "."
	//imagePath := "/assets/logo.png"
	//readiness := "/healthz"

	const port = "8080"

	mux := http.NewServeMux()
	mux.HandleFunc("/", serverReady)
	mux.Handle("/app/", (*apiConfig).middlewareMetrics((http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))))
	mux.Handle("/app/", (http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	//mux.Handle("assets/logo.png", http.FileServer(http.Dir(imagePath)))
	server := &http.Server{
		Addr:        ":" + port,
		Handler:     mux,
		ReadTimeout: 10 * time.Second,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
