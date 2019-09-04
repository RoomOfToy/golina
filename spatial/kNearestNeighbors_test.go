package spatial

import (
	"golina/matrix"
	"math"
	"strconv"
	"testing"
)

func TestKNearestNeighbors(t *testing.T) {
	dataSet := matrix.GenerateRandomMatrix(10, 3)
	v := matrix.GenerateRandomVector(3)
	knn := KNearestNeighbors(dataSet, v, 2, EuclideanDistance)
	if v.Sub(knn.Row(0)).Norm() > v.Sub(knn.Row(1)).Norm() {
		t.Fail()
	}
}

func TestKNearestNeighborsWithDistance(t *testing.T) {
	dataSet := matrix.GenerateRandomMatrix(10, 3)
	v := matrix.GenerateRandomVector(3)
	knn := KNearestNeighborsWithDistance(dataSet, v, 2, EuclideanDistance)
	if knn.Row(0).At(-1) > knn.Row(1).At(-1) {
		t.Fail()
	}
}

func BenchmarkKNearestNeighbors(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n)+"x3", func(b *testing.B) {
			dataSet := matrix.GenerateRandomMatrix(n, 3)
			v := matrix.GenerateRandomVector(3)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				KNearestNeighbors(dataSet, v, n/2, EuclideanDistance)
			}
		})
	}
}
