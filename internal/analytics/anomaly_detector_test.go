package analytics

import "testing"

func TestAnomalyDetector_Warmup_NoAnomaly(t *testing.T) {
	ad := NewAnomalyDetector(5, 2.0)

	// пока окно не заполнено — аномалии не считаем
	for i := 0; i < 4; i++ {
		if isA := ad.Add(100); isA {
			t.Fatalf("expected no anomaly during warmup, got anomaly at i=%d", i)
		}
	}
}

func TestAnomalyDetector_DetectsSpike(t *testing.T) {
	ad := NewAnomalyDetector(50, 2.0)

	// заполняем окно стабильными значениями
	for i := 0; i < 50; i++ {
		if isA := ad.Add(100); isA {
			t.Fatalf("expected no anomaly in stable data at i=%d", i)
		}
	}

	// выброс
	if isA := ad.Add(2000); !isA {
		t.Fatalf("expected anomaly for spike")
	}
}

func TestAnomalyDetector_Reset(t *testing.T) {
	ad := NewAnomalyDetector(10, 2.0)
	for i := 0; i < 10; i++ {
		_ = ad.Add(100)
	}

	ad.Reset()
	mean, std, count := ad.GetStats()
	if count != 0 {
		t.Fatalf("expected count 0 after reset, got %d", count)
	}
	if mean != 0 || std != 0 {
		t.Fatalf("expected mean/std 0 after reset, got mean=%v std=%v", mean, std)
	}
}
