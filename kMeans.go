package golina

import (
	"math/rand"
	"time"
)

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

// K-means
// https://en.wikipedia.org/wiki/K-means_clustering
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

// K-means++
// https://en.wikipedia.org/wiki/K-means%2B%2B
//	Not use RandomMeans but follow the following initialization process:
//	1. Choose one center uniformly at random from among the data points.
//	2. For each data point x, compute D(x), the distance between x and the nearest center that has already been chosen.
//	3. Choose one new data point at random as a new center, using a weighted probability distribution where a point x is chosen with probability proportional to D(x)2.
//	4. Repeat Steps 2 and 3 until k centers have been chosen.
//	5. Now that the initial centers have been chosen, proceed using standard k-means clustering.
func KMeansPP(dataSet *Matrix, k int, distFunc DistFunc, iterLimit int) (ClusteredObservationSet, []int, []int, int) {
	means := PPMeans(dataSet, k, distFunc)
	return KMeans(dataSet, means, distFunc, iterLimit)
}

func PPMeans(dataSet *Matrix, k int, distFunc DistFunc) *Matrix {
	dataLen := len(dataSet._array)
	means := ZeroMatrix(k, dataLen)
	rand.Seed(time.Now().UnixNano())
	// step 1
	means._array[0] = dataSet._array[rand.Intn(dataLen)]
	// step 2
	dx2 := make([]float64, dataLen)
	sum := 0.
	for i := 1; i < k; i++ {
		sum = 0.
		for j, d := range dataSet._array {
			_, minDistance := nearestMean(new(Matrix).Init(means._array[:i]), ObservationWithClusterID{observation: d}, distFunc)
			dx2[j] = minDistance * minDistance
			sum += dx2[j]
		}
		// step 3
		target := rand.Float64() * sum
		idx := 0
		for sum = dx2[0]; sum < target; sum += dx2[idx] {
			idx++
		}
		means._array[i] = dataSet._array[idx]
	}
	return means
}
