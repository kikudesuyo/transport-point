package v1

import (
	"github.com/go-chi/chi/v5"

	"github.com/kikudesuyo/point-hub/handler"
	"github.com/kikudesuyo/point-hub/service"
)

func RegisterRoutes(r chi.Router, ps *service.PointService) {
	tokyuHandler := handler.NewTokyuHandler(ps)
	metpoHandler := handler.NewMetpoHandler(ps)
	toeiHandler := handler.NewToeiHandler(ps)
	sotetsuHandler := handler.NewSotetsuHandler(ps)
	keikyuHandler := handler.NewKeikyuHandler(ps)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/tokyu", tokyuHandler.GetPoints)
		r.Get("/metpo", metpoHandler.GetPoints)
		r.Get("/toei", toeiHandler.GetPoints)
		r.Get("/sotetsu", sotetsuHandler.GetPoints)
		r.Get("/keikyu", keikyuHandler.GetPoints)
	})
}
