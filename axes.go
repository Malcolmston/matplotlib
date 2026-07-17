package matplotlib

import (
	"math"
)

// Axes is a single set of coordinate axes within a Figure. It holds the plot
// items (series) drawn against it together with titles, labels and display
// options. Create one with Figure.AddAxes or Figure.Subplots.
type Axes struct {
	Title      string
	XLabel     string
	YLabel     string
	ShowLegend bool
	ShowGrid   bool

	fig   *Figure
	pos   [4]float64 // normalized [x0,y0,x1,y1] within the figure
	items []plotItem

	categories    []string // non-nil => categorical axis
	horizontalCat bool     // categories on the y axis (horizontal bars)
	isPie         bool

	xlimSet, ylimSet       bool
	xminUser, xmaxUser     float64
	yminUser, ymaxUser     float64
	xmin, xmax, ymin, ymax float64
	xticks, yticks         []float64
	xtickLabels            []string
	ytickLabels            []string

	// plot rectangle in pixels, filled during render
	px0, py0, px1, py1 float64
}

// SetTitle sets the axes title and returns the axes for chaining.
func (ax *Axes) SetTitle(s string) *Axes { ax.Title = s; return ax }

// SetXLabel sets the x-axis label and returns the axes for chaining.
func (ax *Axes) SetXLabel(s string) *Axes { ax.XLabel = s; return ax }

// SetYLabel sets the y-axis label and returns the axes for chaining.
func (ax *Axes) SetYLabel(s string) *Axes { ax.YLabel = s; return ax }

// Legend enables the legend and returns the axes for chaining.
func (ax *Axes) Legend() *Axes { ax.ShowLegend = true; return ax }

// Grid enables or disables the background grid and returns the axes.
func (ax *Axes) Grid(on bool) *Axes { ax.ShowGrid = on; return ax }

// SetXLim fixes the x-axis limits and returns the axes for chaining.
func (ax *Axes) SetXLim(lo, hi float64) *Axes {
	ax.xlimSet, ax.xminUser, ax.xmaxUser = true, lo, hi
	return ax
}

// SetYLim fixes the y-axis limits and returns the axes for chaining.
func (ax *Axes) SetYLim(lo, hi float64) *Axes {
	ax.ylimSet, ax.yminUser, ax.ymaxUser = true, lo, hi
	return ax
}

// Plot adds a line series through the (x[i], y[i]) points.
func (ax *Axes) Plot(x, y []float64) *LinePlot {
	l := &LinePlot{X: x, Y: y, Color: ax.fig.nextColor(), LineWidth: 1.5}
	ax.items = append(ax.items, l)
	return l
}

// Scatter adds an unconnected marker series at the (x[i], y[i]) points.
func (ax *Axes) Scatter(x, y []float64) *ScatterPlot {
	s := &ScatterPlot{X: x, Y: y, Color: ax.fig.nextColor(), Size: 7, Marker: MarkerCircle}
	ax.items = append(ax.items, s)
	return s
}

// Bar adds a vertical bar chart with the given category labels and values.
func (ax *Axes) Bar(labels []string, values []float64) *BarChart {
	b := &BarChart{Labels: labels, Values: values, Color: ax.fig.nextColor()}
	ax.items = append(ax.items, b)
	ax.categories = labels
	ax.horizontalCat = false
	return b
}

// BarH adds a horizontal bar chart with the given category labels and values.
func (ax *Axes) BarH(labels []string, values []float64) *BarChart {
	b := &BarChart{Labels: labels, Values: values, Color: ax.fig.nextColor(), Horizontal: true}
	ax.items = append(ax.items, b)
	ax.categories = labels
	ax.horizontalCat = true
	return b
}

// Hist adds a histogram of data using the given number of bins (<=0 selects a
// default of 10 bins).
func (ax *Axes) Hist(data []float64, bins int) *Histogram {
	h := &Histogram{Data: data, Bins: bins, Color: ax.fig.nextColor()}
	ax.items = append(ax.items, h)
	return h
}

// Pie adds a pie chart. Labels may be nil. Colors are taken from the figure's
// color cycle.
func (ax *Axes) Pie(values []float64, labels []string) *PieChart {
	cols := make([]Color, len(values))
	for i := range values {
		cols[i] = ax.fig.nextColor()
	}
	p := &PieChart{Values: values, Labels: labels, Colors: cols}
	ax.items = append(ax.items, p)
	ax.isPie = true
	return p
}

// tx maps a data-space x value to a pixel x coordinate.
func (ax *Axes) tx(x float64) float64 {
	if ax.xmax == ax.xmin {
		return ax.px0
	}
	return ax.px0 + (x-ax.xmin)/(ax.xmax-ax.xmin)*(ax.px1-ax.px0)
}

// ty maps a data-space y value to a pixel y coordinate (y grows downward).
func (ax *Axes) ty(y float64) float64 {
	if ax.ymax == ax.ymin {
		return ax.py1
	}
	return ax.py1 - (y-ax.ymin)/(ax.ymax-ax.ymin)*(ax.py1-ax.py0)
}

// combinedBounds merges the data bounds of every item.
func (ax *Axes) combinedBounds() (xlo, xhi, ylo, yhi float64, ok bool) {
	first := true
	for _, it := range ax.items {
		bxlo, bxhi, bylo, byhi, bok := it.bounds()
		if !bok {
			continue
		}
		if first {
			xlo, xhi, ylo, yhi = bxlo, bxhi, bylo, byhi
			first = false
			continue
		}
		xlo = math.Min(xlo, bxlo)
		xhi = math.Max(xhi, bxhi)
		ylo = math.Min(ylo, bylo)
		yhi = math.Max(yhi, byhi)
	}
	return xlo, xhi, ylo, yhi, !first
}

// finalize computes axis limits, ticks and tick labels for the current items.
func (ax *Axes) finalize() {
	if ax.isPie {
		return
	}
	bxlo, bxhi, bylo, byhi, ok := ax.combinedBounds()
	if !ok {
		bxlo, bxhi, bylo, byhi = 0, 1, 0, 1
	}

	if ax.categories != nil && !ax.horizontalCat {
		n := len(ax.categories)
		ax.xmin, ax.xmax = -0.5, float64(n)-0.5
		ax.xticks = intTicks(n)
		ax.xtickLabels = ax.categories
		ax.setNumericY(bylo, byhi)
		return
	}
	if ax.horizontalCat {
		n := len(ax.categories)
		ax.ymin, ax.ymax = -0.5, float64(n)-0.5
		ax.yticks = intTicks(n)
		ax.ytickLabels = ax.categories
		ax.setNumericX(bxlo, bxhi)
		return
	}
	ax.setNumericX(bxlo, bxhi)
	ax.setNumericY(bylo, byhi)
}

func (ax *Axes) setNumericX(lo, hi float64) {
	if ax.xlimSet {
		ax.xmin, ax.xmax = ax.xminUser, ax.xmaxUser
		ax.xticks = ticksWithin(ax.xmin, ax.xmax)
	} else {
		t := niceTicks(lo, hi, 5)
		ax.xmin, ax.xmax = t[0], t[len(t)-1]
		ax.xticks = t
	}
	ax.xtickLabels = numericLabels(ax.xticks)
}

func (ax *Axes) setNumericY(lo, hi float64) {
	if ax.ylimSet {
		ax.ymin, ax.ymax = ax.yminUser, ax.ymaxUser
		ax.yticks = ticksWithin(ax.ymin, ax.ymax)
	} else {
		t := niceTicks(lo, hi, 5)
		ax.ymin, ax.ymax = t[0], t[len(t)-1]
		ax.yticks = t
	}
	ax.ytickLabels = numericLabels(ax.yticks)
}

// render draws the axes and all its items onto cv.
func (ax *Axes) render(cv canvas) {
	fw, fh := cv.size()
	ax0 := ax.pos[0] * float64(fw)
	ay0 := ax.pos[1] * float64(fh)
	ax1 := ax.pos[2] * float64(fw)
	ay1 := ax.pos[3] * float64(fh)

	if ax.isPie {
		ax.renderPie(cv, ax0, ay0, ax1, ay1)
		return
	}

	ax.finalize()

	left := 48.0
	if ax.YLabel != "" {
		left += 16
	}
	if ax.horizontalCat {
		left = 78
		if ax.YLabel != "" {
			left += 16
		}
	}
	bottom := 34.0
	if ax.XLabel != "" {
		bottom += 18
	}
	top := 14.0
	if ax.Title != "" {
		top += 22
	}
	right := 18.0

	ax.px0 = ax0 + left
	ax.py0 = ay0 + top
	ax.px1 = ax1 - right
	ax.py1 = ay1 - bottom

	// Grid lines.
	if ax.ShowGrid {
		for _, xt := range ax.xticks {
			px := ax.tx(xt)
			cv.strokeLine(px, ax.py0, px, ax.py1, gridColor, 1)
		}
		for _, yt := range ax.yticks {
			py := ax.ty(yt)
			cv.strokeLine(ax.px0, py, ax.px1, py, gridColor, 1)
		}
	}

	// Series.
	for _, it := range ax.items {
		it.render(cv, ax)
	}

	// Axis frame.
	cv.strokeLine(ax.px0, ax.py1, ax.px1, ax.py1, axisColor, 1.5)
	cv.strokeLine(ax.px0, ax.py0, ax.px0, ax.py1, axisColor, 1.5)

	// X ticks + labels.
	for i, xt := range ax.xticks {
		px := ax.tx(xt)
		if px < ax.px0-0.5 || px > ax.px1+0.5 {
			continue
		}
		cv.strokeLine(px, ax.py1, px, ax.py1+5, axisColor, 1)
		label := ""
		if i < len(ax.xtickLabels) {
			label = ax.xtickLabels[i]
		}
		cv.drawText(px, ax.py1+15, label, textColor, AnchorMiddle, 11, false)
	}

	// Y ticks + labels.
	for i, yt := range ax.yticks {
		py := ax.ty(yt)
		if py < ax.py0-0.5 || py > ax.py1+0.5 {
			continue
		}
		cv.strokeLine(ax.px0-5, py, ax.px0, py, axisColor, 1)
		label := ""
		if i < len(ax.ytickLabels) {
			label = ax.ytickLabels[i]
		}
		cv.drawText(ax.px0-8, py, label, textColor, AnchorEnd, 11, false)
	}

	// Axis labels and title.
	if ax.XLabel != "" {
		cv.drawText((ax.px0+ax.px1)/2, ax.py1+bottom-6, ax.XLabel, textColor, AnchorMiddle, 12, false)
	}
	if ax.YLabel != "" {
		cv.drawText(ax0+14, (ax.py0+ax.py1)/2, ax.YLabel, textColor, AnchorMiddle, 12, true)
	}
	if ax.Title != "" {
		cv.drawText((ax.px0+ax.px1)/2, ay0+top/2, ax.Title, textColor, AnchorMiddle, 15, false)
	}

	if ax.ShowLegend {
		ax.renderLegend(cv)
	}
}

// renderLegend draws a legend box in the top-right corner of the plot area.
func (ax *Axes) renderLegend(cv canvas) {
	type entry struct {
		label  string
		col    Color
		line   bool
		marker Marker
	}
	var entries []entry
	for _, it := range ax.items {
		label, col, line, marker, ok := it.legendEntry()
		if !ok {
			continue
		}
		entries = append(entries, entry{label, col, line, marker})
	}
	if len(entries) == 0 {
		return
	}
	rowH := 16.0
	boxW := 110.0
	boxH := rowH*float64(len(entries)) + 8
	bx := ax.px1 - boxW - 6
	by := ax.py0 + 6
	cv.fillRect(bx, by, boxW, boxH, Color{255, 255, 255, 235})
	cv.strokeLine(bx, by, bx+boxW, by, axisColor, 1)
	cv.strokeLine(bx, by+boxH, bx+boxW, by+boxH, axisColor, 1)
	cv.strokeLine(bx, by, bx, by+boxH, axisColor, 1)
	cv.strokeLine(bx+boxW, by, bx+boxW, by+boxH, axisColor, 1)
	for i, e := range entries {
		cy := by + 8 + rowH*float64(i) + rowH/2 - 4
		if e.line {
			cv.strokeLine(bx+8, cy, bx+26, cy, e.col, 2)
		}
		if e.marker != MarkerNone {
			drawMarker(cv, e.marker, bx+17, cy, 6, e.col)
		}
		if !e.line && e.marker == MarkerNone {
			cv.fillRect(bx+8, cy-4, 18, 8, e.col)
		}
		cv.drawText(bx+32, cy, e.label, textColor, AnchorStart, 11, false)
	}
}

// renderPie draws a pie chart within the given axes pixel rectangle.
func (ax *Axes) renderPie(cv canvas, ax0, ay0, ax1, ay1 float64) {
	pie, ok := firstPie(ax.items)
	if !ok {
		return
	}
	top := 14.0
	if ax.Title != "" {
		top += 22
	}
	total := 0.0
	for _, v := range pie.Values {
		if v > 0 {
			total += v
		}
	}
	if total <= 0 {
		return
	}
	cx := (ax0 + ax1) / 2
	cy := (ay0 + top + ay1) / 2
	radius := math.Min(ax1-ax0, ay1-ay0-top)/2 - 30
	if radius < 5 {
		radius = 5
	}

	angle := -math.Pi / 2 // start at the top
	const seg = 60
	for i, v := range pie.Values {
		if v <= 0 {
			continue
		}
		sweep := v / total * 2 * math.Pi
		steps := int(math.Ceil(float64(seg) * sweep / (2 * math.Pi)))
		if steps < 1 {
			steps = 1
		}
		pts := []Point{{cx, cy}}
		for s := 0; s <= steps; s++ {
			a := angle + sweep*float64(s)/float64(steps)
			pts = append(pts, Point{cx + radius*math.Cos(a), cy + radius*math.Sin(a)})
		}
		cv.polygon(pts, pie.Colors[i%len(pie.Colors)])
		// Label at the wedge midpoint.
		mid := angle + sweep/2
		lx := cx + (radius+16)*math.Cos(mid)
		ly := cy + (radius+16)*math.Sin(mid)
		if i < len(pie.Labels) {
			anchor := AnchorStart
			if math.Cos(mid) < 0 {
				anchor = AnchorEnd
			}
			cv.drawText(lx, ly, pie.Labels[i], textColor, anchor, 11, false)
		}
		angle += sweep
	}
	if ax.Title != "" {
		cv.drawText(cx, ay0+top/2, ax.Title, textColor, AnchorMiddle, 15, false)
	}
}

func firstPie(items []plotItem) (*PieChart, bool) {
	for _, it := range items {
		if p, ok := it.(*PieChart); ok {
			return p, true
		}
	}
	return nil, false
}

// intTicks returns tick positions 0,1,...,n-1.
func intTicks(n int) []float64 {
	t := make([]float64, n)
	for i := range t {
		t[i] = float64(i)
	}
	return t
}

// ticksWithin returns nice ticks clipped to [lo,hi].
func ticksWithin(lo, hi float64) []float64 {
	all := niceTicks(lo, hi, 5)
	var out []float64
	for _, v := range all {
		if v >= lo-1e-9 && v <= hi+1e-9 {
			out = append(out, v)
		}
	}
	if len(out) == 0 {
		out = []float64{lo, hi}
	}
	return out
}

func numericLabels(ticks []float64) []string {
	out := make([]string, len(ticks))
	for i, v := range ticks {
		out[i] = fmtTick(v)
	}
	return out
}
