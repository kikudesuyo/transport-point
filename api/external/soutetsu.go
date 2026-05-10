package external

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	sotetsuLoginPageURL = "https://mypage.sotetsu-point.jp/PISM010_00"
	sotetsuLoginPostURL = "https://mypage.sotetsu-point.jp/PISM010_01"
	sotetsuMyPageURL    = "https://mypage.sotetsu-point.jp/PISM020_00"
)

type SotetsuClient struct{ client *http.Client }
type SotetsuData struct {
	Name        string
	Point       int
	Mile        int
	Rank        string
	PointExpiry string
	MileExpiry  string
}

func NewSotetsuClient() (*SotetsuClient, error) {
	jar, _ := cookiejar.New(nil)
	return &SotetsuClient{client: &http.Client{Jar: jar}}, nil
}

// Login は標準のhttpリクエストで認証を行う
func (s *SotetsuClient) Login(userId, password string) error {
	// 1. ログインページからトークン取得
	resp, err := s.client.Get(sotetsuLoginPageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	token, _ := doc.Find("input[name='jp.hitachisoft.message.TOKEN']").Attr("value")

	// 2. ログインPOST
	postResp, err := s.client.PostForm(sotetsuLoginPostURL, url.Values{
		"userId":                       {userId},
		"passWord":                     {password},
		"jp.hitachisoft.message.TOKEN": {token},
	})
	if err != nil {
		return err
	}
	defer postResp.Body.Close()

	// 成功判定（マイページ PISM020 にリダイレクトされたか）
	if !strings.Contains(postResp.Request.URL.String(), "PISM020") {
		return fmt.Errorf("ログイン失敗")
	}
	return nil
}

// FetchAll はマイページからデータを取得する
func (s *SotetsuClient) FetchAll() (*SotetsuData, error) {
	resp, err := s.client.Get(sotetsuMyPageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	data := &SotetsuData{}

	// 名前
	data.Name = strings.TrimSpace(strings.TrimSuffix(doc.Find("h1.parts-title03").Text(), " 様"))

	// 各ステータス（ポイント、マイル、ランク）の抽出
	doc.Find(".mypage-status__whbord").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find("h2").First().Text())
		valText := strings.TrimSpace(s.Find(".mypage-status__situation__txt").Text())
		expiry := strings.TrimSpace(s.Find("span:contains('有効期限')").NextAll().Text())

		switch {
		case strings.Contains(title, "ポイント"):
			v, _ := strconv.Atoi(strings.ReplaceAll(valText, ",", ""))
			data.Point = v
			data.PointExpiry = expiry
		case strings.Contains(title, "マイル"):
			v, _ := strconv.Atoi(strings.ReplaceAll(valText, ",", ""))
			data.Mile = v
			data.MileExpiry = expiry
		case strings.Contains(title, "ランク"):
			data.Rank = valText
		}
	})

	return data, nil
}
