package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := struct {
		Email string `json:"email"`
	}{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	log.Printf("Attempting to create user with email: %s", params.Email)
	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Database error: %v", err) // Add this line
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	person := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
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
