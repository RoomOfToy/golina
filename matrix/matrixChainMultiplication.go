package matrix

import (
	"fmt"
	"math"
	"strconv"
)

// Matrix Chain Multiplication
//	https://en.wikipedia.org/wiki/Matrix_chain_multiplication
//	Now use the naive O(n^3) way
//	TODO: try Hu & Shing O(nlogn) way
func MatrixChainMultiplication(matrices ...*Matrix) *Matrix {
	switch len(matrices) {
	case 0:
		panic("not enough matrices for calculation")
	case 1:
		return Copy(matrices[0])
	case 2:
		return matrices[0].Mul(matrices[1])
	}
	dims := getChainDims(matrices)
	_, s := getChainOrder(dims)
	resM := new(Matrix)
	//inAResult := make([]bool, len(matrices))
	//multiplySubChainTest(s, 0, len(matrices) - 1, inAResult)
	resM = multiplySubChain(s, 0, len(matrices)-1, matrices)
	return resM
}

func getChainDims(matrices []*Matrix) []int {
	l := len(matrices)
	dims := make([]int, l+1)
	r, cc := matrices[0].Dims()
	dims[0] = r
	for i, m := range matrices[1:] {
		mr, mc := m.Dims()
		dims[i+1] = mr
		if mr != cc {
			panic("invalid matrix dimension for matrix " + strconv.Itoa(i) + " in matrices")
		}
		cc = mc
	}
	_, lc := matrices[l-1].Dims()
	dims[l] = lc
	return dims
}

func getChainOrder(dims []int) ([][]int, [][]int) {
	n := len(dims) - 1
	m := make([][]int, n)
	s := make([][]int, n)
	for i := range m {
		m[i] = make([]int, n)
		s[i] = make([]int, n)
	}
	for lenMinusOne := 1; lenMinusOne < n; lenMinusOne++ {
		for i := 0; i < n-lenMinusOne; i++ {
			j := i + lenMinusOne
			m[i][j] = math.MaxInt64
			for k := i; k < j; k++ {
				cost := m[i][k] + m[k+1][j] + dims[i]*dims[k+1]*dims[j+1]
				if cost < m[i][j] {
					m[i][j] = cost
					s[i][j] = k
				}
			}
		}
	}
	fmt.Println("the optimized cost is " + strconv.Itoa(m[0][n-1]))
	return m, s
}

// For test aim
func multiplySubChainTest(s [][]int, i, j int, inAResult []bool) *Matrix {
	if i != j {
		multiplySubChainTest(s, i, s[i][j], inAResult)
		multiplySubChainTest(s, s[i][j]+1, j, inAResult)
		istr := Ternary(inAResult[i], "_result", " ").(string)
		jstr := Ternary(inAResult[j], "_result", " ").(string)
		fmt.Println("A_" + strconv.Itoa(i) + istr + "* A_" + strconv.Itoa(j) + jstr)
		fmt.Println(inAResult)
		inAResult[i] = true
		inAResult[j] = true
		fmt.Println(i, j)
	}
	return nil
}

func multiplySubChain(s [][]int, i, j int, matrices []*Matrix) *Matrix {
	if i != j {
		multiplySubChain(s, i, s[i][j], matrices)
		multiplySubChain(s, s[i][j]+1, j, matrices)
		tmp := matrices[i].Mul(matrices[j])
		matrices[i] = tmp
		matrices[j] = tmp
	}
	return matrices[0]
}
