package matrix

import (
	"fmt"
	"math"
	"sort"
	"sync"
)

type Vector []float64 // 1D array

// Vector
// 	get vector element at index n
func (v *Vector) At(n int) float64 {
	l := len(*v)
	if AbsInt(n) > l {
		panic("index out of range")
	}
	if n < 0 {
		n = l + n
	}
	return (*v)[n]
}

// add two vectors
func (v *Vector) Add(v1 *Vector) *Vector {
	if len(*v) != len(*v1) {
		panic("add requires equal-length vectors")
	}
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] + (*v1)[i]
	}
	return &res
}

// vector add number
func (v *Vector) AddNum(n interface{}) *Vector {
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] + GetFloat64(n)
	}
	return &res
}

// vector subtract vector
func (v *Vector) Sub(v1 *Vector) *Vector {
	if len(*v) != len(*v1) {
		panic("sub requires equal-length vectors")
	}
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] - (*v1)[i]
	}
	return &res
}

// vector subtract number
func (v *Vector) SubNum(n interface{}) *Vector {
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] - GetFloat64(n)
	}
	return &res
}

// vector multiply number
func (v *Vector) MulNum(n interface{}) *Vector {
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] * GetFloat64(n)
	}
	return &res
}

// vector dot production
func (v *Vector) Dot(v1 *Vector) float64 {
	if len(*v) != len(*v1) {
		panic("dot product requires equal-length vectors")
	}
	res := 0.
	for i := range *v {
		res += (*v)[i] * (*v1)[i]
	}
	return res
}

// vector outer product: v1, v2 -> matrix
// https://en.wikipedia.org/wiki/Outer_product
func (v *Vector) OuterProduct(v1 *Vector) *Matrix {
	row, col := len(*v), len(*v1)
	res := ZeroMatrix(row, col)
	for i := range res.Data {
		for j := range res.Data[i] {
			res.Set(i, j, (*v)[i]*(*v)[j])
		}
	}
	return res
}

// vector cross product, 3D only
func (v *Vector) Cross(v1 *Vector) *Vector {
	if len(*v) != len(*v1) || len(*v) != 3 {
		panic("cross product requires 3d vectors in 3d space!")
	}
	return &Vector{(*v)[1]*(*v1)[2] - (*v)[2]*(*v1)[1], (*v)[2]*(*v1)[0] - (*v)[0]*(*v1)[2], (*v)[0]*(*v1)[1] - (*v)[1]*(*v1)[0]}
}

// vector elements square sum
func (v *Vector) SquareSum() float64 {
	// dot is almost 50% faster than pow by benchmark
	return v.Dot(v)
}

// vector norm
func (v *Vector) Norm() float64 {
	return math.Sqrt(v.SquareSum())
}

// normalize vector
func (v *Vector) Normalize() *Vector {
	n := v.Norm()
	if n == 0 {
		panic("invalid input vector with norm equal to 0")
	}
	res := make(Vector, len(*v))
	for i := range *v {
		res[i] = (*v)[i] / n
	}
	return &res
}

// vector to matrix, row-wise
func (v *Vector) ToMatrix(rows, cols int) *Matrix {
	if len(*v) != rows*cols {
		panic(fmt.Sprintf("invalid target matrix dimensions (%d x %d) with vector length %d\n", rows, cols, len(*v)))
	}
	nt := ZeroMatrix(rows, cols)
	for r := range nt.Data {
		for c := range nt.Data[r] {
			nt.Data[r][c] = (*v)[r*cols+c]
		}
	}
	return nt
}

// sum of vector's elements
func (v *Vector) Sum() float64 {
	s := 0.
	for _, e := range *v {
		s += e
	}
	return s
}

// sum of vector elements' absolute value
func (v *Vector) AbsSum() float64 {
	s := 0.
	for _, e := range *v {
		s += math.Abs(e)
	}
	return s
}

// mean value of vector
func (v *Vector) Mean() float64 {
	return v.Sum() / float64(len(*v))
}

func (v *Vector) Variance() float64 {
	return v.SubNum(v.Mean()).SquareSum() / float64(v.Length())
}

func (v *Vector) StandardDeviation() float64 {
	return math.Sqrt(v.Variance())
}

// tile vector alone certain dimension into matrix, 0 -> vector as row, 1 -> vector as column
func (v *Vector) Tile(dim, n int) *Matrix {
	switch dim {
	case 0:
		d := make(Data, n)
		for i := range d {
			d[i] = *v
		}
		return new(Matrix).Init(d)
	case 1:
		d := make(Data, n)
		for i := range d {
			d[i] = *v
		}
		return new(Matrix).Init(d).T()
	default:
		panic("invalid tile dimension")
	}
}

// length of vector
func (v *Vector) Length() int {
	return len(*v)
}

// max element of vector
func (v *Vector) Max() (int, float64) {
	sortSlice := make(SortPairSlice, v.Length())
	for i, j := range *v {
		sortSlice[i] = SortPair{
			Key:   i,
			Value: j,
		}
	}
	sort.Sort(sortSlice)
	return sortSlice[v.Length()-1].Key, sortSlice[v.Length()-1].Value
}

// min element of vector
func (v *Vector) Min() (int, float64) {
	sortSlice := make(SortPairSlice, v.Length())
	for i, j := range *v {
		sortSlice[i] = SortPair{
			Key:   i,
			Value: j,
		}
	}
	sort.Sort(sortSlice)
	return sortSlice[0].Key, sortSlice[0].Value
}

// Sorted pairs of vector
func (v *Vector) SortedToSortPairSlice() SortPairSlice {
	sortSlice := make(SortPairSlice, v.Length())
	for i, j := range *v {
		sortSlice[i] = SortPair{
			Key:   i,
			Value: j,
		}
	}
	sort.Sort(sortSlice)
	return sortSlice
}

// Sorted return new vector in ascending order
func (v *Vector) SortedAscending() *Vector {
	nv := make(Vector, v.Length())
	copy(nv, *v)
	sort.Float64s(nv)
	return &nv
}

// Sorted return new vector in descending order
func (v *Vector) SortedDescending() *Vector {
	nv := make(Vector, v.Length())
	copy(nv, *v)
	sort.Sort(sort.Reverse(sort.Float64Slice(nv)))
	return &nv
}

// Reversed return new vector with reverse order
func (v *Vector) Reversed() *Vector {
	nv := make(Vector, v.Length())
	copy(nv, *v)
	for i, j := 0, v.Length()-1; i < j; i, j = i+1, j-1 {
		nv[i], nv[j] = nv[j], nv[i]
	}
	return &nv
}

// Find unique elements in a vector
func (v *Vector) Unique() *Vector {
	uSet := map[float64]bool{}
	for _, val := range *v {
		uSet[val] = true
	}
	uv := make(Vector, len(uSet))
	i := 0
	for k := range uSet {
		uv[i] = k
		i++
	}
	return &uv
}

func (v *Vector) UniqueWithCount() map[float64]int {
	uSet := map[float64]int{}
	for _, val := range *v {
		if _, exist := uSet[val]; !exist {
			uSet[val] = 1
		} else {
			uSet[val]++
		}
	}
	return uSet
}

// cross covariance matrix
// 	https://en.wikipedia.org/wiki/Cross-covariance
//	https://www.quora.com/What-is-the-difference-between-cross-correlation-and-cross-covariance
//	Element to Element
//	TODO: am i right on cross covariance and cross correlation???
func CrossCov(u, v *Vector) *Matrix {
	m, n := u.Length(), v.Length()
	mat := ZeroMatrix(m, n)
	um, vm := u.Mean(), v.Mean()
	f := float64(m * n)
	for i := range mat.Data {
		for j := range mat.Data[i] {
			mat.Data[i][j] = (u.At(i) - um) * (v.At(j) - vm) / f
		}
	}
	return mat
}

// Cross Correlation
//	https://en.wikipedia.org/wiki/Cross-correlation
//	https://www.quora.com/What-is-the-difference-between-cross-correlation-and-cross-covariance
//	Element to Element
func CrossCorr(u, v *Vector) *Matrix {
	m, n := u.Length(), v.Length()
	mat := ZeroMatrix(m, n)
	f := float64(m * n)
	for i := range mat.Data {
		for j := range mat.Data[i] {
			mat.Data[i][j] = (u.At(i) * v.At(j)) / f
		}
	}
	return mat
}

func mul(u, v *Vector, k int) (res float64) {
	n := MinInt(k+1, len(*u))
	j := MinInt(k, len(*v)-1)

	for i := k - j; i < n; i, j = i+1, j-1 {
		res += (*u)[i] * (*v)[j]
	}
	return res
}

// Vector convolve
//	Convolve computes w = u * v, where w[k] = Σ u[i]*v[j], i + j = k.
//	Precondition: len(u) > 0, len(v) > 0.
func Convolve(u, v *Vector) *Vector {
	n := len(*u) + len(*v) - 1
	w := make(Vector, n)

	// Divide w into work units that take ~100μs-1ms to compute.
	size := MaxInt(1, 100000/n)

	var wg sync.WaitGroup
	for i, j := 0, size; i < n; i, j = j, j+size {
		if j > n {
			j = n
		}

		// The goroutines share memory, but only for reading.
		wg.Add(1)

		go func(i, j int) {
			for k := i; k < j; k++ {
				w[k] = mul(u, v, k)
			}
			wg.Done()
		}(i, j)
	}

	wg.Wait()

	return &w
}

func (v *Vector) Concatenate(v1 *Vector) *Vector {
	nv := append(*v, *v1...)
	return &nv
}

// pretty-print for vector
func (v *Vector) String() string {
	if v == nil {
		return "{nil}"
	}
	outString := "{"
	maxLen := 0
	vLen := len(*v)
	for i := 0; i < vLen; i++ {
		entry := (*v)[i]
		entryString := fmt.Sprintf("%f", entry)
		maxLen = MaxInt(maxLen, len(entryString))
	}
	for i := 0; i < vLen; i++ {
		entry := (*v)[i]
		entryString := fmt.Sprintf("%f", entry)
		for len(entryString) < maxLen {
			entryString = " " + entryString
		}
		outString += entryString
		if i != vLen-1 {
			outString += ", "
		}
	}
	outString += "}\n"
	return outString
}

func ARRange(start, step, stop int) *Vector {
	l := (stop - start) / step
	v := make(Vector, l)
	for i := start; i < stop; i += step {
		v[i] = float64(i)
	}
	return &v
}

func (v *Vector) MapFloat(f func(float64) float64) *Vector {
	nv := make(Vector, v.Length())
	for i := range *v {
		nv[i] = f((*v)[i])
	}
	return &nv
}

func (v *Vector) MapInt(f func(float64) int) *[]int {
	nv := make([]int, v.Length())
	for i := range *v {
		nv[i] = f((*v)[i])
	}
	return &nv
}

func (v *Vector) Angle(v1 *Vector) float64 {
	if v.Norm() == 0 || v1.Norm() == 0 {
		panic("zero division error: vector normal can NOT be zero")
	}
	return math.Acos(v.Dot(v1) / (v.Norm() * v1.Norm()))
}
