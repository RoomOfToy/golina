package matrix

import (
	"math"
)

// LUPDecompose does LUP decomposition of matrix
//	https://en.wikipedia.org/wiki/LU_decomposition
/* INPUT: t - array of pointers to rows of a square matrix having dimension N
 *        Tol - small tolerance number to detect failure when the matrix is near degenerate
 * OUTPUT: New Matrix nt, it contains both matrices L-E and U as nt=(L-E)+U such that P*nt=L*U.
 *        The permutation matrix is not stored as a matrix, but in an integer vector P of size N+1
 *        containing column indexes where the permutation matrix has "1". The last element P[N]=S+N,
 *        where S is the number of row exchanges needed for determinant computation, det(P)=(-1)^S
 */
func LUPDecompose(t *Matrix, N int, Tol float64) (*Matrix, *[]int) {
	nt := Copy(t)
	P := make([]int, N+1)
	imax := 0
	maxT := 0.
	i, j, k := 0, 0, 0
	tmp := Vector{}

	for i = 0; i <= N; i++ {
		P[i] = i //Unit permutation matrix, P[N] initialized with N
	}
	for i = 0; i < N; i++ {
		imax = i
		for k = i; k < N; k++ {
			if a := math.Abs(nt.At(k, i)); a > maxT {
				maxT = a
				imax = k
			}
		}

		if maxT < Tol { //failure, matrix is degenerate
			return nil, nil
		}

		if imax != i {
			// pivoting P
			j = P[i]
			P[i] = P[imax]
			P[imax] = j

			// pivoting rows of t
			tmp = nt.Data[i]
			nt.Data[i] = nt.Data[imax]
			nt.Data[imax] = tmp

			// counting pivots starting from N (for determinant)
			P[N]++
		}

		for j = i + 1; j < N; j++ {
			nt.Set(j, i, nt.At(j, i)/nt.At(i, i))

			for k = i + 1; k < N; k++ {
				nt.Set(j, k, nt.At(j, k)-nt.At(j, i)*nt.At(i, k))
			}
		}
	}
	return nt, &P
}

// LUPSolve LUP decomposition and solve equations
/* INPUT: A,P filled in LUPDecompose; b - rhs vector; N - dimension
 * OUTPUT: x - solution vector of A*x=b
 */
func LUPSolve(t *Matrix, P *[]int, N int, b *Vector) *Vector {
	x := make(Vector, N)
	for i := 0; i < N; i++ {
		x[i] = (*b)[(*P)[i]]

		for k := 0; k < i; k++ {
			x[i] -= t.At(i, k) * x[k]
		}
	}

	for i := N - 1; i >= 0; i-- {
		for k := i + 1; k < N; k++ {
			x[i] -= t.At(i, k) * x[k]
		}
		x[i] /= t.At(i, i)
	}
	return &x
}

// LUPInvert returns inverse matrix
/* INPUT: A,P filled in LUPDecompose; N - dimension
 * OUTPUT: IA is the inverse of the initial matrix
 */
func LUPInvert(t *Matrix, P *[]int, N int) *Matrix {
	nt := ZeroMatrix(N, N)
	for j := 0; j < N; j++ {
		for i := 0; i < N; i++ {
			if (*P)[i] == j {
				nt.Set(i, j, 1.0)
			} else {
				nt.Set(i, j, 0.0)
			}

			for k := 0; k < i; k++ {
				nt.Set(i, j, nt.At(i, j)-t.At(i, k)*nt.At(k, j))
			}
		}

		for i := N - 1; i >= 0; i-- {
			for k := i + 1; k < N; k++ {
				nt.Set(i, j, nt.At(i, j)-t.At(i, k)*nt.At(k, j))
			}
			nt.Set(i, j, nt.At(i, j)/t.At(i, i))
		}
	}
	return nt
}

// LUPDeterminant returns determinant of matrix
func LUPDeterminant(t *Matrix, P *[]int, N int) float64 {
	det := t.At(0, 0)

	for i := 1; i < N; i++ {
		det *= t.At(i, i)
	}

	if ((*P)[N]-N)%2 == 0 {
		return det
	}
	return -det
}

// LUPRank returns rank of matrix
func LUPRank(t *Matrix, N int) int {
	rank := 0
	for i := 0; i < N; i++ {
		if !FloatEqual(t.At(i, i), 0.) {
			rank++
		}
	}
	return rank
}
