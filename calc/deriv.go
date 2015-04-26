package calc

import (
	"math"
)

func __fa0(x float64, h float64) float64 {
	return (__ff(x+h) - __ff(x-h)) / (2 * h)
}

func __fa1(x float64, h float64) float64 {
	return (4*__fa0(x, h/2) - __fa0(x, h)) / 3
}

func __fa2(x float64) float64 {
	const h0 = 0.0031622776601683793319988935444
	return (16*__fa1(x, h0/4) - __fa1(x, h0)) / 15
}

func __Deriv(x float64) (float64, bool) {
	const h0 = 0.0031622776601683793319988935444
	a1, a2, a3 := __fa0(x, h0), __fa1(x, h0), __fa2(x)
	return a3, math.Abs(a1-a2) >= math.Abs(a2-a3)
}
