package matrix

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"testing"
)

func TestZeroSparseMatrix(t *testing.T) {
	SA := ZeroSparseMatrix(100, 100)
	if len(SA.Data) != 0 {
		t.Fail()
	}
}

func TestNewSparseMatrix(t *testing.T) {
	SA := NewSparseMatrix(map[int]float64{1: 0.1, 50: 0.2}, 100, 100)
	if len(SA.Data) != 2 {
		t.Fail()
	}
}

func TestSparseMatrix_GetSubSparseMatrix(t *testing.T) {
	SA := NewSparseMatrix(map[int]float64{1: 0.1, 50: 0.2}, 100, 100)
	SB := SA.GetSubSparseMatrix(0, 0, 25, 25)
	if SB.AtIndex(1) != 0.1 {
		t.Fail()
	}
}

func TestSparseMatrix_RowColToIndex(t *testing.T) {
	SA := ZeroSparseMatrix(10, 10)
	if SA.RowColToIndex(5, 5) != 55 {
		t.Fail()
	}
}

func TestSparseMatrix_IndexToRowCol(t *testing.T) {
	SA := ZeroSparseMatrix(10, 10)
	r, c := SA.IndexToRowCol(55)
	if r != 5 || c != 5 {
		t.Fail()
	}
}

func TestSparseMatrix_At(t *testing.T) {
	SA := NewSparseMatrix(map[int]float64{1: 0.1, 500: 0.2}, 100, 100)
	if SA.At(0, 1) != 0.1 || SA.At(5, 0) != 0.2 {
		t.Fail()
	}
}

func TestSparseMatrix_AtIndex(t *testing.T) {
	SA := NewSparseMatrix(map[int]float64{1: 0.1, 500: 0.2}, 100, 100)
	if SA.AtIndex(1) != 0.1 || SA.AtIndex(500) != 0.2 {
		t.Fail()
	}
}

func TestSparseMatrix_Set(t *testing.T) {
	SA := ZeroSparseMatrix(50, 50)
	SA.Set(10, 10, 0.3)
	if SA.At(10, 10) != 0.3 {
		t.Fail()
	}
}

func TestSparseMatrix_SetIndex(t *testing.T) {
	SA := ZeroSparseMatrix(50, 50)
	SA.SetIndex(500, 0.3)
	if SA.AtIndex(500) != 0.3 {
		t.Fail()
	}
}

func TestSparseMatrix_GetAllIndexes(t *testing.T) {
	SA := NewSparseMatrix(map[int]float64{1: 0.1, 500: 0.2}, 100, 100)
	s := SA.GetAllIndexes()
	sort.Ints(s)
	if s[0] != 1 && s[1] != 500 {
		t.Fail()
	}
}

func TestSparseMatrix_Row(t *testing.T) {
	SA := NewSparseMatrix(map[int]float64{1: 0.1, 500: 0.2}, 100, 100)
	v := make(Vector, 100)
	v[1] = 0.1
	if !VEqual(SA.Row(0), &v) {
		t.Fail()
	}
}

func TestSparseMatrix_Col(t *testing.T) {
	SA := NewSparseMatrix(map[int]float64{1: 0.1, 500: 0.2}, 100, 100)
	v := make(Vector, 100)
	v[0] = 0.1
	if !VEqual(SA.Col(1), &v) {
		t.Fail()
	}
}

func TestSparseMatrix_FindFirstNonZeroInSubMatrix(t *testing.T) {
	SA := NewSparseMatrix(map[int]float64{50: 0.1, 600: 0.2}, 100, 100)
	if SA.FindFirstNonZeroInSubMatrix(25) != 50 || SA.FindFirstNonZeroInSubMatrix(500) != 600 {
		t.Fail()
	}
}

func TestSparseMatrix_Copy(t *testing.T) {
	SA := GenerateRandomSparseMatrix(20, 20, 10)
	SB := SA.Copy()
	if SA.Rows != SB.Rows || SA.Cols != SB.Cols || SA.Offset != SB.Offset {
		t.Fail()
	}
	for i := range SA.Data {
		if SB.Data[i] != SA.Data[i] {
			t.Fail()
		}
	}
}

func TestSparseMatrix_ToMatrix(t *testing.T) {
	SA := NewSparseMatrix(map[int]float64{5: 0.1, 26: 0.2}, 20, 20)
	ma := SA.ToMatrix()
	r, c := ma.Dims()
	if r != 20 || c != 20 {
		t.Fail()
	}
	if ma.At(0, 5) != 0.1 || ma.At(1, 6) != 0.2 {
		t.Fail()
	}
}

func TestSparseMatrix_T(t *testing.T) {
	SA := GenerateRandomSparseMatrix(20, 20, 10)
	SB := SA.T()
	for idx, value := range SA.Data {
		r, c := SA.IndexToRowCol(idx)
		if value != SB.At(c, r) {
			t.Fail()
		}
	}
}

func TestSparseMatrix_Add(t *testing.T) {
	SA := GenerateRandomSparseMatrix(20, 20, 10)
	SB := GenerateRandomSparseMatrix(20, 20, 10)
	SC := SA.Add(SB)
	for idx, value := range SC.Data {
		if value == 0. || !FloatEqual(value, SA.AtIndex(idx)+SB.AtIndex(idx)) {
			t.Fail()
		}
	}
}

func TestSparseMatrix_AddNum(t *testing.T) {
	SA := GenerateRandomSparseMatrix(20, 20, 10)
	n := GenerateRandomFloat()
	SB := SA.AddNum(n)
	for idx, value := range SB.Data {
		if value == 0. || !FloatEqual(value, SA.AtIndex(idx)+n) {
			t.Fail()
		}
	}
}

func TestSparseMatrix_Mul(t *testing.T) {
	SA := GenerateRandomSparseMatrix(30, 20, 10)
	SB := GenerateRandomSparseMatrix(20, 10, 10)
	SC := SA.Mul(SB)
	if SC.Rows != 30 || SC.Cols != 10 {
		t.Fail()
	}
	for _, value := range SC.Data {
		if value == 0. {
			t.Fail()
		}
	}
}

func TestSparseMatrix_MulNum(t *testing.T) {
	SA := GenerateRandomSparseMatrix(20, 20, 10)
	n := GenerateRandomFloat()
	SB := SA.MulNum(n)
	for idx, value := range SB.Data {
		if value == 0. || !FloatEqual(value, SA.AtIndex(idx)*n) {
			t.Fail()
		}
	}
}

func TestSparseMatrix_MulVec(t *testing.T) {
	SA := GenerateRandomSparseMatrix(30, 20, 10)
	v := GenerateRandomVector(20)
	nv := SA.MulVec(v)

	if nv.Length() != SA.Rows {
		t.Fail()
	}
}

func TestSparseMatrix_Det(t *testing.T) {
	SA := GenerateRandomSparseMatrix(10, 10, 50)
	mat := SA.ToMatrix()
	fmt.Println(mat.Dims())
	fmt.Println(mat)
	fmt.Println(SA.Det())
}

func BenchmarkGenerateRandomSparseMatrix(b *testing.B) {
	for k := 5; k <= 8; k++ {
		n := int(math.Pow(10, float64(k)))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			for i := 1; i < b.N; i++ {
				GenerateRandomSparseMatrix(n, n, n/1000)
			}
		})
	}
}

func BenchmarkSparseMatrix_Mul(b *testing.B) {
	for k := 5; k <= 6; k++ {
		n := int(math.Pow(10, float64(k)))
		b.Run("size-"+strconv.Itoa(n)+"x"+strconv.Itoa(n), func(b *testing.B) {
			SA := GenerateRandomSparseMatrix(n, n, n/1000)
			SB := GenerateRandomSparseMatrix(n, n, n/1000)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				SA.Mul(SB)
			}
		})
	}
}

func BenchmarkSparseMatrix_MulVec(b *testing.B) {
	for k := 5; k <= 8; k++ {
		n := int(math.Pow(10, float64(k)))
		b.Run("size-"+strconv.Itoa(n)+"x"+strconv.Itoa(n), func(b *testing.B) {
			SA := GenerateRandomSparseMatrix(n, n, n/1000)
			v := GenerateRandomVector(n)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				SA.MulVec(v)
			}
		})
	}
}
