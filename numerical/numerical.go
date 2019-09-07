package numerical

type FloatFunc func(float64) float64

// first order differential of function f
func FuncFirstOrderDiff(f FloatFunc, h float64) func(float64) float64 {
	if h == 0. {
		h = 1e-3
	}
	return func(x float64) float64 {
		return (f(x+h) - f(x-h)) / (2 * h)
	}
}

// Some low-order quadrature rules over [-1, 1]
func GetGaussianQuadraturePointWeight(numOfPoint int) (points, weights []float64) {
	switch numOfPoint {
	case 1:
		points = []float64{0}
		weights = []float64{2}
	case 2:
		points = []float64{-0.5773502691896257, +0.5773502691896257}
		weights = []float64{1, 1}
	case 3:
		points = []float64{0, -0.7745966692414834, 0.7745966692414834}
		weights = []float64{0.8888888888888888, 0.5555555555555556, 0.5555555555555556}
	case 4:
		points = []float64{-0.3399810435848563, 0.3399810435848563, -0.8611363115940526, 0.8611363115940526}
		weights = []float64{0.6521451548625461, 0.6521451548625461, 0.3478548451374538, 0.3478548451374538}
	case 5:
		points = []float64{0, -0.5384693101056831, 0.5384693101056831, -0.9061798459386640, 0.9061798459386640}
		weights = []float64{0.5688888888888888, 0.4786286704993665, 0.4786286704993665, 0.2369268850561891, 0.2369268850561891}
	default:
		panic("Sorry but wikipedia only gives these five pairs...")
	}
	return
}

// Change of interval
func ChangeInterval(f FloatFunc, a, b float64) FloatFunc {
	return func(x float64) float64 {
		return (b - a) / 2 * f((b-a)/2*x+(a+b)/2)
	}
}

// Gaussian Quadrature (simplest case: Gauss–Legendre quadrature)
//	https://en.wikipedia.org/wiki/Gaussian_quadrature
func GaussianQuadrature(f FloatFunc, a, b float64, numOfPoint int) float64 {
	points, weights := GetGaussianQuadraturePointWeight(numOfPoint)
	nf := ChangeInterval(f, a, b)
	integral := 0.
	for i := 0; i < numOfPoint; i++ {
		integral += weights[i] * nf(points[i])
	}
	return integral
}
