package main

import (
	"fmt"
	"golina"
)

func fibonacci(n int) int {
	// Series: An = An-2 + An-1
	// State: Sn = [Sn-1, Sn] = [Sn-1, Sn-2 + Sn-1] = [Sn-2 * 0 + Sn-1 * 1, Sn-2 * 1 + Sn-1 * 1]
	// Transfer function: An = An-2 + An-1 => Fn = Fn-1 * [[0, 1], [1, 1]]
	// F0 = [0, 1] => [[1, 0], [0, 1]]
	// Fn = [[1, 0], [0, 1]] * ([[0, 1], [1, 1]]).Power(n)
	f0 := new(golina.Matrix).Init(golina.Data{{1, 0}, {0, 1}})
	stateTransMatrix := new(golina.Matrix).Init(golina.Data{{0, 1}, {1, 1}})
	fn := f0.Mul(stateTransMatrix.Pow(n))
	return int(fn.At(1, 1))
}

func main() {
	for i := 0; i < 50; i++ {
		fmt.Println(fibonacci(i))
	}
}
