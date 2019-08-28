package golina

import (
	"math"
)

// https://en.wikipedia.org/wiki/Singular_value_decomposition
// 	A: m * n (m >= n)
// 	U: m * n orthogonal matrix, U * U.T() = Im
// 	S: n * n diagonal matrix
// 	V: n * n orthogonal matrix, V * V.T() = Im
//
// Step by step solution: https://atozmath.com/MatrixEv.aspx?q=svd
//
// 	Code from `Jama`, derived from LINPACK code
// 	https://github.com/fiji/Jama/blob/master/src/main/java/Jama/SingularValueDecomposition.java
func SVD(t *Matrix) (U, S, V *Matrix) {
	// Initialize
	A := Copy(t)._array
	m, n := t.Dims()
	if m < n {
		panic("error dimensions of input matrix, rows should be larger or equal to columns")
	}

	nu := min(m, n)
	s := make(Vector, min(m+1, n))

	u := make(Data, m)
	v := make(Data, n)
	for i := 0; i < m; i++ {
		u[i] = make(Vector, nu)
	}
	for i := 0; i < n; i++ {
		v[i] = make(Vector, n)
	}

	e := make(Vector, n)
	work := make(Vector, m)
	wantu := true
	wantv := true

	// Reduce A to bi-diagonal form, storing the diagonal elements
	// in s and the super-diagonal elements in e.
	nct := min(m-1, n)
	nrt := max(0, min(n-2, m))
	for k := 0; k < max(nct, nrt); k++ {
		if k < nct {

			// Compute the transformation for the k-th column and
			// place the k-th diagonal in s[k].
			// Compute 2-norm of k-th column without under/overflow.
			s[k] = 0
			for i := k; i < m; i++ {
				s[k] = math.Hypot(s[k], A[i][k])
			}
			if s[k] != 0.0 {
				if A[k][k] < 0.0 {
					s[k] = -s[k]
				}
				for i := k; i < m; i++ {
					A[i][k] /= s[k]
				}
				A[k][k] += 1.0
			}
			s[k] = -s[k]
		}
		for j := k + 1; j < n; j++ {
			if (k < nct) && (s[k] != 0.0) {

				// Apply the transformation.

				t := float64(0)
				for i := k; i < m; i++ {
					t += A[i][k] * A[i][j]
				}
				t = -t / A[k][k]
				for i := k; i < m; i++ {
					A[i][j] += t * A[i][k]
				}
			}

			// Place the k-th row of A into e for the
			// subsequent calculation of the row transformation.

			e[j] = A[k][j]
		}
		if wantu && (k < nct) {

			// Place the transformation in u for subsequent back
			// multiplication.

			for i := k; i < m; i++ {
				u[i][k] = A[i][k]
			}
		}
		if k < nrt {

			// Compute the k-th row transformation and place the
			// k-th super-diagonal in e[k].
			// Compute 2-norm without under/overflow.
			e[k] = 0
			for i := k + 1; i < n; i++ {
				e[k] = math.Hypot(e[k], e[i])
			}
			if e[k] != 0.0 {
				if e[k+1] < 0.0 {
					e[k] = -e[k]
				}
				for i := k + 1; i < n; i++ {
					e[i] /= e[k]
				}
				e[k+1] += 1.0
			}
			e[k] = -e[k]
			if (k+1 < m) && (e[k] != 0.0) {

				// Apply the transformation.

				for i := k + 1; i < m; i++ {
					work[i] = 0.0
				}
				for j := k + 1; j < n; j++ {
					for i := k + 1; i < m; i++ {
						work[i] += e[j] * A[i][j]
					}
				}
				for j := k + 1; j < n; j++ {
					t := -e[j] / e[k+1]
					for i := k + 1; i < m; i++ {
						A[i][j] += t * work[i]
					}
				}
			}
			if wantv {

				// Place the transformation in v for subsequent
				// back multiplication.

				for i := k + 1; i < n; i++ {
					v[i][k] = e[i]
				}
			}
		}
	}

	// Set up the final bidiagonal matrix or order p.

	p := min(n, m+1)
	if nct < n {
		s[nct] = A[nct][nct]
	}
	if m < p {
		s[p-1] = 0.0
	}
	if nrt+1 < p {
		e[nrt] = A[nrt][p-1]
	}
	e[p-1] = 0.0

	// If required, generate u.

	if wantu {
		for j := nct; j < nu; j++ {
			for i := 0; i < m; i++ {
				u[i][j] = 0.0
			}
			u[j][j] = 1.0
		}
		for k := nct - 1; k >= 0; k-- {
			if s[k] != 0.0 {
				for j := k + 1; j < nu; j++ {
					t := float64(0)
					for i := k; i < m; i++ {
						t += u[i][k] * u[i][j]
					}
					t = -t / u[k][k]
					for i := k; i < m; i++ {
						u[i][j] += t * u[i][k]
					}
				}
				for i := k; i < m; i++ {
					u[i][k] = -u[i][k]
				}
				u[k][k] = 1.0 + u[k][k]
				for i := 0; i < k-1; i++ {
					u[i][k] = 0.0
				}
			} else {
				for i := 0; i < m; i++ {
					u[i][k] = 0.0
				}
				u[k][k] = 1.0
			}
		}
	}

	// If required, generate v.

	if wantv {
		for k := n - 1; k >= 0; k-- {
			if (k < nrt) && (e[k] != 0.0) {
				for j := k + 1; j < nu; j++ {
					t := float64(0)
					for i := k + 1; i < n; i++ {
						t += v[i][k] * v[i][j]
					}
					t = -t / v[k+1][k]
					for i := k + 1; i < n; i++ {
						v[i][j] += t * v[i][k]
					}
				}
			}
			for i := 0; i < n; i++ {
				v[i][k] = 0.0
			}
			v[k][k] = 1.0
		}
	}

	// Main iteration loop for the singular values.

	pp := p - 1
	iter := 0
	eps := math.Pow(2.0, -52.0)
	tiny := math.Pow(2.0, -966.0)
	for p > 0 {
		var k, kase int

		// Here is where a test for too many iterations would go.

		// This section of the program inspects for
		// negligible elements in the s and e arrays.  On
		// completion the variables kase and k are set as follows.

		// kase = 1     if s(p) and e[k-1] are negligible and k<p
		// kase = 2     if s(k) is negligible and k<p
		// kase = 3     if e[k-1] is negligible, k<p, and
		//              s(k), ..., s(p) are not negligible (qr step).
		// kase = 4     if e(p-1) is negligible (convergence).

		for k = p - 2; k >= -1; k-- {
			if k == -1 {
				break
			}
			if math.Abs(e[k]) <=
				tiny+eps*(math.Abs(s[k])+math.Abs(s[k+1])) {
				e[k] = 0.0
				break
			}
		}
		if k == p-2 {
			kase = 4
		} else {
			var ks int
			for ks = p - 1; ks >= k; ks-- {
				if ks == k {
					break
				}
				t := float64(0)
				if ks != p {
					t = math.Abs(e[ks])
				}
				if ks != k+1 {
					t += math.Abs(e[ks-1])
				}
				//double t = (ks != p ? Math.abs(e[ks]) : 0.) +
				//           (ks != k+1 ? Math.abs(e[ks-1]) : 0.);
				if math.Abs(s[ks]) <= tiny+eps*t {
					s[ks] = 0.0
					break
				}
			}
			if ks == k {
				kase = 3
			} else if ks == p-1 {
				kase = 1
			} else {
				kase = 2
				k = ks
			}
		}
		k++

		// Perform the task indicated by kase.

		switch kase {

		// Deflate negligible s(p).

		case 1:
			{
				f := e[p-2]
				e[p-2] = 0.0
				for j := p - 2; j >= k; j-- {
					t := math.Hypot(s[j], f)
					cs := s[j] / t
					sn := f / t
					s[j] = t
					if j != k {
						f = -sn * e[j-1]
						e[j-1] = cs * e[j-1]
					}
					if wantv {
						for i := 0; i < n; i++ {
							t = cs*v[i][j] + sn*v[i][p-1]
							v[i][p-1] = -sn*v[i][j] + cs*v[i][p-1]
							v[i][j] = t
						}
					}
				}
			}
			break

		// Split at negligible s(k).

		case 2:
			{
				f := e[k-1]
				e[k-1] = 0.0
				for j := k; j < p; j++ {
					t := math.Hypot(s[j], f)
					cs := s[j] / t
					sn := f / t
					s[j] = t
					f = -sn * e[j]
					e[j] = cs * e[j]
					if wantu {
						for i := 0; i < m; i++ {
							t = cs*u[i][j] + sn*u[i][k-1]
							u[i][k-1] = -sn*u[i][j] + cs*u[i][k-1]
							u[i][j] = t
						}
					}
				}
			}
			break

		// Perform one qr step.

		case 3:
			{

				// Calculate the shift.

				scale := math.Max(math.Max(math.Max(math.Max(
					math.Abs(s[p-1]), math.Abs(s[p-2])),
					math.Abs(e[p-2])),
					math.Abs(s[k])),
					math.Abs(e[k]))
				sp := s[p-1] / scale
				spm1 := s[p-2] / scale
				epm1 := e[p-2] / scale
				sk := s[k] / scale
				ek := e[k] / scale
				b := ((spm1+sp)*(spm1-sp) + epm1*epm1) / 2.0
				c := (sp * epm1) * (sp * epm1)
				shift := float64(0)
				if (b != 0.0) || (c != 0.0) {
					shift = math.Sqrt(b*b + c)
					if b < 0.0 {
						shift = -shift
					}
					shift = c / (b + shift)
				}
				f := (sk+sp)*(sk-sp) + shift
				g := sk * ek

				// Chase zeros.

				for j := k; j < p-1; j++ {
					t := math.Hypot(f, g)
					cs := f / t
					sn := g / t
					if j != k {
						e[j-1] = t
					}
					f = cs*s[j] + sn*e[j]
					e[j] = cs*e[j] - sn*s[j]
					g = sn * s[j+1]
					s[j+1] = cs * s[j+1]
					if wantv {
						for i := 0; i < n; i++ {
							t = cs*v[i][j] + sn*v[i][j+1]
							v[i][j+1] = -sn*v[i][j] + cs*v[i][j+1]
							v[i][j] = t
						}
					}
					t = math.Hypot(f, g)
					cs = f / t
					sn = g / t
					s[j] = t
					f = cs*e[j] + sn*s[j+1]
					s[j+1] = -sn*e[j] + cs*s[j+1]
					g = sn * e[j+1]
					e[j+1] = cs * e[j+1]
					if wantu && (j < m-1) {
						for i := 0; i < m; i++ {
							t = cs*u[i][j] + sn*u[i][j+1]
							u[i][j+1] = -sn*u[i][j] + cs*u[i][j+1]
							u[i][j] = t
						}
					}
				}
				e[p-2] = f
				iter = iter + 1
			}
			break

		// Convergence.

		case 4:
			{

				// Make the singular values positive.

				if s[k] <= 0.0 {
					if s[k] < 0.0 {
						s[k] = -s[k]
					} else {
						s[k] = 0
					}
					if wantv {
						for i := 0; i <= pp; i++ {
							v[i][k] = -v[i][k]
						}
					}
				}

				// Order the singular values.

				for k < pp {
					if s[k] >= s[k+1] {
						break
					}
					t := s[k]
					s[k] = s[k+1]
					s[k+1] = t
					if wantv && (k < n-1) {
						for i := 0; i < n; i++ {
							t = v[i][k+1]
							v[i][k+1] = v[i][k]
							v[i][k] = t
						}
					}
					if wantu && (k < m-1) {
						for i := 0; i < m; i++ {
							t = u[i][k+1]
							u[i][k+1] = u[i][k]
							u[i][k] = t
						}
					}
					k++
				}
				iter = 0
				p--
			}
			break
		}
	}

	U = new(Matrix).Init(u).GetSubMatrix(0, 0, m, min(m+1, n))
	ls := len(s)
	S = ZeroMatrix(ls, ls)
	for i := 0; i < ls; i++ {
		S._array[i][i] = s[i]
	}
	V = new(Matrix).Init(v)

	return
}
