package analytics

import (
	"math"
	"sync"
)

// AnomalyDetector detects anomalies using z-score method
type AnomalyDetector struct {
	windowSize int
	values     []float64
	threshold  float64 // z-score threshold (default 2.0)

	lastZ         float64
	lastIsAnomaly bool

	mu sync.RWMutex
}

// NewAnomalyDetector creates a new AnomalyDetector
func NewAnomalyDetector(windowSize int, threshold float64) *AnomalyDetector {
	if windowSize < 1 {
		windowSize = 50
	}
	if threshold <= 0 {
		threshold = 2.0
	}

	return &AnomalyDetector{
		windowSize: windowSize,
		values:     make([]float64, 0, windowSize),
		threshold:  threshold,
	}
}

func meanStd(values []float64) (mean, std float64) {
	n := float64(len(values))
	if n == 0 {
		return 0, 0
	}

	for _, v := range values {
		mean += v
	}
	mean /= n

	for _, v := range values {
		diff := v - mean
		std += diff * diff
	}
	std = math.Sqrt(std / n)

	return
}

// Add adds a new value and returns if it's an anomaly
func (a *AnomalyDetector) Add(value float64) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	// добавляем значение
	a.values = append(a.values, value)
	if len(a.values) > a.windowSize {
		a.values = a.values[1:]
	}

	// если данных мало — аномалии не считаем
	if len(a.values) < 2 {
		a.lastZ = 0
		a.lastIsAnomaly = false
		return false
	}

	mean, std := meanStd(a.values)
	if std == 0 {
		a.lastZ = 0
		a.lastIsAnomaly = false
		return false
	}

	z := (value - mean) / std
	a.lastZ = z
	a.lastIsAnomaly = math.Abs(z) > a.threshold

	return a.lastIsAnomaly
}

// calculateStats calculates mean and standard deviation
func (ad *AnomalyDetector) calculateStats() (mean, stdDev float64) {
	if len(ad.values) == 0 {
		return 0.0, 0.0
	}

	// Calculate mean
	sum := 0.0
	for _, v := range ad.values {
		sum += v
	}
	mean = sum / float64(len(ad.values))

	// Calculate standard deviation
	variance := 0.0
	for _, v := range ad.values {
		variance += math.Pow(v-mean, 2)
	}
	variance /= float64(len(ad.values))
	stdDev = math.Sqrt(variance)

	return mean, stdDev
}

// GetStats returns current statistics
func (a *AnomalyDetector) GetStats() (mean, std float64, count int) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	count = len(a.values)
	if count == 0 {
		return 0, 0, 0
	}

	mean, std = meanStd(a.values)
	return
}

// Reset clears all values
func (ad *AnomalyDetector) Reset() {
	ad.mu.Lock()
	defer ad.mu.Unlock()
	ad.values = make([]float64, 0, ad.windowSize)
}

func (a *AnomalyDetector) GetWindowSize() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.windowSize
}

func (a *AnomalyDetector) GetThreshold() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.threshold
}

func (a *AnomalyDetector) GetLastDecision() (z float64, isAnomaly bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.lastZ, a.lastIsAnomaly
}
