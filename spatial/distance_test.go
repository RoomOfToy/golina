package spatial

import (
	"fmt"
	"golina/matrix"
	"math"
	"strconv"
	"testing"
)

func TestPointToPointDistance(t *testing.T) {
	p1 := &matrix.Vector{1, 2, 3}
	p2 := &matrix.Vector{5, 8, 6}
	if !matrix.FloatEqual(PointToPointDistance(p1, p2), 7.810249675906654) {
		t.Fail()
	}
}

func TestPointToLineDistance(t *testing.T) {
	pt := &matrix.Vector{1, 2, 3}
	linePt := &matrix.Vector{2, 2, 2}
	lineDir := &matrix.Vector{1, 1, 1}
	if !matrix.FloatEqual(PointToLineDistance(pt, linePt, lineDir), 1.4142135623730951) {
		t.Fail()
	}
}

func TestPointToPlaneDistance(t *testing.T) {
	pt := &matrix.Vector{1, 2, 3}
	planeCenter := &matrix.Vector{0, 1, 0}
	planeNormal := &matrix.Vector{0, 1, 0}
	if !matrix.FloatEqual(PointToPlaneDistance(pt, planeCenter, planeNormal), 1) {
		t.Fail()
	}
}

func TestDirectedHausdorffDistance(t *testing.T) {
	a := matrix.Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	pts1 := new(matrix.Matrix).Init(a)
	b := matrix.Data{{32, 12, 3}, {6, 3, 52}, {9, 2, 15}}
	pts2 := new(matrix.Matrix).Init(b)
	hd := DirectedHausdorffDistance(pts1, pts2)
	if !matrix.FloatEqual(hd.Distance, 43.474130238568314) || hd.LIndex != 1 || hd.RIndex != 2 {
		t.Fail()
	}
}

func TestDirectedHausdorffDistanceBasedOnKNN(t *testing.T) {
	a := matrix.Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	pts1 := new(matrix.Matrix).Init(a)
	b := matrix.Data{{32, 12, 3}, {6, 3, 52}, {9, 2, 15}}
	pts2 := new(matrix.Matrix).Init(b)
	hd := DirectedHausdorffDistanceBasedOnKNN(pts1, pts2)
	if !matrix.FloatEqual(hd.Distance, 43.474130238568314) || hd.LIndex != 1 || hd.RIndex != 2 {
		fmt.Println(hd)
		t.Fail()
	}
}

func TestTaxicabDistance(t *testing.T) {
	v1 := &matrix.Vector{1, 2, -3}
	v2 := &matrix.Vector{5, -8, 6}
	if !matrix.FloatEqual(TaxicabDistance(v1, v2), 23) {
		t.Fail()
	}
}

func TestEuclideanDistance(t *testing.T) {
	v1 := &matrix.Vector{1, 2, -3}
	v2 := &matrix.Vector{5, -8, 6}
	if !matrix.FloatEqual(EuclideanDistance(v1, v2), 14.035668847618199) {
		t.Fail()
	}
}

func TestSquaredEuclideanDistance(t *testing.T) {
	v1 := &matrix.Vector{1, 2, -3}
	v2 := &matrix.Vector{5, -8, 6}
	if !matrix.FloatEqual(SquaredEuclideanDistance(v1, v2), 197) {
		t.Fail()
	}
}

func TestMinkowskiDistance(t *testing.T) {
	v1 := &matrix.Vector{1, 2, -3}
	v2 := &matrix.Vector{5, -8, 6}
	if !matrix.FloatEqual(MinkowskiDistance(v1, v2, 1), TaxicabDistance(v1, v2)) {
		t.Fail()
	}
	if !matrix.FloatEqual(MinkowskiDistance(v1, v2, 2), EuclideanDistance(v1, v2)) {
		t.Fail()
	}
	if !matrix.FloatEqual(MinkowskiDistance(v1, v2, 3), 12.148614834158321) {
		t.Fail()
	}
}

func TestChebyshevDistance(t *testing.T) {
	v1 := &matrix.Vector{1, 2, -3}
	v2 := &matrix.Vector{5, -8, 6}
	if !matrix.FloatEqual(ChebyshevDistance(v1, v2), 10) {
		t.Fail()
	}
}

func TestHammingDistance(t *testing.T) {
	v1 := &matrix.Vector{1, 1, 1}
	v2 := &matrix.Vector{0, 1, 0}
	if !matrix.FloatEqual(HammingDistance(v1, v2), 2) {
		t.Fail()
	}
}

func TestCanberraDistance(t *testing.T) {
	v1 := &matrix.Vector{0, 3, 4}
	v2 := &matrix.Vector{7, 6, 3}
	if !matrix.FloatEqual(CanberraDistance(v1, v2), 1.476190476190476) {
		t.Fail()
	}
}

func BenchmarkDirectedHausdorffDistance(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n)+"x3", func(b *testing.B) {
			pts1 := matrix.GenerateRandomMatrix(n, 3)
			pts2 := matrix.GenerateRandomMatrix(n, 3)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				DirectedHausdorffDistance(pts1, pts2)
			}
		})
	}
}

func BenchmarkDirectedHausdorffDistanceBasedOnKNN(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n)+"x3", func(b *testing.B) {
			pts1 := matrix.GenerateRandomMatrix(n, 3)
			pts2 := matrix.GenerateRandomMatrix(n, 3)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				DirectedHausdorffDistanceBasedOnKNN(pts1, pts2)
			}
		})
	}
}

/*
BenchmarkDirectedHausdorffDistance/size-10x3-8                           5000000               467 ns/op
BenchmarkDirectedHausdorffDistance/size-100x3-8                           200000             21126 ns/op
BenchmarkDirectedHausdorffDistance/size-1000x3-8                            3000            428309 ns/op
BenchmarkDirectedHausdorffDistanceBasedOnKNN/size-10x3-8                   50000             31167 ns/op
BenchmarkDirectedHausdorffDistanceBasedOnKNN/size-100x3-8                   2000            882729 ns/op
BenchmarkDirectedHausdorffDistanceBasedOnKNN/size-1000x3-8                   100          89443641 ns/op
*/
