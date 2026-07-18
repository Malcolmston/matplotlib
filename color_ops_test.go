package matplotlib

import (
	"math"
	"testing"
)

func TestLighten(t *testing.T) {
	if got := Black.Lighten(0.5); got != (Color{128, 128, 128, 255}) {
		t.Errorf("Black.Lighten(0.5) = %v, want mid-gray", got)
	}
	if got := Black.Lighten(1); got != (Color{255, 255, 255, 255}) {
		t.Errorf("Black.Lighten(1) = %v, want white", got)
	}
	if got := White.Lighten(0.5); got != White {
		t.Errorf("White.Lighten(0.5) = %v, want white", got)
	}
	// Alpha preserved.
	if got := RGBA(0, 0, 0, 100).Lighten(0.5); got.A != 100 {
		t.Errorf("Lighten changed alpha: %v", got)
	}
}

func TestDarken(t *testing.T) {
	if got := White.Darken(0.5); got != (Color{128, 128, 128, 255}) {
		t.Errorf("White.Darken(0.5) = %v, want mid-gray", got)
	}
	if got := White.Darken(1); got != (Color{0, 0, 0, 255}) {
		t.Errorf("White.Darken(1) = %v, want black", got)
	}
	if got := Black.Darken(0.5); got != Black {
		t.Errorf("Black.Darken(0.5) = %v, want black", got)
	}
}

func TestWithAlpha(t *testing.T) {
	if got := Red.WithAlpha(128); got.A != 128 || got.R != Red.R {
		t.Errorf("WithAlpha = %v", got)
	}
}

func TestLuminance(t *testing.T) {
	if got := White.Luminance(); math.Abs(got-1) > 1e-9 {
		t.Errorf("White.Luminance = %v, want 1", got)
	}
	if got := Black.Luminance(); math.Abs(got) > 1e-9 {
		t.Errorf("Black.Luminance = %v, want 0", got)
	}
	// Pure green is the brightest primary under Rec.709 weights.
	g := RGB(0, 255, 0).Luminance()
	r := RGB(255, 0, 0).Luminance()
	b := RGB(0, 0, 255).Luminance()
	if !(g > r && r > b) {
		t.Errorf("luminance ordering wrong: g=%v r=%v b=%v", g, r, b)
	}
}

func TestGrayscale(t *testing.T) {
	if got := White.Grayscale(); got != White {
		t.Errorf("White.Grayscale = %v", got)
	}
	if got := Black.Grayscale(); got != Black {
		t.Errorf("Black.Grayscale = %v", got)
	}
	gray := RGB(255, 255, 255).Grayscale()
	if gray.R != gray.G || gray.G != gray.B {
		t.Errorf("Grayscale not equal channels: %v", gray)
	}
}

func TestColorHex(t *testing.T) {
	tests := []struct {
		c    Color
		want string
	}{
		{Red, "#d62728"},
		{Black, "#000000"},
		{White, "#ffffff"},
		{RGB(31, 119, 180), "#1f77b4"},
	}
	for _, tt := range tests {
		if got := tt.c.Hex(); got != tt.want {
			t.Errorf("%v.Hex() = %q, want %q", tt.c, got, tt.want)
		}
	}
}

func TestBlend(t *testing.T) {
	if got := Blend(Black, White, 0.5); got != (Color{128, 128, 128, 255}) {
		t.Errorf("Blend(Black,White,0.5) = %v", got)
	}
	if got := Blend(Black, White, 0); got != Black {
		t.Errorf("Blend t=0 = %v, want Black", got)
	}
	if got := Blend(Black, White, 1); got != White {
		t.Errorf("Blend t=1 = %v, want White", got)
	}
	// Alpha is interpolated too.
	got := Blend(RGBA(0, 0, 0, 0), RGBA(0, 0, 0, 255), 0.5)
	if got.A != 128 {
		t.Errorf("Blend alpha = %d, want 128", got.A)
	}
}
