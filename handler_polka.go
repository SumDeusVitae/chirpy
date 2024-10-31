package main

import (
	"encoding/json"
	"net/http"

	"github.com/SumDeusVitae/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) polkaHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "NO API KEY IN HEADER")
		return
	}
	if cfg.polka_api != apiKey {
		respondWithError(w, http.StatusUnauthorized, "Wrong API KEY")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}{}

	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}
	parsedUUID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, 404, "Error parsing UUID")
		return
	}
	_, err = cfg.db.UpgradeUser(r.Context(), parsedUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Wrong user ID")
	}

	w.WriteHeader(204)
	return

}
