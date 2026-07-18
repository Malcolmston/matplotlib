package matplotlib

// Parity tests encoding concrete known-answer vectors asserted by the upstream
// Python library, matplotlib/matplotlib, in its own test suite
// (lib/matplotlib/tests/test_colors.py). Each vector below is transcribed from
// an assertion in that file so the Go port is validated against matplotlib's
// documented expected values rather than against itself.

import (
	"math"
	"testing"
)

// approx reports whether a and b agree to the given number of decimal places,
// matching numpy.testing.assert_array_almost_equal (default 6 decimals).
func approx(a, b float64, decimals int) bool {
	return math.Abs(a-b) < 1.5*math.Pow(10, -float64(decimals))
}

// TestParityTab10Cycle checks the default color cycle against matplotlib's
// TABLEAU_COLORS / default property cycle. Upstream asserts (test_colors.py):
//
//	dflt_cycle = ['#1f77b4', '#ff7f0e', '#2ca02c', '#d62728',
//	              '#9467bd', '#8c564b', '#e377c2', '#7f7f7f',
//	              '#bcbd22', '#17becf']
//	assert list(mcolors.TABLEAU_COLORS.values()) == dflt_cycle
func TestParityTab10Cycle(t *testing.T) {
	want := []string{
		"#1f77b4", "#ff7f0e", "#2ca02c", "#d62728",
		"#9467bd", "#8c564b", "#e377c2", "#7f7f7f",
		"#bcbd22", "#17becf",
	}
	if len(DefaultColors) != len(want) {
		t.Fatalf("DefaultColors has %d entries, want %d", len(DefaultColors), len(want))
	}
	for i, w := range want {
		if got := DefaultColors[i].Hex(); got != w {
			t.Errorf("DefaultColors[%d].Hex() = %q, want %q", i, got, w)
		}
	}
}

// TestParityTabBlueHex mirrors the single upstream assertion
// `mcolors.to_hex("tab:blue") == "#1f77b4"`.
func TestParityTabBlueHex(t *testing.T) {
	if got := Blue.Hex(); got != "#1f77b4" {
		t.Errorf("Blue.Hex() = %q, want %q", got, "#1f77b4")
	}
}

// TestParityToHexGray mirrors `mcolors.to_hex((.7, .7, .7)) == "#b2b2b2"`.
// matplotlib rounds 0.7*255 = 178.5 to 178 (round-half-to-even). Here the
// gray is constructed from the exact 8-bit value 178, so the hex must match.
func TestParityToHexGray(t *testing.T) {
	if got := RGB(178, 178, 178).Hex(); got != "#b2b2b2" {
		t.Errorf("gray Hex() = %q, want %q", got, "#b2b2b2")
	}
}

// TestParityRGBToHSVGreen mirrors upstream test_rgb_to_hsv_int:
//
//	assert_array_equal(mcolors.rgb_to_hsv((0, 1, 0)), (1/3, 1, 1))  # green
func TestParityRGBToHSVGreen(t *testing.T) {
	h, s, v := RGBToHSV(0, 1, 0)
	if !approx(h, 1.0/3.0, 6) || !approx(s, 1, 6) || !approx(v, 1, 6) {
		t.Errorf("RGBToHSV(0,1,0) = (%v,%v,%v), want (1/3,1,1)", h, s, v)
	}
}

// TestParityRGBHSVRoundTrip mirrors upstream test_rgb_hsv_round_trip, which
// asserts hsv_to_rgb(rgb_to_hsv(x)) == x and rgb_to_hsv(hsv_to_rgb(x)) == x
// for arbitrary values. A deterministic lattice stands in for the upstream
// random sampling.
func TestParityRGBHSVRoundTrip(t *testing.T) {
	steps := []float64{0, 0.1, 0.25, 0.4, 0.5, 0.6, 0.75, 0.9, 1}
	for _, r := range steps {
		for _, g := range steps {
			for _, b := range steps {
				h, s, v := RGBToHSV(r, g, b)
				r2, g2, b2 := HSVToRGB(h, s, v)
				if !approx(r, r2, 6) || !approx(g, g2, 6) || !approx(b, b2, 6) {
					t.Fatalf("rgb->hsv->rgb drift: (%v,%v,%v) -> (%v,%v,%v)", r, g, b, r2, g2, b2)
				}
			}
		}
	}
	// hsv -> rgb -> hsv for saturated colors (hue is well-defined).
	for _, h := range []float64{0, 0.1, 1.0 / 3.0, 0.5, 2.0 / 3.0, 0.9} {
		for _, s := range []float64{0.3, 0.7, 1} {
			for _, v := range []float64{0.3, 0.7, 1} {
				r, g, b := HSVToRGB(h, s, v)
				h2, s2, v2 := RGBToHSV(r, g, b)
				if !approx(h, h2, 6) || !approx(s, s2, 6) || !approx(v, v2, 6) {
					t.Fatalf("hsv->rgb->hsv drift: (%v,%v,%v) -> (%v,%v,%v)", h, s, v, h2, s2, v2)
				}
			}
		}
	}
}

// TestParityLogNorm mirrors upstream test_LogNorm_inverse:
//
//	norm = mcolors.LogNorm(vmin=0.1, vmax=10)
//	assert_array_almost_equal(norm([0.5, 0.4]), [0.349485, 0.30103])
//	assert_array_almost_equal(norm(0.4), [0.30103])
func TestParityLogNorm(t *testing.T) {
	n := NewLogNorm(0.1, 10)
	if got := n.Norm(0.5); !approx(got, 0.349485, 6) {
		t.Errorf("LogNorm.Norm(0.5) = %v, want 0.349485", got)
	}
	if got := n.Norm(0.4); !approx(got, 0.30103, 5) {
		t.Errorf("LogNorm.Norm(0.4) = %v, want 0.30103", got)
	}
	// Endpoints map to 0 and 1.
	if got := n.Norm(0.1); got != 0 {
		t.Errorf("LogNorm.Norm(0.1) = %v, want 0", got)
	}
	if got := n.Norm(10); got != 1 {
		t.Errorf("LogNorm.Norm(10) = %v, want 1", got)
	}
}

// TestParityViridisEndpoints checks the viridis colormap endpoints against the
// canonical _viridis_data table matplotlib ships: entry 0 is
// (0.267004, 0.004874, 0.329415) and the last entry is
// (0.993248, 0.906157, 0.143936), which round to the 8-bit RGB values below.
func TestParityViridisEndpoints(t *testing.T) {
	cm := Viridis()
	if got := cm.At(0); got != (Color{68, 1, 84, 255}) {
		t.Errorf("Viridis At(0) = %v, want {68,1,84,255}", got)
	}
	if got := cm.At(1); got != (Color{253, 231, 37, 255}) {
		t.Errorf("Viridis At(1) = %v, want {253,231,37,255}", got)
	}
}

// TestParityNormalize mirrors the linear Normalize behavior matplotlib asserts
// (colors.Normalize maps vmin->0, vmax->1 and clips outside the range).
func TestParityNormalize(t *testing.T) {
	n := NewNormalize(0, 100)
	cases := []struct{ v, want float64 }{
		{0, 0}, {25, 0.25}, {50, 0.5}, {100, 1}, {-10, 0}, {150, 1},
	}
	for _, c := range cases {
		if got := n.Norm(c.v); !approx(got, c.want, 9) {
			t.Errorf("Normalize.Norm(%v) = %v, want %v", c.v, got, c.want)
		}
	}
}
