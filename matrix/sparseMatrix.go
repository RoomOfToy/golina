package matrix

import (
	"math"
)

// Sparse Matrix
//	https://en.wikipedia.org/wiki/Sparse_matrix
//	Dictionary of keys (DOK)
//	Row-wise
type SparseMatrix struct {
	Rows, Cols int
	Data       map[int]float64 // key = row * Cols + col - Offset + Offset % Cols + Offset / Cols
	Offset     int             // start position
}

func ZeroSparseMatrix(rows, cols int) *SparseMatrix {
	return &SparseMatrix{
		Rows:   rows,
		Cols:   cols,
		Data:   map[int]float64{},
		Offset: 0,
	}
}

func NewSparseMatrix(data map[int]float64, rows, cols int) *SparseMatrix {
	return &SparseMatrix{
		Rows:   rows,
		Cols:   cols,
		Data:   data,
		Offset: 0,
	}
}

func (sm *SparseMatrix) RowColToIndex(row, col int) (idx int) {
	return row*sm.Cols + col - sm.Offset + sm.Offset%sm.Cols + sm.Offset/sm.Cols
}

func (sm *SparseMatrix) IndexToRowCol(idx int) (row, col int) {
	return (idx + sm.Offset - sm.Offset/sm.Cols - sm.Offset%sm.Cols) / sm.Cols, (idx + sm.Offset - sm.Offset/sm.Cols - sm.Offset%sm.Cols) % sm.Cols
}

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

// value == 0. indicates deletion
func (sm *SparseMatrix) Set(row, col int, value float64) {
	if value == 0. {
		delete(sm.Data, sm.RowColToIndex(row, col))
	} else {
		sm.Data[sm.RowColToIndex(row, col)] = value
	}
}

func (sm *SparseMatrix) SetIndex(idx int, value float64) {
	if value == 0. {
		delete(sm.Data, idx)
	} else {
		sm.Data[idx] = value
	}
}

func (sm *SparseMatrix) GetAllIndexes() (idxes []int) {
	idxes = make([]int, 0, len(sm.Data))
	for idx := range sm.Data {
		idxes = append(idxes, idx)
	}
	return
}

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

func (sm *SparseMatrix) Copy() *SparseMatrix {
	nsm := ZeroSparseMatrix(sm.Rows, sm.Cols)
	for idx, value := range sm.Data {
		nsm.Data[idx] = value
	}
	return nsm
}

func (sm *SparseMatrix) ToMatrix() *Matrix {
	nm := ZeroMatrix(sm.Rows, sm.Cols)
	for idx, value := range sm.Data {
		i, j := sm.IndexToRowCol(idx)
		nm.Set(i, j, value)
	}
	return nm
}

func (sm *SparseMatrix) T() *SparseMatrix {
	nsm := ZeroSparseMatrix(sm.Cols, sm.Rows)
	for idx, value := range sm.Data {
		nc, nr := sm.IndexToRowCol(idx)
		nsm.Set(nr, nc, value)
	}
	return nsm
}

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

func (sm *SparseMatrix) AddNum(n float64) *SparseMatrix {
	nsm := ZeroSparseMatrix(sm.Rows, sm.Cols)
	for idx, value := range sm.Data {
		if value+n != 0. {
			nsm.SetIndex(idx, value+n)
		}
	}
	return nsm
}

// Sparse Matrix Multiplication
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

// Determinant of Sparse Matrix
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
