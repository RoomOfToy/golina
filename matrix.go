package golina

import (
	"math"
	"runtime"
	"sync"
)

// init function to set CPU usage
func init() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus) // Try to use all available CPUs.
}

// EPS for float number comparision
const EPS float64 = 1E-9

type Vector []float64 // 1D array
type Data []Vector    // 2D array -> backend of Matrix

// matrix entry
type Entry struct {
	value    float64
	row, col int
}

// _Matrix interface

type _Matrix interface {
	// dimensions
	Dims() (row, col int)

	// value at index(row i, col j), panic if not access
	At(i, j int) float64

	// set value at index(row i, col j), panic if not access
	Set(i, j int, value float64)

	// transpose matrix
	T() _Matrix

	// get row
	Row(i int) *Vector

	// get column
	Col(i int) *Vector

	// max entry
	Max() *Entry

	// min entry
	Min() *Entry

	// rank
	Rank() int
}

// Transpose struct implementing _Matrix interface and return transpose of input _Matrix
type Matrix struct {
	_Matrix
	_array Data // row-wise
}

func (t *Matrix) Init(array Data) *Matrix {
	return &Matrix{_array: array}
}

func (t *Matrix) Dims() (row, col int) {
	return len(t._array), len(t._array[0])
}

func (t *Matrix) At(i, j int) float64 {
	return t._array[i][j]
}

func (t *Matrix) Set(i, j int, value float64) {
	t._array[i][j] = value
}

func (t *Matrix) T() *Matrix {
	row, col := t.Dims()
	ntArray := make(Data, col)
	for i := 0; i < col; i++ {
		ntArray[i] = make([]float64, row)
		for j := 0; j < row; j++ {
			ntArray[i][j] = t._array[j][i]
		}
	}
	nt := new(Matrix).Init(ntArray)
	return nt
}

func (t *Matrix) Row(m int) *Vector {
	row, _ := t.Dims()
	if m > -1 && m < row {
		return &t._array[m]
	}
	panic("row index out of range")
}

func (t *Matrix) Col(n int) *Vector {
	_, col := t.Dims()
	if n > -1 && n < col {
		return &t.T()._array[n]
	}
	panic("column index out of range")
}

func VEqual(v1, v2 *Vector) bool {
	if len(*v1) != len(*v2) {
		return false
	}
	for i, v := range *v1 {
		if v != (*v2)[i] {
			return false
		}
	}
	return true
}

// https://stackoverflow.com/questions/37884152/how-do-i-check-the-equality-of-three-values-elegantly
func Equal(mat1, mat2 *Matrix) bool {
	row1, col1 := mat1.Dims()
	row2, col2 := mat2.Dims()
	if [2]int{row1, col1} != [2]int{row2, col2} {
		return false
	}
	for i, col := range mat1._array {
		if !VEqual(&col, mat2.Row(i)) {
			return false
		}
	}
	return true
}

func Copy(t *Matrix) *Matrix {
	nt := Matrix{_array: make([]Vector, len(t._array))}
	for i := range t._array {
		nt._array[i] = make(Vector, len(t._array[i]))
		copy(nt._array[i], t._array[i])
	}
	return &nt
}

func Empty(t *Matrix) *Matrix {
	row, col := t.Dims()
	nt := Matrix{_array: make([]Vector, row)}
	for i := range t._array {
		nt._array[i] = make(Vector, col)
	}
	return &nt
}

// nil entries
func EmptyMatrix(row, col int) *Matrix {
	nt := Matrix{_array: make([]Vector, row)}
	for i := range nt._array {
		nt._array[i] = make(Vector, col)
	}
	return &nt
}

func ZeroMatrix(row, col int) *Matrix {
	nt := EmptyMatrix(row, col)
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			nt.Set(i, j, 0)
		}
	}
	return nt
}

func IdentityMatrix(n int) *Matrix {
	nt := ZeroMatrix(n, n)
	for i := 0; i < n; i++ {
		nt.Set(i, i, 1)
	}
	return nt
}

// TODO: how to find all max or min?
func (t *Matrix) Max() *Entry {
	entry := Entry{}
	entry.value = math.Inf(-1)
	for r, i := range t._array {
		for c, j := range i {
			if j > entry.value {
				entry.value = j
				entry.row, entry.col = r, c
			}
		}
	}
	return &entry
}

func (t *Matrix) Min() *Entry {
	entry := Entry{}
	entry.value = math.Inf(1)
	for r, i := range t._array {
		for c, j := range i {
			if j < entry.value {
				entry.value = j
				entry.row, entry.col = r, c
			}
		}
	}
	return &entry
}

// TODO: need optimize
// Gaussian elimination (row echelon form)
func (t *Matrix) Rank() (rank int) {
	mat := Copy(t)
	rowN, colN := mat.Dims()
	rank = colN
	for row := 0; row < rank; row++ {
		// diagonal entry is not zero
		if mat.At(row, row) != 0 {
			for col := 0; col < rowN; col++ {
				if col != row {
					// makes all entries of current column as 0 except entry `mat[row][row]`
					multipler := mat.At(col, row) / mat.At(row, row)
					for i := 0; i < rank; i++ {
						mat.Set(col, i, mat.At(col, i)-multipler*mat.At(row, i))
					}
				}
			}
		} else {
			// diagonal entry is already zero, now two cases
			// 1) if there is a row below it with non-zero entry, then swap this row with that row and process that row
			// 2) if all elements in current column below mat[row][row] are 0,
			// 	  then remove this column by swapping it with last column and reducing rank by 1
			reduce := true

			for i := row + 1; i < rowN; i++ {
				// swap the row with non-zero entry with this row
				if mat.At(i, row) > EPS {
					SwapRow(mat, row, i)
					reduce = false
					break
				}
			}

			// if no row with non-zero entry in current column, then all values in this column are 0
			if reduce {
				// reduce rank
				rank--
				// copy the last column here
				for i := 0; i < rowN; i++ {
					mat.Set(i, row, mat.At(i, rank))
				}
			}

			// process this row again
			row--
		}
	}
	return rank
}

func SwapRow(t *Matrix, row1, row2 int) {
	t._array[row1], t._array[row2] = *t.Row(row2), *t.Row(row1)
}

// TODO: need optimize
// Determinant of N x N matrix recursively
func (t *Matrix) Det() float64 {
	row, col := t.Dims()
	if row != col {
		panic("need N x N matrix for determinant calculation")
	}
	return _det(t, row)
}
func _det(t *Matrix, n int) float64 {
	det := 0.

	// base case: if matrix only contains one entry
	if n == 1 {
		return t.At(0, 0)
	}

	// template matrix to store coefficients
	matTmp := Empty(t)
	// sign of multiplier
	sign := 1.
	// iterate for each entry of first row
	for f := 0; f < n; f++ {
		// get coefficient of mat[0][f]
		getCoeff(t, matTmp, 0, f, n)
		det += sign * t.At(0, f) * _det(matTmp, n-1)
		sign = -sign
	}
	return det
}

// func to get coefficients of mat[p][q] in matTmp, n is dimension of current matrix (to avoid re-calculation)
func getCoeff(t, matTmp *Matrix, p, q, n int) {
	i, j := 0, 0

	// looping for each entries of the matrix
	for row := 0; row < n; row++ {
		for col := 0; col < n; col++ {
			// fill template matrix
			if row != p && col != q {
				matTmp.Set(i, j, t.At(row, col))
				j++
				// row is filled, so increase row index and rest col index
				if j == n-1 {
					j = 0
					i++
				}
			}
		}
	}
}

// Adjugate Matrix
// https://en.wikipedia.org/wiki/Adjugate_matrix
func (t *Matrix) Adj() (adj *Matrix) {
	row, col := t.Dims()
	if row != col {
		panic("need N x N matrix for adjugate calculation")
	}

	adj = Empty(t)

	if row == 1 {
		adj.Set(0, 0, 1)
		return
	}

	// temp to store coefficients
	matTmp := Empty(t)
	sign := 1.

	for i := 0; i < row; i++ {
		for j := 0; j < row; j++ {
			// get coefficient of t[i][j]
			getCoeff(t, matTmp, i, j, row)
			// sign of adj[j][i] is positive if sum of row and column indexes is even
			sign = Ternary((i+j)%2 == 0, 1., -1.).(float64)
			// interchanging rows and columns to get transpose
			adj.Set(j, i, sign*_det(matTmp, row-1))
		}
	}
	return
}

// Inverse Matrix
// inverse(t) = adj(t) / det(t)
func (t *Matrix) Inverse() *Matrix {
	det := t.Det()
	if det == 0 {
		panic("this matrix is not invertible")
	}
	adj := t.Adj()
	inverse := Empty(t)
	n, _ := t.Dims()
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			inverse.Set(i, j, adj.At(i, j)/det)
		}
	}
	return inverse
}

func (t *Matrix) Add(mat2 *Matrix) *Matrix {
	row1, col1 := t.Dims()
	row2, col2 := mat2.Dims()
	if [2]int{row1, col1} != [2]int{row2, col2} {
		panic("both matrices should have the same dimension")
	}
	nt := Empty(t)
	for r, i := range t._array {
		for c, j := range i {
			nt.Set(r, c, j+mat2.At(r, c))
		}
	}
	return nt
}

func (t *Matrix) Sub(mat2 *Matrix) *Matrix {
	row1, col1 := t.Dims()
	row2, col2 := mat2.Dims()
	if [2]int{row1, col1} != [2]int{row2, col2} {
		panic("both matrices should have the same dimension")
	}
	nt := Empty(t)
	for r, i := range t._array {
		for c, j := range i {
			nt.Set(r, c, j-mat2.At(r, c))
		}
	}
	return nt
}

// TODO: need optimize
// Matrix multiplication (dot | inner)
// https://en.wikipedia.org/wiki/Matrix_multiplication
func (t *Matrix) Mul(mat2 *Matrix) *Matrix {
	row1, col1 := t.Dims()
	row2, col2 := mat2.Dims()
	if col1 != row2 {
		panic("matrix multiplication need M x N and N x L matrices to get M x L matrix")
	}
	out := ZeroMatrix(row1, col2)
	for i := 0; i < row1; i++ {
		for j := 0; j < col2; j++ {
			for k := 0; k < row2; k++ {
				out.Set(i, j, out.At(i, j)+t.At(i, k)*mat2.At(k, j))
			}
		}
	}
	return out
}

// TODO: need optimize and deal with negative condition (invertible)
// Matrix Power of square matrix
// Precondition: n >= 0
func (t *Matrix) Pow(n int) *Matrix {
	row, col := t.Dims()
	if row != col {
		panic("only square matrix has power")
	}
	if n == 0 {
		return IdentityMatrix(row)
	} else if n == 1 {
		return t
	} else {
		nt := Copy(t)
		for n > 1 {
			nt = nt.Mul(Copy(t))
			n--
		}
		return nt
	}
}

// Vector convolve

func Ternary(statement bool, a, b interface{}) interface{} {
	if statement {
		return a
	}
	return b
}

func min(x, y int) int {
	return Ternary(x > y, y, x).(int)
}

func max(x, y int) int {
	return Ternary(x > y, x, y).(int)
}

func mul(u, v *Vector, k int) (res float64) {
	n := min(k+1, len(*u))
	j := min(k, len(*v)-1)

	for i := k - j; i < n; i, j = i+1, j-1 {
		res += (*u)[i] * (*v)[j]
	}
	return res
}

// Convolve computes w = u * v, where w[k] = Σ u[i]*v[j], i + j = k.
// Precondition: len(u) > 0, len(v) > 0.
func Convolve(u, v *Vector) *Vector {
	n := len(*u) + len(*v) - 1
	w := make(Vector, n)

	// Divide w into work units that take ~100μs-1ms to compute.
	size := max(1, 100000/n)

	var wg sync.WaitGroup
	for i, j := 0, size; i < n; i, j = j, j+size {
		if j > n {
			j = n
		}

		// The goroutines share memory, but only for reading.
		wg.Add(1)

		go func(i, j int) {
			for k := i; k < j; k++ {
				w[k] = mul(u, v, k)
			}
			wg.Done()
		}(i, j)
	}

	wg.Wait()

	return &w
}
