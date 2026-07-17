package matplotlib_test

import (
	"fmt"
	"strings"

	"github.com/malcolmston/matplotlib"
)

// Example renders a tiny line plot to SVG and prints its dimensions and a
// snippet of the output, demonstrating the basic Figure/Axes workflow.
func Example() {
	fig := matplotlib.NewFigure(320, 240)
	ax := fig.AddAxes()
	ax.Plot([]float64{0, 1, 2, 3}, []float64{0, 1, 4, 9}).SetLabel("x^2")
	ax.SetTitle("Squares").SetXLabel("x").SetYLabel("y").Legend()

	svg := fig.RenderSVG()
	fmt.Printf("size: %dx%d\n", fig.Width, fig.Height)
	fmt.Println("has polyline:", strings.Contains(svg, "<polyline"))
	fmt.Println("has title:", strings.Contains(svg, "Squares"))
	// Output:
	// size: 320x240
	// has polyline: true
	// has title: true
}

// ExampleFigure_Save writes a bar chart to a PNG file using the pyplot-style
// convenience API.
func ExampleFigure_Save() {
	fig := matplotlib.NewFigure(400, 300)
	ax := fig.AddAxes()
	ax.Bar([]string{"a", "b", "c"}, []float64{3, 7, 2})
	ax.SetTitle("Counts")

	img := fig.RenderImage()
	fmt.Printf("%dx%d\n", img.Bounds().Dx(), img.Bounds().Dy())
	// Output:
	// 400x300
}
