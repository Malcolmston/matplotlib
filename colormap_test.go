package matplotlib

import "testing"

func TestColormapEndpoints(t *testing.T) {
	cm := NewColormap([]Color{Black, White})
	if got := cm.At(0); got != Black {
		t.Errorf("At(0) = %v, want %v", got, Black)
	}
	if got := cm.At(1); got != White {
		t.Errorf("At(1) = %v, want %v", got, White)
	}
	// Clamping outside [0,1].
	if got := cm.At(-5); got != Black {
		t.Errorf("At(-5) = %v, want Black", got)
	}
	if got := cm.At(2); got != White {
		t.Errorf("At(2) = %v, want White", got)
	}
}

func TestColormapInterpolation(t *testing.T) {
	cm := NewColormap([]Color{Black, White})
	got := cm.At(0.5)
	want := Color{128, 128, 128, 255}
	if got != want {
		t.Errorf("At(0.5) = %v, want %v", got, want)
	}
	// Three-anchor map: midpoint of first segment.
	cm3 := NewColormap([]Color{{0, 0, 0, 255}, {100, 200, 40, 255}, {255, 255, 255, 255}})
	// t=0.25 is the middle of the first of two segments (segment covers [0,0.5]).
	got3 := cm3.At(0.25)
	want3 := Color{50, 100, 20, 255}
	if got3 != want3 {
		t.Errorf("At(0.25) = %v, want %v", got3, want3)
	}
	if mid := cm3.At(0.5); mid != (Color{100, 200, 40, 255}) {
		t.Errorf("At(0.5) = %v, want middle anchor", mid)
	}
}

func TestColormapReversed(t *testing.T) {
	cm := Grays().Reversed()
	if got := cm.At(0); got != White {
		t.Errorf("reversed At(0) = %v, want White", got)
	}
	if got := cm.At(1); got != Black {
		t.Errorf("reversed At(1) = %v, want Black", got)
	}
}

func TestNamedColormapEndpoints(t *testing.T) {
	tests := []struct {
		name   string
		cm     Colormap
		lo, hi Color
	}{
		{"viridis", Viridis(), Color{68, 1, 84, 255}, Color{253, 231, 37, 255}},
		{"plasma", Plasma(), Color{13, 8, 135, 255}, Color{240, 249, 33, 255}},
		{"jet", Jet(), Color{0, 0, 131, 255}, Color{131, 0, 0, 255}},
		{"coolwarm", Coolwarm(), Color{59, 76, 192, 255}, Color{180, 4, 38, 255}},
		{"grays", Grays(), Black, White},
	}
	for _, tt := range tests {
		if got := tt.cm.At(0); got != tt.lo {
			t.Errorf("%s At(0) = %v, want %v", tt.name, got, tt.lo)
		}
		if got := tt.cm.At(1); got != tt.hi {
			t.Errorf("%s At(1) = %v, want %v", tt.name, got, tt.hi)
		}
	}
}

func TestNormalize(t *testing.T) {
	n := NewNormalize(0, 100)
	cases := []struct {
		v    float64
		want float64
	}{
		{0, 0}, {50, 0.5}, {100, 1}, {-10, 0}, {150, 1},
	}
	for _, c := range cases {
		if got := n.Norm(c.v); got != c.want {
			t.Errorf("Norm(%v) = %v, want %v", c.v, got, c.want)
		}
	}
	// Degenerate range.
	if got := NewNormalize(5, 5).Norm(5); got != 0 {
		t.Errorf("degenerate Norm = %v, want 0", got)
	}
	// Map through a grayscale colormap.
	if got := n.Map(Grays(), 50); got != (Color{128, 128, 128, 255}) {
		t.Errorf("Map(50) = %v, want mid-gray", got)
	}
}

func BenchmarkColormapAt(b *testing.B) {
	cm := Viridis()
	for i := 0; i < b.N; i++ {
		_ = cm.At(float64(i%100) / 100)
	}
}
