package matplotlib

import "math"

// This file adds color-space conversions and a logarithmic normalization that
// matplotlib exposes through its colors module (matplotlib.colors.rgb_to_hsv,
// hsv_to_rgb and LogNorm). All operations are deterministic and work on
// floating-point channel values in the closed interval [0, 1], matching the
// conventions of the upstream functions so their documented known-answer
// vectors reproduce exactly.

// RGBToHSV converts an RGB color, with each component given as a float in
// [0, 1], to the HSV color space. The returned hue, saturation and value are
// each in [0, 1] (hue is the fraction of a full turn, so 1/3 is green and 2/3
// is blue). It mirrors matplotlib.colors.rgb_to_hsv exactly, including its
// convention that a fully unsaturated color reports a hue of 0.
func RGBToHSV(r, g, b float64) (h, s, v float64) {
	maxc := math.Max(r, math.Max(g, b))
	minc := math.Min(r, math.Min(g, b))
	v = maxc
	delta := maxc - minc
	if maxc != 0 {
		s = delta / maxc
	}
	if delta == 0 {
		return 0, s, v
	}
	switch {
	case r == maxc:
		h = (g - b) / delta
	case g == maxc:
		h = 2 + (b-r)/delta
	default: // b == maxc
		h = 4 + (r-g)/delta
	}
	h = math.Mod(h/6, 1)
	if h < 0 {
		h += 1
	}
	return h, s, v
}

// HSVToRGB converts an HSV color, with each component given as a float in
// [0, 1], back to RGB. The returned red, green and blue are each in [0, 1]. It
// mirrors matplotlib.colors.hsv_to_rgb and is the inverse of RGBToHSV, so
// round-tripping a color through both reproduces the original channels.
func HSVToRGB(h, s, v float64) (r, g, b float64) {
	if s == 0 {
		return v, v, v
	}
	i := int(math.Floor(h * 6))
	f := h*6 - float64(i)
	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))
	switch ((i % 6) + 6) % 6 {
	case 0:
		return v, t, p
	case 1:
		return q, v, p
	case 2:
		return p, v, t
	case 3:
		return p, q, v
	case 4:
		return t, p, v
	default: // 5
		return v, p, q
	}
}

// HSV returns the color's red, green and blue channels converted to HSV, each
// as a float in [0, 1]. The alpha channel is ignored. It is a convenience
// wrapper over RGBToHSV for the 8-bit Color type.
func (c Color) HSV() (h, s, v float64) {
	return RGBToHSV(float64(c.R)/255, float64(c.G)/255, float64(c.B)/255)
}

// ColorFromHSV builds an opaque 8-bit Color from HSV components in [0, 1],
// rounding each resulting channel to the nearest 8-bit value. It is the
// Color-typed counterpart of HSVToRGB.
func ColorFromHSV(h, s, v float64) Color {
	r, g, b := HSVToRGB(h, s, v)
	return Color{R: clamp8(r * 255), G: clamp8(g * 255), B: clamp8(b * 255), A: 255}
}

// LogNorm maps positive data values onto the unit interval [0, 1] on a
// logarithmic (base-10) scale, the counterpart of matplotlib's colors.LogNorm.
// It is typically paired with a Colormap to color data that spans several
// orders of magnitude. Both Min and Max must be strictly positive.
type LogNorm struct {
	Min, Max float64
}

// NewLogNorm returns a LogNorm spanning [min, max]. Both bounds should be
// strictly positive; a non-positive or degenerate range makes Norm return 0.
func NewLogNorm(min, max float64) LogNorm {
	return LogNorm{Min: min, Max: max}
}

// Norm maps v to [0, 1] by the transform
// (log10(v)-log10(Min))/(log10(Max)-log10(Min)), clamped to the unit interval.
// Non-positive v, or a degenerate/invalid range, returns 0. This reproduces
// matplotlib.colors.LogNorm for in-range values.
func (n LogNorm) Norm(v float64) float64 {
	if n.Min <= 0 || n.Max <= 0 || n.Min == n.Max || v <= 0 {
		return 0
	}
	lmin := math.Log10(n.Min)
	lmax := math.Log10(n.Max)
	t := (math.Log10(v) - lmin) / (lmax - lmin)
	if t < 0 {
		return 0
	}
	if t > 1 {
		return 1
	}
	return t
}

// Map normalizes v with n and looks the result up in cm, returning the color a
// logarithmic scalar-mappable would assign to v.
func (n LogNorm) Map(cm Colormap, v float64) Color {
	return cm.At(n.Norm(v))
}
