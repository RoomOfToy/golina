package golina

// Array interface

type Array interface {
	// dimensions
	Dims() (row, col int)

	// value at index(row i, col j), panic if not access
	At(i, j int) float64

	// set value at index(row i, col j), panic if not access
	Set(i, j int, value float64)

	// transpose array
	T() Array
}

// Transpose struct implementing Array interface and return transpose of input Array
type Matrix struct {
	Array
	_array [][]float64  // row-wise
}

func (t *Matrix) New(array [][]float64) {
	t._array = array
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

func (t *Matrix) T() Matrix {
	row, col := t.Dims()
	ntArray := make([][]float64, col)
	for i := 0; i < col; i ++ {
		ntArray[i] = make([]float64, row)
		for j := 0; j < row; j ++ {
			ntArray[i][j] = t._array[j][i]
		}
	}
	copy(ntArray, t._array)
	var nt Matrix
	nt.New(ntArray)
	return nt
}
