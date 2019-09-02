package golina

import (
	"fmt"
	"testing"
)

func TestSVC_Predict(t *testing.T) {
	clss := SVC{
		w:          nil,
		a:          nil,
		b:          0,
		C:          0.1,
		Tolerance:  1e-10,
		Kernel:     RBFKernel,
		KernelArgs: 3.5,
		bOffset:    0,
		sv:         nil,
	}
	dataSet, err := Load3DToMatrix("data.txt")
	if err != nil {
		panic(err)
	}
	// println(dataSet.String())
	r, c := dataSet.Dims()
	data := dataSet.GetSubMatrix(0, 0, r, c-1)
	labels := dataSet.Col(c - 1)
	clss.Fit(data)
	fmt.Printf("%+v\n", clss)
	tt := clss.Predict(data)
	// fmt.Println(tt)
	_err := labels.Sub(tt).SquareSum() / float64(r)
	fmt.Printf("Error: %f\n", _err)

	data = data.MulNum(25).AddNum(50)

	err = WriteMatrixToTxt("res_.txt", data)
	if err != nil {
		fmt.Println(err)
	}

	if !FloatEqual(_err, 0.5) {
		t.Fail()
	}
}
