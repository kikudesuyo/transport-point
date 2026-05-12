package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kikudesuyo/point-hub/service"
)

type OdakyuHandler struct {
	pointService *service.PointService
}

func NewOdakyuHandler(ps *service.PointService) *OdakyuHandler {
	return &OdakyuHandler{pointService: ps}
}

func (h *OdakyuHandler) GetPoints(w http.ResponseWriter, r *http.Request) {
	points, err := h.pointService.FetchOdakyu()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(points)
}
