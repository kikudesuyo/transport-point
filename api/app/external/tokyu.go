package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const (
	tokyuRootURL    = "https://plus.tokyu.co.jp/"
	tokyuDetailURL  = "https://plus.tokyu.co.jp/my/point/detail"
	tokyuCookieFile = "tokyu_cookie.json"
)

type TokyuClient struct {
	client *http.Client
}

type ExpiryItem struct {
	Balance int
	Date    string
}

type TokyuData struct {
	Point       int
	PointExpiry string
	Expiries    []ExpiryItem
}

// Next.js RSC Payload structs
type TokyuExpirationYear struct {
	Balance    int    `json:"balance"`
	ExpiryDate string `json:"expiryDate"`
}

type TokyuExpirationDate struct {
	CurrentYear     TokyuExpirationYear `json:"currentYear"`
	LastYear        TokyuExpirationYear `json:"lastYear"`
	PrePreviousYear TokyuExpirationYear `json:"prePreviousYear"`
}

type TokyuPointBalances struct {
	Point          int                 `json:"point"`
	ExpirationDate TokyuExpirationDate `json:"expirationDate"`
}

type TokyuRSCPayload struct {
	PointBalances TokyuPointBalances `json:"pointBalances"`
}

func NewTokyuClient() (*TokyuClient, error) {
	jar, _ := cookiejar.New(nil)
	t := &TokyuClient{
		client: &http.Client{Jar: jar},
	}
	t.loadCookies()
	return t, nil
}

func (t *TokyuClient) loadCookies() {
	cookies := make(map[string]string)
	data, err := os.ReadFile(tokyuCookieFile)
	if err == nil {
		json.Unmarshal(data, &cookies)
		t.SetCookies(cookies)
	}
}

func (t *TokyuClient) saveCookies() {
	cookies := t.GetCookies()
	data, _ := json.MarshalIndent(cookies, "", "  ")
	os.WriteFile(tokyuCookieFile, data, 0644)
}

// SetCookies は提供された全ての認証Cookieをセットする
func (t *TokyuClient) SetCookies(cookies map[string]string) {
	uPlus, _ := url.Parse("https://plus.tokyu.co.jp")
	uDot, _ := url.Parse("https://tokyu.co.jp")

	var plusCookies []*http.Cookie
	var dotCookies []*http.Cookie

	for name, value := range cookies {
		if value == "" {
			continue
		}
		c := &http.Cookie{
			Name:  name,
			Value: value,
			Path:  "/",
		}

		// __Host- プレフィックスのクッキーは特定のドメイン属性を持てない
		if strings.HasPrefix(name, "__Host-") {
			plusCookies = append(plusCookies, c)
			continue
		}

		// ドメインの振り分け
		if strings.HasPrefix(name, "_ga") || name == "_clck" || name == "_gcl_au" || name == "_fbp" {
			c.Domain = ".tokyu.co.jp"
			dotCookies = append(dotCookies, c)
		} else {
			c.Domain = "plus.tokyu.co.jp"
			plusCookies = append(plusCookies, c)
		}
	}
	t.client.Jar.SetCookies(uPlus, plusCookies)
	t.client.Jar.SetCookies(uDot, dotCookies)
}

// GetCookies は現在のメモリ上のクッキーをマップ形式で返す
func (t *TokyuClient) GetCookies() map[string]string {
	uPlus, _ := url.Parse("https://plus.tokyu.co.jp")
	uDot, _ := url.Parse("https://tokyu.co.jp")

	cookies := make(map[string]string)

	// 両方のドメインのクッキーを統合
	for _, c := range t.client.Jar.Cookies(uPlus) {
		cookies[c.Name] = c.Value
	}
	for _, c := range t.client.Jar.Cookies(uDot) {
		cookies[c.Name] = c.Value
	}

	return cookies
}

func (t *TokyuClient) setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ja;q=0.8")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("Referer", "https://plus.tokyu.co.jp/")
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
}

func (t *TokyuClient) FetchAll() (*TokyuData, error) {
	data := &TokyuData{}

	req, _ := http.NewRequest("GET", tokyuDetailURL, nil)
	t.setHeaders(req)
	
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// ログイン確認（HTML内を確認）
	if !strings.Contains(string(body), "ログアウト") {
		return nil, fmt.Errorf("認証失敗: セッションが切れているか、Cookieが不足しています")
	}

	// JSON抽出
	payload, err := t.extractJSON(string(body))
	if err != nil {
		return nil, fmt.Errorf("JSON抽出失敗: %v", err)
	}

	data.Point = payload.PointBalances.Point

	// 有効期限の処理
	exp := payload.PointBalances.ExpirationDate

	// 3年分のデータを整理
	items := []struct {
		Balance int
		RawDate string
	}{
		{exp.PrePreviousYear.Balance, exp.PrePreviousYear.ExpiryDate},
		{exp.LastYear.Balance, exp.LastYear.ExpiryDate},
		{exp.CurrentYear.Balance, exp.CurrentYear.ExpiryDate},
	}

	for _, item := range items {
		dateStr := strings.TrimPrefix(item.RawDate, "$D")
		if idx := strings.Index(dateStr, "T"); idx != -1 {
			dateStr = dateStr[:idx]
		}

		data.Expiries = append(data.Expiries, ExpiryItem{
			Balance: item.Balance,
			Date:    dateStr,
		})

		// PointExpiry には残高がある最短の期限をセット（まだセットされていない場合）
		if item.Balance > 0 && data.PointExpiry == "" {
			data.PointExpiry = dateStr
		}
	}

	// すべて0の場合は、便宜上最も古い期限を PointExpiry に入れる
	if data.PointExpiry == "" && len(data.Expiries) > 0 {
		data.PointExpiry = data.Expiries[0].Date
	}

	return data, nil
}

func (t *TokyuClient) extractJSON(html string) (*TokyuRSCPayload, error) {
	// self.__next_f.push([1,"..."]) の中身を抽出
	// pointBalances が含まれるチャンクを探す
	re := regexp.MustCompile(`self\.__next_f\.push\(\[1,"(.*?)"\]\)`)
	matches := re.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		content := match[1]
		if strings.Contains(content, "pointBalances") {
			// エスケープされた文字列をクリーンアップ
			// \" -> " , \\ -> \
			content = strings.ReplaceAll(content, `\"`, `"`)
			content = strings.ReplaceAll(content, `\\`, `\`)

			// JSONの開始位置を探す { "pointBalances": ...
			startIdx := strings.Index(content, `{"currentMonth"`)
			if startIdx == -1 {
				startIdx = strings.Index(content, `{"pointBalances"`)
			}

			if startIdx != -1 {
				jsonStr := content[startIdx:]
				// 閉じ括弧を探す（簡易的だが、RSC形式なら大体最後がオブジェクト）
				if endIdx := strings.LastIndex(jsonStr, "}]"); endIdx != -1 {
					jsonStr = jsonStr[:endIdx+1]
				} else if endIdx := strings.LastIndex(jsonStr, "}"); endIdx != -1 {
					jsonStr = jsonStr[:endIdx+1]
				}

				var payload TokyuRSCPayload
				if err := json.Unmarshal([]byte(jsonStr), &payload); err == nil {
					return &payload, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("pointBalances found in script tags but could not be parsed")
}
