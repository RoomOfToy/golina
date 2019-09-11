package stats

import (
	"golina/matrix"
	"math"
	"sort"
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

// Mode: Most common value
//	https://en.wikipedia.org/wiki/Mode_(statistics)
//	since it may exist many modes, the return value is a Vector
//	TODO: any efficient way?
func Mode(x *matrix.Vector) *matrix.Vector {
	cnt := x.UniqueWithCount()
	type kv struct {
		key   float64
		value int
	}
	ss := make([]kv, len(cnt))
	i := 0
	for k, v := range cnt {
		ss[i] = kv{
			key:   k,
			value: v,
		}
		i++
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].value > ss[j].value
	})
	tmp := matrix.Vector{}
	for i := range ss {
		if ss[i].value != ss[0].value {
			break
		}
		tmp = append(tmp, ss[i].key)
	}
	return &tmp
}

func Variance(x *matrix.Vector) float64 {
	return x.SubNum(x.Mean()).SquareSum() / float64(x.Length())
}

func StandardDeviation(x *matrix.Vector) float64 {
	return math.Sqrt(Variance(x))
}

// Coefficient of Variance (CV) or Relative Standard Deviation (RSD)
//	https://en.wikipedia.org/wiki/Coefficient_of_variation
func CoefficientOfVariance(x *matrix.Vector) float64 {
	return StandardDeviation(x) / x.Mean()
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

// Correlation Coefficient
//	https://en.wikipedia.org/wiki/Correlation_and_dependence
func CorrelationCoefficient(x, y *matrix.Vector) float64 {
	if x.Length() != y.Length() {
		panic("x, y length mismatch")
	}
	// return Covariance(x, y) / math.Sqrt(Variance(x) * Variance(y))
	ex, ey := x.SubNum(x.Mean()), y.SubNum(y.Mean())
	return ex.Dot(ey) / math.Sqrt(ex.SquareSum()*ey.SquareSum())
}

// Bin number from variable formula
//	n: number of data points
//	choice: Sqrt, Sturges, Rice
//	TODO: add more choices: Doanes, ScottNormalReference, FreedmanDiaconis, MinCrossValidationEstiSquaredErr, ShimazakiShinomoto
func GetBinNum(n int, choice string) int {
	if n < 1 {
		panic("in valid data points number")
	}
	switch choice {
	case "Sqrt":
		return int(math.Ceil(math.Sqrt(float64(n))))
	case "Sturges":
		return int(math.Ceil(math.Log2(float64(n)))) + 1
	case "Rice":
		return int(math.Ceil(2 * math.Pow(float64(n), 1./3.)))
	default:
		panic("Sorry, only Sqrt, Sturges, Rice are supported up to now")
	}
}

func GetEqualBinWidth(binNum int, data *matrix.Vector) float64 {
	_, max := data.Max()
	_, min := data.Min()
	return (max - min) / float64(binNum)
}

// Histogram
//	https://en.wikipedia.org/wiki/Histogram
//	histogram of certain bin = Σw[i] * data[i] for certain divider range
//	divider: left close right open
func Histogram(dividers, data, weights *matrix.Vector) *matrix.Vector {
	if dividers.Length() < 2 {
		panic("histogram requires 2 dividers (lower, upper range) at least")
	}
	if !sort.Float64sAreSorted(*dividers) {
		dividers = dividers.SortedAscending()
	}
	if !sort.Float64sAreSorted(*data) {
		data = data.SortedAscending()
	}
	if dividers.At(0) > data.At(0) || dividers.At(-1) <= data.At(-1) { // left close right open
		panic("data range should be within divider range")
	}
	if dividers.Length() == 2 {
		return data
	}
	idx, comp := 0, dividers.At(1)
	cnt := make(matrix.Vector, dividers.Length()-1)
	w := 1.
	for i, x := range *data {
		if x < comp {
			if weights != nil {
				w = weights.At(i)
			}
			cnt[idx] += w
			continue
		}

		// in case of dividers has equal elements inside
		for j := idx + 1; j < dividers.Length(); j++ {
			if x < dividers.At(j+1) {
				idx = j
				comp = dividers.At(j + 1)
				break
			}
		}

		if weights != nil {
			w = weights.At(i)
		}
		cnt[idx] += w
	}
	return &cnt
}
