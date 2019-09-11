package stats

import (
	"golina/matrix"
	"testing"
)

func TestMahalanobisDistance(t *testing.T) {
	dataSet := &matrix.Matrix{Data: matrix.Data{{4, 3, 5}, {2, 6, 8}, {5, 3, 7}}}
	if !matrix.FloatEqual(MahalanobisDistance(dataSet.Row(0), nil, dataSet), 0.5773502691896252) {
		t.Fail()
	}
	if !matrix.FloatEqual(MahalanobisDistance(dataSet.Row(0), dataSet.Row(1), dataSet), 3.7416573867739413) {
		t.Fail()
	}
}

func TestMahalanobisDistanceXYVI(t *testing.T) {
	x, y := &matrix.Vector{1, 0, 0}, &matrix.Vector{0, 1, 0}
	vi := &matrix.Matrix{
		Data: matrix.Data{{1, 0.5, 0.5}, {0.5, 1, 0.5}, {0.5, 0.5, 1}},
	}
	if !matrix.FloatEqual(MahalanobisDistanceXYVI(x, y, vi), 1) {
		t.Fail()
	}
	x, y = &matrix.Vector{2, 0, 0}, &matrix.Vector{0, 1, 0}
	if !matrix.FloatEqual(MahalanobisDistanceXYVI(x, y, vi), 1.7320508075688772) {
		t.Fail()
	}
}
