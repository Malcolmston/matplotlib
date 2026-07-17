package matplotlib

import (
	"math"
	"strconv"
	"strings"
)

// Point is a 2-D coordinate in device (pixel) space, with the origin at the
// top-left corner and y increasing downward.
type Point struct {
	X, Y float64
}

// TextAnchor controls horizontal alignment of rendered text relative to its
// anchor point.
type TextAnchor int

const (
	// AnchorStart left-aligns text at the anchor point.
	AnchorStart TextAnchor = iota
	// AnchorMiddle centers text on the anchor point.
	AnchorMiddle
	// AnchorEnd right-aligns text at the anchor point.
	AnchorEnd
)

// Marker identifies the glyph drawn at data points for line and scatter plots.
type Marker int

const (
	// MarkerNone draws no marker.
	MarkerNone Marker = iota
	// MarkerCircle draws a filled circle.
	MarkerCircle
	// MarkerSquare draws a filled square.
	MarkerSquare
	// MarkerTriangle draws a filled upward triangle.
	MarkerTriangle
	// MarkerCross draws an "x" shaped cross.
	MarkerCross
	// MarkerPlus draws a "+" shaped cross.
	MarkerPlus
)

// canvas is the drawing surface abstraction implemented by both the PNG raster
// backend and the SVG vector backend. All coordinates are in pixels.
type canvas interface {
	size() (w, h int)
	fillBackground(c Color)
	fillRect(x, y, w, h float64, c Color)
	strokeLine(x0, y0, x1, y1 float64, c Color, lw float64)
	polyline(pts []Point, c Color, lw float64)
	polygon(pts []Point, fill Color)
	fillCircle(cx, cy, r float64, c Color)
	drawText(x, y float64, s string, c Color, anchor TextAnchor, size float64, vertical bool)
}

// drawMarker renders marker m centered at (x,y) with the given size (full
// width/height in pixels) using primitives available on any canvas.
func drawMarker(cv canvas, m Marker, x, y, size float64, c Color) {
	h := size / 2
	switch m {
	case MarkerNone:
		return
	case MarkerCircle:
		cv.fillCircle(x, y, h, c)
	case MarkerSquare:
		cv.fillRect(x-h, y-h, size, size, c)
	case MarkerTriangle:
		cv.polygon([]Point{{x, y - h}, {x + h, y + h}, {x - h, y + h}}, c)
	case MarkerCross:
		cv.strokeLine(x-h, y-h, x+h, y+h, c, 1.5)
		cv.strokeLine(x-h, y+h, x+h, y-h, c, 1.5)
	case MarkerPlus:
		cv.strokeLine(x-h, y, x+h, y, c, 1.5)
		cv.strokeLine(x, y-h, x, y+h, c, 1.5)
	}
}

// fnum formats a float for SVG attributes, deterministically, with up to two
// decimal places and trailing zeros trimmed.
func fnum(v float64) string {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return "0"
	}
	s := strconv.FormatFloat(v, 'f', 2, 64)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	if s == "-0" {
		s = "0"
	}
	return s
}

// niceNum returns a "nice" number approximately equal to x. If round is true
// the value is rounded, otherwise the ceiling is taken. This is the classic
// Graphics Gems algorithm for labeling graph axes.
func niceNum(x float64, round bool) float64 {
	if x <= 0 {
		return 1
	}
	exp := math.Floor(math.Log10(x))
	f := x / math.Pow(10, exp)
	var nf float64
	if round {
		switch {
		case f < 1.5:
			nf = 1
		case f < 3:
			nf = 2
		case f < 7:
			nf = 5
		default:
			nf = 10
		}
	} else {
		switch {
		case f <= 1:
			nf = 1
		case f <= 2:
			nf = 2
		case f <= 5:
			nf = 5
		default:
			nf = 10
		}
	}
	return nf * math.Pow(10, exp)
}

// niceTicks returns approximately n evenly spaced "nice" tick values covering
// the range [lo, hi]. The returned slice always has at least two entries.
func niceTicks(lo, hi float64, n int) []float64 {
	if n < 2 {
		n = 2
	}
	if hi < lo {
		lo, hi = hi, lo
	}
	if hi == lo {
		hi = lo + 1
	}
	rng := niceNum(hi-lo, false)
	d := niceNum(rng/float64(n-1), true)
	glo := math.Floor(lo/d) * d
	ghi := math.Ceil(hi/d) * d
	var t []float64
	for v := glo; v < ghi+0.5*d; v += d {
		if math.Abs(v) < d*1e-9 {
			v = 0
		}
		t = append(t, v)
	}
	if len(t) < 2 {
		t = []float64{glo, glo + d}
	}
	return t
}

// fmtTick formats a numeric tick label deterministically.
func fmtTick(v float64) string {
	if v == 0 {
		return "0"
	}
	a := math.Abs(v)
	if a >= 1e5 || a < 1e-4 {
		return strconv.FormatFloat(v, 'e', 1, 64)
	}
	s := strconv.FormatFloat(v, 'f', 4, 64)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}
	return s
}
