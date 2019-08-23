package golina

import (
	"testing"
)

func TestKNearestNeighbors(t *testing.T) {
	dataSet := GenerateRandomMatrix(10, 3)
	v := GenerateRandomVector(3)
	knn := KNearestNeighbors(dataSet, v, 2, EuclideanDistance)
	if v.Sub(knn.Row(0)).Norm() > v.Sub(knn.Row(1)).Norm() {
		t.Fail()
	}
}
