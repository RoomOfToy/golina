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
