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

// Covariance
//	biased -> divide by x.Length(), unbiased -> divide by x.Length() - 1
//	biased here
func Covariance(x, y *matrix.Vector) float64 {
	if x.Length() != y.Length() {
		panic("x, y length mismatch")
	}
	return x.SubNum(x.Mean()).Dot(y.SubNum(y.Mean())) / float64(x.Length())
}

func Variance(x *matrix.Vector) float64 {
	return x.SubNum(x.Mean()).SquareSum() / float64(x.Length())
}

func Correlation(x, y *matrix.Vector) float64 {
	if x.Length() != y.Length() {
		panic("x, y length mismatch")
	}
	// return Covariance(x, y) / math.Sqrt(Variance(x) * Variance(y))
	ex, ey := x.SubNum(x.Mean()), y.SubNum(y.Mean())
	return ex.Dot(ey) / math.Sqrt(ex.SquareSum()*ey.SquareSum())
}
