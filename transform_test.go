package golina

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
	if !Equal(Stretch(matA, 2, 7), new(Matrix).Init(b)) {
		t.Fail()
	}
	// 2D
	a = Data{{1, 2}, {5, 6}, {7, 8}}
	matA = new(Matrix).Init(a)
	b = Data{{2, 2}, {10, 6}, {14, 8}}
	if !Equal(Stretch(matA, 2, 1), new(Matrix).Init(b)) {
		t.Fail()
	}
	// 3D
	a = Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA = new(Matrix).Init(a)
	b = Data{{2, 2, 9}, {8, 5, 18}, {14, 8, 27}}
	if !Equal(Stretch(matA, 2, 1, 3), new(Matrix).Init(b)) {
		t.Fail()
	}
	// 4D
	a = Data{{1, 2, 3, 4}, {4, 5, 6, 7}, {7, 8, 9, 10}}
	matA = new(Matrix).Init(a)
	b = Data{{2, 2, 9, 4}, {8, 5, 18, 7}, {14, 8, 27, 10}}
	if !Equal(Stretch(matA, 2, 1, 3), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestRotate2D(t *testing.T) {
	a := Data{{1, 2}, {5, 6}, {7, 8}}
	matA := new(Matrix).Init(a)
	b := Data{{-2, 1}, {-6, 5}, {-8, 7}}
	if !Equal(Rotate2D(matA, 90), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestRotate3D(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	b := Data{{1, 3, -2}, {4, 6, -5}, {7, 9, -8}}
	if !Equal(Rotate3D(matA, 90, &Vector{1, 0, 0}), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestTranslate(t *testing.T) {
	// 2D
	a := Data{{1, 2}, {5, 6}, {7, 8}}
	matA := new(Matrix).Init(a)
	b := Data{{2, 4}, {6, 8}, {8, 10}}
	if !Equal(Translate(matA, 1, 2), new(Matrix).Init(b)) {
		t.Fail()
	}
	// 3D
	a = Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA = new(Matrix).Init(a)
	b = Data{{3, 3, 6}, {6, 6, 9}, {9, 9, 12}}
	if !Equal(Translate(matA, 2, 1, 3), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestShear2D(t *testing.T) {
	a := Data{{1, 2}, {5, 6}, {7, 8}}
	matA := new(Matrix).Init(a)
	b := Data{{3, 4}, {11, 16}, {15, 22}}
	if !Equal(Shear2D(matA, 1, 2), new(Matrix).Init(b)) {
		t.Fail()
	}
}

func TestShear3D(t *testing.T) {
	a := Data{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	matA := new(Matrix).Init(a)
	b := Data{{8, 5, 3}, {20, 17, 6}, {32, 29, 9}}
	if !Equal(Shear3D(matA, 2, 1, 3), new(Matrix).Init(b)) {
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
