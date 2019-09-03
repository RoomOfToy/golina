package golina

// Sparse Matrix
//	https://en.wikipedia.org/wiki/Sparse_matrix
//	Dictionary of keys (DOK)
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
		panic("invalid index")
	}
	return entry
}

func (sm *SparseMatrix) AtIndex(idx int) float64 {
	entry, ok := sm.Data[idx]
	if !ok {
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
		i = max(0, i)
		j = max(0, j)
		rows = min(sm.Rows-i, rows)
		rows = min(sm.Cols-j, cols)
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
