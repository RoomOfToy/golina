package golina

import (
	"math"
	"strconv"
	"testing"
)

func TestMatrix_CholeskyDecomposition(t *testing.T) {
	a := Data{{4, 12, -16}, {12, 37, -43}, {-16, -43, 98}}
	matA := new(Matrix).Init(a)
	b := Data{{2, 0, 0}, {6, 1, 0}, {-8, 5, 3}}
	if !MEqual(CholeskyDecomposition(matA), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func BenchmarkCholeskyDecomposition(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			m := GenerateRandomSquareMatrix(n)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				CholeskyDecomposition(m)
			}
		})
	}
}
