package matplotlib

// This file provides a small pyplot-style stateful convenience API layered on
// top of the object model. It operates on a lazily created "current figure"
// and "current axes", mirroring matplotlib.pyplot for quick, script-style use.
//
// The API is not safe for concurrent use; it is intended for simple, linear
// plotting scripts. For anything else, use NewFigure and the Axes methods
// directly.

var (
	currentFig *Figure
	currentAx  *Axes
)

// FigSize starts a fresh current figure with the given pixel dimensions,
// discarding any previous global state.
func FigSize(width, height int) *Figure {
	currentFig = NewFigure(width, height)
	currentAx = currentFig.AddAxes()
	return currentFig
}

// Clf clears the current figure, resetting the global state.
func Clf() { currentFig, currentAx = nil, nil }

// Gcf returns the current figure, creating a default 640x480 one if needed.
func Gcf() *Figure {
	if currentFig == nil {
		FigSize(640, 480)
	}
	return currentFig
}

// Gca returns the current axes, creating a default figure if needed.
func Gca() *Axes {
	if currentAx == nil {
		Gcf()
	}
	return currentAx
}

// Plot adds a line series to the current axes.
func Plot(x, y []float64) *LinePlot { return Gca().Plot(x, y) }

// Scatter adds a scatter series to the current axes.
func Scatter(x, y []float64) *ScatterPlot { return Gca().Scatter(x, y) }

// Bar adds a vertical bar chart to the current axes.
func Bar(labels []string, values []float64) *BarChart { return Gca().Bar(labels, values) }

// BarH adds a horizontal bar chart to the current axes.
func BarH(labels []string, values []float64) *BarChart { return Gca().BarH(labels, values) }

// Hist adds a histogram to the current axes.
func Hist(data []float64, bins int) *Histogram { return Gca().Hist(data, bins) }

// Pie adds a pie chart to the current axes.
func Pie(values []float64, labels []string) *PieChart { return Gca().Pie(values, labels) }

// Title sets the current axes title.
func Title(s string) { Gca().SetTitle(s) }

// Xlabel sets the current axes x-axis label.
func Xlabel(s string) { Gca().SetXLabel(s) }

// Ylabel sets the current axes y-axis label.
func Ylabel(s string) { Gca().SetYLabel(s) }

// Legend enables the legend on the current axes.
func Legend() { Gca().Legend() }

// Grid toggles the grid on the current axes.
func Grid(on bool) { Gca().Grid(on) }

// Xlim fixes the current axes x-limits.
func Xlim(lo, hi float64) { Gca().SetXLim(lo, hi) }

// Ylim fixes the current axes y-limits.
func Ylim(lo, hi float64) { Gca().SetYLim(lo, hi) }

// Save renders the current figure to path (format chosen by extension).
func Save(path string) error { return Gcf().Save(path) }
