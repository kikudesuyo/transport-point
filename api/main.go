package main

import (
	"fmt"
	"log"
	"os"

	"hoge/external"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".envの読み込みに失敗:", err)
	}

	email := os.Getenv("SOTETSU_EMAIL")
	password := os.Getenv("SOTETSU_PASSWORD")
	if email == "" || password == "" {
		log.Fatal("SOTETSU_EMAIL / SOTETSU_PASSWORD が未設定です")
	}

	client, err := external.NewSotetsuClient()
	if err != nil {
		log.Fatal("クライアント初期化失敗:", err)
	}

	fmt.Println("ログイン中...")
	if err := client.Login(email, password); err != nil {
		log.Fatal("ログイン失敗:", err)
	}
	fmt.Println("✅ ログイン成功")

	fmt.Println("データ取得中...")
	data, err := client.FetchAll()
	if err != nil {
		log.Fatal("データ取得失敗:", err)
	}

	fmt.Printf("\n── 会員情報 ──────────────────────\n")
	fmt.Printf("名前          : %s\n", data.Name)
	fmt.Printf("ランク        : %s\n", data.Rank)
	fmt.Printf("\n── ポイント ──────────────────────\n")
	fmt.Printf("保有ポイント  : %d pt\n", data.Point)
	fmt.Printf("有効期限      : %s\n", data.PointExpiry)
	fmt.Printf("\n── マイル ────────────────────────\n")
	fmt.Printf("保有マイル    : %d mile\n", data.Mile)
	fmt.Printf("有効期限      : %s\n", data.MileExpiry)
}
