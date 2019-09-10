package stats

import (
	"golina/matrix"
	"math"
)

// Canonical Correlation
//	https://en.wikipedia.org/wiki/Canonical_correlation
//	https://ww2.mathworks.cn/help/stats/canoncorr.html
//	http://numerical.recipes/whp/notes/CanonCorrBySVD.pdf
//	TODO: SVD has sign uncertainty, will it affect?
func CanonicalCorrelation(X, Y *matrix.Matrix) (A, B *matrix.Matrix, r *matrix.Vector) {
	xm, _ := X.Dims()
	ym, _ := Y.Dims()
	if xm != ym {
		panic("X, Y should have the same number of rows (observations)")
	}
	Ux, Sx, Vx := matrix.SVD(X.Sub(X.Mean(0).Tile(0, xm)))
	Uy, Sy, Vy := matrix.SVD(Y.Sub(Y.Mean(0).Tile(0, xm)))

	/*
		// Method 1
		U, S, V := matrix.SVD(Vx.Mul(Ux.T()).Mul(Uy).Mul(Vy.T()))

		for i := range Vx.Data {
			for j := range Vx.Data[i] {
				Vx.Data[i][j] *= math.Sqrt(1. / Sx.Data[j][j])
			}
		}
		for i := range Vy.Data {
			for j := range Vy.Data[i] {
				Vy.Data[i][j] *= math.Sqrt(1. / Sy.Data[j][j])
			}
		}
		A = Vx.Mul(Vx.T()).Mul(U).MulNum(math.Sqrt(float64(xm - 1)))
		B = Vy.Mul(Vy.T()).Mul(V).MulNum(math.Sqrt(float64(xm - 1)))

		/*/
	// Method 2
	U, S, V := matrix.SVD(Ux.T().Mul(Uy))
	A, B = Vx.Mul(Sx.Inverse()).Mul(U).MulNum(math.Sqrt(float64(xm-1))), Vy.Mul(Sy.Inverse()).Mul(V).MulNum(math.Sqrt(float64(xm-1)))
	//*/

	d, _ := S.Dims()
	rr := make(matrix.Vector, d)
	for i := range S.Data {
		rr[i] = S.Data[i][i]
	}
	r = &rr
	return
}
