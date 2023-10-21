package main

import (
	"encoding/json"
	"net/http"

	"github.com/aduatgit/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	if apiKey != cfg.polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := request{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusOK)
		return
	}

	err = cfg.DB.UpdateChirpyRedStatus(params.Data.UserID, true)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(200)

}
