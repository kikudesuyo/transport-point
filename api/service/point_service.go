package service

import (
	"fmt"
	"hoge/external"
	"os"
	"sort"
	"strings"
)

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
	tokyuClient      *external.TokyuClient
	metpoClient      *external.MetpoClient
	toeiClient       *external.ToeiMetroClient
	sotetsuClient    *external.SotetsuClient
	keikyuClient     *external.KeikyuClient
}

func NewPointService() (*PointService, error) {
	tokyu, _ := external.NewTokyuClient()
	metpo, _ := external.NewMetpoClient()
	toei, _ := external.NewToeiMetroClient()
	sotetsu, _ := external.NewSotetsuClient()
	keikyu, _ := external.NewKeikyuClient()

	return &PointService{
		tokyuClient:   tokyu,
		metpoClient:   metpo,
		toeiClient:    toei,
		sotetsuClient: sotetsu,
		keikyuClient:  keikyu,
	}, nil
}

func (s *PointService) FetchAll() (*PointReport, error) {
	report := &PointReport{
		Details: []UnifiedPoint{},
	}

	// Tokyu
	if os.Getenv("TOKYU_SESSION_TOKEN") != "" {
		s.fetchTokyu(report)
	}

	// Metpo
	if os.Getenv("TOKYO_METRO_EMAIL") != "" {
		s.fetchMetpo(report)
	}

	// Toei
	if os.Getenv("TOEI_USER_ID") != "" {
		s.fetchToei(report)
	}

	// Sotetsu
	if os.Getenv("SOTETSU_EMAIL") != "" {
		s.fetchSotetsu(report)
	}

	// Keikyu
	if os.Getenv("KEIKYU_LOGIN_ID") != "" {
		s.fetchKeikyu(report)
	}

	// Calculate total
	for _, d := range report.Details {
		report.TotalBalance += d.Balance
	}

	return report, nil
}

func (s *PointService) fetchTokyu(report *PointReport) {
	cookies := map[string]string{
		"__Host-plus.sessionToken": os.Getenv("TOKYU_SESSION_TOKEN"),
		"nToken":                   os.Getenv("TOKYU_SESSION_TOKEN"),
		"s.sessionToken":           os.Getenv("TOKYU_SESSION_TOKEN"),
		"onToken":                  os.Getenv("TOKYU_SESSION_TOKEN"),
		"_clck":                    os.Getenv("TOKYU_CLCK"),
		"_clsk":                    os.Getenv("TOKYU_CLSK"),
		"_ga":                      os.Getenv("TOKYU_GA"),
		"_ga_B0V3646TYC":           os.Getenv("TOKYU_GA_B0"),
		"_ga_XD2N3Y0135":           os.Getenv("TOKYU_GA_XD"),
		"_ga_Y86R0E9JVH":           os.Getenv("TOKYU_GA_Y8"),
		"_gcl_au":                  os.Getenv("TOKYU_GCL_AU"),
		"_rslgvry":                 os.Getenv("TOKYU_RSLGVRY"),
		"_yjsu_yjad":               os.Getenv("TOKYU_YJSU"),
		"krt_rewrite_uid":          os.Getenv("TOKYU_KRT"),
		"withdesk-id":              os.Getenv("TOKYU_WITHDESK"),
	}
	s.tokyuClient.SetCookies(cookies)

	data, err := s.tokyuClient.FetchAll()
	if err != nil {
		fmt.Printf("Tokyu fetch error: %v\n", err)
		return
	}

	up := UnifiedPoint{
		Provider: "Tokyu",
		Balance:  data.Point,
	}

	for _, exp := range data.Expiries {
		if exp.Balance > 0 {
			up.ExpiryList = append(up.ExpiryList, ExpiryInfo{
				Points: exp.Balance,
				Date:   exp.Date,
			})
		}
	}
	// Sort by date and pick nearest
	if len(up.ExpiryList) > 0 {
		sort.Slice(up.ExpiryList, func(i, j int) bool {
			return up.ExpiryList[i].Date < up.ExpiryList[j].Date
		})
		up.ExpiryDate = up.ExpiryList[0].Date
	}

	report.Details = append(report.Details, up)

	// Update cookies if changed
	updated := s.tokyuClient.GetCookies()
	newToken := updated["__Host-plus.sessionToken"]
	if newToken == "" {
		newToken = updated["nToken"]
	}
	if newToken != "" && newToken != os.Getenv("TOKYU_SESSION_TOKEN") {
		if report.UpdatedCookies == nil {
			report.UpdatedCookies = make(map[string]string)
		}
		report.UpdatedCookies["TOKYU_SESSION_TOKEN"] = newToken
		if clck := updated["_clck"]; clck != "" {
			report.UpdatedCookies["TOKYU_CLCK"] = clck
		}
		if clsk := updated["_clsk"]; clsk != "" {
			report.UpdatedCookies["TOKYU_CLSK"] = clsk
		}
	}
}

func (s *PointService) fetchMetpo(report *PointReport) {
	err := s.metpoClient.Login(os.Getenv("TOKYO_METRO_EMAIL"), os.Getenv("TOKYO_METRO_PASSWORD"))
	if err != nil {
		fmt.Printf("Metpo login error: %v\n", err)
		return
	}

	data, err := s.metpoClient.FetchAll()
	if err != nil {
		fmt.Printf("Metpo fetch error: %v\n", err)
		return
	}

	up := UnifiedPoint{
		Provider: "Tokyo Metro (Metpo)",
		Balance:  data.Point.HoldingPoint,
	}

	if data.Point.NormalExpiryPoint > 0 {
		up.ExpiryList = append(up.ExpiryList, ExpiryInfo{
			Points: data.Point.NormalExpiryPoint,
			Date:   data.Point.NormalExpiry,
		})
	}
	if data.Point.ChargeExpiryPoint > 0 {
		up.ExpiryList = append(up.ExpiryList, ExpiryInfo{
			Points: data.Point.ChargeExpiryPoint,
			Date:   data.Point.ChargeExpiry,
		})
	}
	
	if len(up.ExpiryList) > 0 {
		sort.Slice(up.ExpiryList, func(i, j int) bool {
			return up.ExpiryList[i].Date < up.ExpiryList[j].Date
		})
		up.ExpiryDate = up.ExpiryList[0].Date
	}

	report.Details = append(report.Details, up)
}

func (s *PointService) fetchToei(report *PointReport) {
	err := s.toeiClient.Login(os.Getenv("TOEI_USER_ID"), os.Getenv("TOEI_METRO_PASSWORD"))
	if err != nil {
		fmt.Printf("Toei login error: %v\n", err)
		return
	}

	data, err := s.toeiClient.FetchAll()
	if err != nil {
		fmt.Printf("Toei fetch error: %v\n", err)
		return
	}

	report.Details = append(report.Details, UnifiedPoint{
		Provider: "Toei Metro (Tokopo)",
		Balance:  data.Point,
	})
}

func (s *PointService) fetchSotetsu(report *PointReport) {
	err := s.sotetsuClient.Login(os.Getenv("SOTETSU_EMAIL"), os.Getenv("SOTETSU_PASSWORD"))
	if err != nil {
		fmt.Printf("Sotetsu login error: %v\n", err)
		return
	}

	data, err := s.sotetsuClient.FetchAll()
	if err != nil {
		fmt.Printf("Sotetsu fetch error: %v\n", err)
		return
	}

	up := UnifiedPoint{
		Provider:   "Sotetsu",
		Balance:    data.Point,
		ExpiryDate: data.PointExpiry,
	}
	report.Details = append(report.Details, up)
	
	// If Mile exists, add it as well?
	if data.Mile > 0 {
		report.Details = append(report.Details, UnifiedPoint{
			Provider:   "Sotetsu (Mile)",
			Balance:    data.Mile,
			ExpiryDate: data.MileExpiry,
		})
	}
}

func (s *PointService) fetchKeikyu(report *PointReport) {
	err := s.keikyuClient.Login(os.Getenv("KEIKYU_LOGIN_ID"), os.Getenv("KEIKYU_PASSWORD"))
	if err != nil {
		fmt.Printf("Keikyu login error: %v\n", err)
		return
	}

	data, err := s.keikyuClient.FetchAll()
	if err != nil {
		fmt.Printf("Keikyu fetch error: %v\n", err)
		return
	}

	report.Details = append(report.Details, UnifiedPoint{
		Provider: "Keikyu",
		Balance:  data.AvailablePoint + data.LimitedPoint,
		// Keikyu expiry info is textual (RevocationInfo)
		ExpiryDate: strings.TrimSpace(data.RevocationInfo),
	})
}
