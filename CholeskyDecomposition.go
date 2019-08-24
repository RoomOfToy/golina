package golina

import "math"

// https://en.wikipedia.org/wiki/Cholesky_decomposition
// A = L.Mul(L.T())
// Ljj = sqrt(Ajj - sum((Ljk) ** 2)_from_k=1_to_j-1)
// Lij = (1 / Ljj) * (Aij - sum(Lik * Ljk)_from_k=1_to_j-1)
func CholeskyDecomposition(t *Matrix) *Matrix {
	row, col := t.Dims()
	if row != col {
		panic("square matrix only")
	}
	L := ZeroMatrix(row, col)
	sum := 0.
	// Ljj first
	for i := 0; i < row; i++ {
		for j := 0; j <= i; j++ {
			sum = 0.
			if j == i {
				for k := 0; k < j; k++ {
					sum += L.At(j, k) * L.At(j, k)
				}
				L.Set(j, j, math.Sqrt(t.At(j, j)-sum))
			} else { // Lij
				for k := 0; k < j; k++ {
					sum += L.At(i, k) * L.At(j, k)
				}
				L.Set(i, j, (t.At(i, j)-sum)/L.At(j, j))
			}
		}
	}
	return L
}
