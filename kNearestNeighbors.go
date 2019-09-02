package golina

import "sort"

type SortPair struct {
	Key   int
	Value float64
}

type SortPairSlice []SortPair

func (sps SortPairSlice) Swap(i, j int) {
	sps[i], sps[j] = sps[j], sps[i]
}

func (sps SortPairSlice) Len() int {
	return len(sps)
}

func (sps SortPairSlice) Less(i, j int) bool {
	return sps[i].Value < sps[j].Value
}

// k-nearest-neighbors of some vector to all vectors in dataSet
func KNearestNeighbors(dataSet *Matrix, v *Vector, k int, distFunc func(v1, v2 *Vector) float64) *Matrix {
	row, _ := dataSet.Dims()
	if k > row {
		k = row
	}
	distSlice := make(SortPairSlice, row)
	for i, vd := range dataSet.Data {
		distSlice[i] = SortPair{i, distFunc(v, &vd)}
	}
	sort.Sort(distSlice) // sort default is ascending
	retM := ZeroMatrix(k, len(*v))
	for i := range retM.Data {
		retM.Data[i] = dataSet.Data[distSlice[i].Key]
	}
	return retM
}

func KNearestNeighborsWithDistance(dataSet *Matrix, v *Vector, k int, distFunc func(v1, v2 *Vector) float64) *Matrix {
	row, _ := dataSet.Dims()
	if k > row {
		k = row
	}
	distSlice := make(SortPairSlice, row)
	for i, vd := range dataSet.Data {
		distSlice[i] = SortPair{i, distFunc(v, &vd)}
	}
	sort.Sort(distSlice) // sort default is ascending
	retM := ZeroMatrix(k, len(*v))
	for i := range retM.Data {
		retM.Data[i] = append(dataSet.Data[distSlice[i].Key], float64(distSlice[i].Key), distSlice[i].Value) // output idx, distance for observation
	}
	return retM
}
