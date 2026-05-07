package external

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	keikyuLoginPageURL = "https://kqpoint-portal.keikyu-point.jp/mypage/auth/login_form"
	keikyuLoginPostURL = "https://kqpoint-portal.keikyu-point.jp/mypage/auth/login"
	keikyuMyPageURL    = "https://kqpoint-portal.keikyu-point.jp/mypage/"
	keikyuUserAgent    = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36"
)

type KeikyuClient struct {
	httpClient *http.Client
}

type KeikyuData struct {
	Name           string
	MemberNo       string
	AvailablePoint int
	LimitedPoint   int
	RevocationInfo string
}

func NewKeikyuClient() (*KeikyuClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Jar: jar,
	}
	return &KeikyuClient{httpClient: client}, nil
}

func (k *KeikyuClient) getMGToken() (string, error) {
	req, _ := http.NewRequest("GET", keikyuLoginPageURL, nil)
	req.Header.Set("User-Agent", keikyuUserAgent)

	resp, err := k.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("ログインページ取得失敗: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("HTMLパース失敗: %w", err)
	}

	token, _ := doc.Find("input[name='mg_token']").Attr("value")
	if token == "" {
		return "", fmt.Errorf("mg_tokenが見つかりません")
	}
	return token, nil
}

// Login は標準のhttpリクエストで認証を行う
func (k *KeikyuClient) Login(loginID, password string) error {
	// 1. ログインページにアクセスしてトークンを取得
	token, err := k.getMGToken()
	if err != nil {
		return err
	}

	// 2. POSTリクエストでログイン
	form := url.Values{}
	form.Add("mg_token", token)
	form.Add("ninsyo_id", loginID)
	form.Add("ninsyo_password", password)

	postReq, _ := http.NewRequest("POST", keikyuLoginPostURL, strings.NewReader(form.Encode()))
	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	postReq.Header.Set("User-Agent", keikyuUserAgent)
	postReq.Header.Set("Referer", keikyuLoginPageURL)
	postReq.Header.Set("Origin", "https://kqpoint-portal.keikyu-point.jp")

	postResp, err := k.httpClient.Do(postReq)
	if err != nil {
		return fmt.Errorf("ログインPOST失敗: %w", err)
	}
	defer postResp.Body.Close()

	// ログイン後のURLを確認（auth/login以外なら成功とみなす）
	finalURL := postResp.Request.URL.String()
	if strings.Contains(finalURL, "auth/login") {
		return fmt.Errorf("ログイン失敗: 認証情報を確認してください")
	}

	return nil
}

// FetchAll はマイページからデータを取得する
func (k *KeikyuClient) FetchAll() (*KeikyuData, error) {
	req, _ := http.NewRequest("GET", keikyuMyPageURL, nil)
	req.Header.Set("User-Agent", keikyuUserAgent)

	resp, err := k.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("マイページ取得失敗: %w", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTMLパース失敗: %w", err)
	}

	data := &KeikyuData{}
	// 名前
	data.Name = strings.TrimSpace(doc.Find(".c-information-head-name span").Text())
	// 会員No
	memberNoRaw := doc.Find(".c-information-head-number").Text()
	data.MemberNo = strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(memberNoRaw), "会員No :"))
	// 利用可能ポイント
	availablePointStr := doc.Find(".c-information-body-detail-available strong").Text()
	data.AvailablePoint, _ = strconv.Atoi(strings.ReplaceAll(availablePointStr, ",", ""))
	// 期間限定ポイント
	limitedPointStr := doc.Find(".c-information-body-detail-limited strong").Text()
	data.LimitedPoint, _ = strconv.Atoi(strings.ReplaceAll(limitedPointStr, ",", ""))
	// 失効情報
	data.RevocationInfo = strings.TrimSpace(doc.Find(".c-information-body-detail-revocation").Text())

	return data, nil
}
