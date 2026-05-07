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

	id := os.Getenv("TOKYO_METRO_EMAIL")
	password := os.Getenv("TOKYO_METRO_PASSWORD")
	if id == "" || password == "" {
		log.Fatal("TOKYO_METRO_EMAIL / TOKYO_METRO_PASSWORD が未設定です")
	}

	client, err := external.NewMetpoClient()
	if err != nil {
		log.Fatal("クライアント初期化失敗:", err)
	}

	fmt.Println("ログイン中...")
	if err := client.Login(id, password); err != nil {
		log.Fatal("ログイン失敗:", err)
	}
	fmt.Println("✅ ログイン成功")

	fmt.Println("データ取得中...")
	data, err := client.FetchAll()
	if err != nil {
		log.Fatal("データ取得失敗:", err)
	}

	fmt.Printf("\n── 会員情報 ──────────────────────\n")
	fmt.Printf("名前          : %s\n", data.User.Name)
	fmt.Printf("会員番号      : %s\n", data.User.ID)
	fmt.Printf("\n── ポイント ──────────────────────\n")
	fmt.Printf("合計保有      : %d pt\n", data.Point.HoldingPoint)
	fmt.Printf("通常ポイント  : %d pt (%s: %d pt)\n", data.Point.NormalPoint, data.Point.NormalExpiry, data.Point.NormalExpiryPoint)
	fmt.Printf("チャージ専用  : %d pt (%s: %d pt)\n", data.Point.ChargePoint, data.Point.ChargeExpiry, data.Point.ChargeExpiryPoint)
	fmt.Printf("\n── ランク ────────────────────────\n")
	fmt.Printf("現在ランク    : %s (スコア: %d)\n", data.Score.CurrentRank, data.Score.CurrentScore)
	fmt.Printf("次回更新日    : %s\n", data.Score.NextRankDate)
	fmt.Printf("次ランク目安  : %sまで あと %d スコア\n", data.Score.NextRankName, data.Score.ScoreToNextRank)
}
