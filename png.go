package matplotlib

import (
	"image"
	"image/color"
	"math"
)

// pngCanvas is the raster (PNG) implementation of the canvas interface. It
// draws into an in-memory RGBA image using simple, deterministic routines.
type pngCanvas struct {
	img *image.RGBA
	w   int
	h   int
}

func newPNGCanvas(w, h int) *pngCanvas {
	return &pngCanvas{img: image.NewRGBA(image.Rect(0, 0, w, h)), w: w, h: h}
}

func (p *pngCanvas) size() (int, int) { return p.w, p.h }

// blend draws src over the existing pixel at (x,y) using its alpha value.
func (p *pngCanvas) blend(x, y int, c Color) {
	if x < 0 || y < 0 || x >= p.w || y >= p.h {
		return
	}
	if c.A == 0 {
		return
	}
	if c.A == 255 {
		p.img.SetRGBA(x, y, toRGBA(c))
		return
	}
	i := p.img.PixOffset(x, y)
	dr := float64(p.img.Pix[i])
	dg := float64(p.img.Pix[i+1])
	db := float64(p.img.Pix[i+2])
	a := float64(c.A) / 255
	p.img.Pix[i] = uint8(float64(c.R)*a + dr*(1-a))
	p.img.Pix[i+1] = uint8(float64(c.G)*a + dg*(1-a))
	p.img.Pix[i+2] = uint8(float64(c.B)*a + db*(1-a))
	p.img.Pix[i+3] = 255
}

func (p *pngCanvas) fillBackground(c Color) {
	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			p.img.SetRGBA(x, y, toRGBA(c))
		}
	}
}

func (p *pngCanvas) fillRect(x, y, w, h float64, c Color) {
	if c.isTransparent() {
		return
	}
	x0 := int(math.Round(x))
	y0 := int(math.Round(y))
	x1 := int(math.Round(x + w))
	y1 := int(math.Round(y + h))
	if x1 < x0 {
		x0, x1 = x1, x0
	}
	if y1 < y0 {
		y0, y1 = y1, y0
	}
	for yy := y0; yy < y1; yy++ {
		for xx := x0; xx < x1; xx++ {
			p.blend(xx, yy, c)
		}
	}
}

// plotThick fills a square of side lw centered at (x,y).
func (p *pngCanvas) plotThick(x, y int, lw float64, c Color) {
	r := int(lw / 2)
	if r <= 0 {
		p.blend(x, y, c)
		return
	}
	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			p.blend(x+dx, y+dy, c)
		}
	}
}

func (p *pngCanvas) strokeLine(x0, y0, x1, y1 float64, c Color, lw float64) {
	if c.isTransparent() {
		return
	}
	ix0 := int(math.Round(x0))
	iy0 := int(math.Round(y0))
	ix1 := int(math.Round(x1))
	iy1 := int(math.Round(y1))
	dx := abs(ix1 - ix0)
	dy := -abs(iy1 - iy0)
	sx := 1
	if ix0 > ix1 {
		sx = -1
	}
	sy := 1
	if iy0 > iy1 {
		sy = -1
	}
	err := dx + dy
	for {
		p.plotThick(ix0, iy0, lw, c)
		if ix0 == ix1 && iy0 == iy1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			ix0 += sx
		}
		if e2 <= dx {
			err += dx
			iy0 += sy
		}
	}
}

func (p *pngCanvas) polyline(pts []Point, c Color, lw float64) {
	for i := 1; i < len(pts); i++ {
		p.strokeLine(pts[i-1].X, pts[i-1].Y, pts[i].X, pts[i].Y, c, lw)
	}
}

func (p *pngCanvas) fillCircle(cx, cy, r float64, c Color) {
	if c.isTransparent() || r <= 0 {
		return
	}
	x0 := int(math.Floor(cx - r))
	x1 := int(math.Ceil(cx + r))
	y0 := int(math.Floor(cy - r))
	y1 := int(math.Ceil(cy + r))
	r2 := r * r
	for yy := y0; yy <= y1; yy++ {
		for xx := x0; xx <= x1; xx++ {
			ddx := float64(xx) + 0.5 - cx
			ddy := float64(yy) + 0.5 - cy
			if ddx*ddx+ddy*ddy <= r2 {
				p.blend(xx, yy, c)
			}
		}
	}
}

// polygon fills an arbitrary simple polygon using a scanline algorithm.
func (p *pngCanvas) polygon(pts []Point, fill Color) {
	if fill.isTransparent() || len(pts) < 3 {
		return
	}
	minY := pts[0].Y
	maxY := pts[0].Y
	for _, pt := range pts {
		if pt.Y < minY {
			minY = pt.Y
		}
		if pt.Y > maxY {
			maxY = pt.Y
		}
	}
	y0 := int(math.Floor(minY))
	y1 := int(math.Ceil(maxY))
	for y := y0; y <= y1; y++ {
		yc := float64(y) + 0.5
		var xs []float64
		n := len(pts)
		for i := 0; i < n; i++ {
			a := pts[i]
			b := pts[(i+1)%n]
			if (a.Y <= yc && b.Y > yc) || (b.Y <= yc && a.Y > yc) {
				t := (yc - a.Y) / (b.Y - a.Y)
				xs = append(xs, a.X+t*(b.X-a.X))
			}
		}
		if len(xs) < 2 {
			continue
		}
		sortFloats(xs)
		for i := 0; i+1 < len(xs); i += 2 {
			xa := int(math.Round(xs[i]))
			xb := int(math.Round(xs[i+1]))
			for x := xa; x < xb; x++ {
				p.blend(x, y, fill)
			}
		}
	}
}

func (p *pngCanvas) drawText(x, y float64, s string, c Color, anchor TextAnchor, size float64, vertical bool) {
	if c.isTransparent() || s == "" {
		return
	}
	scale := int(math.Round(size / float64(glyphH)))
	if scale < 1 {
		scale = 1
	}
	runes := []rune(s)
	textLen := len(runes)*glyphAdv*scale - scale // width along the text baseline
	textH := glyphH * scale

	if vertical {
		// Render rotated 90 degrees counter-clockwise, centered on (x,y).
		startY := y + float64(textLen)/2
		for ci, r := range runes {
			g, _ := glyphFor(r)
			base := startY - float64(ci*glyphAdv*scale)
			for row := 0; row < glyphH; row++ {
				line := g[row]
				for col := 0; col < glyphW; col++ {
					if col < len(line) && line[col] == '#' {
						px := x - float64(textH)/2 + float64(row*scale)
						py := base - float64(col*scale)
						for sy := 0; sy < scale; sy++ {
							for sx := 0; sx < scale; sx++ {
								p.blend(int(px)+sx, int(py)+sy, c)
							}
						}
					}
				}
			}
		}
		return
	}

	var startX float64
	switch anchor {
	case AnchorMiddle:
		startX = x - float64(textLen)/2
	case AnchorEnd:
		startX = x - float64(textLen)
	default:
		startX = x
	}
	startYPix := y - float64(textH)/2
	for ci, r := range runes {
		g, _ := glyphFor(r)
		gx := startX + float64(ci*glyphAdv*scale)
		for row := 0; row < glyphH; row++ {
			line := g[row]
			for col := 0; col < glyphW; col++ {
				if col < len(line) && line[col] == '#' {
					px := gx + float64(col*scale)
					py := startYPix + float64(row*scale)
					for sy := 0; sy < scale; sy++ {
						for sx := 0; sx < scale; sx++ {
							p.blend(int(px)+sx, int(py)+sy, c)
						}
					}
				}
			}
		}
	}
}

func toRGBA(c Color) color.RGBA {
	return color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func sortFloats(a []float64) {
	for i := 1; i < len(a); i++ {
		for j := i; j > 0 && a[j-1] > a[j]; j-- {
			a[j-1], a[j] = a[j], a[j-1]
		}
	}
}
