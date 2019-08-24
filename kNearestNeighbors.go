package golina

import "sort"

type SortPair struct {
	key   int
	value float64
}

type SortPairSlice []SortPair

func (sps SortPairSlice) Swap(i, j int) {
	sps[i], sps[j] = sps[j], sps[i]
}

func (sps SortPairSlice) Len() int {
	return len(sps)
}

func (sps SortPairSlice) Less(i, j int) bool {
	return sps[i].value < sps[j].value
}

func KNearestNeighbors(dataSet *Matrix, v *Vector, k int, distFunc func(v1, v2 *Vector) float64) *Matrix {
	row, _ := dataSet.Dims()
	if k > row {
		k = row
	}
	distSlice := make(SortPairSlice, row)
	for i, vd := range dataSet._array {
		distSlice[i] = SortPair{i, distFunc(v, &vd)}
	}
	sort.Sort(distSlice) // sort default is ascending
	retM := ZeroMatrix(k, len(*v))
	for i := range retM._array {
		retM._array[i] = dataSet._array[distSlice[i].key]
		// retM._array[i] = append(dataSet._array[distSlice[i].key], distSlice[i].value)  // output distance for observation
	}
	return retM
}
