package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func isServerReady(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain-text; charset=utf-8") // normal header
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))

}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareMetricsWrite(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	hits := cfg.fileserverHits.Load()
	data := fmt.Sprintf(`<html>
						<body>
							<h1>Welcome, Chirpy Admin</h1>
							<p>Chirpy has been visited %d times!</p>
						</body>
						</html>`, (hits))
	w.Write([]byte(data))

}

func (cfg *apiConfig) middlewareMetricsReset(w http.ResponseWriter, req *http.Request) {

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}
