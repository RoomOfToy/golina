package stats

import (
	"golina/matrix"
	"math"
)

// a: Row/Col i, b: Row/Col j
func LinearKernel(a, b *matrix.Vector) float64 {
	return a.Dot(b)
}

func PolyKernel(a, b *matrix.Vector, gamma, coef0, degree float64) float64 {
	return math.Pow(gamma*a.Dot(b)+coef0, degree)
}

// Radial basis function kernel
func RBFKernel(a, b *matrix.Vector, gamma float64) float64 {
	return math.Exp(-gamma * a.Sub(b).SquareSum())
}

func SigmoidKernel(a, b *matrix.Vector, gamma, coef0 float64) float64 {
	return math.Tanh(gamma*a.Dot(b) + coef0)
}
