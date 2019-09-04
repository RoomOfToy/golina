package matrix

import (
	"math"
	"strconv"
	"testing"
)

func TestSVD(t *testing.T) {
	a := Data{{8, -6, 2}, {-6, 7, -4}, {2, -4, 3}}
	matA := new(Matrix).Init(a)
	U, S, V := SVD(matA)
	if !MEqual(matA, U.Mul(S).Mul(V.T())) {
		t.Fail()
	}
}

func BenchmarkSVD(b *testing.B) {
	for k := 1.0; k <= 2; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			m := GenerateRandomMatrix(n, n)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				SVD(m)
			}
		})
	}
}
