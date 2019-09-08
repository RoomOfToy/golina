package stats

import (
	"golina/matrix"
	"math"
	"strconv"
	"testing"
)

func TestPrincipalComponents(t *testing.T) {
	points := matrix.GenerateRandomMatrix(10, 3)
	_, S, V := matrix.SVD(points.Sub(points.Mean(0).Tile(0, 10)))
	pc, weights := PrincipalComponents(points, nil)
	for i := range pc.Data {
		if !matrix.VEqual(pc.Col(i), V.Col(i)) && !matrix.VEqual(pc.Col(i).MulNum(-1), V.Col(i)) {
			t.Fail()
		}
	}
	for i := range *weights {
		f := S.At(i, i) * S.At(i, i) / float64(weights.Length()-1)
		if !matrix.FloatEqual(weights.At(i), f) && !matrix.FloatEqual(weights.At(i)*(-1), f) {
			t.Fail()
		}
	}
}

func BenchmarkPrincipalComponents(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n)+"x3", func(b *testing.B) {
			points := matrix.GenerateRandomMatrix(n, 3)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				PrincipalComponents(points, nil)
			}
		})
	}
}
