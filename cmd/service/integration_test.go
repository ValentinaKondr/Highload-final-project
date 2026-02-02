package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Интеграционный тест: POST /metrics -> GET /analyze
func TestIntegration_MetricsAndAnalyze(t *testing.T) {
	service := NewTestService()
	router := service.setupRoutes()

	server := httptest.NewServer(router)
	defer server.Close()

	metric := Metric{
		Timestamp: time.Now().Unix(),
		CPU:       42.0,
		RPS:       100.0,
	}

	body, err := json.Marshal(metric)
	if err != nil {
		t.Fatalf("failed to marshal metric: %v", err)
	}

	// --- Act: POST /metrics ---
	resp, err := http.Post(
		server.URL+"/metrics",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		t.Fatalf("POST /metrics failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	// --- Act: GET /analyze ---
	resp, err = http.Get(server.URL + "/analyze")
	if err != nil {
		t.Fatalf("GET /analyze failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var result struct {
		RollingAverage float64 `json:"rolling_average"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode analyze response: %v", err)
	}

	// --- Assert ---
	if result.RollingAverage <= 0 {
		t.Fatalf("expected rolling_average > 0, got %.2f", result.RollingAverage)
	}
}
