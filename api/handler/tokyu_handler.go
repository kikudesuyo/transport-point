package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kikudesuyo/point-hub/service"
)

type TokyuHandler struct {
	ps *service.PointService
}

func NewTokyuHandler(ps *service.PointService) *TokyuHandler {
	return &TokyuHandler{ps: ps}
}

func (h *TokyuHandler) GetPoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	report, err := h.ps.FetchTokyu()
	if err != nil {
		log.Printf("Tokyu fetch error: %v", err)
		http.Error(w, `{"error": "Failed to fetch Tokyu points"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(report)
}
