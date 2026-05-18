package route

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kikudesuyo/point-hub/app/handler"
)

func RunHTTPServer(w http.ResponseWriter, r *http.Request) {
	mux := NewMux()
	mux.ServeHTTP(w, r)
}

func NewMux() http.Handler {
	// Cloud Functions does not need .env normally,
	// but we keep it for local and fallback compatibility.
	_ = godotenv.Load()
	if _, err := os.Stat("../../.env"); err == nil {
		_ = godotenv.Load("../../.env")
	} else if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load(".env")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/tokyu", handler.GetTokyuPoints)
	mux.HandleFunc("/api/v1/metpo", handler.GetMetpoPoints)
	mux.HandleFunc("/api/v1/toei", handler.GetToeiPoints)
	mux.HandleFunc("/api/v1/sotetsu", handler.GetSotetsuPoints)
	mux.HandleFunc("/api/v1/keikyu", handler.GetKeikyuPoints)
	mux.HandleFunc("/api/v1/odakyu", handler.GetOdakyuPoints)
	mux.HandleFunc("/api/v1/tobu", handler.GetTobuPoints)
	return mux
}
