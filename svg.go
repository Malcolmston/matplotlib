package matplotlib

import (
	"fmt"
	"strings"
)

// svgCanvas is the vector (SVG) implementation of the canvas interface. It
// accumulates SVG elements into a buffer that is wrapped by the figure into a
// complete document.
type svgCanvas struct {
	w   int
	h   int
	buf strings.Builder
}

func newSVGCanvas(w, h int) *svgCanvas {
	return &svgCanvas{w: w, h: h}
}

func (s *svgCanvas) size() (int, int) { return s.w, s.h }

func (s *svgCanvas) fillBackground(c Color) {
	fmt.Fprintf(&s.buf, "<rect x=\"0\" y=\"0\" width=\"%d\" height=\"%d\" fill=\"%s\"/>\n",
		s.w, s.h, c.hexString())
}

func (s *svgCanvas) fillRect(x, y, w, h float64, c Color) {
	if c.isTransparent() {
		return
	}
	if w < 0 {
		x += w
		w = -w
	}
	if h < 0 {
		y += h
		h = -h
	}
	fmt.Fprintf(&s.buf, "<rect x=\"%s\" y=\"%s\" width=\"%s\" height=\"%s\" fill=\"%s\"/>\n",
		fnum(x), fnum(y), fnum(w), fnum(h), c.hexString())
}

func (s *svgCanvas) strokeLine(x0, y0, x1, y1 float64, c Color, lw float64) {
	if c.isTransparent() {
		return
	}
	fmt.Fprintf(&s.buf, "<line x1=\"%s\" y1=\"%s\" x2=\"%s\" y2=\"%s\" stroke=\"%s\" stroke-width=\"%s\"/>\n",
		fnum(x0), fnum(y0), fnum(x1), fnum(y1), c.hexString(), fnum(lw))
}

func (s *svgCanvas) polyline(pts []Point, c Color, lw float64) {
	if c.isTransparent() || len(pts) < 2 {
		return
	}
	var b strings.Builder
	for i, p := range pts {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%s,%s", fnum(p.X), fnum(p.Y))
	}
	fmt.Fprintf(&s.buf, "<polyline points=\"%s\" fill=\"none\" stroke=\"%s\" stroke-width=\"%s\"/>\n",
		b.String(), c.hexString(), fnum(lw))
}

func (s *svgCanvas) polygon(pts []Point, fill Color) {
	if fill.isTransparent() || len(pts) < 3 {
		return
	}
	var b strings.Builder
	for i, p := range pts {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%s,%s", fnum(p.X), fnum(p.Y))
	}
	fmt.Fprintf(&s.buf, "<polygon points=\"%s\" fill=\"%s\"/>\n", b.String(), fill.hexString())
}

func (s *svgCanvas) fillCircle(cx, cy, r float64, c Color) {
	if c.isTransparent() || r <= 0 {
		return
	}
	fmt.Fprintf(&s.buf, "<circle cx=\"%s\" cy=\"%s\" r=\"%s\" fill=\"%s\"/>\n",
		fnum(cx), fnum(cy), fnum(r), c.hexString())
}

func (s *svgCanvas) drawText(x, y float64, str string, c Color, anchor TextAnchor, size float64, vertical bool) {
	if c.isTransparent() || str == "" {
		return
	}
	var ta string
	switch anchor {
	case AnchorMiddle:
		ta = "middle"
	case AnchorEnd:
		ta = "end"
	default:
		ta = "start"
	}
	// Approximate vertical centering on the anchor point.
	baseline := y + size*0.35
	transform := ""
	if vertical {
		transform = fmt.Sprintf(" transform=\"rotate(-90 %s %s)\"", fnum(x), fnum(y))
		baseline = y + size*0.35
		ta = "middle"
	}
	fmt.Fprintf(&s.buf,
		"<text x=\"%s\" y=\"%s\" font-family=\"monospace\" font-size=\"%s\" fill=\"%s\" text-anchor=\"%s\"%s>%s</text>\n",
		fnum(x), fnum(baseline), fnum(size), c.hexString(), ta, transform, escapeXML(str))
}

func (s *svgCanvas) document() string {
	var out strings.Builder
	out.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	fmt.Fprintf(&out,
		"<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"%d\" height=\"%d\" viewBox=\"0 0 %d %d\">\n",
		s.w, s.h, s.w, s.h)
	out.WriteString(s.buf.String())
	out.WriteString("</svg>\n")
	return out.String()
}

func escapeXML(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&apos;",
	)
	return r.Replace(s)
}
