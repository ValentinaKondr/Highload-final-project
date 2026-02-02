package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/highload-service/internal/analytics"
	"github.com/highload-service/internal/cache"
	"github.com/highload-service/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getenvInt(name string, def int) int {
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}

func getenvFloat(name string, def float64) float64 {
	v := os.Getenv(name)
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return def
	}
	return f
}

type Metric struct {
	Timestamp int64   `json:"timestamp"`
	CPU       float64 `json:"cpu"`
	RPS       float64 `json:"rps"`
}

type Service struct {
	cache             *cache.RedisCache
	rollingAvg        *analytics.RollingAverage
	anomalyDetector   *analytics.AnomalyDetector
	rpsCounter        int64
	anomalyCounter    int64
	lastRPSUpdate     time.Time
	lastAnomalyUpdate time.Time
}

func NewService() (*Service, error) {
	windowSize := getenvInt("WINDOW_SIZE", 50)
	anomalyThreshold := getenvFloat("ANOMALY_THRESHOLD", 2.0)
	redisDB := getenvInt("REDIS_DB", 0)

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		redisPassword = ""
	}

	redisCache, err := cache.NewRedisCache(redisAddr, redisPassword, redisDB)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	return &Service{
		cache:             redisCache,
		rollingAvg:        analytics.NewRollingAverage(windowSize),
		anomalyDetector:   analytics.NewAnomalyDetector(windowSize, anomalyThreshold),
		lastRPSUpdate:     time.Now(),
		lastAnomalyUpdate: time.Now(),
	}, nil
}

func (s *Service) handleMetrics(w http.ResponseWriter, r *http.Request) {
	var metric Metric
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		metrics.RequestTotal.WithLabelValues(r.Method, "/metrics", "400").Inc()
		return
	}

	// Update timestamp if not provided
	if metric.Timestamp == 0 {
		metric.Timestamp = time.Now().Unix()
	}

	// Store in Redis cache
	cacheKey := fmt.Sprintf("metric:%d", metric.Timestamp)
	if err := s.cache.Set(cacheKey, metric, 5*time.Minute); err != nil {
		log.Printf("Failed to cache metric: %v", err)
	}

	// Update rolling average with RPS
	s.rollingAvg.Add(metric.RPS)
	avg := s.rollingAvg.GetAverage()
	metrics.RollingAverageValue.Set(avg)

	// Update CPU metric
	metrics.CPUMetric.Set(metric.CPU)

	// Detect anomalies
	isAnomaly := s.anomalyDetector.Add(metric.RPS)
	if isAnomaly {
		s.anomalyCounter++
		metrics.AnomalyCount.Inc()
		log.Printf("Anomaly detected: RPS=%.2f, Timestamp=%d", metric.RPS, metric.Timestamp)
	}

	// Update RPS counter
	s.rpsCounter++
	now := time.Now()
	elapsed := now.Sub(s.lastRPSUpdate).Seconds()
	if elapsed >= 1.0 {
		currentRPS := float64(s.rpsCounter) / elapsed
		metrics.RPSRate.Set(currentRPS)
		s.rpsCounter = 0
		s.lastRPSUpdate = now
	}

	// Update anomaly rate per minute
	anomalyElapsed := now.Sub(s.lastAnomalyUpdate).Minutes()
	if anomalyElapsed >= 1.0 {
		anomalyRate := float64(s.anomalyCounter) / anomalyElapsed
		metrics.AnomalyRate.Set(anomalyRate)
		s.anomalyCounter = 0
		s.lastAnomalyUpdate = now
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":          "ok",
		"rolling_average": avg,
		"is_anomaly":      isAnomaly,
	})

	metrics.RequestTotal.WithLabelValues(r.Method, "/metrics", "200").Inc()
}

func (s *Service) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	avg := s.rollingAvg.GetAverage()
	mean, stdDev, count := s.anomalyDetector.GetStats()

	response := map[string]interface{}{
		"rolling_average": avg,
		"anomaly_stats": map[string]interface{}{
			"mean":        mean,
			"std_dev":     stdDev,
			"threshold":   getenvFloat("ANOMALY_THRESHOLD", 2.0),
			"window_size": count,
		},
		"timestamp": time.Now().Unix(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	metrics.RequestTotal.WithLabelValues(r.Method, "/analyze", "200").Inc()
}

func (s *Service) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
	metrics.RequestTotal.WithLabelValues(r.Method, "/health", "200").Inc()
}

func (s *Service) setupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Metrics endpoint for Prometheus
	r.Path("/metrics").Methods("GET").Handler(promhttp.Handler())

	// API endpoints
	r.HandleFunc("/metrics", s.handleMetrics).Methods("POST")
	r.HandleFunc("/analyze", s.handleAnalyze).Methods("GET")
	r.HandleFunc("/health", s.handleHealth).Methods("GET")

	return r
}

func main() {
	service, err := NewService()
	if err != nil {
		log.Fatalf("Failed to initialize service: %v", err)
	}
	defer service.cache.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := service.setupRoutes()

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
