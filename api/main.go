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

	loginID := os.Getenv("KEIKYU_LOGIN_ID")
	password := os.Getenv("KEIKYU_PASSWORD")
	if loginID == "" || password == "" {
		log.Fatal("KEIKYU_LOGIN_ID / KEIKYU_PASSWORD が未設定です")
	}

	client, err := external.NewKeikyuClient()
	if err != nil {
		log.Fatal("クライアント初期化失敗:", err)
	}

	fmt.Println("ログイン中...")
	if err := client.Login(loginID, password); err != nil {
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
	fmt.Printf("会員No        : %s\n", data.MemberNo)
	fmt.Printf("\n── ポイント ──────────────────────\n")
	fmt.Printf("利用可能ポイント: %d pt\n", data.AvailablePoint)
	fmt.Printf("期間限定ポイント: %d pt\n", data.LimitedPoint)
	fmt.Printf("失効情報        : %s\n", data.RevocationInfo)
}
