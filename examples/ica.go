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
	W, S, K, X := stats.FastICA(2, 1e-2, 200, true, stats.FuncLogcosh, dataSet)
	fmt.Println(W.Dims())
	fmt.Println(S.Dims())
	fmt.Println(K.Dims())
	fmt.Println(X.Dims())
	fmt.Println(W)
	fmt.Println(K)
	M, _ := dataSet.Dims()
	_ = matrix.WriteMatrixToTxt("ica_res.txt", W.Mul(K).Mul(dataSet.Sub(dataSet.Mean(0).Tile(0, M)).T()).T())
}
