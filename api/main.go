package main

import (
	"fmt"
	"log"
	"hoge/service"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".envの読み込みに失敗:", err)
	}

	ps, err := service.NewPointService()
	if err != nil {
		log.Fatal("サービス初期化失敗:", err)
	}

	fmt.Println("各社ポイント収集中...")
	report, err := ps.FetchAll()
	if err != nil {
		log.Fatal("データ取得失敗:", err)
	}

	fmt.Println("\n========================================")
	fmt.Printf("合計ポイント: %d pt\n", report.TotalBalance)
	fmt.Println("========================================")

	for _, d := range report.Details {
		fmt.Printf("- %-20s: %6d pt", d.Provider, d.Balance)
		if d.ExpiryDate != "" {
			fmt.Printf(" (最短失効: %s)", d.ExpiryDate)
		}
		fmt.Println()
		
		for _, exp := range d.ExpiryList {
			fmt.Printf("    └ %s: %d pt\n", exp.Date, exp.Points)
		}
	}
	fmt.Println("========================================\n")

	// セッション情報の自動更新
	if len(report.UpdatedCookies) > 0 {
		fmt.Println("セッション更新を検知しました。 .env を更新します...")
		if err := updateEnvFile(".env", report.UpdatedCookies); err != nil {
			fmt.Printf("⚠️ .envの更新に失敗: %v\n", err)
		} else {
			fmt.Println("✅ .envを最新のセッション情報で更新しました")
		}
	}
}

// updateEnvFile は指定されたファイルの環境変数を書き換える
func updateEnvFile(filename string, updates map[string]string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	changed := false

	for i, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		if newVal, ok := updates[key]; ok {
			lines[i] = fmt.Sprintf("%s=%s", key, newVal)
			delete(updates, key) // 処理済み
			changed = true
		}
	}

	// もしファイルになかった新規キーがあれば追記
	for key, val := range updates {
		lines = append(lines, fmt.Sprintf("%s=%s", key, val))
		changed = true
	}

	if changed {
		return os.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
	}
	return nil
}
