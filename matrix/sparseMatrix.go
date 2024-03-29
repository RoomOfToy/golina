package matrix

import (
	"math"
)

// SparseMatrix struct
//	https://en.wikipedia.org/wiki/Sparse_matrix
//	Dictionary of keys (DOK)
//	Row-wise
type SparseMatrix struct {
	Rows, Cols int
	Data       map[int]float64 // key = row * Cols + col - Offset + Offset % Cols + Offset / Cols
	Offset     int             // start position
}

// ZeroSparseMatrix returns a sparse matrix with all elements equal to zero
func ZeroSparseMatrix(rows, cols int) *SparseMatrix {
	return &SparseMatrix{
		Rows:   rows,
		Cols:   cols,
		Data:   map[int]float64{},
		Offset: 0,
	}
}

// NewSparseMatrix generates a new sparse matrix with input data and dims
func NewSparseMatrix(data map[int]float64, rows, cols int) *SparseMatrix {
	return &SparseMatrix{
		Rows:   rows,
		Cols:   cols,
		Data:   data,
		Offset: 0,
	}
}

// RowColToIndex transfers row, col idx into internal data map idx
func (sm *SparseMatrix) RowColToIndex(row, col int) (idx int) {
	return row*sm.Cols + col - sm.Offset + sm.Offset%sm.Cols + sm.Offset/sm.Cols
}

// IndexToRowCol transfers internal data map idx into row, col idx
func (sm *SparseMatrix) IndexToRowCol(idx int) (row, col int) {
	return (idx + sm.Offset - sm.Offset/sm.Cols - sm.Offset%sm.Cols) / sm.Cols, (idx + sm.Offset - sm.Offset/sm.Cols - sm.Offset%sm.Cols) % sm.Cols
}

// At returns elements value at row, col idx
func (sm *SparseMatrix) At(row, col int) float64 {
	entry, ok := sm.Data[sm.RowColToIndex(row, col)]
	if !ok {
		if row < sm.Rows && col < sm.Cols && row >= 0 && col >= 0 {
			return 0.
		}
		panic("invalid index")
	}
	return entry
}

// AtIndex returns elements value at internal data map idx
func (sm *SparseMatrix) AtIndex(idx int) float64 {
	entry, ok := sm.Data[idx]
	if !ok {
		if idx < sm.Rows*sm.Cols && idx >= 0 {
			return 0.
		}
		panic("invalid index")
	}
	return entry
}

// Set sets value at row, col idx
//	notice: value == 0. indicates deletion
func (sm *SparseMatrix) Set(row, col int, value float64) {
	if value == 0. {
		delete(sm.Data, sm.RowColToIndex(row, col))
	} else {
		sm.Data[sm.RowColToIndex(row, col)] = value
	}
}

// SetIndex sets value at internal data map idx
//	notice: value == 0. indicates deletion
func (sm *SparseMatrix) SetIndex(idx int, value float64) {
	if value == 0. {
		delete(sm.Data, idx)
	} else {
		sm.Data[idx] = value
	}
}

// GetAllIndexes returns all idx
//	notice: since map is unordered, it needs sort.Ints(idxes) if wants to keep ascending order
func (sm *SparseMatrix) GetAllIndexes() (idxes []int) {
	idxes = make([]int, len(sm.Data))
	i := 0
	for idx := range sm.Data {
		idxes[i] = idx
		i++
	}
	return
}

// Row constructs and returns a row vector
func (sm *SparseMatrix) Row(n int) *Vector {
	v := make(Vector, sm.Cols)
	for i := 0; i < sm.Cols; i++ {
		entry, ok := sm.Data[sm.RowColToIndex(n, i)]
		if !ok {
			v[i] = 0.
		} else {
			v[i] = entry
		}
	}
	return &v
}

// Col constructs and returns a column vector
func (sm *SparseMatrix) Col(n int) *Vector {
	v := make(Vector, sm.Cols)
	for i := 0; i < sm.Cols; i++ {
		entry, ok := sm.Data[sm.RowColToIndex(i, n)]
		if !ok {
			v[i] = 0.
		} else {
			v[i] = entry
		}
	}
	return &v
}

// FindFirstNonZeroInSubMatrix returns the first sub-matrix starts with non-zero value, it searches with a input starting idx
func (sm *SparseMatrix) FindFirstNonZeroInSubMatrix(startIdx int) (idx int) {
	row, col := sm.IndexToRowCol(startIdx)
	for i := row; i < sm.Rows; i++ {
		for j := col; j < sm.Cols; j++ {
			_, ok := sm.Data[sm.RowColToIndex(i, j)]
			if ok {
				return sm.RowColToIndex(i, j)
			}
		}
	}
	return -1
}

// GetSubSparseMatrix returns a sub-sparse-matrix with input idx and dims
func (sm *SparseMatrix) GetSubSparseMatrix(i, j, rows, cols int) *SparseMatrix {
	if i < 0 || j < 0 || i+rows > sm.Rows || j+cols > sm.Cols {
		i = MaxInt(0, i)
		j = MaxInt(0, j)
		rows = MinInt(sm.Rows-i, rows)
		rows = MinInt(sm.Cols-j, cols)
	}

	return &SparseMatrix{
		Rows:   rows,
		Cols:   cols,
		Data:   sm.Data,
		Offset: i*sm.Cols + sm.Offset + j + sm.Offset%sm.Cols, //sm.FindFirstNonZeroInSubMatrix(sm.RowColToIndex(i, j))
	}
}

// Copy returns a deep copy of sparse matrix
func (sm *SparseMatrix) Copy() *SparseMatrix {
	nsm := ZeroSparseMatrix(sm.Rows, sm.Cols)
	for idx, value := range sm.Data {
		nsm.Data[idx] = value
	}
	return nsm
}

// ToMatrix transfers sparse matrix into matrix (dense)
func (sm *SparseMatrix) ToMatrix() *Matrix {
	nm := ZeroMatrix(sm.Rows, sm.Cols)
	for idx, value := range sm.Data {
		i, j := sm.IndexToRowCol(idx)
		nm.Set(i, j, value)
	}
	return nm
}

// T returns a sparse transpose matrix
func (sm *SparseMatrix) T() *SparseMatrix {
	nsm := ZeroSparseMatrix(sm.Cols, sm.Rows)
	for idx, value := range sm.Data {
		nc, nr := sm.IndexToRowCol(idx)
		nsm.Set(nr, nc, value)
	}
	return nsm
}

// Add sums two sparse matrices and returns a new sparse matrix
func (sm *SparseMatrix) Add(sm2 *SparseMatrix) *SparseMatrix {
	if sm.Rows != sm2.Rows || sm.Cols != sm2.Cols {
		panic("Dimension mismatch")
	}
	nsm := ZeroSparseMatrix(sm.Rows, sm.Cols)
	for idx, value := range sm.Data {
		if _, ok := sm2.Data[idx]; !ok {
			nsm.SetIndex(idx, value)
		} else {
			nsm.SetIndex(idx, value+sm2.Data[idx])
		}
	}
	return nsm
}

// AddNum adds input number to all elements inside the sparse matrix and returns a new sparse matrix
func (sm *SparseMatrix) AddNum(n float64) *SparseMatrix {
	nsm := ZeroSparseMatrix(sm.Rows, sm.Cols)
	for idx, value := range sm.Data {
		if value+n != 0. {
			nsm.SetIndex(idx, value+n)
		}
	}
	return nsm
}

// Mul does sparse matrix multiplication
//	TODO: Need Optimize
func (sm *SparseMatrix) Mul(sm2 *SparseMatrix) *SparseMatrix {
	if sm.Cols != sm2.Rows {
		panic("Dimension mismatch")
	}
	nsm := ZeroSparseMatrix(sm.Rows, sm2.Cols)
	for idx, value := range sm.Data {
		i, j := sm.IndexToRowCol(idx)
		for k := 0; k < sm2.Cols; k++ {
			nv := nsm.At(i, k) + value*sm2.At(j, k)
			if nv != 0. {
				nsm.Set(i, k, nv)
			}
		}
	}
	return nsm
}

// MulVec multiplies sparse matrix with input vector and returns a new vector
func (sm *SparseMatrix) MulVec(v *Vector) *Vector {
	if sm.Cols != v.Length() {
		panic("Dimension mismatch")
	}
	nVec := make(Vector, sm.Rows)
	for idx, value := range sm.Data {
		r, c := sm.IndexToRowCol(idx)
		nVec[r] += value * v.At(c)
	}
	return &nVec
}

// MulNum multiplies sparse matrix elements with input number (float64) and returns a new sparse matrix
func (sm *SparseMatrix) MulNum(n float64) *SparseMatrix {
	nsm := ZeroSparseMatrix(sm.Rows, sm.Cols)
	if n == 0. {
		return nsm
	}
	for idx, value := range sm.Data {
		nsm.SetIndex(idx, value*n)
	}
	return nsm
}

// Det returns determinant of sparse matrix
//	transfer to dense matrix and use LUP decomposition
//	Notice: easy OOM and slow...
//	TODO: Any better way???
func (sm *SparseMatrix) Det() float64 {
	defer func() {
		if r := recover(); r != nil { // Can Not do LUP Decomposition with tolerance 1e-6
			// return zero value of the specified return type
		}
	}()
	det := sm.ToMatrix().Det()
	if math.IsNaN(det) {
		return 0.
	}
	return det
}
