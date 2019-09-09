package stats

import (
	"golina/matrix"
)

// Linear Regression
//	https://en.wikipedia.org/wiki/Linear_regression
//	https://en.wikipedia.org/wiki/Simple_linear_regression
//	y = x * beta + alpha, linear regression tries to minimize Σw[i]*(y[i] - alpha - x[i]*beta)^2
//	rSquared = 1 - Σw[i]*(y[i] - alpha - x[i]*beta)^2 / Σw[i]*(y[i] - y.Mean())^2
//	if origin -> alpha = 0
func SimpleLinearRegression(x, y, weights *matrix.Vector, origin, calRSquared bool) (alpha, beta, rSquared float64) {
	if x.Length() != y.Length() {
		panic("x, y length mismatch")
	}
	if weights != nil && weights.Length() != x.Length() {
		panic("x, weights length mismatch")
	}
	w := 1.
	if origin { // y = x * beta
		alpha = 0.
		xx, xy := 0., 0.
		for i, xi := range *x {
			if weights != nil {
				w = weights.At(i)
			}
			xx += w * xi * xi
			xy += w * xi * y.At(i)
		}
		beta = xy / xx
		return
	}

	beta = Covariance(x, y) / Variance(x)
	alpha = y.Mean() - beta*x.Mean()

	if !calRSquared {
		return
	}

	w = 1.
	ym := y.Mean()
	numerator, denominator := 0., 0.
	for i, xi := range *x {
		if weights != nil {
			w = weights.At(i)
		}
		numerator += w * (y.At(i) - alpha - beta*xi) * (y.At(i) - alpha - beta*xi)
		denominator += w * (y.At(i) - ym) * (y.At(i) - ym)
	}
	rSquared = 1 - numerator/denominator
	return
}
