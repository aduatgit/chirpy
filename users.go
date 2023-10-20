package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID    int    `json:"id"`
	EMail string `json:"email"`
}

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		EMail string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := requestBody{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode parameters")
		return
	}

	email := params.EMail

	user, err := cfg.DB.CreateUser(email)
	if err != nil {
		respondWithError(w, 400, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:    user.ID,
		EMail: user.EMail,
	})

}
