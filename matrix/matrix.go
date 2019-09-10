package matrix

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"log"
	"math"
	"runtime"
)

// init function to set CPU usage
func init() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus) // Try to use all available CPUs.
}

// EPS for float number comparision
const EPS float64 = 1E-6

type Data []Vector // 2D array -> backend of Matrix

// matrix entry
type Entry struct {
	Value    float64
	Row, Col int
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

// Matrix struct
type Matrix struct {
	_Matrix      // basic interface
	Data    Data // row-wise
}

// generate matrix struct from 2D array
func (t *Matrix) Init(array Data) *Matrix {
	return &Matrix{Data: array}
}

// matrix dimensions in row, col
func (t *Matrix) Dims() (row, col int) {
	return len(t.Data), len(t.Data[0])
}

// get element at row i, column j of matrix
func (t *Matrix) At(i, j int) float64 {
	return t.Data[i][j]
}

// set element at row i, column j of matrix
func (t *Matrix) Set(i, j int, value float64) {
	t.Data[i][j] = value
}

// transpose matrix
func (t *Matrix) T() *Matrix {
	row, col := t.Dims()
	ntArray := make(Data, col)
	for i := 0; i < col; i++ {
		ntArray[i] = make([]float64, row)
		for j := 0; j < row; j++ {
			ntArray[i][j] = t.Data[j][i]
		}
	}
	nt := new(Matrix).Init(ntArray)
	return nt
}

// row vector of matrix row m
func (t *Matrix) Row(m int) *Vector {
	row, _ := t.Dims()
	if m > -1 && m < row {
		return &t.Data[m]
	}
	panic("row index out of range")
}

// column vector of matrix column n
func (t *Matrix) Col(n int) *Vector {
	_, col := t.Dims()
	if n > -1 && n < col {
		return &t.T().Data[n]
	}
	panic("column index out of range")
}

// make a copy of matrix
func Copy(t *Matrix) *Matrix {
	nt := Matrix{Data: make([]Vector, len(t.Data))}
	for i := range t.Data {
		nt.Data[i] = make(Vector, len(t.Data[i]))
		copy(nt.Data[i], t.Data[i])
	}
	return &nt
}

// empty matrix
// 	golang make slice has zero value in default, so empty matrix == zero matrix
func Empty(t *Matrix) *Matrix {
	row, col := t.Dims()
	nt := Matrix{Data: make([]Vector, row)}
	for i := range t.Data {
		nt.Data[i] = make(Vector, col)
	}
	return &nt // nt is a zero matrix
}

// generate matrix with all elements are zero
func ZeroMatrix(row, col int) *Matrix {
	nt := Matrix{Data: make([]Vector, row)}
	for i := range nt.Data {
		nt.Data[i] = make(Vector, col)
	}
	return &nt
}

// generate matrix with all elements are one
func OneMatrix(row, col int) *Matrix {
	nt := Matrix{Data: make([]Vector, row)}
	r := make(Vector, col)
	for i := range r {
		r[i] = 1
	}
	for i := range nt.Data {
		nt.Data[i] = make(Vector, col)
		copy(nt.Data[i], r)
	}
	return &nt
}

// generate diagonal matrix with ones as diagonal elements, like `eye` in other libs
func IdentityMatrix(n int) *Matrix {
	nt := ZeroMatrix(n, n)
	for i := 0; i < n; i++ {
		nt.Set(i, i, 1)
	}
	return nt
}

// find the first max entry
func (t *Matrix) Max() *Entry {
	// TODO: how to find all max or min?
	entry := Entry{}
	entry.Value = math.Inf(-1)
	for r, i := range t.Data {
		for c, j := range i {
			if j > entry.Value {
				entry.Value = j
				entry.Row, entry.Col = r, c
			}
		}
	}
	return &entry
}

// find the first min entry
func (t *Matrix) Min() *Entry {
	entry := Entry{}
	entry.Value = math.Inf(1)
	for r, i := range t.Data {
		for c, j := range i {
			if j < entry.Value {
				entry.Value = j
				entry.Row, entry.Col = r, c
			}
		}
	}
	return &entry
}

// get matrix rank through Gaussian elimination (row echelon form)
func (t *Matrix) Rank() (rank int) {
	mat := Copy(t)
	rowN, colN := mat.Dims()
	if rowN == colN {
		nt, _ := LUPDecompose(t, rowN, EPS)
		rank = LUPRank(nt, rowN)
		return
	}
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

// swap two rows
func SwapRow(t *Matrix, row1, row2 int) {
	t.Data[row1], t.Data[row2] = *t.Row(row2), *t.Row(row1)
}

// Determinant of N x N matrix based on LU Decomposition
func (t *Matrix) Det() float64 {
	row, col := t.Dims()
	if row != col {
		panic("need N x N matrix for determinant calculation")
	}
	nt, P := LUPDecompose(t, row, EPS)
	if nt == nil {
		panic("Can Not do LUP Decomposition with tolerance 1e-6")
	}
	return LUPDeterminant(nt, P, row)
}

func (t *Matrix) Inverse() *Matrix {
	row, col := t.Dims()
	if row != col {
		panic("only inverse only support square matrix, left/right is not supported")
	}
	nt, P := LUPDecompose(t, row, EPS)
	if nt == nil {
		panic("Can Not do LUP Decomposition with tolerance 1e-6")
	}
	return LUPInvert(nt, P, row)
}

// Determinant of N x N matrix recursively
func NaiveDet(t *Matrix) float64 {
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
func NaiveAdj(t *Matrix) (adj *Matrix) {
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
// 	inverse(t) = adj(t) / det(t)
func NaiveInverse(t *Matrix) *Matrix {
	det := NaiveDet(t)
	if det == 0 {
		panic("this matrix is not invertible")
	}
	adj := NaiveAdj(t)
	inverse := Empty(t)
	n, _ := t.Dims()
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			inverse.Set(i, j, adj.At(i, j)/det)
		}
	}
	return inverse
}

// Add two matrices
func (t *Matrix) Add(mat2 *Matrix) *Matrix {
	row1, col1 := t.Dims()
	row2, col2 := mat2.Dims()
	if [2]int{row1, col1} != [2]int{row2, col2} {
		panic("both matrices should have the same dimension")
	}
	nt := Empty(t)
	for r, i := range t.Data {
		for c, j := range i {
			nt.Set(r, c, j+mat2.At(r, c))
		}
	}
	return nt
}

// Add number
func (t *Matrix) AddNum(n interface{}) *Matrix {
	nt := Empty(t)
	for r, i := range t.Data {
		for c, j := range i {
			nt.Set(r, c, j+GetFloat64(n))
		}
	}
	return nt
}

// Subtract matrix
func (t *Matrix) Sub(mat2 *Matrix) *Matrix {
	row1, col1 := t.Dims()
	row2, col2 := mat2.Dims()
	if [2]int{row1, col1} != [2]int{row2, col2} {
		panic("both matrices should have the same dimension")
	}
	nt := Empty(t)
	for r, i := range t.Data {
		for c, j := range i {
			nt.Set(r, c, j-mat2.At(r, c))
		}
	}
	return nt
}

// Matrix multiplication (dot | inner)
// https://en.wikipedia.org/wiki/Matrix_multiplication
func (t *Matrix) Mul(mat2 *Matrix) *Matrix {
	// TODO: need optimize
	row1, col1 := t.Dims()
	row2, col2 := mat2.Dims()
	if col1 != row2 {
		panic("matrix multiplication need M x N and N x L matrices to get M x L matrix")
	}
	out := ZeroMatrix(row1, col2)
	if row1 <= 80 {
		for i := 0; i < row1; i++ {
			for j := 0; j < col2; j++ {
				for k := 0; k < row2; k++ {
					out.Set(i, j, out.At(i, j)+t.At(i, k)*mat2.At(k, j))
				}
			}
		}
	} else {
		sem := semaphore.NewWeighted(int64(runtime.NumCPU()))
		ctx := context.TODO()
		sum := 0.
		for i := 0; i < row1; i++ {
			if err := sem.Acquire(ctx, 1); err != nil {
				log.Printf("Failed to acquire semaphore: %v", err)
				break
			}
			go func(i int) {
				defer sem.Release(1)
				for j := 0; j < col2; j++ {
					sum = 0.
					for k := 0; k < row2; k++ {
						sum += out.At(i, j) + t.At(i, k)*mat2.At(k, j)
					}
					out.Set(i, j, sum)
				}
			}(i)
		}
		if err := sem.Acquire(ctx, int64(runtime.NumCPU())); err != nil {
			log.Printf("Failed to acquire semaphore: %v", err)
		}
	}
	return out
}

// matrix multiply vector, please notice all vectors in this package is row vector
func (t *Matrix) MulVec(v *Vector) *Vector {
	return t.Mul(new(Matrix).Init(Data{*v}).T()).T().Row(0)
}

func (t *Matrix) MulNum(n interface{}) *Matrix {
	multiplier := GetFloat64(n)
	row, col := t.Dims()
	out := ZeroMatrix(row, col)
	for i := range t.Data {
		for j, v := range t.Data[i] {
			out.Set(i, j, v*multiplier)
		}
	}
	return out
}

func (t *Matrix) GetDiagonalElements() *Vector {
	row, _ := t.Dims()
	v := make(Vector, row)
	for i := range t.Data {
		v[i] = t.Data[i][i]
	}
	return &v
}

// Matrix Power of square matrix
// 	Precondition: n >= 0
func (t *Matrix) Pow(n int) *Matrix {
	// TODO: need deal with negative condition (invertible)
	row, col := t.Dims()
	if row != col {
		panic("only square matrix has power")
	}
	if n == 0 {
		return IdentityMatrix(row)
	} else if n == 1 {
		return t
	} else {
		V, D := EigenDecompose(t)
		for i := range D.Data {
			D.Data[i][i] = math.Pow(D.Data[i][i], float64(n))
		}
		return V.Mul(D).Mul(V.Inverse()) // change NaiveInverse to LUPDecomposeInvert
	}
}

// ten times slower than LUDecompose + EigenDecompose ways for 100 x 100 matrix
func NaivePow(t *Matrix, n int) *Matrix {
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

// Trace: sum of all diagonal values
// https://en.wikipedia.org/wiki/Trace_(linear_algebra)
func (t *Matrix) Trace() float64 {
	row, col := t.Dims()
	if row != col {
		panic("square matrix only")
	}
	res := 0.
	for i := range t.Data {
		res += t.Data[i][i]
	}
	return res
}

// Frobenius norm
func (t *Matrix) Norm() float64 {
	fr := 0.
	for _, i := range t.Data {
		fr += i.SquareSum()
	}
	return math.Sqrt(fr)
}

// matrix to row vector
func (t *Matrix) Flat() *Vector {
	m, n := t.Dims()
	v := make(Vector, m*n)
	for i, j := range t.Data {
		copy(v[i*n:i*n+n], j)
	}
	return &v
}

// Get a sub-matrix starting at i, j with rows rows and cols columns.
func (t *Matrix) GetSubMatrix(i, j, rows, cols int) *Matrix {
	nt := ZeroMatrix(rows, cols)
	for k := range nt.Data {
		copy(nt.Data[k], t.Data[i+k][j:j+cols])
	}
	return nt
}

// Set a sub-matrix starting at i, j with input matrix, please take care the dimension matching conditions
// 	notice: in-place change
func (t *Matrix) SetSubMatrix(i, j int, mat *Matrix) {
	m, n := mat.Dims()
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			t.Set(i+r, j+c, mat.At(r, c))
		}
	}
}

// sum one column of matrix
func (t *Matrix) SumCol(col int) float64 {
	s := 0.
	for _, e := range *(t.Col(col)) {
		s += e
	}
	return s
}

// sum one row of matrix
func (t *Matrix) SumRow(row int) float64 {
	s := 0.
	for _, e := range *(t.Row(row)) {
		s += e
	}
	return s
}

// sum the matrix along certain dimension, 0 -> sum all rows into one vector, 1 -> sum all columns into one vector
func (t *Matrix) Sum(dim int) *Vector {
	row, col := t.Dims()
	switch dim {
	case 0:
		v := make(Vector, col)
		for i := range v {
			v[i] = t.SumCol(i)
		}
		return &v
	case 1:
		v := make(Vector, row)
		for i := range v {
			v[i] = t.SumRow(i)
		}
		return &v
	case -1:
		d := Ternary(row > col, 1, 0).(int)
		return &Vector{t.Sum(d).Sum()}
	default:
		panic("invalid sum dimension")
	}
}

// mean vector of matrix along certain dimension, 0 -> row, 1 -> column
func (t *Matrix) Mean(dim int) *Vector {
	row, col := t.Dims()
	switch dim {
	case 0:
		return t.Sum(0).MulNum(1. / float64(row)) // Notice: without float64(col), this will be int, e.g. col=3, 1./col=0
	case 1:
		return t.Sum(1).MulNum(1. / float64(col))
	case -1:
		return t.Sum(-1).MulNum(1. / float64(row*col))
	default:
		panic("invalid mean dimension")
	}
}

// covariance matrix
func (t *Matrix) CovMatrix() *Matrix {
	row, col := t.Dims()
	x := t.Sub(t.Mean(0).Tile(0, row))
	cov := x.T().Mul(x).MulNum(1. / float64(col-1))
	return cov
}

// cross covariance matrix
// 	https://en.wikipedia.org/wiki/Cross-covariance
// 	https://en.wikipedia.org/wiki/Cross-covariance_matrix
func CrossCovMatrix(mat1, mat2 *Matrix) *Matrix {
	r1, c1 := mat1.Dims()
	r2, c2 := mat2.Dims()
	if r1 != r2 || c1 != c2 {
		panic("both matrix should have the same dimensions")
	}
	return mat1.Sub(mat1.Mean(0).Tile(0, r1)).T().Mul(mat2.Sub(mat2.Mean(0).Tile(0, r1))).MulNum(1. / float64(c1-1))
}

// check whether matrix is symmetric
func (t *Matrix) IsSymmetric() bool {
	m, n := t.Dims()
	if m != n {
		return false
	}
	for i := range t.Data {
		for j := range t.Data[i] {
			if t.Data[i][j] != t.Data[j][i] {
				return false
			}
		}
	}
	return true
}

// get unique elements in matrix
//	it need to loop the whole matrix
func (t *Matrix) Unique() *Vector {
	uSet := map[float64]bool{}
	for _, r := range t.Data {
		for _, val := range r {
			uSet[val] = true
		}
	}
	uv := make(Vector, len(uSet))
	i := 0
	for k := range uSet {
		uv[i] = k
		i++
	}
	return &uv
}

func (t *Matrix) UniqueWithCount() map[float64]int {
	uSet := map[float64]int{}
	for _, r := range t.Data {
		for _, val := range r {
			if _, exist := uSet[val]; !exist {
				uSet[val] = 1
			} else {
				uSet[val]++
			}
		}
	}
	return uSet
}

// pretty-print for matrix
func (t *Matrix) String() string {
	if t == nil {
		return "{nil}"
	}
	outString := "{"
	maxLen := 0
	row, col := t.Dims()
	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			entry := t.At(i, j)
			entryString := fmt.Sprintf("%f", entry)

			maxLen = MaxInt(maxLen, len(entryString))
		}
	}

	for i := 0; i < row; i++ {
		for j := 0; j < col; j++ {
			entry := t.At(i, j)

			entryString := fmt.Sprintf("%f", entry)

			for len(entryString) < maxLen {
				entryString = " " + entryString
			}
			outString += entryString
			if i != row-1 || j != col-1 {
				outString += ","
			}
			if j != col-1 {
				outString += " "
			}
		}
		if i != row-1 {
			outString += "\n "
		}
	}
	outString += "}\n"
	return outString
}
