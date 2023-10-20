package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	EMail    string `json:"email"`
	Password []byte `json:"password"`
}

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Password string `json:"password"`
		EMail    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := requestBody{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "couldn't decode parameters")
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(params.Password), 10)
	if err != nil {
		respondWithError(w, 400, "Couldn't generate password")
		return
	}
	email := params.EMail

	user, err := cfg.DB.CreateUser(email, password)
	if err != nil {
		respondWithError(w, 400, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:       user.ID,
		EMail:    user.EMail,
		Password: user.Password,
	})

}
