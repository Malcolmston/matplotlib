package matplotlib

import "math"

// Linspace returns num evenly spaced values covering the closed interval
// [start, stop]. It mirrors numpy.linspace, which matplotlib users routinely
// reach for when constructing plot coordinates.
//
// The endpoints are always included exactly: the first element equals start
// and, when num > 1, the last element equals stop. If num <= 0 an empty slice
// is returned; if num == 1 the single-element slice {start} is returned.
func Linspace(start, stop float64, num int) []float64 {
	if num <= 0 {
		return []float64{}
	}
	if num == 1 {
		return []float64{start}
	}
	out := make([]float64, num)
	step := (stop - start) / float64(num-1)
	for i := 0; i < num; i++ {
		out[i] = start + step*float64(i)
	}
	// Pin the final value to stop to avoid floating-point drift.
	out[num-1] = stop
	return out
}

// Arange returns evenly spaced values over the half-open interval [start, stop)
// using the given step, mirroring numpy.arange. A positive step walks upward
// until reaching or passing stop; a negative step walks downward. A zero step
// (or a step whose sign cannot reach stop) yields an empty slice.
func Arange(start, stop, step float64) []float64 {
	out := []float64{}
	if step == 0 {
		return out
	}
	if step > 0 {
		for v := start; v < stop; v += step {
			out = append(out, v)
		}
	} else {
		for v := start; v > stop; v += step {
			out = append(out, v)
		}
	}
	return out
}

// Logspace returns num values spaced evenly on a logarithmic scale, mirroring
// numpy.logspace. The returned values are base raised to the exponents
// Linspace(start, stop, num); that is, the exponents (not the values
// themselves) are linearly spaced. A base <= 0 is treated as 10.
func Logspace(start, stop float64, num int, base float64) []float64 {
	if base <= 0 {
		base = 10
	}
	exps := Linspace(start, stop, num)
	out := make([]float64, len(exps))
	for i, e := range exps {
		out[i] = math.Pow(base, e)
	}
	return out
}
