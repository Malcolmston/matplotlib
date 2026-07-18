package matplotlib

import "math"

// This file adds chart types that bring the library closer to matplotlib's
// feature set: step plots, filled areas between curves, stem plots, error bars,
// reference lines and text annotations. Each implements the internal plotItem
// interface so it composes with the existing Axes rendering, legend and
// autoscaling machinery.

// StepPlot is a piecewise-constant ("staircase") line, the equivalent of
// matplotlib's Axes.step. The Where field selects where the step transition
// occurs relative to the data points: "pre", "post" (the default) or "mid".
type StepPlot struct {
	X, Y      []float64
	Color     Color
	LineWidth float64
	Where     string
	Label     string
}

// SetColor sets the step color and returns the plot for chaining.
func (s *StepPlot) SetColor(c Color) *StepPlot { s.Color = c; return s }

// SetLabel sets the legend label and returns the plot for chaining.
func (s *StepPlot) SetLabel(l string) *StepPlot { s.Label = l; return s }

// SetWhere sets the step placement ("pre", "post" or "mid") and returns the
// plot for chaining. Unrecognized values are treated as "post".
func (s *StepPlot) SetWhere(w string) *StepPlot { s.Where = w; return s }

func (s *StepPlot) bounds() (xlo, xhi, ylo, yhi float64, ok bool) {
	return xyBounds(s.X, s.Y)
}

// stepPoints returns the staircase vertices in data space.
func (s *StepPlot) stepPoints() []Point {
	n := min(len(s.X), len(s.Y))
	if n == 0 {
		return nil
	}
	pts := make([]Point, 0, 2*n)
	switch s.Where {
	case "pre":
		pts = append(pts, Point{s.X[0], s.Y[0]})
		for i := 1; i < n; i++ {
			pts = append(pts, Point{s.X[i-1], s.Y[i]})
			pts = append(pts, Point{s.X[i], s.Y[i]})
		}
	case "mid":
		pts = append(pts, Point{s.X[0], s.Y[0]})
		for i := 1; i < n; i++ {
			mid := (s.X[i-1] + s.X[i]) / 2
			pts = append(pts, Point{mid, s.Y[i-1]})
			pts = append(pts, Point{mid, s.Y[i]})
		}
		pts = append(pts, Point{s.X[n-1], s.Y[n-1]})
	default: // "post"
		for i := 0; i < n-1; i++ {
			pts = append(pts, Point{s.X[i], s.Y[i]})
			pts = append(pts, Point{s.X[i+1], s.Y[i]})
		}
		pts = append(pts, Point{s.X[n-1], s.Y[n-1]})
	}
	return pts
}

func (s *StepPlot) render(cv canvas, ax *Axes) {
	dp := s.stepPoints()
	if len(dp) == 0 {
		return
	}
	pts := make([]Point, len(dp))
	for i, p := range dp {
		pts[i] = Point{ax.tx(p.X), ax.ty(p.Y)}
	}
	lw := s.LineWidth
	if lw <= 0 {
		lw = 1.5
	}
	cv.polyline(pts, s.Color, lw)
}

func (s *StepPlot) legendEntry() (string, Color, bool, Marker, bool) {
	return s.Label, s.Color, true, MarkerNone, s.Label != ""
}

// FillBetween shades the region between two curves y1 and y2 sampled at the
// same x positions, the equivalent of matplotlib's Axes.fill_between. When Y2
// is nil the region is filled down to zero.
type FillBetween struct {
	X, Y1, Y2 []float64
	Color     Color
	Label     string
}

// SetColor sets the fill color and returns the area for chaining. Use
// Color.WithAlpha to make the fill translucent.
func (f *FillBetween) SetColor(c Color) *FillBetween { f.Color = c; return f }

// SetLabel sets the legend label and returns the area for chaining.
func (f *FillBetween) SetLabel(l string) *FillBetween { f.Label = l; return f }

// y2At returns the lower-curve value at index i, treating a nil Y2 as zero.
func (f *FillBetween) y2At(i int) float64 {
	if i < len(f.Y2) {
		return f.Y2[i]
	}
	return 0
}

func (f *FillBetween) bounds() (xlo, xhi, ylo, yhi float64, ok bool) {
	n := min(len(f.X), len(f.Y1))
	if n == 0 {
		return 0, 0, 0, 0, false
	}
	xlo, xhi = f.X[0], f.X[0]
	ylo, yhi = f.Y1[0], f.Y1[0]
	for i := 0; i < n; i++ {
		xlo = math.Min(xlo, f.X[i])
		xhi = math.Max(xhi, f.X[i])
		ylo = math.Min(ylo, math.Min(f.Y1[i], f.y2At(i)))
		yhi = math.Max(yhi, math.Max(f.Y1[i], f.y2At(i)))
	}
	return xlo, xhi, ylo, yhi, true
}

func (f *FillBetween) render(cv canvas, ax *Axes) {
	n := min(len(f.X), len(f.Y1))
	if n < 2 {
		return
	}
	poly := make([]Point, 0, 2*n)
	for i := 0; i < n; i++ {
		poly = append(poly, Point{ax.tx(f.X[i]), ax.ty(f.Y1[i])})
	}
	for i := n - 1; i >= 0; i-- {
		poly = append(poly, Point{ax.tx(f.X[i]), ax.ty(f.y2At(i))})
	}
	cv.polygon(poly, f.Color)
}

func (f *FillBetween) legendEntry() (string, Color, bool, Marker, bool) {
	return f.Label, f.Color, false, MarkerSquare, f.Label != ""
}

// StemPlot draws a vertical stem from a baseline at y = 0 up to each data
// point, with a marker at the tip, mirroring matplotlib's Axes.stem.
type StemPlot struct {
	X, Y   []float64
	Color  Color
	Marker Marker
	Label  string
}

// SetColor sets the stem color and returns the plot for chaining.
func (s *StemPlot) SetColor(c Color) *StemPlot { s.Color = c; return s }

// SetLabel sets the legend label and returns the plot for chaining.
func (s *StemPlot) SetLabel(l string) *StemPlot { s.Label = l; return s }

func (s *StemPlot) bounds() (xlo, xhi, ylo, yhi float64, ok bool) {
	xlo, xhi, ylo, yhi, ok = xyBounds(s.X, s.Y)
	if !ok {
		return
	}
	ylo = math.Min(ylo, 0)
	yhi = math.Max(yhi, 0)
	return
}

func (s *StemPlot) render(cv canvas, ax *Axes) {
	n := min(len(s.X), len(s.Y))
	if n == 0 {
		return
	}
	m := s.Marker
	if m == MarkerNone {
		m = MarkerCircle
	}
	baseY := ax.ty(0)
	// Baseline across the data span.
	cv.strokeLine(ax.tx(s.X[0]), baseY, ax.tx(s.X[n-1]), baseY, axisColor, 1)
	for i := 0; i < n; i++ {
		px := ax.tx(s.X[i])
		cv.strokeLine(px, baseY, px, ax.ty(s.Y[i]), s.Color, 1.5)
		drawMarker(cv, m, px, ax.ty(s.Y[i]), 6, s.Color)
	}
}

func (s *StemPlot) legendEntry() (string, Color, bool, Marker, bool) {
	m := s.Marker
	if m == MarkerNone {
		m = MarkerCircle
	}
	return s.Label, s.Color, true, m, s.Label != ""
}

// ErrorBar draws markers at (x, y) with symmetric vertical error bars of
// half-height yerr[i] and horizontal caps, connected by a line, mirroring
// matplotlib's Axes.errorbar.
type ErrorBar struct {
	X, Y, YErr []float64
	Color      Color
	Marker     Marker
	Label      string
}

// SetColor sets the series color and returns the plot for chaining.
func (e *ErrorBar) SetColor(c Color) *ErrorBar { e.Color = c; return e }

// SetLabel sets the legend label and returns the plot for chaining.
func (e *ErrorBar) SetLabel(l string) *ErrorBar { e.Label = l; return e }

// errAt returns the error magnitude at index i, or zero when unspecified.
func (e *ErrorBar) errAt(i int) float64 {
	if i < len(e.YErr) {
		return math.Abs(e.YErr[i])
	}
	return 0
}

func (e *ErrorBar) bounds() (xlo, xhi, ylo, yhi float64, ok bool) {
	n := min(len(e.X), len(e.Y))
	if n == 0 {
		return 0, 0, 0, 0, false
	}
	xlo, xhi = e.X[0], e.X[0]
	ylo, yhi = e.Y[0], e.Y[0]
	for i := 0; i < n; i++ {
		er := e.errAt(i)
		xlo = math.Min(xlo, e.X[i])
		xhi = math.Max(xhi, e.X[i])
		ylo = math.Min(ylo, e.Y[i]-er)
		yhi = math.Max(yhi, e.Y[i]+er)
	}
	return xlo, xhi, ylo, yhi, true
}

func (e *ErrorBar) render(cv canvas, ax *Axes) {
	n := min(len(e.X), len(e.Y))
	if n == 0 {
		return
	}
	m := e.Marker
	if m == MarkerNone {
		m = MarkerCircle
	}
	line := make([]Point, n)
	const capW = 4.0
	for i := 0; i < n; i++ {
		px := ax.tx(e.X[i])
		py := ax.ty(e.Y[i])
		line[i] = Point{px, py}
		er := e.errAt(i)
		if er > 0 {
			ptop := ax.ty(e.Y[i] + er)
			pbot := ax.ty(e.Y[i] - er)
			cv.strokeLine(px, ptop, px, pbot, e.Color, 1.5)
			cv.strokeLine(px-capW, ptop, px+capW, ptop, e.Color, 1.5)
			cv.strokeLine(px-capW, pbot, px+capW, pbot, e.Color, 1.5)
		}
	}
	cv.polyline(line, e.Color, 1.5)
	for _, p := range line {
		drawMarker(cv, m, p.X, p.Y, 6, e.Color)
	}
}

func (e *ErrorBar) legendEntry() (string, Color, bool, Marker, bool) {
	m := e.Marker
	if m == MarkerNone {
		m = MarkerCircle
	}
	return e.Label, e.Color, true, m, e.Label != ""
}

// RefLine is a full-width horizontal or full-height vertical reference line at
// a fixed data coordinate, the equivalent of matplotlib's axhline and axvline.
// A reference line does not participate in autoscaling.
type RefLine struct {
	Value     float64
	Vertical  bool
	Color     Color
	LineWidth float64
	Label     string
}

// SetColor sets the line color and returns the reference line for chaining.
func (r *RefLine) SetColor(c Color) *RefLine { r.Color = c; return r }

// SetLabel sets the legend label and returns the reference line for chaining.
func (r *RefLine) SetLabel(l string) *RefLine { r.Label = l; return r }

func (r *RefLine) bounds() (xlo, xhi, ylo, yhi float64, ok bool) {
	return 0, 0, 0, 0, false
}

func (r *RefLine) render(cv canvas, ax *Axes) {
	lw := r.LineWidth
	if lw <= 0 {
		lw = 1.5
	}
	if r.Vertical {
		if r.Value < ax.xmin || r.Value > ax.xmax {
			return
		}
		px := ax.tx(r.Value)
		cv.strokeLine(px, ax.py0, px, ax.py1, r.Color, lw)
	} else {
		if r.Value < ax.ymin || r.Value > ax.ymax {
			return
		}
		py := ax.ty(r.Value)
		cv.strokeLine(ax.px0, py, ax.px1, py, r.Color, lw)
	}
}

func (r *RefLine) legendEntry() (string, Color, bool, Marker, bool) {
	return r.Label, r.Color, true, MarkerNone, r.Label != ""
}

// Annotation is a text label placed at a data-space coordinate, the equivalent
// of matplotlib's Axes.text. It does not participate in autoscaling.
type Annotation struct {
	X, Y   float64
	Text   string
	Color  Color
	Size   float64
	Anchor TextAnchor
}

func (a *Annotation) bounds() (xlo, xhi, ylo, yhi float64, ok bool) {
	return 0, 0, 0, 0, false
}

func (a *Annotation) render(cv canvas, ax *Axes) {
	if a.Text == "" {
		return
	}
	size := a.Size
	if size <= 0 {
		size = 12
	}
	col := a.Color
	if col.A == 0 {
		col = textColor
	}
	cv.drawText(ax.tx(a.X), ax.ty(a.Y), a.Text, col, a.Anchor, size, false)
}

func (a *Annotation) legendEntry() (string, Color, bool, Marker, bool) {
	return "", Color{}, false, MarkerNone, false
}

// Step adds a piecewise-constant step series through the (x[i], y[i]) points.
func (ax *Axes) Step(x, y []float64) *StepPlot {
	s := &StepPlot{X: x, Y: y, Color: ax.fig.nextColor(), LineWidth: 1.5, Where: "post"}
	ax.items = append(ax.items, s)
	return s
}

// FillBetween shades the area between the curves y1 and y2 over the x
// positions. Passing nil for y2 fills down to zero.
func (ax *Axes) FillBetween(x, y1, y2 []float64) *FillBetween {
	f := &FillBetween{X: x, Y1: y1, Y2: y2, Color: ax.fig.nextColor().WithAlpha(160)}
	ax.items = append(ax.items, f)
	return f
}

// Stem adds a stem plot with vertical lines from y = 0 to each (x[i], y[i]).
func (ax *Axes) Stem(x, y []float64) *StemPlot {
	s := &StemPlot{X: x, Y: y, Color: ax.fig.nextColor(), Marker: MarkerCircle}
	ax.items = append(ax.items, s)
	return s
}

// ErrorBar adds a series at (x[i], y[i]) with symmetric vertical error bars of
// half-height yerr[i]. yerr may be shorter than x/y; missing entries are zero.
func (ax *Axes) ErrorBar(x, y, yerr []float64) *ErrorBar {
	e := &ErrorBar{X: x, Y: y, YErr: yerr, Color: ax.fig.nextColor(), Marker: MarkerCircle}
	ax.items = append(ax.items, e)
	return e
}

// AxHLine adds a horizontal reference line spanning the axes at data y.
func (ax *Axes) AxHLine(y float64) *RefLine {
	r := &RefLine{Value: y, Vertical: false, Color: axisColor, LineWidth: 1}
	ax.items = append(ax.items, r)
	return r
}

// AxVLine adds a vertical reference line spanning the axes at data x.
func (ax *Axes) AxVLine(x float64) *RefLine {
	r := &RefLine{Value: x, Vertical: true, Color: axisColor, LineWidth: 1}
	ax.items = append(ax.items, r)
	return r
}

// Text places the given string at data-space coordinate (x, y) on the axes.
func (ax *Axes) Text(x, y float64, s string) *Annotation {
	a := &Annotation{X: x, Y: y, Text: s, Color: textColor, Size: 12, Anchor: AnchorStart}
	ax.items = append(ax.items, a)
	return a
}
