package external

import (
	"fmt"
	"net/http/cookiejar"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

const (
	loginPageURL = "https://www.metro-point-club.jp/mp/contents/wb1004/ws0010101/view"
	loginPostURL = "https://www.metro-point-club.jp/mp/contents/wb1004/ws0010101/login"
	topPageURL   = "https://www.metro-point-club.jp/mp/contents/wb1005/ws0030101/view"
	userAgent    = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
)

type MetpoClient struct {
	jar *cookiejar.Jar
}

type UserInfo struct {
	Name string
	ID   string
}

type PointInfo struct {
	HoldingPoint      int
	NormalPoint       int
	NormalExpiry      string // 例: "2027年03月失効ポイント"
	NormalExpiryPoint int
	ChargePoint       int
	ChargeExpiry      string
	ChargeExpiryPoint int
}

type ScoreInfo struct {
	CurrentScore    int
	CurrentRank     string
	NextRankDate    string
	NextRankName    string
	ScoreToNextRank int
}

type MetpoData struct {
	User  UserInfo
	Point PointInfo
	Score ScoreInfo
}

func NewMetpoClient() (*MetpoClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &MetpoClient{jar: jar}, nil
}

func (m *MetpoClient) newCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowURLRevisit(),
		colly.UserAgent(userAgent),
	)
	c.SetCookieJar(m.jar)
	return c
}

// Login はメトポにログインしてセッションCookieを取得する
func (m *MetpoClient) Login(email, password string) error {
	c := m.newCollector()
	var loginErr error
	loggedIn := false

	c.OnHTML("form#loginModel", func(e *colly.HTMLElement) {
		if err := e.Request.Post(loginPostURL, map[string]string{
			"customerNumber": email,
			"webPassword":    password,
		}); err != nil {
			loginErr = fmt.Errorf("POSTリクエスト失敗: %w", err)
		}
	})

	c.OnResponse(func(r *colly.Response) {
		if strings.Contains(r.Request.URL.String(), "wb1005") {
			loggedIn = true
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		loginErr = err
	})

	if err := c.Visit(loginPageURL); err != nil {
		return fmt.Errorf("ログインページアクセス失敗: %w", err)
	}
	if loginErr != nil {
		return loginErr
	}
	if !loggedIn {
		return fmt.Errorf("ログイン失敗: 認証情報を確認してください")
	}
	return nil
}

// FetchAll はトップページから全データを取得する
func (m *MetpoClient) FetchAll() (*MetpoData, error) {
	c := m.newCollector()
	data := &MetpoData{}
	var fetchErr error

	c.OnHTML("body", func(e *colly.HTMLElement) {
		data.User = parseUserInfo(e)
		data.Point = parsePointInfo(e)
		data.Score = parseScoreInfo(e)
	})

	c.OnError(func(r *colly.Response, err error) {
		fetchErr = err
	})

	if err := c.Visit(topPageURL); err != nil {
		return nil, fmt.Errorf("トップページアクセス失敗: %w", err)
	}
	if fetchErr != nil {
		return nil, fetchErr
	}
	return data, nil
}

// parseUserInfo はユーザー名とIDを取得する
func parseUserInfo(e *colly.HTMLElement) UserInfo {
	// .user-name 内の "様" スパンを除いてテキストを取得
	name := e.DOM.Find(".user-name").First().Clone().
		Find(".unit").Remove().End().Text()
	id := strings.TrimSpace(e.DOM.Find(".user-id").First().Text())
	return UserInfo{
		Name: strings.TrimSpace(name),
		ID:   id,
	}
}

// parsePointInfo はポイント詳細モーダルから全ポイント情報を取得する
func parsePointInfo(e *colly.HTMLElement) PointInfo {
	info := PointInfo{}

	// 保有ポイント合計
	holdingStr := strings.TrimSpace(e.DOM.Find(".holding-point").First().Text())
	info.HoldingPoint, _ = strconv.Atoi(holdingStr)

	// モーダル内 dt/dd ペアを順番に取得
	var labels []string
	var values []int

	e.DOM.Find("#modal dl.menu-modal-text dt").Each(func(_ int, s *goquery.Selection) {
		labels = append(labels, strings.TrimSpace(s.Text()))
		raw := strings.TrimSpace(s.Next().Text()) // 隣の dd
		raw = strings.ReplaceAll(raw, "pt", "")
		raw = strings.ReplaceAll(raw, " ", "")
		val, _ := strconv.Atoi(strings.TrimSpace(raw))
		values = append(values, val)
	})

	// 順序: [0]通常ポイント [1]通常失効 [2]チャージ専用 [3]チャージ失効
	if len(labels) >= 4 {
		info.NormalPoint = values[0]
		info.NormalExpiry = labels[1]
		info.NormalExpiryPoint = values[1]
		info.ChargePoint = values[2]
		info.ChargeExpiry = labels[3]
		info.ChargeExpiryPoint = values[3]
	}

	return info
}

// parseScoreInfo はランク・スコア情報を取得する
func parseScoreInfo(e *colly.HTMLElement) ScoreInfo {
	info := ScoreInfo{}

	// 現在スコア
	scoreStr := strings.TrimSpace(e.DOM.Find(".score-meter-now-value").Text())
	info.CurrentScore, _ = strconv.Atoi(scoreStr)

	// 現在ランク: icon-point-rank の追加クラスから取得 (例: "icon-point-rank silver")
	rankClass, _ := e.DOM.Find(".icon-point-rank").Attr("class")
	for _, cls := range strings.Fields(rankClass) {
		if cls != "icon-point-rank" {
			info.CurrentRank = cls // "silver", "gold", "platinum", "regular"
			break
		}
	}

	// 次回ランク更新日
	info.NextRankDate, _ = e.DOM.Find(".score-rank-message time").Attr("datetime")

	// 次のランクまでのスコアとランク名
	// 例: "シルバー会員まであと600スコア"
	nextText := strings.TrimSpace(e.DOM.Find(".score-meter-next").Text())
	if idx := strings.Index(nextText, "会員"); idx > 0 {
		info.NextRankName = nextText[:idx]
	}
	nextValStr := strings.TrimSpace(e.DOM.Find(".score-meter-next-value").Text())
	info.ScoreToNextRank, _ = strconv.Atoi(nextValStr)

	return info
}
