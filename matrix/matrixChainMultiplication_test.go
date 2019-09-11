package matrix

import (
	"testing"
)

func TestMatrixChainMultiplication(t *testing.T) {
	g := GenerateRandomMatrix
	A, B, C, D := g(40, 20), g(20, 30), g(30, 10), g(10, 30)
	res := MatrixChainMultiplication(A, B, C, D)
	res1 := A.Mul(B).Mul(C).Mul(D)
	if !MEqual(res, res1) {
		t.Fail()
	}

	A, B, C, D = g(10, 20), g(20, 30), g(30, 40), g(40, 30)
	res = MatrixChainMultiplication(A, B, C, D)
	res1 = A.Mul(B).Mul(C).Mul(D)
	if !MEqual(res, res1) {
		t.Fail()
	}
}
