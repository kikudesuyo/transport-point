package external

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
)

const (
	metpoLoginPageURL = "https://www.metro-point-club.jp/mp/contents/wb1004/ws0010101/view"
	metpoLoginPostURL = "https://www.metro-point-club.jp/mp/contents/wb1004/ws0010101/login"
	metpoTopPageURL   = "https://www.metro-point-club.jp/mp/contents/wb1005/ws0030101/view"
)

type MetpoClient struct{ client *http.Client }

type MetpoData struct {
	User  MetpoUserInfo
	Point MetpoPointInfo
	Score MetpoScoreInfo
}

type MetpoUserInfo struct {
	Name string
	ID   string
}

type MetpoPointInfo struct {
	HoldingPoint      int
	NormalPoint       int
	NormalExpiry      string
	NormalExpiryPoint int
	ChargePoint       int
	ChargeExpiry      string
	ChargeExpiryPoint int
}

type MetpoScoreInfo struct {
	CurrentScore    int
	CurrentRank     string
	NextRankDate    string
	NextRankName    string
	ScoreToNextRank int
}

func NewMetpoClient() (*MetpoClient, error) {
	jar, _ := cookiejar.New(nil)
	return &MetpoClient{client: &http.Client{Jar: jar}}, nil
}

// Login はメトポにログインしてセッションCookieを取得する
func (m *MetpoClient) Login(customerNumber, webPassword string) error {
	resp, err := m.client.PostForm(metpoLoginPostURL, url.Values{
		"customerNumber": {customerNumber},
		"webPassword":    {webPassword},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 成功判定（URLに wb1005 が含まれているか）
	if !strings.Contains(resp.Request.URL.String(), "wb1005") {
		return fmt.Errorf("ログイン失敗: 認証情報を確認してください")
	}
	return nil
}

// FetchAll はトップページから全データを取得する
func (m *MetpoClient) FetchAll() (*MetpoData, error) {
	resp, err := m.client.Get(metpoTopPageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Shift-JISをUTF-8に変換
	utf8Body, _ := charset.NewReader(resp.Body, "shift_jis")
	doc, _ := goquery.NewDocumentFromReader(utf8Body)
	data := &MetpoData{}

	// ユーザー情報
	nameNode := doc.Find(".user-name").First().Clone()
	nameNode.Find(".unit").Remove()
	data.User = MetpoUserInfo{
		Name: strings.TrimSpace(nameNode.Text()),
		ID:   strings.TrimSpace(doc.Find(".user-id").First().Text()),
	}

	// ポイント情報
	p := MetpoPointInfo{}
	p.HoldingPoint, _ = strconv.Atoi(strings.TrimSpace(doc.Find(".holding-point").First().Text()))

	var labels []string
	var values []int
	doc.Find("#modal dl.menu-modal-text dt").Each(func(_ int, s *goquery.Selection) {
		labels = append(labels, strings.TrimSpace(s.Text()))
		valStr := strings.NewReplacer("pt", "", " ", "", ",", "").Replace(s.Next().Text())
		val, _ := strconv.Atoi(valStr)
		values = append(values, val)
	})
	if len(values) >= 4 {
		p.NormalPoint, p.NormalExpiry, p.NormalExpiryPoint = values[0], labels[1], values[1]
		p.ChargePoint, p.ChargeExpiry, p.ChargeExpiryPoint = values[2], labels[3], values[3]
	}
	data.Point = p

	// スコア・ランク情報
	s := MetpoScoreInfo{}
	s.CurrentScore, _ = strconv.Atoi(strings.TrimSpace(doc.Find(".score-meter-now-value").Text()))
	rankClass, _ := doc.Find(".icon-point-rank").Attr("class")
	for _, cls := range strings.Fields(rankClass) {
		if cls != "icon-point-rank" {
			s.CurrentRank = cls
			break
		}
	}
	s.NextRankDate, _ = doc.Find(".score-rank-message time").Attr("datetime")
	nextText := strings.TrimSpace(doc.Find(".score-meter-next").Text())
	if idx := strings.Index(nextText, "会員"); idx > 0 {
		s.NextRankName = nextText[:idx]
	}
	s.ScoreToNextRank, _ = strconv.Atoi(strings.TrimSpace(doc.Find(".score-meter-next-value").Text()))
	data.Score = s

	return data, nil
}
