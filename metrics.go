package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if _, err := fmt.Fprintf(w, `<html>
	<body>
    	<h1>Welcome, Chirpy Admin</h1>
    	<p>Chirpy has been visited %d times!</p>
  	</body>
	</html>`, cfg.fileserverHits.Load()); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}

	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(response))
}

func (cfg *apiConfig) resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hit counter reset to 0"))
}
