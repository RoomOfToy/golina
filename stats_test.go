package golina

import "testing"

func TestLinearKernel(t *testing.T) {
	if LinearKernel(&Vector{1, 2, 3}, &Vector{3, 2, 1}) != 10. {
		t.Fail()
	}
}

func TestPolyKernel(t *testing.T) {
	if PolyKernel(&Vector{1, 2, 3}, &Vector{3, 2, 1}, -1, 2, 3) != -512. {
		t.Fail()
	}
}

func TestRBFKernel(t *testing.T) {
	if !FloatEqual(RBFKernel(&Vector{1, 2, 3}, &Vector{3, 2, 1}, 0.5), 0.018315639) {
		t.Fail()
	}
}

func TestSigmoidKernel(t *testing.T) {
	if !FloatEqual(SigmoidKernel(&Vector{1, 2, 3}, &Vector{3, 2, 1}, 0.2, 0.5), 0.986614298) {
		t.Fail()
	}
}
