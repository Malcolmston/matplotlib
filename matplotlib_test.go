package matplotlib

import (
	"bytes"
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// renderPNG renders a figure and returns the decoded image, failing on error.
func renderPNG(t *testing.T, f *Figure) image.Image {
	t.Helper()
	var buf bytes.Buffer
	if err := f.WritePNG(&buf); err != nil {
		t.Fatalf("WritePNG: %v", err)
	}
	img, err := png.Decode(&buf)
	if err != nil {
		t.Fatalf("png.Decode: %v", err)
	}
	return img
}

// nonBlank reports whether the image contains at least two distinct colors.
func nonBlank(img image.Image) bool {
	b := img.Bounds()
	var first uint32
	got := false
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, _ := img.At(x, y).RGBA()
			key := r<<20 ^ g<<10 ^ bl
			if !got {
				first = key
				got = true
				continue
			}
			if key != first {
				return true
			}
		}
	}
	return false
}

func TestLinePlotBothFormats(t *testing.T) {
	fig := NewFigure(400, 300)
	ax := fig.AddAxes()
	ax.Plot([]float64{0, 1, 2, 3}, []float64{0, 1, 4, 9}).SetLabel("sq").SetMarker(MarkerCircle)
	ax.Plot([]float64{0, 1, 2, 3}, []float64{0, 2, 4, 6}).SetLabel("lin")
	ax.SetTitle("Lines").SetXLabel("x").SetYLabel("y").Legend().Grid(true)

	img := renderPNG(t, fig)
	if img.Bounds().Dx() != 400 || img.Bounds().Dy() != 300 {
		t.Fatalf("PNG dims = %v, want 400x300", img.Bounds())
	}
	if !nonBlank(img) {
		t.Fatal("PNG appears blank")
	}

	svg := fig.RenderSVG()
	for _, want := range []string{"<svg", "<polyline", "<text", "Lines", "</svg>"} {
		if !strings.Contains(svg, want) {
			t.Errorf("SVG missing %q", want)
		}
	}
}

func TestScatter(t *testing.T) {
	fig := NewFigure(300, 200)
	ax := fig.AddAxes()
	ax.Scatter([]float64{1, 2, 3}, []float64{3, 1, 2}).SetMarker(MarkerSquare).SetLabel("pts").SetSize(8)
	ax.Legend()
	svg := fig.RenderSVG()
	if !strings.Contains(svg, "<rect") {
		t.Error("scatter square SVG missing <rect>")
	}
	if !nonBlank(renderPNG(t, fig)) {
		t.Error("scatter PNG blank")
	}
}

func TestBarVerticalAndHorizontal(t *testing.T) {
	labels := []string{"a", "b", "c"}
	vals := []float64{3, 7, 2}

	fig := NewFigure(320, 240)
	ax := fig.AddAxes()
	ax.Bar(labels, vals).SetLabel("v")
	ax.SetTitle("Bar").Legend()
	svg := fig.RenderSVG()
	if !strings.Contains(svg, "<rect") {
		t.Error("bar SVG missing rect")
	}
	if !nonBlank(renderPNG(t, fig)) {
		t.Error("bar PNG blank")
	}

	fig2 := NewFigure(320, 240)
	ax2 := fig2.AddAxes()
	ax2.BarH(labels, vals)
	if !nonBlank(renderPNG(t, fig2)) {
		t.Error("hbar PNG blank")
	}
	if !strings.Contains(fig2.RenderSVG(), "<rect") {
		t.Error("hbar SVG missing rect")
	}
}

func TestHistogram(t *testing.T) {
	data := make([]float64, 100)
	for i := range data {
		data[i] = math.Sin(float64(i)) * 10
	}
	fig := NewFigure(320, 240)
	ax := fig.AddAxes()
	h := ax.Hist(data, 8)
	if !nonBlank(renderPNG(t, fig)) {
		t.Error("hist PNG blank")
	}
	if len(h.counts) != 8 {
		t.Fatalf("bins = %d, want 8", len(h.counts))
	}
	total := 0
	for _, c := range h.counts {
		total += c
	}
	if total != len(data) {
		t.Errorf("hist counted %d, want %d", total, len(data))
	}
	if !strings.Contains(fig.RenderSVG(), "<rect") {
		t.Error("hist SVG missing rect")
	}
}

func TestHistDefaultsAndEmpty(t *testing.T) {
	fig := NewFigure(200, 200)
	ax := fig.AddAxes()
	ax.Hist([]float64{}, 0) // empty data, default bins
	_ = fig.RenderSVG()
	renderPNG(t, fig)
}

func TestPie(t *testing.T) {
	fig := NewFigure(300, 300)
	ax := fig.AddAxes()
	ax.Pie([]float64{1, 2, 3, 4}, []string{"w", "x", "y", "z"})
	ax.SetTitle("Pie")
	svg := fig.RenderSVG()
	if !strings.Contains(svg, "<polygon") {
		t.Error("pie SVG missing polygon")
	}
	if !nonBlank(renderPNG(t, fig)) {
		t.Error("pie PNG blank")
	}
}

func TestSubplots(t *testing.T) {
	fig := NewFigure(400, 400)
	axs := fig.Subplots(2, 2)
	if len(axs) != 4 {
		t.Fatalf("got %d axes, want 4", len(axs))
	}
	axs[0].Plot([]float64{0, 1}, []float64{0, 1})
	axs[1].Scatter([]float64{0, 1}, []float64{1, 0})
	axs[2].Bar([]string{"a", "b"}, []float64{1, 2})
	axs[3].Pie([]float64{1, 1}, nil)
	if !nonBlank(renderPNG(t, fig)) {
		t.Error("subplots PNG blank")
	}
}

func TestDeterministic(t *testing.T) {
	build := func() *Figure {
		f := NewFigure(300, 200)
		a := f.AddAxes()
		a.Plot([]float64{0, 1, 2}, []float64{1, 3, 2}).SetLabel("s")
		a.SetTitle("T").Legend().Grid(true)
		return f
	}
	a := build().RenderSVG()
	b := build().RenderSVG()
	if a != b {
		t.Error("SVG not deterministic")
	}
	p1, _ := build().PNGBytes()
	p2, _ := build().PNGBytes()
	if !bytes.Equal(p1, p2) {
		t.Error("PNG not deterministic")
	}
}

func TestSaveFiles(t *testing.T) {
	dir := t.TempDir()
	fig := NewFigure(200, 150)
	fig.AddAxes().Plot([]float64{0, 1, 2}, []float64{0, 1, 0})

	pngPath := filepath.Join(dir, "out.png")
	svgPath := filepath.Join(dir, "out.SVG") // uppercase to test fold
	if err := fig.Save(pngPath); err != nil {
		t.Fatal(err)
	}
	if err := fig.Save(svgPath); err != nil {
		t.Fatal(err)
	}
	pb, err := os.ReadFile(pngPath)
	if err != nil || len(pb) == 0 {
		t.Fatalf("png file: %v len=%d", err, len(pb))
	}
	if _, err := png.Decode(bytes.NewReader(pb)); err != nil {
		t.Fatalf("saved png not decodable: %v", err)
	}
	sb, err := os.ReadFile(svgPath)
	if err != nil || !bytes.Contains(sb, []byte("<svg")) {
		t.Fatalf("svg file bad: %v", err)
	}
}

func TestXYLimits(t *testing.T) {
	fig := NewFigure(300, 200)
	ax := fig.AddAxes()
	ax.Plot([]float64{0, 5, 10}, []float64{0, 5, 10})
	ax.SetXLim(0, 10).SetYLim(0, 10)
	ax.finalize()
	if ax.xmin != 0 || ax.xmax != 10 {
		t.Errorf("xlim = %v..%v", ax.xmin, ax.xmax)
	}
}

func TestPyplotAPI(t *testing.T) {
	Clf()
	FigSize(320, 240)
	Plot([]float64{0, 1, 2}, []float64{2, 1, 3}).SetLabel("p")
	Scatter([]float64{0, 1}, []float64{1, 2})
	Title("PP")
	Xlabel("x")
	Ylabel("y")
	Legend()
	Grid(true)
	Xlim(0, 2)
	Ylim(0, 3)
	dir := t.TempDir()
	if err := Save(filepath.Join(dir, "p.png")); err != nil {
		t.Fatal(err)
	}
	if Gcf().Width != 320 {
		t.Error("Gcf width wrong")
	}
	if Gca() == nil {
		t.Error("Gca nil")
	}
	Clf()
	// Lazy creation path.
	Bar([]string{"a"}, []float64{1})
	BarH([]string{"b"}, []float64{2})
	Hist([]float64{1, 2, 3}, 3)
	Pie([]float64{1, 1}, nil)
	_ = Gcf().RenderSVG()
	Clf()
}

func TestColorHelpers(t *testing.T) {
	if RGB(1, 2, 3).A != 255 {
		t.Error("RGB alpha")
	}
	if RGBA(1, 2, 3, 4).A != 4 {
		t.Error("RGBA alpha")
	}
	c, err := Hex("#ff8800")
	if err != nil || c.R != 255 || c.G != 136 || c.B != 0 {
		t.Errorf("Hex 6 = %+v err=%v", c, err)
	}
	c, err = Hex("f80")
	if err != nil || c.R != 255 || c.G != 136 || c.B != 0 {
		t.Errorf("Hex 3 = %+v err=%v", c, err)
	}
	if _, err := Hex("nope"); err == nil {
		t.Error("Hex bad should error")
	}
	if Black.hexString() != "#000000" {
		t.Error("hexString")
	}
	if !Transparent.isTransparent() {
		t.Error("isTransparent")
	}
}

func TestNiceTicksAndFormat(t *testing.T) {
	ticks := niceTicks(0, 100, 5)
	if len(ticks) < 2 {
		t.Fatal("too few ticks")
	}
	if ticks[0] != 0 {
		t.Errorf("first tick = %v", ticks[0])
	}
	// Degenerate range.
	if len(niceTicks(5, 5, 5)) < 2 {
		t.Error("degenerate ticks")
	}
	// reversed
	if len(niceTicks(10, 0, 5)) < 2 {
		t.Error("reversed ticks")
	}
	if fmtTick(0) != "0" {
		t.Error("fmtTick 0")
	}
	if fmtTick(1.5) != "1.5" {
		t.Errorf("fmtTick 1.5 = %q", fmtTick(1.5))
	}
	if !strings.Contains(fmtTick(1e6), "e") {
		t.Errorf("fmtTick big = %q", fmtTick(1e6))
	}
	if fnum(1.230) != "1.23" {
		t.Errorf("fnum = %q", fnum(1.23))
	}
	if fnum(math.NaN()) != "0" {
		t.Error("fnum NaN")
	}
	if fnum(-0.0001) != "0" {
		t.Errorf("fnum -0 = %q", fnum(-0.0001))
	}
}

func TestFontCoverage(t *testing.T) {
	if _, ok := glyphFor('A'); !ok {
		t.Error("A not found")
	}
	if _, ok := glyphFor('☃'); ok {
		t.Error("snowman should be unknown")
	}
	// Render all printable ASCII through the PNG text path.
	var sb strings.Builder
	for r := rune(32); r < 127; r++ {
		sb.WriteRune(r)
	}
	fig := NewFigure(700, 60)
	ax := fig.AddAxes()
	ax.SetTitle(sb.String())
	renderPNG(t, fig)
}

func TestVerticalText(t *testing.T) {
	cv := newPNGCanvas(60, 120)
	cv.fillBackground(White)
	cv.drawText(30, 60, "Value", Black, AnchorMiddle, 12, true)
	if !nonBlank(cv.img) {
		t.Error("vertical text produced blank image")
	}
	s := newSVGCanvas(60, 120)
	s.drawText(30, 60, "Value", Black, AnchorMiddle, 12, true)
	if !strings.Contains(s.document(), "rotate(") {
		t.Error("SVG vertical text missing rotate")
	}
}

func TestMarkersAndPrimitives(t *testing.T) {
	cv := newPNGCanvas(50, 50)
	cv.fillBackground(White)
	for _, m := range []Marker{MarkerCircle, MarkerSquare, MarkerTriangle, MarkerCross, MarkerPlus, MarkerNone} {
		drawMarker(cv, m, 25, 25, 10, Blue)
	}
	cv.polygon([]Point{{5, 5}, {45, 5}, {25, 45}}, Green)
	cv.polyline([]Point{{0, 0}, {49, 49}}, Red, 2)
	if !nonBlank(cv.img) {
		t.Error("primitives blank")
	}

	sv := newSVGCanvas(50, 50)
	sv.fillBackground(White)
	drawMarker(sv, MarkerTriangle, 25, 25, 10, Blue)
	sv.fillCircle(25, 25, 5, Red)
	sv.fillRect(-5, -5, 10, 10, Green) // negative offset path
	doc := sv.document()
	for _, w := range []string{"<polygon", "<circle", "<rect"} {
		if !strings.Contains(doc, w) {
			t.Errorf("SVG missing %s", w)
		}
	}
}

func TestEmptyAxes(t *testing.T) {
	fig := NewFigure(200, 200)
	fig.AddAxes() // no data
	if !strings.Contains(fig.RenderSVG(), "<svg") {
		t.Error("empty axes svg")
	}
	renderPNG(t, fig)
}

func TestAxesAccessorAndColorCycle(t *testing.T) {
	fig := NewFigure(200, 200)
	first := fig.Axes()
	if first != fig.Axes() {
		t.Error("Axes() should be stable")
	}
	c0 := fig.nextColor()
	c1 := fig.nextColor()
	if c0 == c1 {
		t.Error("color cycle not advancing")
	}
}
