package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	r := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	rApi := chi.NewRouter()
	rApi.Get("/healthz", handlerReadiness)
	rApi.Get("/reset", apiCfg.handlerMetricsReset)
	r.Mount("/api", rApi)

	rAdmin := chi.NewRouter()
	rAdmin.Get("/metrics", apiCfg.handlerMetrics)
	r.Mount("/admin", rAdmin)
	corsMux := middlewareCors(r)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from%s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

/* Structs */

type apiConfig struct {
	fileserverHits int
}
