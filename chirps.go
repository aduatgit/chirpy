package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func checkProfanity(chirp string) string {
	s := strings.Split(chirp, " ")
	profanity := []string{"kerfuffle", "sharbert", "fornax"}
	for i := range s {
		if slices.Contains(profanity, strings.ToLower(s[i])) {
			s[i] = "****"
		}
	}
	return strings.Join(s, " ")
}

func (cfg *apiConfig) handlerChirpGetById(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "chirpid"))
	if err != nil {
		respondWithError(w, 400, "Couldn't convert chirpid to int")
		return
	}

	type requestBody struct {
		ID   int    `json:"id"`
		Body string `json:"body"`
	}

	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, 400, "Couldn't load Chirps")
		return
	}
	if len(chirps) < id {
		respondWithError(w, 404, "Chirp not found")
		return
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	chirp := chirps[id-1]

	reqBody := requestBody{
		ID:   chirp.ID,
		Body: chirp.Body,
	}

	respondWithJSON(w, http.StatusOK, reqBody)

}

func (cfg *apiConfig) handlerChirpCreate(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := requestBody{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode parameters")
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleaned := checkProfanity(params.Body)

	chirp, err := cfg.DB.CreateChirp(cleaned)
	if err != nil {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:   chirp.ID,
		Body: chirp.Body,
	})

}

func (cfg *apiConfig) handlerChirpGet(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, 400, "Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}
