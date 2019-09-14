package cluster

import (
	"golina/matrix"
	"log"
	"math"
	"math/rand"
)

type Kernel func(x, y *matrix.Vector) float64

// SVM based on SMO (http://cs229.stanford.edu/materials/smo.pdf)
// Input:
//	C: regularization parameter
//	tol: numerical tolerance
//	maxIter: max # of times to iterate over α’s without changing
//	kernel: kernel function (in `stats` package)
//	dataSet: (x(1), y(1)), . . . ,(x(m), y(m)): training data, X features, y labels
// Output:
//	predictions: Vector
// TODO: result not correct...
func SVM(C, tol float64, maxIter int, kernel Kernel, dataSet *matrix.Matrix) *matrix.Vector {
	// train
	// initialize
	if dataSet == nil {
		panic("invalid input data")
	}
	row, col := dataSet.Dims()
	y := dataSet.Col(col - 1)
	X := dataSet.GetSubMatrix(0, 0, row, col-1)
	m := y.Length()
	alphaPrev, alpha, svIdx, b, iter := make(matrix.Vector, m), make(matrix.Vector, m), matrix.ARRange(0, 1, m), 0., 0
	// initialize supporting vector idx
	K := matrix.ZeroMatrix(row, row)
	// calculate kernel matrix
	for i := range K.Data {
		for j := range K.Data[i] {
			K.Data[i][j] = kernel(X.Row(i), X.Row(j))
		}
	}
	for iter < maxIter {
		iter++
		copy(alphaPrev, alpha)
		η := 0.
		for i := 0; i < m; i++ {
			j := pickRandomIdx(i, m)
			η = 2.0*K.At(i, j) - K.At(i, i) - K.At(j, j)
			if η >= 0 {
				continue
			}
			L, H := findBounds(i, j, C, &alpha, y)
			// calculate err
			errI, errJ := predictRow(X, X.Row(i), y, svIdx, &alpha, b, kernel)-y.At(i), predictRow(X, X.Row(j), y, svIdx, &alpha, b, kernel)-y.At(j)
			// save old alpha
			alphaIOld, alphaJOld := alpha[i], alpha[j]
			// update alpha
			alpha[j] -= (y.At(j) * (errI - errJ)) / η
			alpha[j] = matrix.Ternary(alpha[j] > H, H, alpha[j]).(float64)
			alpha[j] = matrix.Ternary(alpha[j] < L, L, alpha[j]).(float64)
			alpha[i] = alpha[i] + y.At(i)*y.At(j)*(alphaJOld-alpha[j])
			// find new b
			b1 := b - errI - y.At(i)*(alpha[i]-alphaJOld)*K.At(i, i) - y.At(j)*(alpha[j]-alphaJOld)*K.At(i, j)
			b2 := b - errJ - y.At(j)*(alpha[j]-alphaJOld)*K.At(j, j) - y.At(i)*(alpha[i]-alphaIOld)*K.At(i, j)
			if alpha[i] > 0 && alpha[i] < C {
				b = b1
			} else if alpha[j] > 0 && alpha[j] < C {
				b = b2
			} else {
				b = 0.5 * (b1 + b2)
			}
		}
		// check convergence
		diff := alpha.Sub(&alphaPrev).Norm()
		if diff < tol {
			break
		}
	}
	log.Printf("Reach Convergence with %d iterations\n", iter)
	// update support vectors index
	nsvIdx := matrix.Vector{}
	for i := range alpha {
		if alpha[i] > 0 {
			nsvIdx = append(nsvIdx, float64(i))
		}
	}
	svIdx = &nsvIdx

	// predict
	result := make(matrix.Vector, row)
	for i := range result {
		result[i] = sign(predictRow(X, X.Row(i), y, svIdx, &alpha, b, kernel))
	}
	return &result
}

func sign(a float64) float64 {
	if a > 0 {
		return 1
	} else if a < 0 {
		return 0
	}
	return -1
}

func predictRow(X *matrix.Matrix, x, y, svIdx *matrix.Vector, alpha *matrix.Vector, b float64, kernel Kernel) float64 {
	svIdx = svIdx.Unique()
	kv := make(matrix.Vector, svIdx.Length())
	for i := range kv {
		kv[i] = kernel(X.Row(int(svIdx.At(i))), x)
	}
	f := make(matrix.Vector, svIdx.Length())
	for i := range f {
		f[i] = alpha.At(int(svIdx.At(i))) * y.At(int(svIdx.At(i)))
	}
	return f.Dot(&kv) + b
}

func findBounds(i, j int, C float64, alpha, y *matrix.Vector) (L, H float64) {
	// L <= alpha <= H while 0 <= alpha <= C
	if y.At(i) != y.At(j) {
		L = math.Max(0, alpha.At(j)-alpha.At(i))
		H = math.Min(C, C-alpha.At(i)+alpha.At(j))
	} else {
		L = math.Max(0, alpha.At(i)+alpha.At(j)-C)
		H = math.Min(C, alpha.At(i)+alpha.At(j))
	}
	return
}

func pickRandomIdx(i, n int) int {
	j := i
	for j == i {
		j = rand.Intn(n)
	}
	return j
}
