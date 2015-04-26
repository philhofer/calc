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
