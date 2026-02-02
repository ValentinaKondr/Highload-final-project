package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/highload-service/internal/analytics"
	"github.com/highload-service/internal/cache"
)

type memCache struct {
	mu sync.Mutex
	m  map[string][]byte
}

func newMemCache() *memCache {
	return &memCache{m: make(map[string][]byte)}
}

func (c *memCache) Set(key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	c.m[key] = b
	return nil
}

func (c *memCache) Close() error { return nil }

var _ cache.Cache = (*memCache)(nil)

func newTestService() *Service {
	return &Service{
		cache:             newMemCache(),
		rollingAvg:        analytics.NewRollingAverage(50),
		anomalyDetector:   analytics.NewAnomalyDetector(50, 2.0),
		lastRPSUpdate:     time.Now(),
		lastAnomalyUpdate: time.Now(),
	}
}

func TestHealth(t *testing.T) {
	s := newTestService()
	ts := httptest.NewServer(s.setupRoutes())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func TestMetricsAndAnalyze(t *testing.T) {
	s := newTestService()
	ts := httptest.NewServer(s.setupRoutes())
	defer ts.Close()

	// прогрев: 50 нормальных точек
	for i := 1; i <= 50; i++ {
		body := []byte(`{"timestamp":` + itoa(int64(i)) + `,"cpu":20,"rps":100}`)
		resp, err := http.Post(ts.URL+"/metrics", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
	}

	// выброс
	resp, err := http.Post(ts.URL+"/metrics", "application/json",
		bytes.NewReader([]byte(`{"timestamp":999,"cpu":95,"rps":2000}`)))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	// анализ
	resp, err = http.Get(ts.URL + "/analyze")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var out map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}

	// минимальные проверки структуры
	stats := out["anomaly_stats"].(map[string]interface{})
	if int(stats["window_size"].(float64)) != 50 {
		t.Fatalf("expected window_size 50, got %v", stats["window_size"])
	}
}

func itoa(x int64) string {
	// маленький helper без strconv чтобы не тянуть лишнего в пример
	if x == 0 {
		return "0"
	}
	sign := ""
	if x < 0 {
		sign = "-"
		x = -x
	}
	var buf [20]byte
	i := len(buf)
	for x > 0 {
		i--
		buf[i] = byte('0' + x%10)
		x /= 10
	}
	return sign + string(buf[i:])
}
