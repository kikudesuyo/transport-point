package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kikudesuyo/point-hub/service"
)

type ToeiHandler struct {
	ps *service.PointService
}

func NewToeiHandler(ps *service.PointService) *ToeiHandler {
	return &ToeiHandler{ps: ps}
}

func (h *ToeiHandler) GetPoints(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	report, err := h.ps.FetchToei()
	if err != nil {
		log.Printf("Toei fetch error: %v", err)
		http.Error(w, `{"error": "Failed to fetch Toei points"}`, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(report)
}
