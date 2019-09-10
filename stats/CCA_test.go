package stats

import (
	"golina/matrix"
	"math"
	"strconv"
	"testing"
)

func TestCanonicalCorrelation(t *testing.T) {
	X := new(matrix.Matrix).Init(matrix.Data{
		{5.1, 3.5},
		{4.9, 3.0},
		{4.7, 3.2},
		{4.6, 3.1},
		{5.0, 3.6},
		{5.4, 3.9},
		{4.6, 3.4},
		{5.0, 3.4},
		{4.4, 2.9},
		{4.9, 3.1},
	})
	Y := new(matrix.Matrix).Init(matrix.Data{
		{1.4, 0.2},
		{1.4, 0.2},
		{1.3, 0.2},
		{1.5, 0.2},
		{1.4, 0.2},
		{1.7, 0.4},
		{1.4, 0.3},
		{1.5, 0.2},
		{1.4, 0.2},
		{1.5, 0.1},
	})
	A, B, r := CanonicalCorrelation(X, Y)
	if !matrix.MEqual(A, &matrix.Matrix{Data: matrix.Data{{-1.9794877596804641, -5.2016325219025124}, {4.5211829944066553, 2.7263663170835697}}}) ||
		!matrix.MEqual(B, &matrix.Matrix{Data: matrix.Data{{-0.0613084818030103, -10.8514169865438941}, {12.7209032660734298, 7.6793888180353775}}}) ||
		!matrix.VEqual(r, &matrix.Vector{0.7250624174504773, 0.5547679185730191}) {
		t.Fail()
	}
	//U := X.Sub(X.Mean(0).Tile(0, 10)).Mul(A)
	//V := Y.Sub(Y.Mean(0).Tile(0, 10)).Mul(B)
	//fmt.Println(U, V)
}

func BenchmarkCanonicalCorrelation(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n)+"x10", func(b *testing.B) {
			X := matrix.GenerateRandomMatrix(n, 10)
			Y := matrix.GenerateRandomMatrix(n, 10)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				CanonicalCorrelation(X, Y)
			}
		})
	}
}

/*
// Method 1
BenchmarkCanonicalCorrelation/size-10x10-8         	   20000	     93910 ns/op
BenchmarkCanonicalCorrelation/size-100x10-8        	    5000	    435422 ns/op
BenchmarkCanonicalCorrelation/size-1000x10-8       	     300	   4151420 ns/op

// Method 2
BenchmarkCanonicalCorrelation/size-10x10-8         	   20000	     94420 ns/op
BenchmarkCanonicalCorrelation/size-100x10-8        	    5000	    395215 ns/op
BenchmarkCanonicalCorrelation/size-1000x10-8       	     300	   3719333 ns/op
*/
