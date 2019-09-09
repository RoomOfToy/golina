package stats

import (
	"golina/matrix"
	"testing"
)

func TestLinearKernel(t *testing.T) {
	if LinearKernel(&matrix.Vector{1, 2, 3}, &matrix.Vector{3, 2, 1}) != 10. {
		t.Fail()
	}
}

func TestPolyKernel(t *testing.T) {
	if PolyKernel(&matrix.Vector{1, 2, 3}, &matrix.Vector{3, 2, 1}, -1, 2, 3) != -512. {
		t.Fail()
	}
}

func TestRBFKernel(t *testing.T) {
	if !matrix.FloatEqual(RBFKernel(&matrix.Vector{1, 2, 3}, &matrix.Vector{3, 2, 1}, 0.5), 0.018315639) {
		t.Fail()
	}
}

func TestSigmoidKernel(t *testing.T) {
	if !matrix.FloatEqual(SigmoidKernel(&matrix.Vector{1, 2, 3}, &matrix.Vector{3, 2, 1}, 0.2, 0.5), 0.986614298) {
		t.Fail()
	}
}

func TestCovariance(t *testing.T) {
	if !matrix.FloatEqual(2.1666666666666665, Covariance(&matrix.Vector{1, 5, 7, 2, 6, 9}, &matrix.Vector{1, 2, 3, 1, 2, 3})) {
		t.Fail()
	}
}

func TestVariance(t *testing.T) {
	if !matrix.FloatEqual(7.666666666666667, Variance(&matrix.Vector{1, 5, 7, 2, 6, 9})) {
		t.Fail()
	}
}

func TestCorrelation(t *testing.T) {
	if !matrix.FloatEqual(0.95837272, Correlation(&matrix.Vector{1, 5, 7, 2, 6, 9}, &matrix.Vector{1, 2, 3, 1, 2, 3})) {
		t.Fail()
	}
}

func TestHistogram(t *testing.T) {
	data := &matrix.Vector{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	dividers := &matrix.Vector{-1, 2, 4, 5, 8, 11}
	hist := Histogram(dividers, data, nil)
	if !matrix.VEqual(hist, &matrix.Vector{2, 2, 1, 3, 3}) {
		t.Fail()
	}
}
