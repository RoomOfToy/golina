package matrix

import (
	"math"
	"strconv"
	"testing"
)

func TestStretch(t *testing.T) {
	// 1D
	a := Data{{1}, {6}, {9}}
	matA := new(Matrix).Init(a)
	b := Data{{2}, {12}, {18}}
	if !MEqual(Stretch(matA, 2, 7), new(Matrix).Init(b)) {
		t.Fail()
	}
	// 2D
	a = Data{{1, 2}, {5, 6}, {7, 8}}
	matA = new(Matrix).Init(a)
	b = Data{{2, 2}, {10, 6}, {14, 8}}
	if !MEqual(Stretch(matA, 2, 1), new(Matrix).Init(b)) {
		t.Fail()
	}
	// 3D
	a = Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA = new(Matrix).Init(a)
	b = Data{{2, 2, 9}, {8, 5, 18}, {14, 8, 27}}
	if !MEqual(Stretch(matA, 2, 1, 3), new(Matrix).Init(b)) {
		t.Fail()
	}
	// 4D
	a = Data{{1, 2, 3, 4}, {4, 5, 6, 7}, {7, 8, 9, 10}}
	matA = new(Matrix).Init(a)
	b = Data{{2, 2, 9, 4}, {8, 5, 18, 7}, {14, 8, 27, 10}}
	if !MEqual(Stretch(matA, 2, 1, 3), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestRotate2D(t *testing.T) {
	a := Data{{1, 2}, {5, 6}, {7, 8}}
	matA := new(Matrix).Init(a)
	b := Data{{-2, 1}, {-6, 5}, {-8, 7}}
	if !MEqual(Rotate2D(matA, 90), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestRotate3D(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	b := Data{{1, -3, 2}, {4, -6, 5}, {7, -9, 8}}
	if !MEqual(Rotate3D(matA, 90, &Vector{1, 0, 0}), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestTranslate(t *testing.T) {
	// 2D
	a := Data{{1, 2}, {5, 6}, {7, 8}}
	matA := new(Matrix).Init(a)
	b := Data{{2, 4}, {6, 8}, {8, 10}}
	if !MEqual(Translate(matA, 1, 2), new(Matrix).Init(b)) {
		t.Fail()
	}
	// 3D
	a = Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA = new(Matrix).Init(a)
	b = Data{{3, 3, 6}, {6, 6, 9}, {9, 9, 12}}
	if !MEqual(Translate(matA, 2, 1, 3), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestShear2D(t *testing.T) {
	a := Data{{1, 2}, {5, 6}, {7, 8}}
	matA := new(Matrix).Init(a)
	b := Data{{3, 4}, {11, 16}, {15, 22}}
	if !MEqual(Shear2D(matA, 1, 2), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestShear3D(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	b := Data{{8, 5, 3}, {20, 17, 6}, {32, 29, 9}}
	if !MEqual(Shear3D(matA, 2, 1, 3), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestToAffineMatrix(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	b := Data{{1, 2, 3, 0}, {4, 5, 6, 0}, {7, 8, 9, 0}, {0, 0, 0, 1}}
	if !MEqual(ToAffineMatrix(matA), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestKabsch(t *testing.T) {
	matA := new(Matrix).Init(Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}})
	matB := Rotate3D(matA, 90, &Vector{1, 1, 1})
	linear, translation := Kabsch(matA, matB)
	// X -> AX + B, A: linear, B: translation
	// first linear transformation, then translate
	if !MEqual(matB, Translate(TransformOnRow(matA, ToAffineMatrix(linear)), translation.At(0), translation.At(1), translation.At(2))) {
		t.Fail()
	}
}

func BenchmarkRotate3D(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n)+"x3", func(b *testing.B) {
			m := GenerateRandomMatrix(n, 3)
			angle := GenerateRandomFloat()
			axis := GenerateRandomVector(3)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				Rotate3D(m, angle, axis)
			}
		})
	}
}

func BenchmarkKabsch(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n)+"x3", func(b *testing.B) {
			P := GenerateRandomMatrix(n, 3)
			Q := GenerateRandomMatrix(n, 3)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				Kabsch(P, Q)
			}
		})
	}
}
