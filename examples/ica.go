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
	W, S, K := stats.FastICA(2, 1e-2, 200, true, stats.FuncLogcosh, dataSet)
	fmt.Println(W.Dims())
	fmt.Println(S.Dims())
	fmt.Println(K.Dims())
	_ = matrix.WriteMatrixToTxt("ica_res.txt", S)
}
