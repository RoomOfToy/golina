package golina

import (
	"math"
	"strconv"
	"testing"
)

func TestSVD(t *testing.T) {
	a := Data{{8, -6, 2}, {-6, 7, -4}, {2, -4, 3}}
	matA := new(Matrix).Init(a)
	U, S, V := SVD(matA)
	u := Data{{2. / 3., -2. / 3., -1. / 3.}, {-2. / 3., -1. / 3., -2. / 3.}, {1. / 3., 2. / 3., -2. / 3.}}
	s := Data{{15, 0, 0}, {0, 3, 0}, {0, 0, 0}}
	v := Data{{2. / 3., -2. / 3., -1. / 3.}, {-2. / 3., -1. / 3., -2. / 3.}, {1. / 3., 2. / 3., -2. / 3.}}
	if !Equal(U, new(Matrix).Init(u)) || !Equal(S, new(Matrix).Init(s)) || !Equal(V, new(Matrix).Init(v)) {
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
