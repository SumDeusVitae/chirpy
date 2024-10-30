package main

import (
	"net/http"
	"time"

	"github.com/SumDeusVitae/chirpy/internal/auth"
)

func (cfg *apiConfig) checkRefreshHandler(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token in header")
		return
	}
	tokenInDB, err := cfg.db.CheckRefreshToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "NO SUCH TOKEN")
		return
	}
	// Check if the token has been revoked
	if tokenInDB.RevokedAt.Valid && !tokenInDB.RevokedAt.Time.IsZero() {
		respondWithError(w, http.StatusUnauthorized, "TOKEN REVOKED")
		return
	}

	// Check if the token has expired
	if !tokenInDB.ExpiresAt.Valid || tokenInDB.ExpiresAt.Time.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "TOKEN EXPIRED")
		return
	}

	type Response struct {
		Token string `json:"token"`
	}

	expiresIn := time.Duration(3600) * time.Second
	token, err := auth.MakeJWT(tokenInDB.UserID, cfg.secret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token")
		return
	}
	response := Response{Token: token}

	respondWithJSON(w, 200, response)
}

func (cfg *apiConfig) revokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token in header")
		return
	}
	err = cfg.db.RevokeToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not located in db")
	}
	w.WriteHeader(http.StatusNoContent)

}
