package stats

import (
	"golina/matrix"
	"testing"
)

func TestSimpleLinearRegression(t *testing.T) {
	x := &matrix.Vector{0, 1, 2, 3}
	y := &matrix.Vector{-1, -0.2, 0.9, 2.1}
	_, _, rSquared := SimpleLinearRegression(x, y, nil, false, true)
	if rSquared < 0.99 {
		t.Fail()
	}
}
