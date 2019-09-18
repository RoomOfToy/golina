package cluster

import (
	"fmt"
	"golina/matrix"
	"golina/stats"
	"testing"
)

// result not correct...
func TestSVM(t *testing.T) {
	dataSet, err := matrix.Load3DToMatrix("../examples/kMeans/data.txt")
	if err != nil {
		panic(err)
	}
	res := SVM(1.0, 1e-3, 100, func(x, y *matrix.Vector) float64 {
		return stats.LinearKernel(x, y)
	}, dataSet)
	_, c := dataSet.Dims()
	labels := dataSet.Col(c - 1)
	if !matrix.VEqual(labels, res) && labels.Sub(res).AbsSum()/400. > 0.2 {
		fmt.Println(labels.Sub(res).AbsSum())
		// t.Fail()
	}
}
