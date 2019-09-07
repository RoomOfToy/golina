package numerical

import (
	"golina/matrix"
	"math"
	"testing"
)

func TestFuncFirstOrderDiff(t *testing.T) {
	if !matrix.FloatEqual(math.Cos(10), FuncFirstOrderDiff(math.Sin, 0)(10)) {
		t.Fail()
	}
}

func TestGaussianQuadrature(t *testing.T) {
	if !matrix.FloatEqual(GaussianQuadrature(math.Sin, -math.Pi/2, math.Pi/2, 1), 0) {
		t.Fail()
	}
}
