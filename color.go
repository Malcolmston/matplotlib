package matplotlib

import (
	"fmt"
	"image/color"
)

// Color is an 8-bit-per-channel RGBA color used throughout the package.
// The zero value is fully transparent black.
type Color struct {
	R, G, B, A uint8
}

// RGB returns an opaque color with the given red, green and blue components.
func RGB(r, g, b uint8) Color { return Color{r, g, b, 255} }

// RGBA returns a color with the given components (alpha included).
func RGBA(r, g, b, a uint8) Color { return Color{r, g, b, a} }

// Hex parses a color from a "#rrggbb" or "#rgb" style string. On any parse
// error it returns opaque black and a non-nil error.
func Hex(s string) (Color, error) {
	if len(s) > 0 && s[0] == '#' {
		s = s[1:]
	}
	switch len(s) {
	case 3:
		var r, g, b uint8
		if _, err := fmt.Sscanf(s, "%1x%1x%1x", &r, &g, &b); err != nil {
			return Color{0, 0, 0, 255}, err
		}
		return Color{r * 17, g * 17, b * 17, 255}, nil
	case 6:
		var r, g, b uint8
		if _, err := fmt.Sscanf(s, "%02x%02x%02x", &r, &g, &b); err != nil {
			return Color{0, 0, 0, 255}, err
		}
		return Color{r, g, b, 255}, nil
	default:
		return Color{0, 0, 0, 255}, fmt.Errorf("matplotlib: invalid hex color %q", s)
	}
}

// RGBA implements the standard library image/color.Color interface so a Color
// can be used directly with the image package.
func (c Color) RGBA() (r, g, b, a uint32) {
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}.RGBA()
}

// hexString renders the color as a "#rrggbb" string for SVG output.
func (c Color) hexString() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
}

// isTransparent reports whether the color is fully transparent.
func (c Color) isTransparent() bool { return c.A == 0 }

// Common named colors.
var (
	Black       = Color{0, 0, 0, 255}
	White       = Color{255, 255, 255, 255}
	Red         = Color{214, 39, 40, 255}
	Green       = Color{44, 160, 44, 255}
	Blue        = Color{31, 119, 180, 255}
	Transparent = Color{0, 0, 0, 0}
	gridColor   = Color{221, 221, 221, 255}
	axisColor   = Color{51, 51, 51, 255}
	textColor   = Color{34, 34, 34, 255}
)

// DefaultColors is the default color cycle, matching matplotlib's "tab10"
// categorical palette. New series are assigned successive colors from this
// slice, wrapping around when exhausted.
var DefaultColors = []Color{
	{31, 119, 180, 255},
	{255, 127, 14, 255},
	{44, 160, 44, 255},
	{214, 39, 40, 255},
	{148, 103, 189, 255},
	{140, 86, 75, 255},
	{227, 119, 194, 255},
	{127, 127, 127, 255},
	{188, 189, 34, 255},
	{23, 190, 207, 255},
}
