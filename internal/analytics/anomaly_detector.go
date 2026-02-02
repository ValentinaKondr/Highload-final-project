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
	mu         sync.RWMutex
}

// NewAnomalyDetector creates a new AnomalyDetector
func NewAnomalyDetector(windowSize int, threshold float64) *AnomalyDetector {
	return &AnomalyDetector{
		windowSize: windowSize,
		values:     make([]float64, 0, windowSize),
		threshold:  threshold,
	}
}

// Add adds a new value and returns if it's an anomaly
func (ad *AnomalyDetector) Add(value float64) bool {
	ad.mu.Lock()
	defer ad.mu.Unlock()

	isAnomaly := false
	if len(ad.values) >= ad.windowSize {
		mean, stdDev := ad.calculateStats()
		if stdDev > 0 {
			zScore := math.Abs((value - mean) / stdDev)
			isAnomaly = zScore > ad.threshold
		}
	}

	ad.values = append(ad.values, value)
	if len(ad.values) > ad.windowSize {
		ad.values = ad.values[1:]
	}

	return isAnomaly
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
func (ad *AnomalyDetector) GetStats() (mean, stdDev float64, count int) {
	ad.mu.RLock()
	defer ad.mu.RUnlock()
	mean, stdDev = ad.calculateStats()
	return mean, stdDev, len(ad.values)
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
