# Changelog

All notable changes to this project are documented here. This project adheres
to semantic versioning.

## [0.2.0]

Feature release expanding coverage of matplotlib's chart types, color tooling
and numeric helpers. All additions are pure standard library, deterministic and
covered by known-answer tests.

### Added

- **New chart types** (`plots_extra.go`), each composing with the existing
  Axes/legend/autoscale machinery:
  - `Axes.Step` / `StepPlot` — piecewise-constant staircase lines with
    `"pre"`, `"post"` and `"mid"` placement (matplotlib `step`).
  - `Axes.FillBetween` / `FillBetween` — shaded area between two curves, or
    down to zero (matplotlib `fill_between`).
  - `Axes.Stem` / `StemPlot` — stem plots with a baseline and tip markers
    (matplotlib `stem`).
  - `Axes.ErrorBar` / `ErrorBar` — points with symmetric vertical error bars
    and caps (matplotlib `errorbar`).
  - `Axes.AxHLine` / `Axes.AxVLine` / `RefLine` — full-span reference lines
    (matplotlib `axhline` / `axvline`).
  - `Axes.Text` / `Annotation` — text labels placed in data coordinates
    (matplotlib `text`).
- **Colormaps** (`colormap.go`): `Colormap` with `NewColormap`, `At` and
  `Reversed`; `Normalize` with `NewNormalize`, `Norm` and `Map`; and the named
  constructors `Viridis`, `Plasma`, `Grays`, `Jet` and `Coolwarm`.
- **Color operations** (`color_ops.go`): `Color.Lighten`, `Color.Darken`,
  `Color.WithAlpha`, `Color.Luminance`, `Color.Grayscale`, `Color.Hex` and the
  package-level `Blend`.
- **Numeric helpers** (`numeric.go`): `Linspace`, `Arange` and `Logspace`,
  mirroring the numpy routines matplotlib users rely on for coordinates.

## [0.1.0]

Initial release: Figure/Axes model, pyplot-style API, line/scatter/bar/barh/
histogram/pie charts, and deterministic PNG and SVG output.
