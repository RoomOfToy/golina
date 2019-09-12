package stats

import (
	"fmt"
	"golina/matrix"
	"math"
)

// Independent Component Analysis
//	https://en.wikipedia.org/wiki/Independent_component_analysis
//	The following explanation are from wiki.
//	Two assumptions:
//		1. The source signals are independent of each other.
//		2. The values in each source signal have non-Gaussian distributions.
//	Three effects of mixing source signals:
//		1. Independence: As per assumption 1, the source signals are independent; however, their signal mixtures are not.
//		This is because the signal mixtures share the same source signals.
//		2. Normality: According to the Central Limit Theorem, the distribution of a sum of independent random variables
//		with finite variance tends towards a Gaussian distribution. Loosely speaking, a sum of two independent random
//		variables usually has a distribution that is closer to Gaussian than any of the two original variables.
//		Here we consider the value of each signal as the random variable.
//		3. Complexity: The temporal complexity of any signal mixture is greater than that of its simplest constituent
//		source signal.
// FastICA (https://en.wikipedia.org/wiki/FastICA)
func FastICA(C int, tol float64, maxIter int, whitening bool, nonLinearFunc func(w *matrix.Vector, X *matrix.Matrix) (wp *matrix.Vector), dataSet *matrix.Matrix) (W, S *matrix.Matrix) {
	N, _ := dataSet.Dims()
	if C > N || C < 0 {
		panic("independent components should be less or equal to observations")
	}

	if whitening {
		X, K := PreWhitening(C, dataSet)
		W = CalW(C, tol, maxIter, nonLinearFunc, X)
		fmt.Println(W.Dims())
		fmt.Println(K.Dims())
		S = W.Mul(K.T()).Mul(dataSet)
	} else {
		W = CalW(C, tol, maxIter, nonLinearFunc, dataSet)
		S = W.Mul(dataSet)
	}
	return
}

// Pre-whitening the data
func PreWhitening(C int, dataSet *matrix.Matrix) (X, K *matrix.Matrix) {
	// step 1: centering
	N, M := dataSet.Dims()
	data := dataSet.Sub(dataSet.Mean(0).Tile(0, N))
	// step 2: whitening

	/*
		// by Eigen
		// too slow for large dataset...
		cov := data.Mul(data.T()).MulNum(1. / float64(N))
		V, D := matrix.EigenDecompose(cov)
		for i := range D.Data {
			D.Data[i][i] = math.Sqrt(1. / D.Data[i][i])
		}
		return matrix.MatrixChainMultiplication(V, D, V.T(), data)
		// return V.Mul(D).Mul(V.T()).Mul(data)

		/*/
	// by SVD
	U, D, _ := matrix.SVD(data)
	K = matrix.ZeroMatrix(C, M)
	for i := range K.Data {
		for j := range K.Data[i] {
			K.Data[i][j] = U.Data[i][j] / D.Data[j][j]
		}
	}
	X = K.T().Mul(data).MulNum(math.Sqrt(float64(M)))
	return
	//*/
}

func CalW(C int, tol float64, maxIter int, nonLinearFunc func(w *matrix.Vector, X *matrix.Matrix) (wp *matrix.Vector), dataSet *matrix.Matrix) *matrix.Matrix {
	N, _ := dataSet.Dims()
	W := matrix.GenerateRandomMatrix(C, N)
	iter := make([]int, C)
	for i := 0; i < C; i++ {
		cnt := 0
		w := W.Row(i)
		for {
			wp := nonLinearFunc(w, dataSet)
			s := make(matrix.Vector, N)
			for j := 1; j < i-1; j++ {
				s.Add(W.Row(j).MulNum(wp.Dot(W.Row(j))))
			}
			wp = wp.Sub(&s)
			wp = wp.MulNum(1. / wp.Norm())
			lim := math.Abs(math.Abs(wp.Sub(w).AbsSum()))
			W.Data[i] = *wp
			cnt++
			if lim < tol || cnt >= maxIter {
				iter[i] = cnt
				break
			}
		}
	}
	fmt.Println("iteration times for each component: ", iter)
	return W
}

// standard non-linear function
//	w: 1xN, X: NxM, wp: 1xM
func FuncLogcosh(w *matrix.Vector, X *matrix.Matrix) (wp *matrix.Vector) {
	N, M := X.Dims()
	wtx := w.ToMatrix(1, N).Mul(X).Row(0)
	g, gg := make(matrix.Vector, M), make(matrix.Vector, M)
	for i := 0; i < M; i++ {
		g[i], gg[i] = logcosh(wtx.At(i))
	}
	wp = X.Mul(g.ToMatrix(M, 1)).Col(0).Sub(w.MulNum(gg.Sum())).MulNum(1. / float64(M))
	return
}

//	f: func, g: first derivative, gg: second derivative
func logcosh(u float64) (g, gg float64) {
	// f = math.Log(math.Cosh(u))
	g = math.Tanh(u)
	gg = 1 - g*g
	return
}

// highly robust
func FuncExp(w *matrix.Vector, X *matrix.Matrix) (wp *matrix.Vector) {
	N, M := X.Dims()
	wtx := w.ToMatrix(1, N).Mul(X).Row(0)
	g, gg := make(matrix.Vector, M), make(matrix.Vector, M)
	for i := 0; i < M; i++ {
		g[i], gg[i] = exp(wtx.At(i))
	}
	wp = X.Mul(g.ToMatrix(M, 1)).Col(0).Sub(w.MulNum(gg.Sum())).MulNum(1. / float64(M))
	return
}

func exp(u float64) (g, gg float64) {
	u2 := u * u
	eu2 := math.Exp(-u2 / 2)
	// f = -eu2
	g = u * eu2
	gg = (1 - u2) * eu2
	return
}
