package golina

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

func TestKMeans(t *testing.T) {
	dataSet := GenerateRandomMatrix(1000, 3)
	k := 10
	means := RandomMeans(dataSet, k)
	fmt.Println(means)
	clusteredData, initDistribution, finalDistribution, cnt := KMeans(dataSet, means, SquaredEuclideanDistance, 200)
	fmt.Printf("initDistribution: %+v\n", initDistribution)
	fmt.Printf("finalDistribution: %+v\n", finalDistribution)
	fmt.Println(cnt)
	fmt.Println(len(clusteredData))
}

func BenchmarkKMeans(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		b.Run("size-"+strconv.Itoa(n)+"x3", func(b *testing.B) {
			dataSet := GenerateRandomMatrix(n, 3)
			means := RandomMeans(dataSet, int(k*10))
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				KMeans(dataSet, means, SquaredEuclideanDistance, int(k*10*2))
			}
		})
	}
}
