package main

import (
	"fmt"
	"golina"
	"os"
)

func main() {
	dataSet, err := golina.Load3DToMatrix("data.txt")
	if err != nil {
		panic(err)
	}
	r, c := dataSet.Dims()
	data := dataSet.GetSubMatrix(0, 0, r, c-1)
	labels := dataSet.Col(c - 1)

	fmt.Println("KMeans: ")
	clusteredData, initDistribution, finalDistribution, cnt := golina.KMeansPP(data, 2, golina.SquaredEuclideanDistance, 200)
	fmt.Printf("initDistribution: %+v\n", initDistribution)
	fmt.Printf("finalDistribution: %+v\n", finalDistribution)
	fmt.Println("Iteration: ", cnt)
	file, err := os.Create("res_k.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	kk := make(golina.Vector, len(clusteredData))
	for i := range clusteredData {
		kk[i] = float64(clusteredData[i].ClusterID)
		_, err = fmt.Fprintf(file, "%v %d\n", clusteredData[i].Observation, clusteredData[i].ClusterID)
	}
	_err := labels.Sub(&kk).MulNum(1. / float64(r)).Norm()
	fmt.Printf("Error of KMeans: %f%%\n", _err*100)
	if err != nil {
		fmt.Println(err)
	}
}
