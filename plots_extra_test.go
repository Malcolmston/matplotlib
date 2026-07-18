package matplotlib

import (
	"bytes"
	"testing"
)

func pointsEqual(a, b []Point, tol float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if fabs(a[i].X-b[i].X) > tol || fabs(a[i].Y-b[i].Y) > tol {
			return false
		}
	}
	return true
}

func fabs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func TestStepPoints(t *testing.T) {
	x := []float64{0, 1, 2}
	y := []float64{0, 1, 2}

	post := (&StepPlot{X: x, Y: y, Where: "post"}).stepPoints()
	wantPost := []Point{{0, 0}, {1, 0}, {1, 1}, {2, 1}, {2, 2}}
	if !pointsEqual(post, wantPost, 1e-12) {
		t.Errorf("post steps = %v, want %v", post, wantPost)
	}

	pre := (&StepPlot{X: x, Y: y, Where: "pre"}).stepPoints()
	wantPre := []Point{{0, 0}, {0, 1}, {1, 1}, {1, 2}, {2, 2}}
	if !pointsEqual(pre, wantPre, 1e-12) {
		t.Errorf("pre steps = %v, want %v", pre, wantPre)
	}

	mid := (&StepPlot{X: x, Y: y, Where: "mid"}).stepPoints()
	wantMid := []Point{{0, 0}, {0.5, 0}, {0.5, 1}, {1.5, 1}, {1.5, 2}, {2, 2}}
	if !pointsEqual(mid, wantMid, 1e-12) {
		t.Errorf("mid steps = %v, want %v", mid, wantMid)
	}
}

func TestFillBetweenBounds(t *testing.T) {
	f := &FillBetween{X: []float64{0, 1, 2}, Y1: []float64{1, 2, 3}}
	xlo, xhi, ylo, yhi, ok := f.bounds()
	if !ok || xlo != 0 || xhi != 2 || ylo != 0 || yhi != 3 {
		t.Errorf("bounds = %v %v %v %v ok=%v", xlo, xhi, ylo, yhi, ok)
	}
	f2 := &FillBetween{X: []float64{0, 1}, Y1: []float64{2, 4}, Y2: []float64{1, 1}}
	_, _, ylo2, yhi2, _ := f2.bounds()
	if ylo2 != 1 || yhi2 != 4 {
		t.Errorf("bounds y = %v..%v, want 1..4", ylo2, yhi2)
	}
}

func TestStemBounds(t *testing.T) {
	s := &StemPlot{X: []float64{1, 2}, Y: []float64{3, 5}}
	_, _, ylo, yhi, ok := s.bounds()
	if !ok || ylo != 0 || yhi != 5 {
		t.Errorf("stem bounds ylo=%v yhi=%v", ylo, yhi)
	}
	// Negative values should extend the baseline downward and cap yhi at 0.
	s2 := &StemPlot{X: []float64{0, 1}, Y: []float64{-2, -4}}
	_, _, ylo2, yhi2, _ := s2.bounds()
	if ylo2 != -4 || yhi2 != 0 {
		t.Errorf("stem neg bounds ylo=%v yhi=%v", ylo2, yhi2)
	}
}

func TestErrorBarBounds(t *testing.T) {
	e := &ErrorBar{X: []float64{0, 1}, Y: []float64{5, 5}, YErr: []float64{1, 2}}
	_, _, ylo, yhi, ok := e.bounds()
	if !ok || ylo != 3 || yhi != 7 {
		t.Errorf("errorbar bounds ylo=%v yhi=%v", ylo, yhi)
	}
}

func TestRefLineNoAutoscale(t *testing.T) {
	r := &RefLine{Value: 100}
	if _, _, _, _, ok := r.bounds(); ok {
		t.Errorf("RefLine should not contribute bounds")
	}
	a := &Annotation{X: 5, Y: 5, Text: "hi"}
	if _, _, _, _, ok := a.bounds(); ok {
		t.Errorf("Annotation should not contribute bounds")
	}
}

// TestExtraPlotsRenderDeterministic exercises every new chart type end to end
// through both backends and checks the raster output is reproducible.
func TestExtraPlotsRenderDeterministic(t *testing.T) {
	build := func() *Figure {
		fig := NewFigure(400, 300)
		ax := fig.AddAxes()
		x := Linspace(0, 10, 11)
		y := make([]float64, len(x))
		for i := range x {
			y[i] = x[i] * x[i] / 10
		}
		ax.Step(x, y).SetLabel("step").SetWhere("mid")
		ax.FillBetween(x, y, nil).SetColor(Blue.WithAlpha(120)).SetLabel("area")
		ax.Stem([]float64{1, 3, 5}, []float64{2, 5, 3}).SetLabel("stem")
		ax.ErrorBar([]float64{2, 4, 6}, []float64{4, 6, 5}, []float64{0.5, 1, 0.5}).SetLabel("err")
		ax.AxHLine(5).SetColor(Red).SetLabel("threshold")
		ax.AxVLine(4).SetColor(Green)
		ax.Text(1, 8, "note")
		ax.SetTitle("extras").Legend().Grid(true)
		return fig
	}
	a, err := build().PNGBytes()
	if err != nil {
		t.Fatalf("render 1: %v", err)
	}
	b, err := build().PNGBytes()
	if err != nil {
		t.Fatalf("render 2: %v", err)
	}
	if !bytes.Equal(a, b) {
		t.Errorf("PNG output not deterministic (%d vs %d bytes)", len(a), len(b))
	}
	// SVG should also be produced without panicking and be non-trivial.
	if svg := build().RenderSVG(); len(svg) < 100 {
		t.Errorf("SVG output too small: %d bytes", len(svg))
	}
}

func BenchmarkStepRender(b *testing.B) {
	x := Linspace(0, 100, 200)
	fig := NewFigure(600, 400)
	ax := fig.AddAxes()
	ax.Step(x, x)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fig.RenderImage()
	}
}
