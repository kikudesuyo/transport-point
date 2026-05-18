package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kikudesuyo/point-hub/app/service"
)

func GetMetpoPoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ps, err := service.NewPointService()
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch points"}`, http.StatusInternalServerError)
		return
	}
	report, err := ps.FetchMetpo()
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch Metpo points"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(report)
}
