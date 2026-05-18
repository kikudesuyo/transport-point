package service

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/kikudesuyo/point-hub/app/external"
)

type ExpiryInfo struct {
	Points int    `json:"points"`
	Date   string `json:"date"` // YYYY-MM-DD
}

type SubPoint struct {
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

type UnifiedPoint struct {
	Provider      string       `json:"provider"`
	Balance       int          `json:"balance"`
	ExpiryDate    string       `json:"expiry_date"` // Nearest expiry date
	ExpiryList    []ExpiryInfo `json:"expiry_list,omitempty"`
	SubPoints     []SubPoint   `json:"sub_points,omitempty"`
	IsMaintenance bool         `json:"is_maintenance"`
}

type PointService struct {
	tokyu   *external.TokyuClient
	metpo   *external.MetpoClient
	toei    *external.ToeiMetroClient
	sotetsu *external.SotetsuClient
	keikyu  *external.KeikyuClient
	odakyu  *external.OdakyuClient
	tobu    *external.TobuClient
}

func NewPointService() (*PointService, error) {
	tokyu, _ := external.NewTokyuClient()
	metpo, _ := external.NewMetpoClient()
	toei, _ := external.NewToeiMetroClient()
	sotetsu, _ := external.NewSotetsuClient()
	keikyu, _ := external.NewKeikyuClient()
	odakyu, _ := external.NewOdakyuClient()
	tobu, _ := external.NewTobuClient()

	return &PointService{
		tokyu:   tokyu,
		metpo:   metpo,
		toei:    toei,
		sotetsu: sotetsu,
		keikyu:  keikyu,
		odakyu:  odakyu,
		tobu:    tobu,
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

	// メンテナンス時間判定
	if s.metpo.IsMaintenance() {
		return []UnifiedPoint{{
			Provider:      "Tokyo Metro (Metpo)",
			IsMaintenance: true,
		}}, nil
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

	up := UnifiedPoint{
		Provider: "Sotetsu",
		Balance:  data.Point + data.Mile,
	}

	s.addExpiry(&up, data.Point, data.PointExpiry)
	s.addExpiry(&up, data.Mile, data.MileExpiry)

	var subPoints []SubPoint
	if data.Point > 0 {
		subPoints = append(subPoints, SubPoint{Name: "相鉄ポイント", Balance: data.Point})
	}
	if data.Mile > 0 {
		subPoints = append(subPoints, SubPoint{Name: "相鉄マイル", Balance: data.Mile})
	}
	up.SubPoints = subPoints

	s.finalizePoint(&up)
	return []UnifiedPoint{up}, nil
}

func (s *PointService) FetchKeikyu() ([]UnifiedPoint, error) {
	data, err := s.keikyu.FetchAll()
	if err != nil {
		return nil, fmt.Errorf("keikyu fetch error: %w", err)
	}

	up := UnifiedPoint{
		Provider:   "Keikyu",
		Balance:    data.AvailablePoint + data.LimitedPoint,
		ExpiryDate: strings.TrimSpace(data.RevocationInfo),
	}

	var subPoints []SubPoint
	if data.AvailablePoint > 0 {
		subPoints = append(subPoints, SubPoint{Name: "通常ポイント", Balance: data.AvailablePoint})
	}
	if data.LimitedPoint > 0 {
		subPoints = append(subPoints, SubPoint{Name: "期間限定ポイント", Balance: data.LimitedPoint})
	}
	up.SubPoints = subPoints

	return []UnifiedPoint{up}, nil
}

func (s *PointService) FetchOdakyu() ([]UnifiedPoint, error) {
	data, err := s.odakyu.FetchAll()
	if err != nil {
		return nil, fmt.Errorf("odakyu fetch error: %w", err)
	}

	up := UnifiedPoint{
		Provider:   "Odakyu",
		Balance:    data.ThisYearBalance + data.LastYearBalance,
		ExpiryDate: data.PointInvalidDate,
	}

	var subPoints []SubPoint
	if data.ThisYearBalance > 0 {
		subPoints = append(subPoints, SubPoint{Name: "今年度ポイント", Balance: data.ThisYearBalance})
	}
	if data.LastYearBalance > 0 {
		subPoints = append(subPoints, SubPoint{Name: "前年度ポイント", Balance: data.LastYearBalance})
	}
	up.SubPoints = subPoints

	return []UnifiedPoint{up}, nil
}

func (s *PointService) FetchTobu() ([]UnifiedPoint, error) {
	data, err := s.tobu.FetchAll()
	if err != nil {
		return nil, fmt.Errorf("tobu fetch error: %w", err)
	}

	up := UnifiedPoint{
		Provider: "Tobu",
		Balance:  data.TotalPoint + data.Miles,
	}

	var subPoints []SubPoint
	if data.NormalPoint > 0 {
		subPoints = append(subPoints, SubPoint{Name: "通常ポイント", Balance: data.NormalPoint})
		s.addExpiry(&up, data.NormalPoint, data.NormalExpiry)
	}
	if data.LimitedPoint > 0 {
		subPoints = append(subPoints, SubPoint{Name: "期間限定ポイント", Balance: data.LimitedPoint})
		s.addExpiry(&up, data.LimitedPoint, data.LimitedExpiry)
	}
	if data.Miles > 0 {
		subPoints = append(subPoints, SubPoint{Name: "トブポマイル", Balance: data.Miles})
	}
	up.SubPoints = subPoints
	s.finalizePoint(&up)

	return []UnifiedPoint{up}, nil
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
