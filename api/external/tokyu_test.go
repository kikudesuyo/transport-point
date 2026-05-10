package external

import (
	"os"
	"testing"
)

func init() {
	// 実行ディレクトリを api/ に変更して、tokyu_cookie.json を正しく読み書きできるようにする
	os.Chdir("..")
}

func TestTokyuClient_FetchAll_MinimalCookies(t *testing.T) {
	client, err := NewTokyuClient()
	if err != nil {
		t.Errorf("NewTokyuClient failed: %v", err)
		return
	}

	data, err := client.FetchAll()
	if err != nil {
		t.Logf("Minimal cookies fetch failed: %v (expected if tracking cookies are mandatory)", err)
	} else {
		t.Logf("Minimal cookies fetch success! Points: %d", data.Point)
	}
}

func TestTokyuClient_FetchAll_FromJSON(t *testing.T) {
	client, err := NewTokyuClient()
	if err != nil {
		t.Errorf("NewTokyuClient failed: %v", err)
		return
	}

	data, err := client.FetchAll()
	if err != nil {
		t.Errorf("Fetch from JSON failed: %v", err)
	} else {
		t.Logf("Fetch from JSON success! Points: %d", data.Point)
	}
}
