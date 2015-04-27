package tests

//go:generate calc -func cube -method root
//go:generate calc -func cube -method deriv
//go:generate calc -func cube -method integral
func cube(x float64) float64 { return x * x * x }
