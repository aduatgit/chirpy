package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aduatgit/chirpy/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_APIKEY")
	const filepathRoot = "."
	const port = "8080"

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *dbg {
		os.Remove("database.json") // please handle the error properly
		fmt.Println("Debug mode enabled, database deleted.")
	} else {
		fmt.Println("Debug mode not enabled, database safe.")
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
		polkaKey:       polkaKey,
	}

	r := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	rApi := chi.NewRouter()
	rApi.Get("/healthz", handlerReadiness)
	rApi.Get("/reset", apiCfg.handlerMetricsReset)

	rApi.Post("/chirps", apiCfg.handlerChirpsCreate)
	rApi.Get("/chirps", apiCfg.handlerChirpsRetrieve)
	rApi.Get("/chirps/{chirpID}", apiCfg.handlerChirpsGet)
	rApi.Delete("/chirps/{chirpID}", apiCfg.handlerChirpsDelete)

	rApi.Post("/users", apiCfg.handlerUsersCreate)
	rApi.Put("/users", apiCfg.handlerUsersUpdate)
	rApi.Post("/login", apiCfg.handlerLogin)

	rApi.Post("/refresh", apiCfg.handlerRefresh)
	rApi.Post("/revoke", apiCfg.handlerRevoke)

	rApi.Post("/polka/webhooks", apiCfg.handlerPolkaWebhook)
	r.Mount("/api", rApi)

	rAdmin := chi.NewRouter()
	rAdmin.Get("/metrics", apiCfg.handlerMetrics)
	r.Mount("/admin", rAdmin)
	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

/* Structs */

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
	polkaKey       string
}
