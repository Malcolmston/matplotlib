package matplotlib

import "math"

// plotItem is the internal interface implemented by every chart type. It knows
// how to report its data-space bounds, render itself onto a canvas using the
// axes coordinate transforms, and describe its legend entry.
type plotItem interface {
	bounds() (xlo, xhi, ylo, yhi float64, ok bool)
	render(cv canvas, ax *Axes)
	legendEntry() (label string, col Color, line bool, marker Marker, ok bool)
}

// LinePlot is a connected line series with optional markers.
type LinePlot struct {
	X, Y      []float64
	Color     Color
	LineWidth float64
	Marker    Marker
	Label     string
}

// SetColor sets the series color and returns the plot for chaining.
func (l *LinePlot) SetColor(c Color) *LinePlot { l.Color = c; return l }

// SetLabel sets the legend label and returns the plot for chaining.
func (l *LinePlot) SetLabel(s string) *LinePlot { l.Label = s; return l }

// SetLineWidth sets the line width in pixels and returns the plot for chaining.
func (l *LinePlot) SetLineWidth(w float64) *LinePlot { l.LineWidth = w; return l }

// SetMarker sets the point marker and returns the plot for chaining.
func (l *LinePlot) SetMarker(m Marker) *LinePlot { l.Marker = m; return l }

func (l *LinePlot) bounds() (xlo, xhi, ylo, yhi float64, ok bool) {
	return xyBounds(l.X, l.Y)
}

func (l *LinePlot) render(cv canvas, ax *Axes) {
	n := min(len(l.X), len(l.Y))
	pts := make([]Point, 0, n)
	for i := 0; i < n; i++ {
		pts = append(pts, Point{ax.tx(l.X[i]), ax.ty(l.Y[i])})
	}
	lw := l.LineWidth
	if lw <= 0 {
		lw = 1.5
	}
	cv.polyline(pts, l.Color, lw)
	if l.Marker != MarkerNone {
		for _, p := range pts {
			drawMarker(cv, l.Marker, p.X, p.Y, 6, l.Color)
		}
	}
}

func (l *LinePlot) legendEntry() (string, Color, bool, Marker, bool) {
	return l.Label, l.Color, true, l.Marker, l.Label != ""
}

// ScatterPlot is a set of unconnected markers.
type ScatterPlot struct {
	X, Y   []float64
	Color  Color
	Size   float64
	Marker Marker
	Label  string
}

// SetColor sets the marker color and returns the plot for chaining.
func (s *ScatterPlot) SetColor(c Color) *ScatterPlot { s.Color = c; return s }

// SetLabel sets the legend label and returns the plot for chaining.
func (s *ScatterPlot) SetLabel(l string) *ScatterPlot { s.Label = l; return s }

// SetSize sets the marker size in pixels and returns the plot for chaining.
func (s *ScatterPlot) SetSize(sz float64) *ScatterPlot { s.Size = sz; return s }

// SetMarker sets the marker glyph and returns the plot for chaining.
func (s *ScatterPlot) SetMarker(m Marker) *ScatterPlot { s.Marker = m; return s }

func (s *ScatterPlot) bounds() (xlo, xhi, ylo, yhi float64, ok bool) {
	return xyBounds(s.X, s.Y)
}

func (s *ScatterPlot) render(cv canvas, ax *Axes) {
	n := min(len(s.X), len(s.Y))
	sz := s.Size
	if sz <= 0 {
		sz = 7
	}
	m := s.Marker
	if m == MarkerNone {
		m = MarkerCircle
	}
	for i := 0; i < n; i++ {
		drawMarker(cv, m, ax.tx(s.X[i]), ax.ty(s.Y[i]), sz, s.Color)
	}
}

func (s *ScatterPlot) legendEntry() (string, Color, bool, Marker, bool) {
	m := s.Marker
	if m == MarkerNone {
		m = MarkerCircle
	}
	return s.Label, s.Color, false, m, s.Label != ""
}

// BarChart is a categorical bar series. When Horizontal is true bars extend
// along the x axis and categories are placed on the y axis.
type BarChart struct {
	Labels     []string
	Values     []float64
	Color      Color
	Horizontal bool
	Label      string
}

// SetColor sets the bar color and returns the chart for chaining.
func (b *BarChart) SetColor(c Color) *BarChart { b.Color = c; return b }

// SetLabel sets the legend label and returns the chart for chaining.
func (b *BarChart) SetLabel(s string) *BarChart { b.Label = s; return b }

func (b *BarChart) bounds() (xlo, xhi, ylo, yhi float64, ok bool) {
	if len(b.Values) == 0 {
		return 0, 0, 0, 0, false
	}
	vlo, vhi := b.Values[0], b.Values[0]
	for _, v := range b.Values {
		vlo = math.Min(vlo, v)
		vhi = math.Max(vhi, v)
	}
	vlo = math.Min(vlo, 0)
	vhi = math.Max(vhi, 0)
	n := float64(len(b.Values))
	if b.Horizontal {
		return vlo, vhi, -0.5, n - 0.5, true
	}
	return -0.5, n - 0.5, vlo, vhi, true
}

func (b *BarChart) render(cv canvas, ax *Axes) {
	const barFrac = 0.8
	for i, v := range b.Values {
		center := float64(i)
		if b.Horizontal {
			y0 := ax.ty(center - barFrac/2)
			y1 := ax.ty(center + barFrac/2)
			x0 := ax.tx(0)
			x1 := ax.tx(v)
			cv.fillRect(math.Min(x0, x1), math.Min(y0, y1), math.Abs(x1-x0), math.Abs(y1-y0), b.Color)
		} else {
			x0 := ax.tx(center - barFrac/2)
			x1 := ax.tx(center + barFrac/2)
			y0 := ax.ty(0)
			y1 := ax.ty(v)
			cv.fillRect(math.Min(x0, x1), math.Min(y0, y1), math.Abs(x1-x0), math.Abs(y1-y0), b.Color)
		}
	}
}

func (b *BarChart) legendEntry() (string, Color, bool, Marker, bool) {
	return b.Label, b.Color, false, MarkerSquare, b.Label != ""
}

// Histogram bins raw data into equal-width bins and draws the counts as bars.
type Histogram struct {
	Data  []float64
	Bins  int
	Color Color
	Label string

	edges  []float64
	counts []int
}

// SetColor sets the histogram color and returns it for chaining.
func (h *Histogram) SetColor(c Color) *Histogram { h.Color = c; return h }

// SetLabel sets the legend label and returns the histogram for chaining.
func (h *Histogram) SetLabel(s string) *Histogram { h.Label = s; return h }

// compute fills the edges and counts fields based on Data and Bins.
func (h *Histogram) compute() {
	if len(h.edges) != 0 {
		return
	}
	bins := h.Bins
	if bins < 1 {
		bins = 10
	}
	if len(h.Data) == 0 {
		h.edges = []float64{0, 1}
		h.counts = []int{0}
		return
	}
	lo, hi := h.Data[0], h.Data[0]
	for _, v := range h.Data {
		lo = math.Min(lo, v)
		hi = math.Max(hi, v)
	}
	if hi == lo {
		hi = lo + 1
	}
	width := (hi - lo) / float64(bins)
	h.edges = make([]float64, bins+1)
	for i := 0; i <= bins; i++ {
		h.edges[i] = lo + float64(i)*width
	}
	h.counts = make([]int, bins)
	for _, v := range h.Data {
		idx := int((v - lo) / width)
		if idx >= bins {
			idx = bins - 1
		}
		if idx < 0 {
			idx = 0
		}
		h.counts[idx]++
	}
}

func (h *Histogram) bounds() (xlo, xhi, ylo, yhi float64, ok bool) {
	h.compute()
	maxC := 0
	for _, c := range h.counts {
		if c > maxC {
			maxC = c
		}
	}
	return h.edges[0], h.edges[len(h.edges)-1], 0, float64(maxC), true
}

func (h *Histogram) render(cv canvas, ax *Axes) {
	h.compute()
	for i, c := range h.counts {
		x0 := ax.tx(h.edges[i])
		x1 := ax.tx(h.edges[i+1])
		y0 := ax.ty(0)
		y1 := ax.ty(float64(c))
		cv.fillRect(math.Min(x0, x1), math.Min(y0, y1), math.Abs(x1-x0), math.Abs(y1-y0), h.Color)
		cv.strokeLine(x0, y0, x0, y1, White, 1)
		cv.strokeLine(x1, y0, x1, y1, White, 1)
	}
}

func (h *Histogram) legendEntry() (string, Color, bool, Marker, bool) {
	return h.Label, h.Color, false, MarkerSquare, h.Label != ""
}

// PieChart draws proportional wedges. It does not use the axis coordinate
// system; the enclosing axes renders it specially.
type PieChart struct {
	Values []float64
	Labels []string
	Colors []Color
}

func (p *PieChart) bounds() (xlo, xhi, ylo, yhi float64, ok bool) {
	return 0, 0, 0, 0, false
}

func (p *PieChart) render(cv canvas, ax *Axes) {}

func (p *PieChart) legendEntry() (string, Color, bool, Marker, bool) {
	return "", Color{}, false, MarkerNone, false
}

// xyBounds computes the min/max of paired x/y slices.
func xyBounds(xs, ys []float64) (xlo, xhi, ylo, yhi float64, ok bool) {
	n := min(len(xs), len(ys))
	if n == 0 {
		return 0, 0, 0, 0, false
	}
	xlo, xhi = xs[0], xs[0]
	ylo, yhi = ys[0], ys[0]
	for i := 0; i < n; i++ {
		xlo = math.Min(xlo, xs[i])
		xhi = math.Max(xhi, xs[i])
		ylo = math.Min(ylo, ys[i])
		yhi = math.Max(yhi, ys[i])
	}
	return xlo, xhi, ylo, yhi, true
}
