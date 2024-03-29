package matrix

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
)

// Ternary simple function for simulating ternary operator
func Ternary(statement bool, a, b interface{}) interface{} {
	if statement {
		return a
	}
	return b
}

// AbsInt returns absolute value of input int
func AbsInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// MinInt returns minimum value of input two ints
func MinInt(x, y int) int {
	return Ternary(x > y, y, x).(int)
}

// MaxInt returns maximum value of input two ints
func MaxInt(x, y int) int {
	return Ternary(x > y, x, y).(int)
}

// GetFloat64 transfers input golang number type into float64
func GetFloat64(x interface{}) float64 {
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

// SortPair struct for sorting according to value
type SortPair struct {
	Key   int
	Value float64
}

// SortPairSlice slice of SortPair
type SortPairSlice []SortPair

func (sps SortPairSlice) Swap(i, j int) {
	sps[i], sps[j] = sps[j], sps[i]
}

func (sps SortPairSlice) Len() int {
	return len(sps)
}

func (sps SortPairSlice) Less(i, j int) bool {
	return sps[i].Value < sps[j].Value
}

// FloatEqual checks whether two float numbers are equal, defined by threshold EPS
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

// VEqual checks whether two vector are equal, based on `FloatEqual`
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

// MEqual checks whether two matrix are equal, based on `VEqual`
func MEqual(mat1, mat2 *Matrix) bool {
	row1, col1 := mat1.Dims()
	row2, col2 := mat2.Dims()
	if [2]int{row1, col1} != [2]int{row2, col2} {
		return false
	}
	for i, col := range mat1.Data {
		if !VEqual(&col, mat2.Row(i)) {
			return false
		}
	}
	return true
}

// GenerateRandomFloat generates a random float64
func GenerateRandomFloat() float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64() - rand.Float64()
}

// GenerateRandomVector generates a vector with random float64
func GenerateRandomVector(size int) *Vector {
	slice := make(Vector, size, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = rand.Float64() - rand.Float64()
	}
	return &slice
}

// GenerateRandomSymmetric33Matrix generates a 3 x 3 matrix with random float64
func GenerateRandomSymmetric33Matrix() *Matrix {
	entries := *GenerateRandomVector(6)
	m := ZeroMatrix(3, 3)
	m.Set(0, 0, entries[0])
	m.Set(1, 1, entries[1])
	m.Set(2, 2, entries[2])
	m.Set(0, 1, entries[3])
	m.Set(1, 0, entries[3])
	m.Set(0, 2, entries[4])
	m.Set(2, 0, entries[4])
	m.Set(1, 2, entries[5])
	m.Set(2, 1, entries[5])
	return m
}

// GenerateRandomSquareMatrix generates a `size x size` square matrix with random float64
func GenerateRandomSquareMatrix(size int) *Matrix {
	return GenerateRandomMatrix(size, size)
}

// GenerateRandomMatrix generates a `row x col` matrix with random float64
func GenerateRandomMatrix(row, col int) *Matrix {
	rows := make(Data, row)
	for i := range rows {
		rows[i] = *GenerateRandomVector(col)
	}
	m := new(Matrix).Init(rows)
	return m
}

// GenerateRandomSparseMatrix generates a `rows x cols` sparse matrix with `entriesNum` elements
func GenerateRandomSparseMatrix(rows, cols, entriesNum int) *SparseMatrix {
	nsm := ZeroSparseMatrix(rows, cols)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < entriesNum; i++ {
		nsm.Set(rand.Intn(rows), rand.Intn(cols), rand.Float64()-rand.Float64())
	}
	return nsm
}

/*
// Vector to iterable
func VectorIter(v *Vector) interface{} {
	return *v
}

// Matrix to row iterable
func MatrixRowIter(t *Matrix) interface{} {
	return t.Data
}

// Matrix to element iterable
// 	row-wise
func MatrixElementIter(t *Matrix) interface{} {
	return *(t.Flat())
}

// Map function on iterable
func Map(input interface{}, mapper func(interface{}) interface{}) (output interface{}) {
	val := reflect.ValueOf(input)
	out := make([]interface{}, val.Len())
	wg := &sync.WaitGroup{}
	for i := 0; i < val.Len(); i++ {
		wg.Add(1)
		go func(i int) {
			wg.Done()
			out[i] = mapper(val.Index(i).Interface())
		}(i)
	}
	wg.Wait()
	return out
}

// Reduce iterable by function
func Reduce(input interface{}, reducer func(interface{}, interface{}) interface{}) interface{} {
	val := reflect.ValueOf(input)
	tmp := val.Index(0).Interface()
	for i := 0; i < val.Len()-1; i++ {
		tmp = reducer(tmp, val.Index(i).Interface())
	}
	return tmp
}

// Filter iterable by function
func Filter(input interface{}, filter func(interface{}) bool) interface{} {
	val := reflect.ValueOf(input)
	out := make([]interface{}, 0, val.Len())
	for i := 0; i < val.Len(); i++ {
		if filter(val.Index(i).Interface()) {
			out = append(out, val.Index(i).Interface())
		}
	}
	return out
}
*/

func getFileSize(filename string) int64 {
	fileStat, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
	fileSize := fileStat.Size()
	return fileSize
}

// Load3DToMatrix reads 3D data into matrix
func Load3DToMatrix(path string) (*Matrix, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileSize := getFileSize(path)
	est := 3 * fileSize / (4*8*3 + 1*2)
	lines := make(Data, 0, est)

	var x, y, z float64
	for {
		rowNum, err := fmt.Fscanln(file, &x, &y, &z)
		if rowNum == 0 || err != nil {
			break
		}
		lines = append(lines, Vector{x, y, z})
	}
	return new(Matrix).Init(lines), err
}

// Load2DToMatrix reads 2D data into matrix
func Load2DToMatrix(path string) (*Matrix, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileSize := getFileSize(path)
	est := 3 * fileSize / (4*8*2 + 1)
	lines := make(Data, 0, est)

	var x, y float64
	for {
		rowNum, err := fmt.Fscanln(file, &x, &y)
		if rowNum == 0 || err != nil {
			break
		}
		lines = append(lines, Vector{x, y})
	}
	return new(Matrix).Init(lines), err
}

// WriteMatrixToTxt writes matrix data into file
func WriteMatrixToTxt(path string, t *Matrix) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, c := t.Dims()

	for i := range t.Data {
		for j := range t.Data[i] {
			_, err = fmt.Fprintf(file, "%f ", t.Data[i][j])
			if j == c-1 {
				_, err = fmt.Fprintf(file, "\n")
			}
		}
	}
	return err
}
