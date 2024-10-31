package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/SumDeusVitae/chirpy/internal/auth"
	"github.com/SumDeusVitae/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) chirpHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := struct {
		Body string `json:"body"`
	}{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "No token in header")
		return
	}
	user_id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Wrong token provided")
		return
	}
	corrected := profanityCheck(params.Body)
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   corrected,
		UserID: user_id,
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
	order := r.URL.Query().Get("sort")
	order = strings.ToUpper(order)
	var dbChirps []database.Chirp
	var err error
	if order == "DESC" {
		dbChirps, err = cfg.db.GetChirpsDESC(r.Context())
	} else {
		dbChirps, err = cfg.db.GetChirpsASC(r.Context())
	}

	if err != nil {
		log.Printf("Database error: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps from database")
		return
	}

	authorID := uuid.Nil
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invaslid author ID")
			return
		}
	}

	var array_chirps []Chirp
	for _, chirp := range dbChirps {
		if authorID != uuid.Nil && chirp.UserID != authorID {
			continue
		}
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
	chirp_id := r.PathValue("chirpID")
	log.Printf("Retrieved chirpID for get: %v", chirp_id)
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

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
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
	chirp_id := r.PathValue("chirpID")
	log.Printf("Retrieved chirpID in DELETE: %v", chirp_id)
	parsedUUID, err := uuid.Parse(chirp_id)
	if err != nil {
		respondWithError(w, 404, "Error parsing UUID")
		return
	}
	_, err = cfg.db.GetChirpById(r.Context(), parsedUUID)
	if err != nil {
		log.Printf("Database error: %v", err)
		respondWithError(w, 404, "Couldn't get chirp by ID from database")
		return
	}
	result, err := cfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     parsedUUID,
		UserID: user_id,
	})
	if err != nil {
		respondWithError(w, 403, "You are not authorized to delete it.")
		return
	}
	if result == [16]byte{} {
		respondWithError(w, 403, "Didn't find chirp to delete")
		return
	}
	w.WriteHeader(204)
	return

}
