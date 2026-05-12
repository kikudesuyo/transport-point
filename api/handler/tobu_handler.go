package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kikudesuyo/point-hub/service"
)

type TobuHandler struct {
	pointService *service.PointService
}

func NewTobuHandler(ps *service.PointService) *TobuHandler {
	return &TobuHandler{pointService: ps}
}

func (h *TobuHandler) GetPoints(w http.ResponseWriter, r *http.Request) {
	points, err := h.pointService.FetchTobu()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(points)
}
