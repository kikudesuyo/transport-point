package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kikudesuyo/point-hub/service"
)

type SotetsuHandler struct {
	ps *service.PointService
}

func NewSotetsuHandler(ps *service.PointService) *SotetsuHandler {
	return &SotetsuHandler{ps: ps}
}

func (h *SotetsuHandler) GetPoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	report, err := h.ps.FetchSotetsu()
	if err != nil {
		log.Printf("Sotetsu fetch error: %v", err)
		http.Error(w, `{"error": "Failed to fetch Sotetsu points"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(report)
}
