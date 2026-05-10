package service

import (
	"encoding/json"
	"fmt"
	"hoge/external"
	"os"
	"sort"
	"strings"
)

const tokyuCookieFile = "tokyu_session.json"

type ExpiryInfo struct {
	Points int    `json:"points"`
	Date   string `json:"date"` // YYYY-MM-DD
}

type UnifiedPoint struct {
	Provider   string       `json:"provider"`
	Balance    int          `json:"balance"`
	ExpiryDate string       `json:"expiry_date"` // Nearest expiry date
	ExpiryList []ExpiryInfo `json:"expiry_list,omitempty"`
}

type PointReport struct {
	TotalBalance   int            `json:"total_balance"`
	Details        []UnifiedPoint `json:"details"`
	UpdatedCookies map[string]string `json:"updated_cookies,omitempty"`
}

type PointService struct {
	tokyu   *external.TokyuClient
	metpo   *external.MetpoClient
	toei    *external.ToeiMetroClient
	sotetsu *external.SotetsuClient
	keikyu  *external.KeikyuClient
}

func NewPointService() (*PointService, error) {
	tokyu, _ := external.NewTokyuClient()
	metpo, _ := external.NewMetpoClient()
	toei, _ := external.NewToeiMetroClient()
	sotetsu, _ := external.NewSotetsuClient()
	keikyu, _ := external.NewKeikyuClient()

	return &PointService{
		tokyu:   tokyu,
		metpo:   metpo,
		toei:    toei,
		sotetsu: sotetsu,
		keikyu:  keikyu,
	}, nil
}

func (s *PointService) FetchAll() (*PointReport, error) {
	report := &PointReport{Details: []UnifiedPoint{}}

	// Each provider fetch logic
	s.fetchTokyu(report)
	s.fetchMetpo(report)
	s.fetchToei(report)
	s.fetchSotetsu(report)
	s.fetchKeikyu(report)

	for _, d := range report.Details {
		report.TotalBalance += d.Balance
	}
	return report, nil
}

func (s *PointService) fetchTokyu(report *PointReport) {
	token := os.Getenv("TOKYU_SESSION_TOKEN")

	// Load existing cookies from file
	cookies := s.loadTokyuCookies()

	// Update session tokens if provided in environment
	if token != "" {
		cookies["__Host-plus.sessionToken"] = token
		cookies["nToken"] = token
		cookies["s.sessionToken"] = token
		cookies["onToken"] = token
	}

	if len(cookies) == 0 {
		return // No way to authenticate
	}

	s.tokyu.SetCookies(cookies)

	data, err := s.tokyu.FetchAll()
	if err != nil {
		fmt.Printf("Tokyu fetch error: %v\n", err)
		return
	}

	up := UnifiedPoint{Provider: "Tokyu", Balance: data.Point}
	for _, exp := range data.Expiries {
		s.addExpiry(&up, exp.Balance, exp.Date)
	}
	s.finalizePoint(&up)
	report.Details = append(report.Details, up)

	// Save all current cookies back to file
	updated := s.tokyu.GetCookies()
	s.saveTokyuCookies(updated)

	// Still check for session token change for .env sync if desired
	s.syncTokyuCookies(report, updated)
}

func (s *PointService) fetchMetpo(report *PointReport) {
	user, pass := os.Getenv("TOKYO_METRO_EMAIL"), os.Getenv("TOKYO_METRO_PASSWORD")
	if user == "" || pass == "" {
		return
	}

	if err := s.metpo.Login(user, pass); err != nil {
		fmt.Printf("Metpo login error: %v\n", err)
		return
	}

	data, err := s.metpo.FetchAll()
	if err != nil {
		fmt.Printf("Metpo fetch error: %v\n", err)
		return
	}

	up := UnifiedPoint{Provider: "Tokyo Metro (Metpo)", Balance: data.Point.HoldingPoint}
	s.addExpiry(&up, data.Point.NormalExpiryPoint, data.Point.NormalExpiry)
	s.addExpiry(&up, data.Point.ChargeExpiryPoint, data.Point.ChargeExpiry)
	s.finalizePoint(&up)
	report.Details = append(report.Details, up)
}

func (s *PointService) fetchToei(report *PointReport) {
	user, pass := os.Getenv("TOEI_USER_ID"), os.Getenv("TOEI_METRO_PASSWORD")
	if user == "" || pass == "" {
		return
	}

	if err := s.toei.Login(user, pass); err != nil {
		fmt.Printf("Toei login error: %v\n", err)
		return
	}

	data, err := s.toei.FetchAll()
	if err != nil {
		fmt.Printf("Toei fetch error: %v\n", err)
		return
	}

	report.Details = append(report.Details, UnifiedPoint{Provider: "Toei Metro (Tokopo)", Balance: data.Point})
}

func (s *PointService) fetchSotetsu(report *PointReport) {
	user, pass := os.Getenv("SOTETSU_EMAIL"), os.Getenv("SOTETSU_PASSWORD")
	if user == "" || pass == "" {
		return
	}

	if err := s.sotetsu.Login(user, pass); err != nil {
		fmt.Printf("Sotetsu login error: %v\n", err)
		return
	}

	data, err := s.sotetsu.FetchAll()
	if err != nil {
		fmt.Printf("Sotetsu fetch error: %v\n", err)
		return
	}

	report.Details = append(report.Details, UnifiedPoint{
		Provider: "Sotetsu", Balance: data.Point, ExpiryDate: data.PointExpiry,
	})
	if data.Mile > 0 {
		report.Details = append(report.Details, UnifiedPoint{
			Provider: "Sotetsu (Mile)", Balance: data.Mile, ExpiryDate: data.MileExpiry,
		})
	}
}

func (s *PointService) fetchKeikyu(report *PointReport) {
	user, pass := os.Getenv("KEIKYU_LOGIN_ID"), os.Getenv("KEIKYU_PASSWORD")
	if user == "" || pass == "" {
		return
	}

	if err := s.keikyu.Login(user, pass); err != nil {
		fmt.Printf("Keikyu login error: %v\n", err)
		return
	}

	data, err := s.keikyu.FetchAll()
	if err != nil {
		fmt.Printf("Keikyu fetch error: %v\n", err)
		return
	}

	report.Details = append(report.Details, UnifiedPoint{
		Provider: "Keikyu", Balance: data.AvailablePoint + data.LimitedPoint,
		ExpiryDate: strings.TrimSpace(data.RevocationInfo),
	})
}

// Helpers
func (s *PointService) addExpiry(up *UnifiedPoint, balance int, date string) {
	if balance > 0 && date != "" {
		up.ExpiryList = append(up.ExpiryList, ExpiryInfo{Points: balance, Date: date})
	}
}

func (s *PointService) finalizePoint(up *UnifiedPoint) {
	if len(up.ExpiryList) == 0 {
		return
	}
	sort.Slice(up.ExpiryList, func(i, j int) bool {
		return up.ExpiryList[i].Date < up.ExpiryList[j].Date
	})
	up.ExpiryDate = up.ExpiryList[0].Date
}

func (s *PointService) syncTokyuCookies(report *PointReport, updated map[string]string) {
	newToken := updated["__Host-plus.sessionToken"]
	if newToken == "" {
		newToken = updated["nToken"]
	}

	if newToken != "" && newToken != os.Getenv("TOKYU_SESSION_TOKEN") {
		if report.UpdatedCookies == nil {
			report.UpdatedCookies = make(map[string]string)
		}
		report.UpdatedCookies["TOKYU_SESSION_TOKEN"] = newToken
	}
}

func (s *PointService) loadTokyuCookies() map[string]string {
	cookies := make(map[string]string)
	data, err := os.ReadFile(tokyuCookieFile)
	if err == nil {
		json.Unmarshal(data, &cookies)
	}
	return cookies
}

func (s *PointService) saveTokyuCookies(cookies map[string]string) {
	data, _ := json.MarshalIndent(cookies, "", "  ")
	os.WriteFile(tokyuCookieFile, data, 0644)
}
