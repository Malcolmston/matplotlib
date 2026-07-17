// Library content for the matplotlib documentation site. Mirrors the shape used
// by the malcolmston/go landing site's data.ts so the sibling sites stay in sync.
export interface Lib {
  id: string; name: string; icon: string; accent: string; pkg: string; node: string;
  repo: string; docs: string; tagline: string; blurb: string; tags: string[];
  features: string[]; node_code: string; go_code: string; integrate: string;
}

export const NODE_ACCENT = '#8cc84b';

export const MATPLOTLIB: Lib = {
  id:"matplotlib", name:"matplotlib", icon:'<i class="fa-solid fa-chart-line"></i>', accent:"#11557c",
  pkg:"github.com/malcolmston/matplotlib", node:"matplotlib/matplotlib",
  repo:"https://github.com/malcolmston/matplotlib", docs:"https://malcolmston.github.io/matplotlib/",
  tagline:"Matplotlib-style plotting and charting in pure Go.",
  blurb:"A small, dependency-free plotting library for Go, inspired by Python's matplotlib. It exposes a "+
    "Figure/Axes drawing model with a pyplot-style convenience API on top, so you can build a chart either "+
    "object-first (via NewFigure/AddAxes) or with package-level calls (Plot, Bar, Title, Save). Six chart "+
    "types cover the common cases — line plots, scatter, vertical and horizontal bars, histograms with "+
    "automatic binning, and pie — each returned as a typed value with fluent per-series styling for color, "+
    "line width, marker and label. Axes handle the rest automatically: data-range to pixel mapping, axis "+
    "lines, numeric ticks and labels, an optional legend and grid, and a tab10 color cycle. Every figure "+
    "renders to two real formats from the standard library alone — PNG (a raster image with a built-in "+
    "bitmap font) and vector SVG — and identical input always produces identical bytes. No cgo, no "+
    "third-party modules; the import path and package are both matplotlib.",
  tags:["Figure/Axes","pyplot API","line plot","scatter","bar / barh","histogram","pie","tab10 cycle","legend & grid","PNG output","SVG output","deterministic","zero deps"],
  features:[
    "<b>Figure + Axes model</b> — a <code>Figure</code> owns size, DPI and the color cycle; <code>Figure.AddAxes</code> / <code>Figure.Subplots</code> add coordinate systems",
    "Line plots via <code>Axes.Plot</code> returning a <code>LinePlot</code> — overlay multiple series with one call each",
    "Scatter, bars and more — <code>Axes.Scatter</code> (<code>ScatterPlot</code>), <code>Axes.Bar</code> / <code>Axes.BarH</code> (<code>BarChart</code>), <code>Axes.Hist</code> (<code>Histogram</code>), <code>Axes.Pie</code> (<code>PieChart</code>)",
    "Fluent per-series styling — <code>SetColor</code>, <code>SetLabel</code>, <code>SetLineWidth</code>, <code>SetMarker</code> and <code>SetSize</code> chain off each plot",
    "Marker styles — <code>MarkerCircle</code>, <code>MarkerSquare</code>, <code>MarkerTriangle</code>, <code>MarkerCross</code>, <code>MarkerPlus</code> (and <code>MarkerNone</code>)",
    "Automatic decorations — <code>SetTitle</code>, <code>SetXLabel</code>/<code>SetYLabel</code>, numeric ticks, <code>Legend</code> and <code>Grid</code>, with <code>SetXLim</code>/<code>SetYLim</code> overrides",
    "tab10 color cycle — new series draw successive colors from <code>DefaultColors</code>; build your own with <code>RGB</code>, <code>RGBA</code> or <code>Hex</code>",
    "pyplot-style convenience API — package-level <code>FigSize</code>, <code>Plot</code>, <code>Bar</code>, <code>Title</code>, <code>Legend</code>, <code>Grid</code>, <code>Save</code> over an implicit current figure (<code>Gcf</code>/<code>Gca</code>/<code>Clf</code>)",
    "PNG output — <code>SavePNG</code> / <code>WritePNG</code> / <code>PNGBytes</code> render a raster <code>*image.RGBA</code> with Bresenham lines, filled shapes and a built-in bitmap font",
    "Vector SVG output — <code>SaveSVG</code> / <code>WriteSVG</code> / <code>RenderSVG</code> emit clean <code>&lt;line&gt;</code>/<code>&lt;rect&gt;</code>/<code>&lt;polyline&gt;</code>/<code>&lt;polygon&gt;</code>/<code>&lt;text&gt;</code> markup",
    "Format-by-extension <code>Save</code> plus <code>RenderImage</code> for direct pixel access",
    "Deterministic — identical input always produces identical bytes",
    "Zero dependencies — pure Go standard library (<code>image</code>, <code>image/png</code>, <code>math</code>), no cgo, nothing to audit but the toolchain"
  ],
  node_code:
`import matplotlib.pyplot as plt

fig, ax = plt.subplots()
ax.plot([0, 1, 2, 3], [0, 1, 4, 9], marker="o", label="x^2")
ax.plot([0, 1, 2, 3], [0, 2, 4, 6], label="2x")
ax.set_title("Demo")
ax.set_xlabel("x"); ax.set_ylabel("y")
ax.legend(); ax.grid(True)
fig.savefig("demo.png")`,
  go_code:
`import "github.com/malcolmston/matplotlib"

fig := matplotlib.NewFigure(640, 480)
ax := fig.AddAxes()
ax.Plot([]float64{0, 1, 2, 3}, []float64{0, 1, 4, 9}).
	SetLabel("x^2").SetMarker(matplotlib.MarkerCircle)
ax.Plot([]float64{0, 1, 2, 3}, []float64{0, 2, 4, 6}).SetLabel("2x")
ax.SetTitle("Demo").SetXLabel("x").SetYLabel("y").Legend().Grid(true)
_ = fig.SavePNG("demo.png")
_ = fig.SaveSVG("demo.svg")`,
  integrate:
`<span class="tok-c">// Build a figure explicitly, then style two overlaid series. New series</span>
<span class="tok-c">// pull successive colors from the tab10 cycle unless you override them.</span>
fig := matplotlib.NewFigure(640, 480)
ax := fig.AddAxes()
ax.Plot(xs, ys).SetLabel("signal").SetLineWidth(2).SetMarker(matplotlib.MarkerCircle)
ax.Scatter(xs, noise).SetLabel("samples").SetColor(matplotlib.RGB(0xd6, 0x27, 0x28)).SetSize(6)
ax.SetTitle("Overlay").SetXLabel("t").SetYLabel("value").Legend().Grid(true)

<span class="tok-c">// A categorical bar chart with an explicit hex color.</span>
blue, _ := matplotlib.Hex("#11557c")
bx := matplotlib.NewFigure(500, 320).AddAxes()
bx.Bar([]string{"a", "b", "c"}, []float64{3, 7, 2}).SetColor(blue)

<span class="tok-c">// The pyplot-style API drives an implicit current figure, then Save</span>
<span class="tok-c">// picks PNG or SVG from the file extension.</span>
matplotlib.FigSize(480, 360)
matplotlib.Hist(data, 20)
matplotlib.Title("Distribution")
_ = matplotlib.Save("hist.svg")`
};
