package matrix

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

// Vector
func TestVector_At(t *testing.T) {
	v := &Vector{1, 2, 3}
	if v.At(0) != 1 || v.At(-1) != 3 {
		t.Fail()
	}
}

func TestVector_Add(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	v2 := &Vector{1, 2, 3}
	if !VEqual(v1.Add(v2), &Vector{2, 4, 6}) {
		t.Fail()
	}
}

func TestVector_AddNum(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	n := 1
	if !VEqual(v1.AddNum(n), &Vector{2, 3, 4}) {
		t.Fail()
	}
}

func TestVector_Sub(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	v2 := &Vector{1, 2, 3}
	if !VEqual(v1.Sub(v2), &Vector{0, 0, 0}) {
		t.Fail()
	}
}

func TestVector_SubNum(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	n := 1
	if !VEqual(v1.SubNum(n), &Vector{0, 1, 2}) {
		t.Fail()
	}
}

func TestVector_MulNum(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	n := 2
	if !VEqual(v1.MulNum(n), &Vector{2, 4, 6}) {
		t.Fail()
	}
}

func TestVector_Dot(t *testing.T) {
	v1 := &Vector{1, 2, 3, 4, 5, 6}
	v2 := &Vector{6, 5, 4, 3, 2, 1}
	if v1.Dot(v2) != 56 || v1.Dot(v2) != v2.Dot(v1) {
		t.Fail()
	}
}

func TestVector_OuterProduct(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	v2 := &Vector{4, 5}
	a := Data{{1, 2}, {2, 4}, {3, 6}}
	if !MEqual(v1.OuterProduct(v2), new(Matrix).Init(a)) {
		t.Fail()
	}
}

func TestVector_Cross(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	v2 := &Vector{5, 8, 6}
	if !VEqual(v1.Cross(v2), &Vector{-12, 9, -2}) {
		t.Fail()
	}
}

func TestVector_SquareSum(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	if v1.SquareSum() != 14 {
		t.Fail()
	}
}

func TestVector_Norm(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	if !FloatEqual(v1.Norm(), 3.741657387) {
		t.Fail()
	}
}

func TestVector_Normalize(t *testing.T) {
	v1 := &Vector{1, 2, 3}
	if !VEqual(v1.Normalize(), &Vector{0.2672612419124244, 0.5345224838248488, 0.8017837257372732}) {
		t.Fail()
	}
}

func TestVector_ToMatrix(t *testing.T) {
	v := &Vector{1, 2, 3, 4, 5, 6}
	m := Data{{1, 2, 3}, {4, 5, 6}}
	if !MEqual(v.ToMatrix(2, 3), new(Matrix).Init(m)) {
		t.Fail()
	}
}

func TestVector_Sum(t *testing.T) {
	v := &Vector{1, 2, 3, 4, 5, 6}
	if !FloatEqual(v.Sum(), 21) {
		t.Fail()
	}
}

func TestVector_AbsSum(t *testing.T) {
	v := &Vector{1, -2, 3, -4, 5, -6}
	if !FloatEqual(v.AbsSum(), 21) {
		t.Fail()
	}
}

func TestVector_Mean(t *testing.T) {
	v := &Vector{1, 2, 3, 4, 5, 6}
	if !FloatEqual(v.Mean(), 3.5) {
		t.Fail()
	}
}

func TestVector_Variance(t *testing.T) {
	if !FloatEqual(7.666666666666667, (&Vector{1, 5, 7, 2, 6, 9}).Variance()) {
		t.Fail()
	}
}

func TestVector_StandardDeviation(t *testing.T) {
	if !FloatEqual(2.768874621, (&Vector{1, 5, 7, 2, 6, 9}).StandardDeviation()) {
		t.Fail()
	}
}

func TestVector_Tile(t *testing.T) {
	v := &Vector{1, 2, 3}
	m := Data{{1, 2, 3}, {1, 2, 3}}
	n := Data{{1, 1}, {2, 2}, {3, 3}}
	if !MEqual(v.Tile(0, 2), new(Matrix).Init(m)) {
		t.Fail()
	}
	if !MEqual(v.Tile(1, 2), new(Matrix).Init(n)) {
		t.Fail()
	}
}

func TestVector_Length(t *testing.T) {
	v := &Vector{1, 2, 3}
	if v.Length() != 3 {
		t.Fail()
	}
}

func TestVector_Max(t *testing.T) {
	v := &Vector{1, 2, 3}
	idx, value := v.Max()
	if idx != 2 || value != v.At(2) {
		t.Fail()
	}
}

func TestVector_Min(t *testing.T) {
	v := &Vector{1, 2, 3}
	idx, value := v.Min()
	if idx != 0 || value != v.At(0) {
		t.Fail()
	}
}

func TestVector_SortedToSortPairSlice(t *testing.T) {
	v := &Vector{1, 2, 3}
	sorted := v.SortedToSortPairSlice()
	for i := 1; i < v.Length(); i++ {
		if sorted[i].Value < sorted[i-1].Value {
			t.Fail()
		}
	}
}

func TestVector_SortedAscending(t *testing.T) {
	v := &Vector{1, 2, 3}
	sorted := v.SortedAscending()
	// sorted function return a new vector
	if v == sorted {
		t.Fail()
	}
	for i := 1; i < sorted.Length(); i++ {
		if sorted.At(i) < sorted.At(i-1) {
			t.Fail()
		}
	}
}

func TestVector_SortedDescending(t *testing.T) {
	v := &Vector{1, 2, 3}
	sorted := v.SortedDescending()
	if v == sorted {
		t.Fail()
	}
	for i := 1; i < sorted.Length(); i++ {
		if sorted.At(i) > sorted.At(i-1) {
			t.Fail()
		}
	}
}

func TestVector_Reversed(t *testing.T) {
	v := &Vector{1, 2, 3}
	if !VEqual(v, v.Reversed().Reversed()) {
		t.Fail()
	}
}

func TestVector_Unique(t *testing.T) {
	v := &Vector{1, 1, 2, 3, 4, 5, 2, 6, 9, 8, 6, 11, 12}
	vv := &Vector{1, 2, 3, 4, 5, 6, 8, 9, 11, 12}
	if !VEqual(v.Unique().SortedAscending(), vv) {
		t.Fail()
	}
}

func TestVector_UniqueWithCount(t *testing.T) {
	v := &Vector{1, 1, 2, 3, 4, 5, 2, 6, 9, 8, 6, 11, 12}
	vv := map[float64]int{1: 2, 2: 2, 3: 1, 4: 1, 5: 1, 6: 2, 8: 1, 9: 1, 11: 1, 12: 1}
	vc := v.UniqueWithCount()
	for k, val := range vc {
		if vv[k] != val {
			t.Fail()
		}
	}
}

func TestCrossCov(t *testing.T) {
	u, v := &Vector{1, 2, 3}, &Vector{4, 5, 6}
	cc := CrossCov(u, v)
	fmt.Println(cc)
}

func TestCrossCorr(t *testing.T) {
	u, v := &Vector{1, 2, 3}, &Vector{4, 5, 6}
	cc := CrossCorr(u, v)
	fmt.Println(cc)
}

// Vector convolve
func TestConvolve(t *testing.T) {
	size := 10000
	u := GenerateRandomVector(size)
	v := GenerateRandomVector(size)

	res := Convolve(u, v)
	if len(*res) != size+size-1 {
		t.Fail()
	}
}

func TestVector_String(t *testing.T) {
	v := &Vector{1, 2, 3}
	if v.String() != "{1.000000, 2.000000, 3.000000}\n" {
		t.Fail()
	}
}

func TestARRange(t *testing.T) {
	if !VEqual(ARRange(0, 1, 6), &Vector{0, 1, 2, 3, 4, 5}) {
		t.Fail()
	}
}

func TestVector_Concatenate(t *testing.T) {
	u, v := &Vector{1, 2, 3}, &Vector{4, 5, 6}
	if !VEqual(u.Concatenate(v), &Vector{1, 2, 3, 4, 5, 6}) {
		t.Fail()
	}
}

func BenchmarkVector_SquareSum(b *testing.B) {
	for k := 1.0; k <= 5; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			v := GenerateRandomVector(n)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				v.SquareSum()
			}
		})
	}
}

func BenchmarkConvolve(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			u := GenerateRandomVector(n)
			v := GenerateRandomVector(n)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				Convolve(u, v)
			}
		})
	}
}
