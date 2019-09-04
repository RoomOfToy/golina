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
