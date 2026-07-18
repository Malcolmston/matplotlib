package matplotlib

import "math"

// Colormap maps a scalar in the closed interval [0, 1] to a Color by piecewise
// linear interpolation between a sequence of evenly spaced anchor colors. It is
// the equivalent of matplotlib's LinearSegmentedColormap and lets scalar data
// (heatmap values, series indices, gradients) be turned into colors.
//
// The zero value is not usable; construct a Colormap with NewColormap or one of
// the named-colormap constructors such as Viridis.
type Colormap struct {
	stops []Color
}

// NewColormap builds a Colormap from two or more anchor colors. The anchors are
// treated as being equally spaced across [0, 1], and Colormap.At interpolates
// linearly between neighbouring anchors. Fewer than two colors are padded so
// the result is always usable; the input slice is copied.
func NewColormap(colors []Color) Colormap {
	cp := make([]Color, len(colors))
	copy(cp, colors)
	switch len(cp) {
	case 0:
		cp = []Color{Black, White}
	case 1:
		cp = []Color{cp[0], cp[0]}
	}
	return Colormap{stops: cp}
}

// At returns the color at position t. Values of t are clamped to [0, 1], so
// t <= 0 yields the first anchor and t >= 1 yields the last. Interpolation is
// performed independently per channel (including alpha) and rounded to the
// nearest 8-bit value.
func (m Colormap) At(t float64) Color {
	n := len(m.stops)
	if n == 0 {
		return Black
	}
	if math.IsNaN(t) || t <= 0 {
		return m.stops[0]
	}
	if t >= 1 {
		return m.stops[n-1]
	}
	// Locate the segment: there are n-1 segments across [0,1].
	scaled := t * float64(n-1)
	i := int(math.Floor(scaled))
	if i >= n-1 {
		i = n - 2
	}
	frac := scaled - float64(i)
	a, b := m.stops[i], m.stops[i+1]
	return Color{
		R: lerp8(a.R, b.R, frac),
		G: lerp8(a.G, b.G, frac),
		B: lerp8(a.B, b.B, frac),
		A: lerp8(a.A, b.A, frac),
	}
}

// Reversed returns a new Colormap whose anchor colors are in the opposite
// order, equivalent to appending "_r" to a matplotlib colormap name.
func (m Colormap) Reversed() Colormap {
	n := len(m.stops)
	rev := make([]Color, n)
	for i, c := range m.stops {
		rev[n-1-i] = c
	}
	return Colormap{stops: rev}
}

// lerp8 linearly interpolates between two 8-bit channel values, rounding to the
// nearest integer.
func lerp8(a, b uint8, t float64) uint8 {
	v := float64(a) + (float64(b)-float64(a))*t
	if v < 0 {
		v = 0
	}
	if v > 255 {
		v = 255
	}
	return uint8(v + 0.5)
}

// Normalize maps data values in a fixed range onto the unit interval [0, 1],
// the counterpart of matplotlib's colors.Normalize. It is typically paired with
// a Colormap to color scalar data.
type Normalize struct {
	Min, Max float64
}

// NewNormalize returns a Normalize spanning [min, max]. If min == max the range
// is treated as degenerate and Norm always returns 0.
func NewNormalize(min, max float64) Normalize {
	return Normalize{Min: min, Max: max}
}

// Norm maps v to [0, 1] by the linear transform (v-Min)/(Max-Min), clamped to
// the unit interval. A degenerate range (Min == Max) returns 0.
func (n Normalize) Norm(v float64) float64 {
	if n.Max == n.Min {
		return 0
	}
	t := (v - n.Min) / (n.Max - n.Min)
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return t
}

// Map normalizes v with n and looks the result up in cm, returning the color
// that a scalar-mappable would assign to v.
func (n Normalize) Map(cm Colormap, v float64) Color {
	return cm.At(n.Norm(v))
}

// Viridis returns matplotlib's default perceptually uniform colormap,
// approximated by five anchor colors from dark purple through teal and green to
// yellow.
func Viridis() Colormap {
	return NewColormap([]Color{
		{68, 1, 84, 255},
		{59, 82, 139, 255},
		{33, 145, 140, 255},
		{94, 201, 98, 255},
		{253, 231, 37, 255},
	})
}

// Plasma returns an approximation of matplotlib's "plasma" colormap, running
// from dark blue through magenta and orange to yellow.
func Plasma() Colormap {
	return NewColormap([]Color{
		{13, 8, 135, 255},
		{126, 3, 168, 255},
		{204, 71, 120, 255},
		{248, 149, 64, 255},
		{240, 249, 33, 255},
	})
}

// Grays returns a simple grayscale colormap from black at 0 to white at 1.
func Grays() Colormap {
	return NewColormap([]Color{Black, White})
}

// Jet returns an approximation of the classic "jet" colormap, running from dark
// blue through cyan, yellow and red to dark red.
func Jet() Colormap {
	return NewColormap([]Color{
		{0, 0, 131, 255},
		{0, 60, 255, 255},
		{0, 255, 255, 255},
		{255, 255, 0, 255},
		{255, 60, 0, 255},
		{131, 0, 0, 255},
	})
}

// Coolwarm returns an approximation of matplotlib's diverging "coolwarm"
// colormap, running from cool blue through neutral gray to warm red.
func Coolwarm() Colormap {
	return NewColormap([]Color{
		{59, 76, 192, 255},
		{221, 221, 221, 255},
		{180, 4, 38, 255},
	})
}
