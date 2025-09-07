package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/theokyle/chirpygo/internal/database"
)

func main() {
	// Import environment variables
	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	port := os.Getenv("PORT")
	platform := os.Getenv("PLATFORM")
	const filepathRoot = "."

	// Connect to sql database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("Error accessing database: %s", err)
	}

	dbQueries := database.New(db)

	mux := http.NewServeMux()
	srv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
	}

	//Admin Routes
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	//User Routes
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	//Chirp Routes
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)

	log.Printf("Serving on port: %s", port)
	log.Fatal(srv.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}
