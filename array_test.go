package golina

import (
	"testing"
)

// https://blog.karenuorteva.fi/go-unit-test-setup-and-teardown-db1601a796f2#.2aherx2z5

func TestMatrix_Init(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	if matA._array == nil {
		t.Fail()
	}
}

func TestMatrix_Dims(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	row, col := matA.Dims()
	if row != 3 || col != 3 {
		t.Fail()
	}
}

func TestMatrix_At(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	if matA.At(1, 1) != 5 {
		t.Fail()
	}
}

func TestMatrix_Set(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	matA.Set(1, 1, 10)
	if matA.At(1, 1) != 10 {
		t.Fail()
	}
}

func TestMatrix_T(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	matAT := matA.T()
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if matAT.At(i, j) != matA.At(j, i) {
				t.Fail()
			}
		}
	}
}
