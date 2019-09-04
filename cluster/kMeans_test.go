package cluster

import (
	"fmt"
	"golina/matrix"
	"golina/spatial"
	"math"
	"strconv"
	"testing"
)

func TestKMeans(t *testing.T) {
	dataSet := matrix.GenerateRandomMatrix(1000, 3)
	k := 10
	means := RandomMeans(dataSet, k)
	fmt.Println(means)
	clusteredData, initDistribution, finalDistribution, cnt := KMeans(dataSet, means, spatial.SquaredEuclideanDistance, 200)
	fmt.Printf("initDistribution: %+v\n", initDistribution)
	fmt.Printf("finalDistribution: %+v\n", finalDistribution)
	fmt.Println(cnt)
	fmt.Println(len(clusteredData))
}

func TestKMeansPP(t *testing.T) {
	dataSet := matrix.GenerateRandomMatrix(1000, 3)
	k := 10
	means := RandomMeans(dataSet, k)
	clusteredData, initDistribution, finalDistribution, cnt := KMeans(dataSet, means, spatial.SquaredEuclideanDistance, 200)
	fmt.Println("KMeans")
	fmt.Printf("initDistribution: %+v\n", initDistribution)
	fmt.Printf("finalDistribution: %+v\n", finalDistribution)
	fmt.Println(cnt)
	fmt.Println(len(clusteredData))
	fmt.Println()
	fmt.Println("KMeansPP")
	clusteredData, initDistribution, finalDistribution, cnt = KMeansPP(dataSet, k, spatial.SquaredEuclideanDistance, 200)
	fmt.Printf("initDistribution: %+v\n", initDistribution)
	fmt.Printf("finalDistribution: %+v\n", finalDistribution)
	fmt.Println(cnt)
	fmt.Println(len(clusteredData))
}

func BenchmarkKMeans(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n)+"x3", func(b *testing.B) {
			dataSet := matrix.GenerateRandomMatrix(n, 3)
			means := RandomMeans(dataSet, int(k*10))
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				KMeans(dataSet, means, spatial.SquaredEuclideanDistance, int(k*10*2))
			}
		})
	}
}

// TODO: why `KMeansPP` not faster than `KMeans`???
func BenchmarkKMeansPP(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n)+"x3", func(b *testing.B) {
			dataSet := matrix.GenerateRandomMatrix(n, 3)
			means := PPMeans(dataSet, int(k*10), spatial.SquaredEuclideanDistance)
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				KMeans(dataSet, means, spatial.SquaredEuclideanDistance, int(k*10*2))
			}
		})
	}
}
