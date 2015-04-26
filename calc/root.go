package calc

import (
	"math"
)

const (
	__rootAcc = 10E-16 // desired accuracy of FindRoot
)

// reverse-quadratic interpolation
func __rqi(a, b, c float64) float64 {
	fa := __ff(a)
	fb := __ff(b)
	fc := __ff(c)
	s1 := a * (fb * fc) / ((fa - fb) * (fa - fc))
	s2 := b * (fa * fc) / ((fb - fa) * (fb - fc))
	s3 := c * (fa * fb) / ((fc - fa) * (fc - fb))
	return s1 + s2 + s3
}

func check(a, b, c, d, s float64, flag bool) bool {
	coef := (3*a + b) / 4
	if (s > b && s < coef && b > coef) || (s < b && s > coef && b < coef) {
		return true
	}
	if flag {
		bc := math.Abs(b - c)
		if math.Abs(s-b) >= (bc/2.0) || bc < __rootAcc {
			return true
		}
		return false
	}
	cd := math.Abs(c - d)
	if math.Abs(s-b) >= (cd/2.0) || cd < __rootAcc {
		return true
	}
	return false
}

// FindRoot implements "Brent's method" for root-finding.
// 'a' and 'b' should be two points such that f(a) and f(b)
// have opposite signs. FindRoot will find a point 'x' between
// 'a' and 'b' such that f(x) = 0.
func __FindRoot(a float64, b float64) (float64, bool) {
	fa, fb := __ff(a), __ff(b)
	if math.Signbit(a) == math.Signbit(b) {
		return math.NaN(), false
	}
	// 'b' should always be closer to the root
	if math.Abs(fa) < math.Abs(fb) {
		a, b = b, a
	}
	s := 0.0
	d := a
	c := a
	flag := true
	for math.Abs(b-a) > __rootAcc {
		if a != c && b != c {
			s = __rqi(a, b, c)
		} else {
			fb := __ff(b)
			s = b - fb*(b-a)/(fb-__ff(a))
		}
		if check(a, b, c, d, s, flag) {
			s = (a + b) / 2
			flag = true
		} else {
			flag = false
		}
		d, c = c, b
		if __ff(a)*__ff(s) < 0 {
			b = s
		} else {
			a = s
		}
		if math.Abs(__ff(a)) < math.Abs(__ff(b)) {
			a, b = b, a
		}
	}
	return s, true
}
