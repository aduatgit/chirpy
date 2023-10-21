package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/aduatgit/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	issuer, err := auth.GetJWTIssuer(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}
	if issuer != "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "Token not an refresh one")
		return
	}

	revoked, err := cfg.DB.CheckTokenRevokeStatus(token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if revoked {
		respondWithError(w, http.StatusUnauthorized, "Token revoked")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	accToken, err := auth.MakeJWT(userIDInt, cfg.jwtSecret, time.Duration(time.Hour), "chirpy-access")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT")
		return
	}

	resp := response{
		Token: accToken,
	}

	respondWithJSON(w, http.StatusOK, resp)
}
