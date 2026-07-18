package matplotlib

// This file adds color manipulation helpers matplotlib exposes through its
// colors module: lightening, darkening, alpha adjustment, luminance and
// blending. All operations are per-channel and deterministic.

// clamp8 rounds and clamps a float to the 0..255 range as a uint8.
func clamp8(v float64) uint8 {
	if v < 0 {
		v = 0
	}
	if v > 255 {
		v = 255
	}
	return uint8(v + 0.5)
}

// Lighten returns a copy of the color moved a fraction f of the way toward
// white, leaving alpha unchanged. f is clamped to [0, 1]; f == 0 returns the
// color unchanged and f == 1 returns white (at the original alpha).
func (c Color) Lighten(f float64) Color {
	if f < 0 {
		f = 0
	}
	if f > 1 {
		f = 1
	}
	return Color{
		R: clamp8(float64(c.R) + (255-float64(c.R))*f),
		G: clamp8(float64(c.G) + (255-float64(c.G))*f),
		B: clamp8(float64(c.B) + (255-float64(c.B))*f),
		A: c.A,
	}
}

// Darken returns a copy of the color scaled a fraction f of the way toward
// black, leaving alpha unchanged. f is clamped to [0, 1]; f == 0 returns the
// color unchanged and f == 1 returns black (at the original alpha).
func (c Color) Darken(f float64) Color {
	if f < 0 {
		f = 0
	}
	if f > 1 {
		f = 1
	}
	s := 1 - f
	return Color{
		R: clamp8(float64(c.R) * s),
		G: clamp8(float64(c.G) * s),
		B: clamp8(float64(c.B) * s),
		A: c.A,
	}
}

// WithAlpha returns a copy of the color with its alpha channel replaced by a.
func (c Color) WithAlpha(a uint8) Color {
	c.A = a
	return c
}

// Luminance returns the relative luminance of the color in [0, 1] using the
// Rec. 709 weights (0.2126 R + 0.7152 G + 0.0722 B) over the sRGB channel
// values. Black returns 0 and white returns 1; alpha is ignored.
func (c Color) Luminance() float64 {
	return (0.2126*float64(c.R) + 0.7152*float64(c.G) + 0.0722*float64(c.B)) / 255
}

// Grayscale returns the color converted to a shade of gray whose intensity is
// its luminance, preserving the original alpha.
func (c Color) Grayscale() Color {
	g := clamp8(c.Luminance() * 255)
	return Color{R: g, G: g, B: g, A: c.A}
}

// Hex returns the color as a "#rrggbb" string (the alpha channel is not
// included), suitable for CSS or SVG.
func (c Color) Hex() string {
	return c.hexString()
}

// Blend linearly interpolates between colors a and b, returning a when t <= 0,
// b when t >= 1 and a per-channel mix (including alpha) in between.
func Blend(a, b Color, t float64) Color {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return Color{
		R: clamp8(float64(a.R) + (float64(b.R)-float64(a.R))*t),
		G: clamp8(float64(a.G) + (float64(b.G)-float64(a.G))*t),
		B: clamp8(float64(a.B) + (float64(b.B)-float64(a.B))*t),
		A: clamp8(float64(a.A) + (float64(b.A)-float64(a.A))*t),
	}
}
