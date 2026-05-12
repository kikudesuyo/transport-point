package external

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"encoding/json"

	"github.com/PuerkitoBio/goquery"
)

const (
	tobuMyPageURL   = "https://history.tobupoint.jp/mypage/PISM020_00"
	tobuUserAgent   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36"
	tobuCookieFile  = "tobu_cookie.json"
)

type TobuClient struct {
	client *http.Client
}

type TobuData struct {
	TotalPoint    int
	NormalPoint   int
	NormalExpiry  string
	LimitedPoint  int
	LimitedExpiry string
	Miles         int
}

func NewTobuClient() (*TobuClient, error) {
	jar, _ := cookiejar.New(nil)
	k := &TobuClient{
		client: &http.Client{Jar: jar},
	}
	k.loadCookies()
	return k, nil
}

func (k *TobuClient) loadCookies() {
	data, err := os.ReadFile(tobuCookieFile)
	if err == nil {
		var cookies map[string]string
		if err := json.Unmarshal(data, &cookies); err == nil {
			k.SetCookies(cookies)
		}
	}
}

func (k *TobuClient) saveCookies() {
	cookies := k.GetCookies()
	data, _ := json.MarshalIndent(cookies, "", "  ")
	os.WriteFile(tobuCookieFile, data, 0644)
}

func (k *TobuClient) SetCookies(cookies map[string]string) {
	u, _ := url.Parse("https://history.tobupoint.jp")
	var httpCookies []*http.Cookie
	for name, value := range cookies {
		httpCookies = append(httpCookies, &http.Cookie{
			Name:   name,
			Value:  value,
			Path:   "/",
			Domain: u.Host,
		})
	}
	k.client.Jar.SetCookies(u, httpCookies)
}

func (k *TobuClient) GetCookies() map[string]string {
	u, _ := url.Parse("https://history.tobupoint.jp")
	cookies := make(map[string]string)
	for _, c := range k.client.Jar.Cookies(u) {
		cookies[c.Name] = c.Value
	}
	return cookies
}

func (k *TobuClient) FetchAll() (*TobuData, error) {
	req, _ := http.NewRequest("GET", tobuMyPageURL, nil)
	req.Header.Set("User-Agent", tobuUserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Referer", "https://auth.tobupoint.jp/")

	resp, err := k.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tobu status error: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	data := &TobuData{}
	
	// 利用可能ポイント
	totalStr := strings.ReplaceAll(doc.Find("span.status_total_point").First().Text(), ",", "")
	fmt.Sscanf(totalStr, "%d", &data.TotalPoint)

	// 通常ポイント
	normalStr := strings.ReplaceAll(doc.Find(".point_normal span.point_detail_value").First().Text(), ",", "")
	normalStr = strings.TrimSpace(strings.ReplaceAll(normalStr, "ポイント", ""))
	fmt.Sscanf(normalStr, "%d", &data.NormalPoint)
	data.NormalExpiry = strings.TrimSpace(doc.Find(".point_normal .point_detail_note span.point_detail_value").First().Text())

	// 期間限定ポイント
	limitedStr := strings.ReplaceAll(doc.Find(".point_limited .point_detail span.point_detail_value").First().Text(), ",", "")
	limitedStr = strings.TrimSpace(strings.ReplaceAll(limitedStr, "ポイント", ""))
	fmt.Sscanf(limitedStr, "%d", &data.LimitedPoint)
	data.LimitedExpiry = strings.TrimSpace(doc.Find(".point_limited p:last-child span.point_detail_value").First().Text())

	// トブポマイル
	mileStr := strings.ReplaceAll(doc.Find("span.status_mile").First().Text(), ",", "")
	fmt.Sscanf(mileStr, "%d", &data.Miles)

	return data, nil
}
