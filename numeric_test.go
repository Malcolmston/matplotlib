package matplotlib

import (
	"math"
	"testing"
)

func floatsEqual(a, b []float64, tol float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > tol {
			return false
		}
	}
	return true
}

func TestLinspace(t *testing.T) {
	tests := []struct {
		start, stop float64
		num         int
		want        []float64
	}{
		{0, 1, 5, []float64{0, 0.25, 0.5, 0.75, 1}},
		{0, 10, 11, []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{-1, 1, 3, []float64{-1, 0, 1}},
		{5, 5, 1, []float64{5}},
		{2, 3, 0, []float64{}},
	}
	for _, tt := range tests {
		got := Linspace(tt.start, tt.stop, tt.num)
		if !floatsEqual(got, tt.want, 1e-12) {
			t.Errorf("Linspace(%v,%v,%d) = %v, want %v", tt.start, tt.stop, tt.num, got, tt.want)
		}
	}
	// Endpoints must be exact.
	got := Linspace(0, 1, 7)
	if got[0] != 0 || got[len(got)-1] != 1 {
		t.Errorf("Linspace endpoints not exact: %v", got)
	}
}

func TestArange(t *testing.T) {
	tests := []struct {
		start, stop, step float64
		want              []float64
	}{
		{0, 5, 1, []float64{0, 1, 2, 3, 4}},
		{0, 1, 0.25, []float64{0, 0.25, 0.5, 0.75}},
		{5, 0, -1, []float64{5, 4, 3, 2, 1}},
		{0, 5, 0, []float64{}},
		{0, -5, 1, []float64{}},
	}
	for _, tt := range tests {
		got := Arange(tt.start, tt.stop, tt.step)
		if !floatsEqual(got, tt.want, 1e-9) {
			t.Errorf("Arange(%v,%v,%v) = %v, want %v", tt.start, tt.stop, tt.step, got, tt.want)
		}
	}
}

func TestLogspace(t *testing.T) {
	got := Logspace(0, 2, 3, 10)
	want := []float64{1, 10, 100}
	if !floatsEqual(got, want, 1e-9) {
		t.Errorf("Logspace(0,2,3,10) = %v, want %v", got, want)
	}
	got2 := Logspace(0, 3, 4, 2)
	want2 := []float64{1, 2, 4, 8}
	if !floatsEqual(got2, want2, 1e-9) {
		t.Errorf("Logspace(0,3,4,2) = %v, want %v", got2, want2)
	}
	// Non-positive base falls back to 10.
	if !floatsEqual(Logspace(0, 1, 2, 0), []float64{1, 10}, 1e-9) {
		t.Errorf("Logspace with base 0 should default to 10")
	}
}

func BenchmarkLinspace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Linspace(0, 1000, 1000)
	}
}
