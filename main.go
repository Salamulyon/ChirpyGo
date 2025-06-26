package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Salamulyon/ChirpyGo.git/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	secret := os.Getenv("secret")

	const port = "8080"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cant access database")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("Platform must be set")
	}

	filepathRoot := "."
	fileServer := http.FileServer(http.Dir(filepathRoot))
	apiCfg := apiConfig{platform: platform,
		secretKey: secret,
		dbQueries: database.New(db)}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", isServerReady)

	mux.HandleFunc("GET /admin/metrics", apiCfg.middlewareMetricsWrite)
	mux.HandleFunc("POST /admin/reset", apiCfg.middlewareMetricsReset)

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	mux.HandleFunc("POST /api/login", apiCfg.handlerUserLogin)

	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileServer)))

	server := &http.Server{
		Addr:        ":" + port,
		Handler:     mux,
		ReadTimeout: 10 * time.Second,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
