package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SumDeusVitae/chirpy/internal/auth"
	"github.com/SumDeusVitae/chirpy/internal/database"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
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
	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Couldn't hash password: %v", err)
	}
	log.Printf("Attempting to create user with email: %s", params.Email)
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashed_password,
	})
	if err != nil {
		log.Printf("Database error: %v", err) // Add this line
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	person := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(w, 201, person)

}

func (cfg *apiConfig) resetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbiden")
	}
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete users")
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Everything reset"))
}

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
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
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token in header")
		return
	}
	user_id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Wrong token provided")
		return
	}
	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}
	updated_user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		HashedPassword: hashed_password,
		Email:          params.Email,
		ID:             user_id,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user in db")
		return
	}
	person := User{
		ID:          updated_user.ID,
		CreatedAt:   updated_user.CreatedAt,
		UpdatedAt:   updated_user.UpdatedAt,
		Email:       updated_user.Email,
		IsChirpyRed: updated_user.IsChirpyRed,
	}
	respondWithJSON(w, 200, person)
}
