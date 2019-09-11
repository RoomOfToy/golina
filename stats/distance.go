package stats

import (
	"golina/matrix"
	"math"
)

// Mahalanobis Distance
//	https://en.wikipedia.org/wiki/Mahalanobis_distance
// From wiki: the ellipsoid that best represents the set's probability distribution can be estimated by building the covariance matrix of the samples.
// The Mahalanobis distance is the distance of the test point from the center of mass divided by the width of the ellipsoid in the direction of the test point.
func MahalanobisDistance(x, y *matrix.Vector, dataSet *matrix.Matrix) float64 {
	cov := dataSet.CovMatrix()
	tmp := &matrix.Vector{}
	if x != nil && y == nil {
		tmp = x.Sub(dataSet.Mean(0))
	} else if x != nil && y != nil {
		tmp = x.Sub(y)
	} else {
		panic("at least input non-nil x as input vector")
	}
	return math.Sqrt((tmp.ToMatrix(1, x.Length()).Mul(cov.Inverse())).Row(0).Dot(tmp))
}

// Keep the same api from scipy: x, y -> Vector, vi -> Inverse of Covariance Matrix
//	https://docs.scipy.org/doc/scipy/reference/generated/scipy.spatial.distance.mahalanobis.html
func MahalanobisDistanceXYVI(x, y *matrix.Vector, vi *matrix.Matrix) float64 {
	if x.Length() != y.Length() {
		panic("x, y should have the same dimension")
	}
	tmp := x.Sub(y)
	return math.Sqrt(tmp.ToMatrix(1, x.Length()).Mul(vi).Row(0).Dot(tmp))
}
