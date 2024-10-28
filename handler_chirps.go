package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/SumDeusVitae/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) chirpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := struct {
		Body string    `json:"body"`
		ID   uuid.UUID `json:"user_id"`
	}{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	corrected := profanityCheck(params.Body)
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   corrected,
		UserID: params.ID,
	})
	if err != nil {
		log.Printf("Database error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't save chirp to database")
		return
	}

	resp_chirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		User_id:   chirp.UserID,
	}
	respondWithJSON(w, 201, resp_chirp)
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("Database error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps from database")
		return
	}
	var array_chirps []Chirp
	for _, chirp := range chirps {
		curr_chirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			User_id:   chirp.UserID,
		}
		array_chirps = append(array_chirps, curr_chirp)
	}
	respondWithJSON(w, 200, array_chirps)
}

func (cfg *apiConfig) getChirpByIdHandler(w http.ResponseWriter, r *http.Request) {
	chirp_id := r.PathValue("chirp_id")
	parsedUUID, err := uuid.Parse(chirp_id)
	if err != nil {
		respondWithError(w, 404, "Error parsing UUID")
		return
	}
	chirp, err := cfg.db.GetChirpById(r.Context(), parsedUUID)
	if err != nil {
		log.Printf("Database error: %v", err)
		respondWithError(w, 404, "Couldn't get chirp by ID from database")
		return
	}
	resp_chirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		User_id:   chirp.UserID,
	}
	respondWithJSON(w, 200, resp_chirp)
}
