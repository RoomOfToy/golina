package golina

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func generateRandomVector(size int) *Vector {
	slice := make(Vector, size, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = rand.Float64() - rand.Float64()
	}
	return &slice
}

// https://blog.karenuorteva.fi/go-unit-test-setup-and-teardown-db1601a796f2#.2aherx2z5

func TestMatrix_Init(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	if matA._array == nil {
		t.Fail()
	}
}

func TestMatrix_Dims(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	row, col := matA.Dims()
	if row != 3 || col != 3 {
		t.Fail()
	}
}

func TestMatrix_At(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	if matA.At(1, 1) != 5 {
		t.Fail()
	}
}

func TestMatrix_Set(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	matA.Set(1, 1, 10)
	if matA.At(1, 1) != 10 {
		t.Fail()
	}
}

func TestMatrix_T(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	matAT := matA.T()
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if matAT.At(i, j) != matA.At(j, i) {
				t.Fail()
			}
		}
	}
}

func TestMatrix_Row(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	row := matA.Row(1)
	if !VEqual(row, &Vector{4, 5, 6}) {
		t.Fail()
	}
}

func TestMatrix_Col(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	col := matA.Col(1)
	if !VEqual(col, &Vector{2, 5, 8}) {
		t.Fail()
	}
}

func TestEqual(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	if !Equal(matA, matA) {
		t.Fail()
	}
}

func TestMatrix_Max(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	maxA := matA.Max()
	if maxA.value != 9 && maxA.row != 2 && maxA.col != 2 {
		t.Fail()
	}
}

func TestMatrix_Min(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	minA := matA.Min()
	if minA.value != 1 && minA.row != 0 && minA.col != 0 {
		t.Fail()
	}
}

func TestCopy(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	matB := Copy(matA)
	if !Equal(matA, matB) {
		t.Fail()
	}
}

func TestEmpty(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	matB := Empty(matA)
	row1, col1 := matA.Dims()
	row2, col2 := matB.Dims()
	if [2]int{row1, col1} != [2]int{row2, col2} {
		t.Fail()
	}
}

func TestZeroMatrix(t *testing.T) {
	a := Data{{0}, {0}}
	matA := ZeroMatrix(2, 1)
	if !Equal(matA, new(Matrix).Init(a)) {
		t.Fail()
	}
}

func TestIdentityMatrix(t *testing.T) {
	a := Data{{1, 0}, {0, 1}}
	matA := IdentityMatrix(2)
	if !Equal(matA, new(Matrix).Init(a)) {
		t.Fail()
	}
}

func TestSwapRow(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	matB := Copy(matA)
	SwapRow(matB, 1, 2)
	if !VEqual(matA.Row(1), matB.Row(2)) || !VEqual(matA.Row(2), matB.Row(1)) {
		t.Fail()
	}
}

func TestMatrix_Rank(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	if matA.Rank() != 2 {
		t.Fail()
	}

	b := Data{{0, 1, 2}, {-1, -2, 1}, {2, 7, 8}}
	matB := new(Matrix).Init(b)
	if matB.Rank() != 3 {
		t.Fail()
	}
}

func TestMatrix_Det(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	if matA.Det() != 0 {
		t.Fail()
	}

	b := Data{{32, 12, 1}, {6, 3, 45}, {9, 2, 1}}
	matB := new(Matrix).Init(b)
	if matB.Det() != 1989 {
		t.Fail()
	}
}

func TestMatrix_Adj(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	b := Data{{-500, 500, 500}, {300, -300, -300}, {-100, 100, 100}}
	if !Equal(matA.Adj(), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestMatrix_Inverse(t *testing.T) {
	a := Data{{32, 12, 1}, {6, 3, 45}, {9, 2, 1}}
	matA := new(Matrix).Init(a)
	b := Data{{-0.04374057315233785821, -0.00502765208647561595, 0.26998491704374057313},
		{0.2006033182503770739, 0.0115635997988939167, -0.7209653092006033182},
		{-0.007541478129713423831, 0.022121669180492709902, 0.012066365007541478129}}
	if !Equal(matA.Inverse(), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestMatrix_Add(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	b := Data{{32, 12, 1}, {6, 3, 45}, {9, 2, 1}}
	matB := new(Matrix).Init(b)
	matC := matA.Add(matB)
	c := Data{{42, 32, 11}, {-14, -27, 55}, {39, 52, 1}}
	if !Equal(matC, new(Matrix).Init(c)) {
		t.Fail()
	}
}

func TestMatrix_Sub(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	b := Data{{32, 12, 1}, {6, 3, 45}, {9, 2, 1}}
	matB := new(Matrix).Init(b)
	matC := matA.Sub(matB)
	c := Data{{-22, 8, 9}, {-26, -33, -35}, {21, 48, -1}}
	if !Equal(matC, new(Matrix).Init(c)) {
		t.Fail()
	}
}

func TestMatrix_Mul(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	b := Data{{32, 12}, {6, 3}, {9, 2}}
	matB := new(Matrix).Init(b)
	matC := matA.Mul(matB)
	c := Data{{530, 200}, {-730, -310}, {1260, 510}}
	if !Equal(matC, new(Matrix).Init(c)) {
		t.Fail()
	}
}

func TestMatrix_Pow(t *testing.T) {
	a := Data{{10, 20, 10}, {-20, -30, 10}, {30, 50, 0}}
	matA := new(Matrix).Init(a)
	if !Equal(matA.Pow(0), IdentityMatrix(3)) {
		t.Fail()
	}
	if !Equal(matA.Pow(1), matA) {
		t.Fail()
	}
	b := Data{{0, 100, 300}, {700, 1000, -500}, {-700, -900, 800}}
	if !Equal(matA.Pow(2), new(Matrix).Init(b)) {
		t.Fail()
	}
}

// Vector convolve
func TestConvolve(t *testing.T) {
	size := 10000
	u := generateRandomVector(size)
	v := generateRandomVector(size)

	res := Convolve(u, v)
	if len(*res) != size+size-1 {
		t.Fail()
	}
}

// https://medium.com/justforfunc/analyzing-the-performance-of-go-functions-with-benchmarks-60b8162e61c6
func BenchmarkConvolve(b *testing.B) {
	convolves := []struct {
		name string
		fun  func(u, v *Vector) *Vector
	}{
		{"size-10", Convolve},
		{"size-100", Convolve},
		{"size-1000", Convolve},
	}

	for _, convolve := range convolves {
		for k := 1.0; k <= 3; k++ {
			n := int(math.Pow(10, k))
			b.Run(convolve.name, func(b *testing.B) {
				for i := 1; i < b.N; i++ {
					u := generateRandomVector(n)
					v := generateRandomVector(n)
					Convolve(u, v)
				}
			})
		}
	}
}
