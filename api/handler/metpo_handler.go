package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kikudesuyo/point-hub/service"
)

type MetpoHandler struct {
	ps *service.PointService
}

func NewMetpoHandler(ps *service.PointService) *MetpoHandler {
	return &MetpoHandler{ps: ps}
}

func (h *MetpoHandler) GetPoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	report, err := h.ps.FetchMetpo()
	if err != nil {
		log.Printf("Metpo fetch error: %v", err)
		http.Error(w, `{"error": "Failed to fetch Metpo points"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(report)
}
