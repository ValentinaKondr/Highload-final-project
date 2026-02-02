package analytics

import "testing"

func TestRollingAverage_Empty(t *testing.T) {
	ra := NewRollingAverage(3)
	if got := ra.GetAverage(); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
	if got := ra.GetCount(); got != 0 {
		t.Fatalf("expected count 0, got %d", got)
	}
}

func TestRollingAverage_SlidingWindow(t *testing.T) {
	ra := NewRollingAverage(3)

	ra.Add(1) // avg=1
	if got := ra.GetAverage(); got != 1 {
		t.Fatalf("expected 1, got %v", got)
	}

	ra.Add(2) // avg=1.5
	if got := ra.GetAverage(); got != 1.5 {
		t.Fatalf("expected 1.5, got %v", got)
	}

	ra.Add(3) // avg=2
	if got := ra.GetAverage(); got != 2 {
		t.Fatalf("expected 2, got %v", got)
	}

	ra.Add(4) // window [2,3,4] avg=3
	if got := ra.GetAverage(); got != 3 {
		t.Fatalf("expected 3, got %v", got)
	}

	if got := ra.GetCount(); got != 3 {
		t.Fatalf("expected count 3, got %d", got)
	}
}

func TestRollingAverage_Reset(t *testing.T) {
	ra := NewRollingAverage(3)
	ra.Add(10)
	ra.Add(20)

	ra.Reset()

	if got := ra.GetCount(); got != 0 {
		t.Fatalf("expected count 0 after reset, got %d", got)
	}
	if got := ra.GetAverage(); got != 0 {
		t.Fatalf("expected avg 0 after reset, got %v", got)
	}
}
