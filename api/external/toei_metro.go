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

type ToeiMetroData struct {
	Point int
}

type ToeiMetroClient struct{ client *http.Client }

func NewToeiMetroClient() (*ToeiMetroClient, error) {
	jar, _ := cookiejar.New(nil)
	return &ToeiMetroClient{client: &http.Client{Jar: jar}}, nil
}

func (t *ToeiMetroClient) Login(cardNo, password string) error {
	resp, err := t.client.PostForm("https://tokopo.jp/gv/pc/login/PcLoginLoginAction.do", url.Values{
		"cardNo": {cardNo}, "password": {password}, "x": {"83"}, "y": {"13"},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// コンテンツ確認による成功判定（IDセレクタならエンコーディングに依存しない）
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	if doc.Find("#logout").Length() == 0 {
		return fmt.Errorf("ログイン失敗")
	}
	return nil
}

func (t *ToeiMetroClient) FetchAll() (*ToeiMetroData, error) {
	resp, err := t.client.Get("https://tokopo.jp/gv/pc/mymenu/MyMenuInitAction.do")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Shift-JISをUTF-8に変換して読み込む
	utf8Body, _ := charset.NewReader(resp.Body, "shift_jis")
	doc, _ := goquery.NewDocumentFromReader(utf8Body)

	data := &ToeiMetroData{}
	// ポイント抽出
	pText := doc.Find("div.boxPoint li.green").Text()
	pStr := strings.NewReplacer("ポイント", "", ",", "").Replace(pText)
	data.Point, _ = strconv.Atoi(strings.TrimSpace(pStr))

	return data, nil
}
