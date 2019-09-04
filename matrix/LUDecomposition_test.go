package matrix

import (
	"math"
	"strconv"
	"testing"
)

func TestLUPDecompose(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 5}, {30, 50, 10}}
	matA := new(Matrix).Init(a)
	nt, P := LUPDecompose(matA, 3, EPS)
	// fmt.Printf("%+v\n", nt)
	// fmt.Printf("%+v\n", P)
	p := []int{2, 1, 0, 4}
	for i := range *P {
		if (*P)[i] != p[i] {
			t.Fail()
		}
	}
	b := Data{{30, 50, 10}, {-2. / 3., 10. / 3., 35. / 3.}, {1. / 3., 1, -5}}
	if !MEqual(nt, new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestLUPSolve(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 5}, {30, 50, 10}}
	matA := new(Matrix).Init(a)
	nt, P := LUPDecompose(matA, 3, EPS)
	x := LUPSolve(nt, P, 3, &Vector{40, -40, 80})
	// fmt.Printf("%v\n", x)
	if !VEqual(x, &Vector{-4, 4, 0}) {
		t.Fail()
	}
}

func TestLUPInvert(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 5}, {30, 50, 10}}
	matA := new(Matrix).Init(a)
	nt, P := LUPDecompose(matA, 3, EPS)
	it := LUPInvert(nt, P, 3)
	// fmt.Printf("%v\n", it)
	b := Data{{-1.1, 0.6, 0.8}, {0.7, -0.4, -0.5}, {-0.2, 0.2, 0.2}}
	if !MEqual(it, new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestLUPDeterminant(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 5}, {30, 50, 10}}
	matA := new(Matrix).Init(a)
	nt, P := LUPDecompose(matA, 3, EPS)
	det := LUPDeterminant(nt, P, 3)
	// fmt.Printf("%v\n", det)
	if !FloatEqual(det, 500) {
		t.Fail()
	}
}

func TestLUPRank(t *testing.T) {
	matA := GenerateRandomSquareMatrix(3)
	nt, _ := LUPDecompose(matA, 3, EPS)
	rank := LUPRank(nt, 3)
	if rank != matA.Rank() {
		t.Fail()
	}
}

func BenchmarkLUPDeterminant(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			m := GenerateRandomSquareMatrix(n)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				nt, P := LUPDecompose(m, n, EPS)
				LUPDeterminant(nt, P, 3)
			}
		})
	}
}

func BenchmarkLUPRank(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			m := GenerateRandomSquareMatrix(n)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				nt, _ := LUPDecompose(m, n, EPS)
				LUPRank(nt, n)
			}
		})
	}
}

func BenchmarkLUPDecompose(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			m := GenerateRandomSquareMatrix(n)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				LUPDecompose(m, n, EPS)
			}
		})
	}
}
