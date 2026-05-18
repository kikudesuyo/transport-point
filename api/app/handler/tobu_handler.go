package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kikudesuyo/point-hub/app/service"
)

func GetTobuPoints(w http.ResponseWriter, r *http.Request) {
	ps, err := service.NewPointService()
	if err != nil {
		http.Error(w, `{"error": "Failed to fetch points"}`, http.StatusInternalServerError)
		return
	}
	points, err := ps.FetchTobu()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(points)
}
