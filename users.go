package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aduatgit/chirpy/internal/auth"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"-"`
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

	password, err := auth.HashPassword(params.Password)
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
		EMail: user.Email,
	})

}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.LookupUserByMail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	defaultExpiration := 60 * 60 * 24
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = defaultExpiration
	} else if params.ExpiresInSeconds > defaultExpiration {
		params.ExpiresInSeconds = defaultExpiration
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(params.ExpiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
		Token: token,
	})
}
