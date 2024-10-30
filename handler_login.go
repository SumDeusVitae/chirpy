package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/SumDeusVitae/chirpy/internal/auth"
)

func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Expires  int64  `json:"expires_in_seconds"`
	}{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("Database error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't receive hash from db")
		return
	}
	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
	} else {
		expiresIn := time.Duration(params.Expires) * time.Second
		if params.Expires < 1 || params.Expires > 216000 {
			expiresIn = time.Duration(216000) * time.Second
		}
		token, err := auth.MakeJWT(user.ID, cfg.secret, expiresIn)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get token")
			return
		}
		person := User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     token,
		}
		respondWithJSON(w, 200, person)
	}
}
