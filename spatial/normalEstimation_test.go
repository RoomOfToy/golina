package spatial

import (
	"golina/matrix"
	"math"
	"strconv"
	"testing"
)

func TestPlanePcaEigen(t *testing.T) {
	a := matrix.Data{{-1, 2, -1}, {2, -1, -2}, {-1, 3, -1}}
	points := new(matrix.Matrix).Init(a)
	// np.linalg.eig(np.cov(a.transpose()))
	b := matrix.Vector{3.16227766e-01, 2.45758907e-15, 9.48683298e-01}
	if !matrix.VEqual(PlanePcaEigen(points), &b) && !matrix.VEqual(PlanePcaEigen(points).MulNum(-1.), &b) {
		t.Fail()
	}
}

func TestPlanePcaSVD(t *testing.T) {
	a := matrix.Data{{-1, 2, -1}, {2, -1, -2}, {-1, 3, -1}}
	points := new(matrix.Matrix).Init(a)
	b := matrix.Vector{3.16227766e-01, 2.45758907e-15, 9.48683298e-01}
	if !matrix.VEqual(PlanePcaSVD(points), &b) && !matrix.VEqual(PlanePcaSVD(points).MulNum(-1.), &b) {
		t.Fail()
	}
}

func TestPlaneLinearSolveWeighted(t *testing.T) {
	a := matrix.Data{{-1, 2, -1}, {2, -1, -2}, {-1, 3, -1}}
	points := new(matrix.Matrix).Init(a)
	b := matrix.Vector{3.16227766e-01, 2.45758907e-15, 9.48683298e-01}
	if !matrix.VEqual(PlaneLinearSolveWeighted(points), &b) {
		t.Fail()
	}
}

func BenchmarkPlanePcaEigen(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n)+"x3", func(b *testing.B) {
			points := matrix.GenerateRandomMatrix(n, 3)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				PlanePcaEigen(points)
			}
		})
	}
}

func BenchmarkPlanePcaSVD(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n)+"x3", func(b *testing.B) {
			points := matrix.GenerateRandomMatrix(n, 3)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				PlanePcaSVD(points)
			}
		})
	}
}

func BenchmarkPlaneLinearSolveWeighted(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n)+"x3", func(b *testing.B) {
			points := matrix.GenerateRandomMatrix(n, 3)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				PlaneLinearSolveWeighted(points)
			}
		})
	}
}
