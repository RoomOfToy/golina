package spatial

import (
	"golina/matrix"
	"sort"
)

// k-nearest-neighbors of some vector to all vectors in dataSet
func KNearestNeighbors(dataSet *matrix.Matrix, v *matrix.Vector, k int, distFunc func(v1, v2 *matrix.Vector) float64) *matrix.Matrix {
	row, _ := dataSet.Dims()
	if k > row {
		k = row
	}
	distSlice := make(matrix.SortPairSlice, row)
	for i, vd := range dataSet.Data {
		distSlice[i] = matrix.SortPair{i, distFunc(v, &vd)}
	}
	sort.Sort(distSlice) // sort default is ascending
	retM := matrix.ZeroMatrix(k, len(*v))
	for i := range retM.Data {
		retM.Data[i] = dataSet.Data[distSlice[i].Key]
	}
	return retM
}

func KNearestNeighborsWithDistance(dataSet *matrix.Matrix, v *matrix.Vector, k int, distFunc func(v1, v2 *matrix.Vector) float64) *matrix.Matrix {
	row, _ := dataSet.Dims()
	if k > row {
		k = row
	}
	distSlice := make(matrix.SortPairSlice, row)
	for i, vd := range dataSet.Data {
		distSlice[i] = matrix.SortPair{i, distFunc(v, &vd)}
	}
	sort.Sort(distSlice) // sort default is ascending
	retM := matrix.ZeroMatrix(k, len(*v))
	for i := range retM.Data {
		retM.Data[i] = append(dataSet.Data[distSlice[i].Key], float64(distSlice[i].Key), distSlice[i].Value) // output idx, distance for observation
	}
	return retM
}
