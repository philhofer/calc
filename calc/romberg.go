package calc

import "math"

const (
	__romK   = 10     // Romberg steps
	__romAcc = 10E-16 // accuracy
)

// this gets compiled twice; once for 'ff' and
// once with 'ftrans'
func __trap(a float64, b float64, n int) float64 {
	if n == 0 {
		return a
	} else if n == 1 {
		return (b - a) * (__ff(b) + __ff(a)) / 2
	}
	h := (b - a) / float64(n) // step size

	first := __ff(a)
	middle := 0.0
	a += h
	for i := 1; i < n; i++ {
		middle += __ff(a)
		a += h
	}
	last := __ff(a)

	return (first + (2 * middle) + last) * (h / 2.0)
}

func __Integral(a float64, b float64) (float64, bool) {
	if a == b {
		return 0.0, true
	}
	if math.IsNaN(a) || math.IsNaN(b) {
		return math.NaN(), false
	}

	if math.IsInf(a, 0) || math.IsInf(b, 0) {
		return __rangeTransform(a, b)
	}

	// actual integration is done here
	// via the trapezoidal rule. the rest
	// is just richardson extrapolation.
	var k0 [__romK]float64
	for i := range k0 {
		k0[i] = __trap(a, b, 1<<uint(i))
	}
	// early check; we may have already converged
	curr := k0[__romK-1]
	if math.Abs(k0[__romK-2]-curr) <= __romAcc {
		return curr, true
	}

	// Ip always holds previous values; Ik holds
	// the unused (more previous) values
	Ip := k0[:]
	Ik := make([]float64, __romK)

	for k := 1; k < __romK; k++ {
		Icur := Ik[:__romK-k]
		lk := 1 << uint(k)
		m := lk * lk // 4^(k)
		for i := range Icur {
			Icur[i] = (float64(m)*Ip[i+1] - Ip[i]) / float64(m-1)
		}
		curr := Icur[__romK-k-1]
		if math.Abs(curr-Ip[__romK-k]) <= __romAcc {
			return curr, true
		}
		Ip, Ik = Icur, Ip[:__romK]
	}
	return Ip[len(Ip)-1], false
}

// a function x(z) such that (-Inf, +Inf) can be mapped
// onto (-1, 1)
func xof(z float64) float64 { return -z / ((z - 1) * (z + 1)) }

func __ftrans(z float64) float64 {
	zsq := z * z
	zsm := zsq - 1
	return __ff(xof(z)) * (zsq + 1) / (zsm * zsm)
}

func rangeSwap(a float64, b float64) (float64, float64) {
	var anew, bnew float64
	if a == 0.0 {
		anew = 0.0
	} else if math.IsInf(a, -1) {
		anew = -1.0 + (1E-15)
	} else {
		anew = (math.Sqrt(4*a*a+1) - 1.0) / (2 * a)
	}
	if b == 0.0 {
		bnew = 0.0
	} else if math.IsInf(b, 1) {
		bnew = 1.0 - (1E-15)
	} else {
		bnew = (math.Sqrt(4*b*b+1) - 1.0) / (2 * b)
	}
	return anew, bnew
}

func __rangeTransform(a float64, b float64) (float64, bool) {
	flipped := false
	if a > b {
		a, b = b, a
		flipped = true
	}
	a, b = rangeSwap(a, b)
	f, ok := __transIntegral(a, b)
	if flipped {
		f *= -1.0
	}
	return f, ok
}

// implemented by code generator; same as Integral, but uses
// alternate version of 'trap'
func __transIntegral(a, b float64) (float64, bool) { return 0.0, false }
