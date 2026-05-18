package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

const (
	odakyuBalanceURL = "https://one-odakyu.com/titan/v1/op/cards/balance?device_no=9999"
	odakyuUserAgent  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36"
	odakyuCookieFile = "odakyu_cookie.json"
)

type OdakyuClient struct {
	client *http.Client
	token  string
}

type OdakyuResponse struct {
	AdmissionDate    string `json:"admission_date"`
	ErrMessage       string `json:"err_message"`
	LastYearBalance  int    `json:"last_year_balance"`
	MemberStatus     string `json:"member_status"`
	PointAccumTerm   string `json:"point_accum_term"`
	PointInvalidDate string `json:"point_invalid_date"`
	PreviousBalance  int    `json:"previous_balance"`
	ReturnStatus     string `json:"return_status"`
	ThisYearBalance  int    `json:"this_year_balance"`
}

type OdakyuSession struct {
	Token   string            `json:"token"`
	Cookies map[string]string `json:"cookies"`
}

func NewOdakyuClient() (*OdakyuClient, error) {
	jar, _ := cookiejar.New(nil)
	k := &OdakyuClient{
		client: &http.Client{Jar: jar},
	}
	k.loadSession()
	return k, nil
}

func (k *OdakyuClient) loadSession() {
	data, err := os.ReadFile(odakyuCookieFile)
	if err == nil {
		var session OdakyuSession
		if err := json.Unmarshal(data, &session); err == nil {
			k.token = session.Token
			k.SetCookies(session.Cookies)
		}
	}
}

func (k *OdakyuClient) saveSession() {
	session := OdakyuSession{
		Token:   k.token,
		Cookies: k.GetCookies(),
	}
	data, _ := json.MarshalIndent(session, "", "  ")
	os.WriteFile(odakyuCookieFile, data, 0644)
}

func (k *OdakyuClient) SetCookies(cookies map[string]string) {
	u, _ := url.Parse("https://one-odakyu.com")
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

func (k *OdakyuClient) GetCookies() map[string]string {
	u, _ := url.Parse("https://one-odakyu.com")
	cookies := make(map[string]string)
	for _, c := range k.client.Jar.Cookies(u) {
		cookies[c.Name] = c.Value
	}
	return cookies
}

func (k *OdakyuClient) FetchAll() (*OdakyuResponse, error) {
	req, _ := http.NewRequest("GET", odakyuBalanceURL, nil)
	req.Header.Set("User-Agent", odakyuUserAgent)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ja;q=0.8")
	req.Header.Set("Referer", "https://one-odakyu.com/odakyu-point?tab=point-history")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Ch-Ua", `"Chromium";v="146", "Not-A.Brand";v="24", "Google Chrome";v="146"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	
	if k.token != "" {
		req.Header.Set("Authorization", "Bearer "+k.token)
	}

	resp, err := k.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("認証失敗: Tokenが期限切れの可能性があります")
	}

	var result OdakyuResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
