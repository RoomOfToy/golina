package golina

import (
	"math"
	"strconv"
	"testing"
)

func TestPlanePCA(t *testing.T) {
	a := Data{{-1, 2, -1}, {2, -1, -2}, {-1, 3, -1}}
	points := new(Matrix).Init(a)
	// np.linalg.eig(np.cov(a.transpose()))
	b := Vector{3.16227766e-01, 2.45758907e-15, 9.48683298e-01}
	if !VEqual(PlanePCA(points), &b) {
		t.Fail()
	}
}

func TestPlaneLinearSolveWeighted(t *testing.T) {
	a := Data{{-1, 2, -1}, {2, -1, -2}, {-1, 3, -1}}
	points := new(Matrix).Init(a)
	b := Vector{3.16227766e-01, 2.45758907e-15, 9.48683298e-01}
	if !VEqual(PlaneLinearSolveWeighted(points), &b) {
		t.Fail()
	}
}

func BenchmarkPlanePCA(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n)+"x3", func(b *testing.B) {
			points := GenerateRandomMatrix(n, 3)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				PlanePCA(points)
			}
		})
	}
}

func BenchmarkPlaneLinearSolveWeighted(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n)+"x3", func(b *testing.B) {
			points := GenerateRandomMatrix(n, 3)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				PlaneLinearSolveWeighted(points)
			}
		})
	}
}
