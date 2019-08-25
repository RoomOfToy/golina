package golina

import (
	"testing"
)

func TestPointToPointDistance(t *testing.T) {
	p1 := &Vector{1, 2, 3}
	p2 := &Vector{5, 8, 6}
	if !FloatEqual(PointToPointDistance(p1, p2), 7.810249675906654) {
		t.Fail()
	}
}

func TestPointToLineDistance(t *testing.T) {
	pt := &Vector{1, 2, 3}
	linePt := &Vector{2, 2, 2}
	lineDir := &Vector{1, 1, 1}
	if !FloatEqual(PointToLineDistance(pt, linePt, lineDir), 1.4142135623730951) {
		t.Fail()
	}
}

func TestPointToPlaneDistance(t *testing.T) {
	pt := &Vector{1, 2, 3}
	planeCenter := &Vector{0, 1, 0}
	planeNormal := &Vector{0, 1, 0}
	if !FloatEqual(PointToPlaneDistance(pt, planeCenter, planeNormal), 1) {
		t.Fail()
	}
}

func TestDirectedHausdorffDistance(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	pts1 := new(Matrix).Init(a)
	b := Data{{32, 12, 3}, {6, 3, 52}, {9, 2, 15}}
	pts2 := new(Matrix).Init(b)
	hd := DirectedHausdorffDistance(pts1, pts2)
	if !FloatEqual(hd.distance, 43.474130238568314) || hd.lIndex != 1 || hd.rIndex != 2 {
		t.Fail()
	}
}

func TestTaxicabDistance(t *testing.T) {
	v1 := &Vector{1, 2, -3}
	v2 := &Vector{5, -8, 6}
	if !FloatEqual(TaxicabDistance(v1, v2), 23) {
		t.Fail()
	}
}

func TestEuclideanDistance(t *testing.T) {
	v1 := &Vector{1, 2, -3}
	v2 := &Vector{5, -8, 6}
	if !FloatEqual(EuclideanDistance(v1, v2), 14.035668847618199) {
		t.Fail()
	}
}

func TestSquaredEuclideanDistance(t *testing.T) {
	v1 := &Vector{1, 2, -3}
	v2 := &Vector{5, -8, 6}
	if !FloatEqual(SquaredEuclideanDistance(v1, v2), 197) {
		t.Fail()
	}
}

func TestMinkowskiDistance(t *testing.T) {
	v1 := &Vector{1, 2, -3}
	v2 := &Vector{5, -8, 6}
	if !FloatEqual(MinkowskiDistance(v1, v2, 1), TaxicabDistance(v1, v2)) {
		t.Fail()
	}
	if !FloatEqual(MinkowskiDistance(v1, v2, 2), EuclideanDistance(v1, v2)) {
		t.Fail()
	}
	if !FloatEqual(MinkowskiDistance(v1, v2, 3), 12.148614834158321) {
		t.Fail()
	}
}

func TestChebyshevDistance(t *testing.T) {
	v1 := &Vector{1, 2, -3}
	v2 := &Vector{5, -8, 6}
	if !FloatEqual(ChebyshevDistance(v1, v2), 10) {
		t.Fail()
	}
}

func TestHammingDistance(t *testing.T) {
	v1 := &Vector{1, 1, 1}
	v2 := &Vector{0, 1, 0}
	if !FloatEqual(HammingDistance(v1, v2), 2) {
		t.Fail()
	}
}

func TestCanberraDistance(t *testing.T) {
	v1 := &Vector{0, 3, 4}
	v2 := &Vector{7, 6, 3}
	if !FloatEqual(CanberraDistance(v1, v2), 1.476190476190476) {
		t.Fail()
	}
}