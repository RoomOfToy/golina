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
	W, S := stats.FastICA(20000, 1e-2, 200, true, FuncLogcosh, dataSet)
	fmt.Println(W.Dims())
	fmt.Println(S.Dims())
	for i := range S.Data {
		S.Data[i] = *(S.Data[i].MulNum(stats.StandardDeviation(&(S.Data[i]))))
	}
	//fmt.Println(S)
	_ = matrix.WriteMatrixToTxt("ica_res.txt", S)
}
