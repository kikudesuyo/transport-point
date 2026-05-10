package external

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func loadTestCookies() map[string]string {
	f, err := os.Open("../tokyu_cookie.json")
	if err != nil {
		return nil
	}
	defer f.Close()
	data, _ := io.ReadAll(f)
	var cookies map[string]string
	json.Unmarshal(data, &cookies)
	return cookies
}

func saveTestCookies(cookies map[string]string) {
	data, _ := json.MarshalIndent(cookies, "", "  ")
	os.WriteFile("../tokyu_cookie.json", data, 0644)
}

func TestTokyuClient_FetchAll_MinimalCookies(t *testing.T) {
	client, err := NewTokyuClient()
	if err != nil {
		t.Errorf("NewTokyuClient failed: %v", err)
		return
	}
	cookies := loadTestCookies()
	if cookies == nil {
		t.Skip("tokyu_cookie.json が見つからないためスキップします")
	}
	client.SetCookies(cookies)
	defer func() {
		saveTestCookies(client.GetCookies())
	}()

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
	cookies := loadTestCookies()
	if cookies == nil {
		t.Skip("tokyu_cookie.json が見つからないためスキップします")
	}
	client.SetCookies(cookies)
	defer func() {
		saveTestCookies(client.GetCookies())
	}()

	data, err := client.FetchAll()
	if err != nil {
		t.Errorf("Fetch from JSON failed: %v", err)
	} else {
		t.Logf("Fetch from JSON success! Points: %d", data.Point)
	}
}
