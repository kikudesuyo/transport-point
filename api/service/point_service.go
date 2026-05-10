package service

import (
	"fmt"
	"github.com/kikudesuyo/point-hub/external"
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

func (s *PointService) FetchTokyu() ([]UnifiedPoint, error) {
	data, err := s.tokyu.FetchAll()
	if err != nil {
		return nil, fmt.Errorf("tokyu fetch error: %w", err)
	}

	up := UnifiedPoint{Provider: "Tokyu", Balance: data.Point}
	for _, exp := range data.Expiries {
		s.addExpiry(&up, exp.Balance, exp.Date)
	}
	s.finalizePoint(&up)
	return []UnifiedPoint{up}, nil
}

func (s *PointService) FetchMetpo() ([]UnifiedPoint, error) {
	user, pass := os.Getenv("TOKYO_METRO_EMAIL"), os.Getenv("TOKYO_METRO_PASSWORD")
	if user == "" || pass == "" {
		return nil, fmt.Errorf("credentials not found")
	}

	if err := s.metpo.Login(user, pass); err != nil {
		return nil, fmt.Errorf("metpo login error: %w", err)
	}

	data, err := s.metpo.FetchAll()
	if err != nil {
		return nil, fmt.Errorf("metpo fetch error: %w", err)
	}

	up := UnifiedPoint{Provider: "Tokyo Metro (Metpo)", Balance: data.Point.HoldingPoint}
	s.addExpiry(&up, data.Point.NormalExpiryPoint, data.Point.NormalExpiry)
	s.addExpiry(&up, data.Point.ChargeExpiryPoint, data.Point.ChargeExpiry)
	s.finalizePoint(&up)
	return []UnifiedPoint{up}, nil
}

func (s *PointService) FetchToei() ([]UnifiedPoint, error) {
	user, pass := os.Getenv("TOEI_USER_ID"), os.Getenv("TOEI_METRO_PASSWORD")
	if user == "" || pass == "" {
		return nil, fmt.Errorf("credentials not found")
	}

	if err := s.toei.Login(user, pass); err != nil {
		return nil, fmt.Errorf("toei login error: %w", err)
	}

	data, err := s.toei.FetchAll()
	if err != nil {
		return nil, fmt.Errorf("toei fetch error: %w", err)
	}

	return []UnifiedPoint{{Provider: "Toei Metro (Tokopo)", Balance: data.Point}}, nil
}

func (s *PointService) FetchSotetsu() ([]UnifiedPoint, error) {
	user, pass := os.Getenv("SOTETSU_EMAIL"), os.Getenv("SOTETSU_PASSWORD")
	if user == "" || pass == "" {
		return nil, fmt.Errorf("credentials not found")
	}

	if err := s.sotetsu.Login(user, pass); err != nil {
		return nil, fmt.Errorf("sotetsu login error: %w", err)
	}

	data, err := s.sotetsu.FetchAll()
	if err != nil {
		return nil, fmt.Errorf("sotetsu fetch error: %w", err)
	}

	var results []UnifiedPoint
	results = append(results, UnifiedPoint{
		Provider: "Sotetsu", Balance: data.Point, ExpiryDate: data.PointExpiry,
	})
	if data.Mile > 0 {
		results = append(results, UnifiedPoint{
			Provider: "Sotetsu (Mile)", Balance: data.Mile, ExpiryDate: data.MileExpiry,
		})
	}
	return results, nil
}

func (s *PointService) FetchKeikyu() ([]UnifiedPoint, error) {
	user, pass := os.Getenv("KEIKYU_LOGIN_ID"), os.Getenv("KEIKYU_PASSWORD")
	if user == "" || pass == "" {
		return nil, fmt.Errorf("credentials not found")
	}

	if err := s.keikyu.Login(user, pass); err != nil {
		return nil, fmt.Errorf("keikyu login error: %w", err)
	}

	data, err := s.keikyu.FetchAll()
	if err != nil {
		return nil, fmt.Errorf("keikyu fetch error: %w", err)
	}

	return []UnifiedPoint{{
		Provider: "Keikyu", Balance: data.AvailablePoint + data.LimitedPoint,
		ExpiryDate: strings.TrimSpace(data.RevocationInfo),
	}}, nil
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
