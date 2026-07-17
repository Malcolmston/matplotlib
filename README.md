# matplotlib

A small, dependency-free plotting and charting library for Go, inspired by
Python's matplotlib. It renders charts to both **PNG** (raster) and **SVG**
(vector) using only the Go standard library — no cgo, no third-party modules.

## Features

- **Figure + Axes model** with a pyplot-style convenience API.
- **Chart types:** line plot (multiple series), scatter, vertical bar,
  horizontal bar, histogram (automatic binning), and pie.
- **Axes rendering:** automatic data-range to pixel mapping, x/y axis lines,
  numeric ticks and labels, axis labels, title, optional legend, optional grid.
- **Styling:** per-series color (a default `tab10` color cycle), line width,
  markers; figure size in pixels and DPI.
- **Two real output formats** from the standard library only:
  - PNG via `image` / `image/png` (Bresenham lines, filled rects/wedges,
    markers, and a built-in 5x7 bitmap font for text).
  - SVG via string building (`<line>`, `<rect>`, `<circle>`, `<polyline>`,
    `<polygon>`, `<text>`).
- **Deterministic:** identical input always produces identical bytes.

## Install

```sh
go get github.com/malcolmston/matplotlib
```

Requires Go 1.24 or newer.

## Quick start

```go
package main

import "github.com/malcolmston/matplotlib"

func main() {
	fig := matplotlib.NewFigure(640, 480)
	ax := fig.AddAxes()

	ax.Plot([]float64{0, 1, 2, 3}, []float64{0, 1, 4, 9}).
		SetLabel("x^2").
		SetMarker(matplotlib.MarkerCircle)
	ax.Plot([]float64{0, 1, 2, 3}, []float64{0, 2, 4, 6}).SetLabel("2x")

	ax.SetTitle("Demo").SetXLabel("x").SetYLabel("y").Legend().Grid(true)

	_ = fig.SavePNG("demo.png")
	_ = fig.SaveSVG("demo.svg")
}
```

### pyplot-style API

```go
matplotlib.FigSize(640, 480)
matplotlib.Bar([]string{"a", "b", "c"}, []float64{3, 7, 2})
matplotlib.Title("Counts")
_ = matplotlib.Save("bar.png") // format chosen by extension
```

## Chart types

| Method          | Chart                    |
| --------------- | ------------------------ |
| `Axes.Plot`     | Line plot (per series)   |
| `Axes.Scatter`  | Scatter plot             |
| `Axes.Bar`      | Vertical bar chart       |
| `Axes.BarH`     | Horizontal bar chart     |
| `Axes.Hist`     | Histogram                |
| `Axes.Pie`      | Pie chart                |

## Output

- `fig.SavePNG(path)` / `fig.SaveSVG(path)` / `fig.Save(path)` (by extension)
- `fig.WritePNG(w)` / `fig.WriteSVG(w)` to any `io.Writer`
- `fig.RenderImage()` returns an `*image.RGBA`; `fig.RenderSVG()` returns the
  SVG document string.

## License

See repository.
