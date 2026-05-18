package external

import (
	"os"
	"sync"
	"testing"
)

func TestKeikyuClient_ConsecutiveRequests(t *testing.T) {
	loginID := os.Getenv("KEIKYU_LOGIN_ID")
	password := os.Getenv("KEIKYU_PASSWORD")
	if loginID == "" || password == "" {
		t.Skip("Credentials not set, skipping Keikyu test")
	}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			t.Logf("--- Request %d ---", i+1)
			client, err := NewKeikyuClient()
			if err != nil {
				t.Errorf("Failed to create client on request %d: %v", i+1, err)
				return
			}

			err = client.Login(loginID, password)
			if err != nil {
				t.Errorf("Login failed on request %d: %v", i+1, err)
				return
			}

			data, err := client.FetchAll()
			if err != nil {
				t.Errorf("FetchAll failed on request %d: %v", i+1, err)
				return
			}

			t.Logf("Request %d success! Available: %d, Limited: %d", i+1, data.AvailablePoint, data.LimitedPoint)
		}(i)
	}
	wg.Wait()
}
