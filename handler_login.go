package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/SumDeusVitae/chirpy/internal/auth"
	"github.com/SumDeusVitae/chirpy/internal/database"
)

func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
	} else {
		expiresIn := time.Duration(3600) * time.Second
		token, err := auth.MakeJWT(user.ID, cfg.secret, expiresIn)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get token")
			return
		}
		refresh_token, err := auth.MakeRefreshToken()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token")
			return
		}
		refresh_token_info, err := cfg.db.StoreRefreshToken(r.Context(), database.StoreRefreshTokenParams{
			Token:  refresh_token,
			UserID: user.ID,
			ExpiresAt: sql.NullTime{
				Time:  time.Now().Add(60 * 24 * time.Hour),
				Valid: true,
			},
			RevokedAt: sql.NullTime{},
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't Store refresh token")
			return
		}
		person := User{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			Token:        token,
			RefreshToken: refresh_token_info.Token,
		}
		respondWithJSON(w, 200, person)
	}
}
