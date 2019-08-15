package golina

import (
	"runtime"
	"sync"
)

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

type Data [][]float64
type Vector []float64

// Transpose struct implementing Array interface and return transpose of input Array
type Matrix struct {
	Array
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

// TODO: Max, Min, Rank, Det, Norm, Sum

// Vector convolve
// https://programming.guide/go/efficient-parallel-computation.html

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
	size := max(1, 1000000/n)

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

func init() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus) // Try to use all available CPUs.
}
