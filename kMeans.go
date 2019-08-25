package golina

import (
	"math/rand"
	"time"
)

// https://en.wikipedia.org/wiki/K-means_clustering
type ObservationWithClusterID struct {
	clusterID   int
	observation Vector
}

type ClusteredObservationSet []ObservationWithClusterID

type DistFunc func(v1, v2 *Vector) float64

func nearestMean(means *Matrix, observation ObservationWithClusterID, distFunc DistFunc) (idx int, minDistance float64) {
	d := 0.
	minDistance = distFunc(&(observation.observation), &(means._array[0]))
	for i := 1; i < len(means._array); i++ {
		d = distFunc(&(observation.observation), &(means._array[i]))
		if d < minDistance {
			minDistance = d
			idx = i
		}
	}
	return
}

func observationInit(rawData *Matrix) ClusteredObservationSet {
	observationSet := make([]ObservationWithClusterID, len(rawData._array))
	for i, d := range rawData._array {
		observationSet[i].observation = d
	}
	return observationSet
}

func RandomMeans(dataSet *Matrix, k int) *Matrix {
	means := ZeroMatrix(k, len(dataSet._array))
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < k; i++ {
		means._array[i] = dataSet._array[rand.Intn(len(dataSet._array))]
	}
	return means
}

func KMeans(dataSet *Matrix, means *Matrix, distFunc DistFunc, iterLimit int) (ClusteredObservationSet, []int, []int, int) {
	data := observationInit(dataSet)
	cnt := 0
	// Assign Step
	for i, d := range data {
		id, _ := nearestMean(means, d, distFunc)
		data[i].clusterID = id
	}
	initDistribution := make([]int, len(means._array))
	for i := range data {
		initDistribution[data[i].clusterID]++
	}
	// Update Step
	setVolume := make([]int, len(means._array))
	change := false
	for {
		// generate new means
		means = Empty(means)
		setVolume = make([]int, len(means._array))
		for _, d := range data {
			means._array[d.clusterID] = *(means._array[d.clusterID].Add(&(d.observation)))
			setVolume[d.clusterID]++
		}
		for i := range means._array {
			means._array[i] = *(means._array[i].MulNum(1. / float64(setVolume[i])))
		}
		change = false
		for i, d := range data {
			if id, _ := nearestMean(means, d, distFunc); id != d.clusterID {
				change = true
				data[i].clusterID = id
			}
		}
		cnt++
		if change == false || cnt > iterLimit {
			finalDistribution := make([]int, len(means._array))
			for i := range data {
				finalDistribution[data[i].clusterID]++
			}
			return data, initDistribution, finalDistribution, cnt
		}
	}
}
