// Package matplotlib is a small, dependency-free plotting and charting library
// for Go, inspired by Python's matplotlib. It renders charts to both PNG
// (raster) and SVG (vector) using only the Go standard library — there is no
// cgo and no third-party code.
//
// # Model
//
// The drawing model has two levels:
//
//   - A [Figure] is the top-level canvas. It owns the output size in pixels,
//     the background color, the DPI and a color cycle shared by its series.
//   - An [Axes] is a coordinate system within a figure. It holds one or more
//     plot items (series) and the surrounding decorations: title, axis labels,
//     numeric ticks, an optional grid and an optional legend.
//
// A figure may contain a single axes (the common case, via [Figure.AddAxes] or
// [Figure.Axes]) or a grid of them (via [Figure.Subplots]).
//
// # Chart types
//
// The following chart types are available as methods on [Axes]:
//
//   - [Axes.Plot]    — line plot, one call per series (multiple series overlay)
//   - [Axes.Scatter] — scatter plot with configurable markers
//   - [Axes.Bar]     — vertical bar chart over string categories
//   - [Axes.BarH]    — horizontal bar chart
//   - [Axes.Hist]    — histogram with automatic equal-width binning
//   - [Axes.Pie]     — pie chart
//
// Each returns a typed value ([LinePlot], [ScatterPlot], [BarChart],
// [Histogram], [PieChart]) whose fields — color, line width, marker, label —
// can be customized with fluent setters.
//
// # Styling
//
// New series are automatically assigned successive colors from [DefaultColors]
// (matplotlib's "tab10" palette). Colors are [Color] values and can be built
// with [RGB], [RGBA] or [Hex]. Line width and marker style are per series.
//
// # Output
//
// Rendering is deterministic: identical input always produces identical bytes.
//
//   - PNG is produced with image and image/png. Lines use a Bresenham routine,
//     bars are filled rectangles, markers and pie wedges are filled shapes, and
//     text is drawn with a built-in 5x7 bitmap font.
//   - SVG is produced by emitting <line>, <rect>, <circle>, <polyline>,
//     <polygon> and <text> elements.
//
// Use [Figure.SavePNG], [Figure.SaveSVG] or the extension-sensitive
// [Figure.Save] to write files; [Figure.WritePNG] and [Figure.WriteSVG] to
// stream to any io.Writer; and [Figure.RenderImage] / [Figure.RenderSVG] to
// obtain the in-memory image or SVG string.
//
// # Quick start
//
//	fig := matplotlib.NewFigure(640, 480)
//	ax := fig.AddAxes()
//	ax.Plot([]float64{0, 1, 2, 3}, []float64{0, 1, 4, 9}).SetLabel("x^2")
//	ax.SetTitle("Demo").SetXLabel("x").SetYLabel("y").Legend().Grid(true)
//	_ = fig.SavePNG("demo.png")
//	_ = fig.SaveSVG("demo.svg")
//
// A pyplot-style stateful API ([Plot], [Bar], [Title], [Save], …) is also
// provided for quick scripts; see pyplot.go.
package matplotlib
