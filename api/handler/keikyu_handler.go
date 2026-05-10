package handler

import (
	"encoding/json"
	"github.com/kikudesuyo/point-hub/service"
	"log"
	"net/http"
)

type KeikyuHandler struct {
	ps *service.PointService
}

func NewKeikyuHandler(ps *service.PointService) *KeikyuHandler {
	return &KeikyuHandler{ps: ps}
}

func (h *KeikyuHandler) GetPoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	report, err := h.ps.FetchKeikyu()
	if err != nil {
		log.Printf("Keikyu fetch error: %v", err)
		http.Error(w, `{"error": "Failed to fetch Keikyu points"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(report)
}
