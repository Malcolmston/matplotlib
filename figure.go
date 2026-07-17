package matplotlib

import (
	"bufio"
	"bytes"
	"image"
	"image/png"
	"io"
	"os"
)

// Figure is the top-level container for one or more Axes. It owns the output
// dimensions (in pixels), background color and the color cycle shared by the
// series it contains.
type Figure struct {
	Width      int
	Height     int
	DPI        float64
	Background Color

	axes     []*Axes
	colorIdx int
}

// NewFigure creates a figure of the given pixel size with a white background
// and 100 DPI. A width or height less than 1 is clamped to a small default.
func NewFigure(width, height int) *Figure {
	if width < 1 {
		width = 640
	}
	if height < 1 {
		height = 480
	}
	return &Figure{
		Width:      width,
		Height:     height,
		DPI:        100,
		Background: White,
	}
}

// nextColor returns the next color from the figure's cycle.
func (f *Figure) nextColor() Color {
	c := DefaultColors[f.colorIdx%len(DefaultColors)]
	f.colorIdx++
	return c
}

// AddAxes appends a new Axes filling the entire figure and returns it.
func (f *Figure) AddAxes() *Axes {
	ax := &Axes{fig: f, pos: [4]float64{0, 0, 1, 1}}
	f.axes = append(f.axes, ax)
	return ax
}

// Axes returns the figure's first Axes, creating one if none exist. This is a
// convenience for the common single-plot case.
func (f *Figure) Axes() *Axes {
	if len(f.axes) == 0 {
		return f.AddAxes()
	}
	return f.axes[0]
}

// Subplots creates a rows x cols grid of Axes and returns them in row-major
// order. Each Axes occupies its cell of the figure.
func (f *Figure) Subplots(rows, cols int) []*Axes {
	if rows < 1 {
		rows = 1
	}
	if cols < 1 {
		cols = 1
	}
	out := make([]*Axes, 0, rows*cols)
	fw := 1.0 / float64(cols)
	fh := 1.0 / float64(rows)
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			ax := &Axes{
				fig: f,
				pos: [4]float64{float64(c) * fw, float64(r) * fh, float64(c+1) * fw, float64(r+1) * fh},
			}
			f.axes = append(f.axes, ax)
			out = append(out, ax)
		}
	}
	return out
}

// renderTo draws the whole figure onto the given canvas.
func (f *Figure) renderTo(cv canvas) {
	cv.fillBackground(f.Background)
	for _, ax := range f.axes {
		ax.render(cv)
	}
}

// RenderImage renders the figure to an in-memory RGBA image.
func (f *Figure) RenderImage() *image.RGBA {
	cv := newPNGCanvas(f.Width, f.Height)
	f.renderTo(cv)
	return cv.img
}

// RenderSVG renders the figure to an SVG document string.
func (f *Figure) RenderSVG() string {
	cv := newSVGCanvas(f.Width, f.Height)
	f.renderTo(cv)
	return cv.document()
}

// WritePNG encodes the figure as PNG to w.
func (f *Figure) WritePNG(w io.Writer) error {
	return png.Encode(w, f.RenderImage())
}

// WriteSVG writes the figure's SVG document to w.
func (f *Figure) WriteSVG(w io.Writer) error {
	_, err := io.WriteString(w, f.RenderSVG())
	return err
}

// SavePNG renders and writes the figure to a PNG file at path.
func (f *Figure) SavePNG(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	bw := bufio.NewWriter(file)
	if err := f.WritePNG(bw); err != nil {
		return err
	}
	return bw.Flush()
}

// SaveSVG renders and writes the figure to an SVG file at path.
func (f *Figure) SaveSVG(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	bw := bufio.NewWriter(file)
	if err := f.WriteSVG(bw); err != nil {
		return err
	}
	return bw.Flush()
}

// Save renders the figure to path, choosing the format from the file
// extension (".svg" for SVG, anything else for PNG).
func (f *Figure) Save(path string) error {
	if hasSuffixFold(path, ".svg") {
		return f.SaveSVG(path)
	}
	return f.SavePNG(path)
}

// PNGBytes returns the figure encoded as a PNG byte slice.
func (f *Figure) PNGBytes() ([]byte, error) {
	var buf bytes.Buffer
	if err := f.WritePNG(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func hasSuffixFold(s, suffix string) bool {
	if len(s) < len(suffix) {
		return false
	}
	s = s[len(s)-len(suffix):]
	for i := 0; i < len(suffix); i++ {
		a := s[i]
		b := suffix[i]
		if 'A' <= a && a <= 'Z' {
			a += 'a' - 'A'
		}
		if 'A' <= b && b <= 'Z' {
			b += 'a' - 'A'
		}
		if a != b {
			return false
		}
	}
	return true
}
