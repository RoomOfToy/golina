package cluster

import (
	"golina/matrix"
	"math/rand"
	"time"
)

type ObservationWithClusterID struct {
	ClusterID   int
	Observation matrix.Vector
}

type ClusteredObservationSet []ObservationWithClusterID

type DistFunc func(v1, v2 *matrix.Vector) float64

func nearestMean(means *matrix.Matrix, observation ObservationWithClusterID, distFunc DistFunc) (idx int, minDistance float64) {
	d := 0.
	minDistance = distFunc(&(observation.Observation), &(means.Data[0]))
	for i := 1; i < len(means.Data); i++ {
		d = distFunc(&(observation.Observation), &(means.Data[i]))
		if d < minDistance {
			minDistance = d
			idx = i
		}
	}
	return
}

func observationInit(rawData *matrix.Matrix) ClusteredObservationSet {
	observationSet := make([]ObservationWithClusterID, len(rawData.Data))
	for i, d := range rawData.Data {
		observationSet[i].Observation = d
	}
	return observationSet
}

func RandomMeans(dataSet *matrix.Matrix, k int) *matrix.Matrix {
	means := matrix.ZeroMatrix(k, len(dataSet.Data))
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < k; i++ {
		means.Data[i] = dataSet.Data[rand.Intn(len(dataSet.Data))]
	}
	return means
}

// K-means
// https://en.wikipedia.org/wiki/K-means_clustering
func KMeans(dataSet *matrix.Matrix, means *matrix.Matrix, distFunc DistFunc, iterLimit int) (ClusteredObservationSet, []int, []int, int) {
	data := observationInit(dataSet)
	cnt := 0
	// Assign Step
	for i, d := range data {
		id, _ := nearestMean(means, d, distFunc)
		data[i].ClusterID = id
	}
	initDistribution := make([]int, len(means.Data))
	for i := range data {
		initDistribution[data[i].ClusterID]++
	}
	// Update Step
	setVolume := make([]int, len(means.Data))
	change := false
	for {
		// generate new means
		means = matrix.Empty(means)
		setVolume = make([]int, len(means.Data))
		for _, d := range data {
			means.Data[d.ClusterID] = *(means.Data[d.ClusterID].Add(&(d.Observation)))
			setVolume[d.ClusterID]++
		}
		for i := range means.Data {
			means.Data[i] = *(means.Data[i].MulNum(1. / float64(setVolume[i])))
		}
		change = false
		for i, d := range data {
			if id, _ := nearestMean(means, d, distFunc); id != d.ClusterID {
				change = true
				data[i].ClusterID = id
			}
		}
		cnt++
		if change == false || cnt > iterLimit {
			finalDistribution := make([]int, len(means.Data))
			for i := range data {
				finalDistribution[data[i].ClusterID]++
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
func KMeansPP(dataSet *matrix.Matrix, k int, distFunc DistFunc, iterLimit int) (ClusteredObservationSet, []int, []int, int) {
	means := PPMeans(dataSet, k, distFunc)
	return KMeans(dataSet, means, distFunc, iterLimit)
}

func PPMeans(dataSet *matrix.Matrix, k int, distFunc DistFunc) *matrix.Matrix {
	dataLen := len(dataSet.Data)
	means := matrix.ZeroMatrix(k, dataLen)
	rand.Seed(time.Now().UnixNano())
	// step 1
	means.Data[0] = dataSet.Data[rand.Intn(dataLen)]
	// step 2
	dx2 := make([]float64, dataLen)
	sum := 0.
	for i := 1; i < k; i++ {
		sum = 0.
		for j, d := range dataSet.Data {
			_, minDistance := nearestMean(new(matrix.Matrix).Init(means.Data[:i]), ObservationWithClusterID{Observation: d}, distFunc)
			dx2[j] = minDistance * minDistance
			sum += dx2[j]
		}
		// step 3
		target := rand.Float64() * sum
		idx := 0
		for sum = dx2[0]; sum < target; sum += dx2[idx] {
			idx++
		}
		means.Data[i] = dataSet.Data[idx]
	}
	return means
}
