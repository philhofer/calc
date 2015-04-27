package tests

import (
	"math"
	"testing"
)

func isapprox(d, f float64) bool {
	return math.Abs(d-f) < 1E-8
}

func TestCubeRoot(t *testing.T) {
	f, ok := cubeFindRoot(-1.0, 1.0)
	if !ok {
		t.Fatal("root did not converge")
	}
	if f != 0.0 {
		t.Fatalf("expected 0.0; got %f", f)
	}
}

func BenchmarkFindCubeRoot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cubeFindRoot(-1.0, 1.0)
	}
}

func TestCubeDeriv(t *testing.T) {
	tests := []float64{0.0, 1.0, 3.0, -3.0, 8.0, 1.3897214}
	for _, f := range tests {
		d, ok := cubeDeriv(f)
		if !ok {
			t.Errorf("cubeDeriv did not converge for value %f", f)
		}
		if !isapprox(d, 3*(f*f)) {
			t.Errorf("cubeDeriv(%f) returned %f; expected %f", f, d, 3*(f*f))
		}
	}
}

func BenchmarkCubeDeriv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cubeDeriv(0.0)
	}
}

// indefinite integral of x^3
func cubeIndefIntegral(x float64) float64 { return (x * x * x * x) / 4 }

func TestCubeIntegral(t *testing.T) {
	tests := [][2]float64{
		{0, 1},
		{-1, 2},
		{-5, 5},
	}
	for _, ab := range tests {
		i, ok := cubeIntegral(ab[0], ab[1])
		if !ok {
			t.Errorf("cubeIntegral(%f, %f) did not converge", ab[0], ab[1])
		}
		if !isapprox(i, cubeIndefIntegral(ab[1])-cubeIndefIntegral(ab[0])) {
			t.Errorf("cubeIntegral(%f, %f) returned %f; expected %f", ab[0], ab[1], i, cubeIndefIntegral(ab[1])-cubeIndefIntegral(ab[0]))
		}
	}
}

func BenchmarkCubeIntegral(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cubeIntegral(-1, 1)
	}
}
