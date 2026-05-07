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

	client, err := external.NewTokyuClient()
	if err != nil {
		log.Fatal("クライアント初期化失敗:", err)
	}

	// 全盛りCookie
	cookies := map[string]string{
		"__Host-plus.sessionToken": os.Getenv("TOKYU_SESSION_TOKEN"),
		"nToken":                   os.Getenv("TOKYU_SESSION_TOKEN"),
		"s.sessionToken":           os.Getenv("TOKYU_SESSION_TOKEN"),
		"onToken":                  os.Getenv("TOKYU_SESSION_TOKEN"),
		"_clck":                    os.Getenv("TOKYU_CLCK"),
		"_clsk":                    os.Getenv("TOKYU_CLSK"),
		"_ga":                      os.Getenv("TOKYU_GA"),
		"_ga_B0V3646TYC":           os.Getenv("TOKYU_GA_B0"),
		"_ga_XD2N3Y0135":           os.Getenv("TOKYU_GA_XD"),
		"_ga_Y86R0E9JVH":           os.Getenv("TOKYU_GA_Y8"),
		"_gcl_au":                  os.Getenv("TOKYU_GCL_AU"),
		"_rslgvry":                 os.Getenv("TOKYU_RSLGVRY"),
		"_yjsu_yjad":               os.Getenv("TOKYU_YJSU"),
		"krt_rewrite_uid":          os.Getenv("TOKYU_KRT"),
		"withdesk-id":              os.Getenv("TOKYU_WITHDESK"),
	}

	client.SetCookies(cookies)

	fmt.Println("東急トップページにアクセス中...")
	data, err := client.FetchAll()
	if err != nil {
		log.Fatal("アクセス失敗:", err)
	}

	fmt.Println("✅ ログイン成功！")
	fmt.Printf("取得データ: %+v\n", data)
}
