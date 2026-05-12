package external

import (
	"encoding/json"
	"os"
	"testing"
)

func TestOdakyuClient_FetchAll(t *testing.T) {
	// 手動でトークンとクッキーをセットする（本来は DB や JSON から読み込む）
	token := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ik1ETTBSREJGUmpVM01qUkNRelk0UmpRNU9EbEZPRVV6T0VKR1JURkdRakkwUlVSR05UWkRPUSJ9.eyJodHRwczovL29uZS1vZGFreXUuY29tL3JvbGVzIjpbIm9uZS11c2VyIl0sImlzcyI6Imh0dHBzOi8vYXV0aC5vbmUtb2Rha3l1LmNvbS8iLCJzdWIiOiJhdXRoMHw2YTAzMzQ1ZGE1ZDhiNmMxNzU5MzVjNWIiLCJhdWQiOlsiaHR0cHM6Ly9vbmUtb2Rha3l1LmF1dGgwLmNvbS9hcGkvdjIvIiwiaHR0cHM6Ly9vbmUtb2Rha3l1LmF1dGgwLmNvbS91c2VyaW5mbyJdLCJpYXQiOjE3Nzg1OTU3MjMsImV4cCI6MTc3ODY4MjEyMywic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCByZWFkOmN1cnJlbnRfdXNlciB1cGRhdGU6Y3VycmVudF91c2VyX21ldGFkYXRhIiwiYXpwIjoidklvamRZRmFSVHBUOUoxZ21OQW1yZXZMT0hUeVk0aTgifQ.I5xeQA8Vmxgvfjg7LS08LxWkBHETA5Q82uojJe-GjdKDlHoHEZzI5si5YU_uzO0sWSoN4REcWcxiip1PIkiv4CpPsKOdPALwQN_TwRhilMMelDQeC2qLcv0w1F6ZeycTKcpnEkoE2dgDA1VggNY0ypT3ytY7KaqibOI2LdzebCRVcvIK4vgsJgITOxkdrlKveE-A4cTw3WHZy5KTTwKTa7_J28HwTVdxHvLRcjAduSnXpfZU6Gie7C8yF5dhjyacyFZA0t72ps0Bz2DUi1Ts5dNZoP8INtLW2Q9wQw0NBcTEElvzkx8eo-DhtZzy5LbYXnmguuyKOsuqHcrSzOJvTg"
	cookies := map[string]string{
		"_yjsu_yjad": "1778594750.604240ca-a771-4cb3-856d-66494b58a649",
		"__lt__cid":  "046f3b12-249c-465a-b7ef-9334df14758a",
		"__lt__sid":  "204ff49a-4e5b56eb",
		"_gid":       "GA1.2.473925414.1778594751",
		"withdesk-id": "db0d96c3-5e3d-41a5-a546-e1041f61a39e",
		"_clck":      "1yeuvy0%5E2%5Eg5z%5E0%5E2323",
		"_gcl_au":    "1.1.1502067103.1778594750.1260494806.1778594768.1778594909",
		"_ga":        "GA1.2.1563838141.1778594751",
	}

	session := OdakyuSession{
		Token:   token,
		Cookies: cookies,
	}
	sessionData, _ := json.MarshalIndent(session, "", "  ")
	os.WriteFile("odakyu_cookie.json", sessionData, 0644)
	defer os.Remove("odakyu_cookie.json")

	client, _ := NewOdakyuClient()
	res, err := client.FetchAll()
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	if res.ReturnStatus != "000" {
		t.Fatalf("API returned failure: %+v", res)
	}

	t.Logf("Fetch success! This Year: %d, Last Year: %d", res.ThisYearBalance, res.LastYearBalance)
}
