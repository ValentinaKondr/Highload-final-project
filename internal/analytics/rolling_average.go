package analytics

import (
	"sync"
)

// RollingAverage calculates rolling average over a sliding window
type RollingAverage struct {
	windowSize int
	values     []float64
	mu         sync.RWMutex
}

// NewRollingAverage creates a new RollingAverage with specified window size
func NewRollingAverage(windowSize int) *RollingAverage {
	return &RollingAverage{
		windowSize: windowSize,
		values:     make([]float64, 0, windowSize),
	}
}

// Add adds a new value to the window
func (ra *RollingAverage) Add(value float64) {
	ra.mu.Lock()
	defer ra.mu.Unlock()

	ra.values = append(ra.values, value)
	if len(ra.values) > ra.windowSize {
		ra.values = ra.values[1:]
	}
}

// GetAverage calculates and returns the current average
func (ra *RollingAverage) GetAverage() float64 {
	ra.mu.RLock()
	defer ra.mu.RUnlock()

	if len(ra.values) == 0 {
		return 0.0
	}

	sum := 0.0
	for _, v := range ra.values {
		sum += v
	}
	return sum / float64(len(ra.values))
}

// GetCount returns the number of values in the window
func (ra *RollingAverage) GetCount() int {
	ra.mu.RLock()
	defer ra.mu.RUnlock()
	return len(ra.values)
}

// Reset clears all values
func (ra *RollingAverage) Reset() {
	ra.mu.Lock()
	defer ra.mu.Unlock()
	ra.values = make([]float64, 0, ra.windowSize)
}
