package main

import (
	"fmt"
	"log"
	"net/http"

	route "github.com/kikudesuyo/point-hub/app/routes/v1"
)

func main() {
	mux := route.NewMux()
	fmt.Println("サーバーを起動します: http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
