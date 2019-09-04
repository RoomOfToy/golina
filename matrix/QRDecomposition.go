package matrix

func Sign(a float64) float64 {
	if a > 0 {
		return 1
	} else if a < 0 {
		return -1
	}
	return 0
}

func UnitVector(n, i int) *Matrix {
	v := ZeroMatrix(n, 1)
	v.Set(0, i, 1)
	return v
}

// Householder
// https://en.wikipedia.org/wiki/QR_decomposition
func Householder(t *Matrix) *Matrix {
	m, _ := t.Dims()
	s := Sign(t.At(0, 0))
	e := UnitVector(m, 0)
	u := t.Add(e.MulNum(t.Norm() * s))
	v := u.MulNum(1. / u.At(0, 0))
	prod := v.T().Mul(v)
	beta := 2. / prod.At(0, 0)

	prod = v.Mul(v.T())
	return IdentityMatrix(m).Sub(prod.MulNum(beta))
}

// QR-Decomposition based on `Householder`
func QRDecomposition(t *Matrix) (*Matrix, *Matrix) {
	// TODO: need optimize, many `Mul` -> optimize `Mul` to run in parallel
	m, n := t.Dims()
	q := IdentityMatrix(m)
	r := Copy(t)

	last := n - 1
	if m == n {
		last--
	}
	for i := 0; i <= last; i++ {
		b := r.GetSubMatrix(i, i, m-i, n-i)
		x := b.Col(0).ToMatrix(n-i, 1)
		h := IdentityMatrix(m)
		h.SetSubMatrix(i, i, Householder(x))
		q = q.Mul(h)
		r = h.Mul(r)
	}
	return q, r
}
