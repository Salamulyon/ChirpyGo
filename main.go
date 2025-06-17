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
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cant access database")
	}

	filepathRoot := "."
	fileServer := http.FileServer(http.Dir(filepathRoot))
	apiCfg := apiConfig{}
	const port = "8080"

	apiCfg.dbQueries = database.New(db)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/healthz", isServerReady)

	mux.HandleFunc("GET /admin/metrics", apiCfg.middlewareMetricsWrite)
	mux.HandleFunc("POST /admin/reset", apiCfg.middlewareMetricsReset)

	mux.HandleFunc("POST /api/validate_chirp", reqHandler)

	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fileServer)))

	server := &http.Server{
		Addr:        ":" + port,
		Handler:     mux,
		ReadTimeout: 10 * time.Second,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())

}
