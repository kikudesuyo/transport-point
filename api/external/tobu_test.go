package external

import (
	"encoding/json"
	"os"
	"testing"
)

func TestTobuClient_FetchAll(t *testing.T) {
	cookies := map[string]string{
		"JSESSIONID": "2E4DC9D043B27357E2ED0E4A51C34079",
		"AWSALB":     "0yZ/QifsyKKTwlQ8KLJQIG8eOOwe3heU+swvsnB1r79KaM3j7VKIYjE4/8jyyckoUoqp0K9JkMaoUigL+e4sD0nrEK45rgbiHVCDyoFL1zLmNvInS6d/oXZLlcc1",
		"AWSALBCORS": "0yZ/QifsyKKTwlQ8KLJQIG8eOOwe3heU+swvsnB1r79KaM3j7VKIYjE4/8jyyckoUoqp0K9JkMaoUigL+e4sD0nrEK45rgbiHVCDyoFL1zLmNvInS6d/oXZLlcc1",
		"_yjsu_yjad": "1778596125.eab12fea-b63e-4a3b-8da0-f57c29b6f91f",
	}

	data, _ := json.MarshalIndent(cookies, "", "  ")
	os.WriteFile("tobu_cookie.json", data, 0644)
	defer os.Remove("tobu_cookie.json")

	client, _ := NewTobuClient()
	res, err := client.FetchAll()
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	t.Logf("Fetch success! Total: %d, Miles: %d", res.TotalPoint, res.Miles)
	t.Logf("Breakdown: Normal=%d (%s), Limited=%d (%s)", res.NormalPoint, res.NormalExpiry, res.LimitedPoint, res.LimitedExpiry)
}
