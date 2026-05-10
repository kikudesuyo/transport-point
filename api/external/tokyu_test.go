package external

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestTokyuClient_FetchAll_MinimalCookies(t *testing.T) {
	// .env から認証情報を読み込み
	_ = godotenv.Load("../.env")
	token := os.Getenv("TOKYU_SESSION_TOKEN")
	if token == "" {
		t.Skip("TOKYU_SESSION_TOKEN が設定されていないためスキップします")
	}

	client, _ := NewTokyuClient()
	
	// 最小限のセッショントークンのみをセット
	cookies := map[string]string{
		"__Host-plus.sessionToken": token,
		"nToken":                   token,
		"s.sessionToken":           token,
		"onToken":                  token,
	}
	client.SetCookies(cookies)

	data, err := client.FetchAll()
	if err != nil {
		t.Logf("Minimal cookies fetch failed: %v (expected if tracking cookies are mandatory)", err)
	} else {
		t.Logf("Minimal cookies fetch success! Points: %d", data.Point)
	}
}

func TestTokyuClient_FetchAll_FullCookies(t *testing.T) {
	_ = godotenv.Load("../.env")
	token := os.Getenv("TOKYU_SESSION_TOKEN")
	if token == "" {
		t.Skip("TOKYU_SESSION_TOKEN が設定されていないためスキップします")
	}

	client, _ := NewTokyuClient()
	
	// すべての環境変数からクッキーを構成（後でサービス層のロジックに合わせる）
	cookies := map[string]string{
		"__Host-plus.sessionToken": token,
		"nToken":                   token,
		"_clck":                    os.Getenv("TOKYU_CLCK"),
		"_clsk":                    os.Getenv("TOKYU_CLSK"),
		"_ga":                      os.Getenv("TOKYU_GA"),
	}
	client.SetCookies(cookies)

	data, err := client.FetchAll()
	if err != nil {
		t.Errorf("Full cookies fetch failed: %v", err)
	} else {
		t.Logf("Full cookies fetch success! Points: %d", data.Point)
	}
}
