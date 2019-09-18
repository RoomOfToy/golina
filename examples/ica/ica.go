package main

import (
	"fmt"
	"golina/matrix"
	"golina/stats"
)

func main() {
	dataSet, err := matrix.Load2DToMatrix("ica_data.txt")
	if err != nil {
		panic(err)
	}
	W, S, K, X := stats.FastICA(2, 1e-3, 200, true, stats.FuncLogcosh, dataSet)
	fmt.Println(W.Dims())
	fmt.Println(S.Dims())
	fmt.Println(K.Dims())
	fmt.Println(X.Dims())
	fmt.Println(W)
	fmt.Println(K)
	recovered := S.Mul(W.T())
	recovered = recovered.MulNum(1. / recovered.StandardDeviation(-1).At(0))
	_ = matrix.WriteMatrixToTxt("ica_res.txt", recovered)
}
