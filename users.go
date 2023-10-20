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
	type responseBody struct {
		ID    int    `json:"id"`
		EMail string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := requestBody{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(params.Password), 4)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}
	email := params.EMail

	user, err := cfg.DB.CreateUser(email, password)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, responseBody{
		ID:    user.ID,
		EMail: user.EMail,
	})

}

func (cfg *apiConfig) handlerUserLogin(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Password string `json:"password"`
		EMail    string `json:"email"`
	}
	type responseBody struct {
		ID    int    `json:"id"`
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
	password := params.Password
	user, err := cfg.DB.LookupUserByMail(email)
	if err != nil {
		respondWithError(w, 401, err.Error())
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		respondWithError(w, 401, err.Error())
	}
	respBody := responseBody{
		ID:    user.ID,
		EMail: user.EMail,
	}

	respondWithJSON(w, http.StatusOK, respBody)

}
