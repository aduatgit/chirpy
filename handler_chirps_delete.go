package main

import (
	"net/http"
	"strconv"

	"github.com/aduatgit/chirpy/internal/auth"
	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	author, err := auth.CheckAuthorization(r.Header, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 403, err.Error())
		return
	}
	authorID, err := strconv.Atoi(author)
	if err != nil {
		respondWithError(w, 403, err.Error())
		return
	}
	chirpIDString := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, 403, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, 403, "Invalid chirp ID")
		return
	}

	if chirp.AuthorID != authorID {
		respondWithError(w, 403, "Not author of Chirp")
		return
	}

	err = cfg.DB.DeleteChirp(chirpID)
	if err != nil {
		respondWithError(w, 403, "Couldn't delete chirps")
		return
	}

	w.WriteHeader(200)
}
