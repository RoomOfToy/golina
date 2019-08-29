// Package `golina` provides primitives for linear algebra calculations on top of pure golang.

package golina

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"log"
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
const EPS float64 = 1E-6

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

// Matrix struct
type Matrix struct {
	_Matrix      // basic interface
	_array  Data // row-wise
}

// generate matrix struct from 2D array
func (t *Matrix) Init(array Data) *Matrix {
	return &Matrix{_array: array}
}

// matrix dimensions in row, col
func (t *Matrix) Dims() (row, col int) {
	return len(t._array), len(t._array[0])
}

// get element at row i, column j of matrix
func (t *Matrix) At(i, j int) float64 {
	return t._array[i][j]
}

// set element at row i, column j of matrix
func (t *Matrix) Set(i, j int, value float64) {
	t._array[i][j] = value
}

// transpose matrix
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

// row vector of matrix row m
func (t *Matrix) Row(m int) *Vector {
	row, _ := t.Dims()
	if m > -1 && m < row {
		return &t._array[m]
	}
	panic("row index out of range")
}

// column vector of matrix column n
func (t *Matrix) Col(n int) *Vector {
	_, col := t.Dims()
	if n > -1 && n < col {
		return &t.T()._array[n]
	}
	panic("column index out of range")
}

// check whether two float numbers are equal, defined by threshold EPS
// https://floating-point-gui.de/errors/comparison/
func FloatEqual(x, y float64) bool {
	diff := math.Abs(x - y)
	mean := math.Abs(x+y) / 2.
	absX := math.Abs(x)
	absY := math.Abs(y)
	if x == y {
		return true
	} else if x == 0 || y == 0 || absX+absY < EPS {
		return diff < EPS
	} else {
		return diff/mean < EPS
	}
}

// check whether two vector are equal, based on `FloatEqual`
func VEqual(v1, v2 *Vector) bool {
	if len(*v1) != len(*v2) {
		return false
	}
	for i, v := range *v1 {
		if v != (*v2)[i] && !FloatEqual(v, (*v2)[i]) {
			return false
		}
	}
	return true
}

// check whether two matrix are equal, based on `VEqual`
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

// make a copy of matrix
func Copy(t *Matrix) *Matrix {
	nt := Matrix{_array: make([]Vector, len(t._array))}
	for i := range t._array {
		nt._array[i] = make(Vector, len(t._array[i]))
		copy(nt._array[i], t._array[i])
	}
	return &nt
}

// empty matrix
// 	golang make slice has zero value in default, so empty matrix == zero matrix
func Empty(t *Matrix) *Matrix {
	row, col := t.Dims()
	nt := Matrix{_array: make([]Vector, row)}
	for i := range t._array {
		nt._array[i] = make(Vector, col)
	}
	return &nt // nt is a zero matrix
}

// generate matrix with all elements are zero
func ZeroMatrix(row, col int) *Matrix {
	nt := Matrix{_array: make([]Vector, row)}
	for i := range nt._array {
		nt._array[i] = make(Vector, col)
	}
	return &nt
}

// generate matrix with all elements are one
func OneMatrix(row, col int) *Matrix {
	nt := Matrix{_array: make([]Vector, row)}
	r := make(Vector, col)
	for i := range r {
		r[i] = 1
	}
	for i := range nt._array {
		nt._array[i] = make(Vector, col)
		copy(nt._array[i], r)
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

// find the first min entry
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
	t._array[row1], t._array[row2] = *t.Row(row2), *t.Row(row1)
}

// Determinant of N x N matrix based on LU Decomposition
func (t *Matrix) Det() float64 {
	row, col := t.Dims()
	if row != col {
		panic("need N x N matrix for determinant calculation")
	}
	nt, P := LUPDecompose(t, 3, EPS)
	return LUPDeterminant(nt, P, 3)
}

func (t *Matrix) Inverse() *Matrix {
	nt, P := LUPDecompose(t, 3, EPS)
	return LUPInvert(nt, P, 3)
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
	for r, i := range t._array {
		for c, j := range i {
			nt.Set(r, c, j+mat2.At(r, c))
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
	for r, i := range t._array {
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
	multiplier := getFloat64(n)
	row, col := t.Dims()
	out := ZeroMatrix(row, col)
	for i := range t._array {
		for j, v := range t._array[i] {
			out.Set(i, j, v*multiplier)
		}
	}
	return out
}

func getFloat64(x interface{}) float64 {
	switch x := x.(type) {
	case uint8:
		return float64(x)
	case int8:
		return float64(x)
	case uint16:
		return float64(x)
	case int16:
		return float64(x)
	case uint32:
		return float64(x)
	case int32:
		return float64(x)
	case uint64:
		return float64(x)
	case int64:
		return float64(x)
	case int:
		return float64(x)
	case float32:
		return float64(x)
	case float64:
		return x
	}
	panic("invalid numeric type of input")
}

// Matrix Power of square matrix
// 	Precondition: n >= 0
func (t *Matrix) Pow(n int) *Matrix {
	// TODO: need optimize and deal with negative condition (invertible)
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
	for i := range t._array {
		res += t._array[i][i]
	}
	return res
}

// Frobenius norm
func (t *Matrix) Norm() float64 {
	fr := 0.
	for _, i := range t._array {
		fr += i.SquareSum()
	}
	return math.Sqrt(fr)
}

// matrix to row vector
func (t *Matrix) Flat() *Vector {
	m, n := t.Dims()
	v := make(Vector, m*n)
	for i, j := range t._array {
		copy(v[i*n:i*n+m], j)
	}
	return &v
}

// Get a sub-matrix starting at i, j with rows rows and cols columns.
func (t *Matrix) GetSubMatrix(i, j, rows, cols int) *Matrix {
	nt := ZeroMatrix(rows, cols)
	for k := range nt._array {
		copy(nt._array[k], t._array[i+k][j:j+cols])
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
		return t.Sum(0).MulNum(1. / float64(col)) // Notice: without float64(col), this will be int, e.g. col=3, 1./col=0
	case 1:
		return t.Sum(1).MulNum(1. / float64(row))
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
	for i := range t._array {
		for j := range t._array[i] {
			if t._array[i][j] != t._array[j][i] {
				return false
			}
		}
	}
	return true
}

// Vector
// 	get vector element at index n
func (v *Vector) At(n int) float64 {
	l := len(*v)
	if abs(n) > l {
		panic("index out of range")
	}
	if n < 0 {
		n = l + n
	}
	return (*v)[n]
}

// add two vectors
func (v *Vector) Add(v1 *Vector) *Vector {
	if len(*v) != len(*v1) {
		panic("add requires equal-length vectors")
	}
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] + (*v1)[i]
	}
	return &res
}

// vector add number
func (v *Vector) AddNum(n interface{}) *Vector {
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] + getFloat64(n)
	}
	return &res
}

// vector subtract vector
func (v *Vector) Sub(v1 *Vector) *Vector {
	if len(*v) != len(*v1) {
		panic("sub requires equal-length vectors")
	}
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] - (*v1)[i]
	}
	return &res
}

// vector subtract number
func (v *Vector) SubNum(n interface{}) *Vector {
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] - getFloat64(n)
	}
	return &res
}

// vector multiply number
func (v *Vector) MulNum(n interface{}) *Vector {
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] * getFloat64(n)
	}
	return &res
}

// vector dot production
func (v *Vector) Dot(v1 *Vector) float64 {
	if len(*v) != len(*v1) {
		panic("dot product requires equal-length vectors")
	}
	res := 0.
	for i := range *v {
		res += (*v)[i] * (*v1)[i]
	}
	return res
}

// vector outer product: v1, v2 -> matrix
// https://en.wikipedia.org/wiki/Outer_product
func (v *Vector) OuterProduct(v1 *Vector) *Matrix {
	row, col := len(*v), len(*v1)
	res := ZeroMatrix(row, col)
	for i := range res._array {
		for j := range res._array[i] {
			res.Set(i, j, (*v)[i]*(*v)[j])
		}
	}
	return res
}

// vector cross product, 3D only
func (v *Vector) Cross(v1 *Vector) *Vector {
	if len(*v) != len(*v1) || len(*v) != 3 {
		panic("cross product requires 3d vectors in 3d space!")
	}
	return &Vector{(*v)[1]*(*v1)[2] - (*v)[2]*(*v1)[1], (*v)[2]*(*v1)[0] - (*v)[0]*(*v1)[2], (*v)[0]*(*v1)[1] - (*v)[1]*(*v1)[0]}
}

// vector elements square sum
func (v *Vector) SquareSum() float64 {
	// dot is almost 50% faster than pow by benchmark
	return v.Dot(v)
}

// vector norm
func (v *Vector) Norm() float64 {
	return math.Sqrt(v.SquareSum())
}

// normalize vector
func (v *Vector) Normalize() *Vector {
	n := v.Norm()
	if n == 0 {
		panic("invalid input vector with norm equal to 0")
	}
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] / n
	}
	return &res
}

// vector to matrix, row-wise
func (v *Vector) ToMatrix(rows, cols int) *Matrix {
	if len(*v) != rows*cols {
		panic(fmt.Sprintf("invalid target matrix dimensions (%d x %d) with vector length %d\n", rows, cols, len(*v)))
	}
	nt := ZeroMatrix(rows, cols)
	for r := range nt._array {
		for c := range nt._array[r] {
			nt._array[r][c] = (*v)[r*cols+c]
		}
	}
	return nt
}

// sum of vector's elements
func (v *Vector) Sum() float64 {
	s := 0.
	for _, e := range *v {
		s += e
	}
	return s
}

// sum of vector elements' absolute value
func (v *Vector) AbsSum() float64 {
	s := 0.
	for _, e := range *v {
		s += math.Abs(e)
	}
	return s
}

// mean value of vector
func (v *Vector) Mean() float64 {
	return v.Sum() / float64(len(*v))
}

// tile vector alone certain dimension into matrix, 0 -> vector as row, 1 -> vector as column
func (v *Vector) Tile(dim, n int) *Matrix {
	switch dim {
	case 0:
		d := make(Data, n)
		for i := range d {
			d[i] = *v
		}
		return new(Matrix).Init(d)
	case 1:
		d := make(Data, n)
		for i := range d {
			d[i] = *v
		}
		return new(Matrix).Init(d).T()
	default:
		panic("invalid tile dimension")
	}
}

// length of vector
func (v *Vector) Length() int {
	return len(*v)
}

// simple function for simulating ternary operator
func Ternary(statement bool, a, b interface{}) interface{} {
	if statement {
		return a
	}
	return b
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
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

// Vector convolve
//	Convolve computes w = u * v, where w[k] = Σ u[i]*v[j], i + j = k.
//	Precondition: len(u) > 0, len(v) > 0.
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

			maxLen = max(maxLen, len(entryString))
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

// pretty-print for vector
func (v *Vector) String() string {
	if v == nil {
		return "{nil}"
	}
	outString := "{"
	maxLen := 0
	vLen := len(*v)
	for i := 0; i < vLen; i++ {
		entry := (*v)[i]
		entryString := fmt.Sprintf("%f", entry)
		maxLen = max(maxLen, len(entryString))
	}
	for i := 0; i < vLen; i++ {
		entry := (*v)[i]
		entryString := fmt.Sprintf("%f", entry)
		for len(entryString) < maxLen {
			entryString = " " + entryString
		}
		outString += entryString
		if i != vLen-1 {
			outString += ", "
		}
	}
	outString += "}\n"
	return outString
}
