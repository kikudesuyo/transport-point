package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	v1 "github.com/kikudesuyo/point-hub/routes/v1"
	"github.com/kikudesuyo/point-hub/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .envの読み込みに失敗しましたが、既存の環境変数で続行します:", err)
	}

	ps, err := service.NewPointService()
	if err != nil {
		log.Fatal("サービス初期化失敗:", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		v1.RegisterRoutes(r, ps)
	})

	port := "8081"
	log.Printf("🚀 サーバー起動: http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("サーバー起動失敗:", err)
	}
}
