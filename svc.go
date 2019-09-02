package golina

// Reference: https://github.com/josiahw/SimpleSVClustering

import (
	"math"
)

// Polynomial kernel
// 	https://en.wikipedia.org/wiki/Polynomial_kernel
func PolyKernel(a, b *Vector, power float64) float64 {
	return math.Pow(a.Dot(b), power)
}

// Radial basis function
// 	https://en.wikipedia.org/wiki/Radial_basis_function_kernel
func RBFKernel(a, b *Vector, gamma float64) float64 {
	return math.Exp(-gamma * a.Sub(b).Norm())
}

type SVC struct {
	w, a      *Vector
	b         float64 // mean of support vectors (bias value)
	C         float64 // SVC cost
	tolerance float64 // gradient descent solution accuracy
	kernel    func(a, v *Vector, x float64) float64
	args      float64 // power or gamma
	bOffset   float64
	sv        *Matrix
}

func (svc *SVC) CheckConnected(a, b *Vector, segments int) bool {
	// segments = 20
	// fmt.Println(a, b)
	for i := 0; i < segments; i++ {
		data := make(Data, 1)
		data[0] = *(a.MulNum(i).Add(b.MulNum(1 - i)))
		if svc.CalRadius(new(Matrix).Init(data)).At(0) > svc.b {
			return false
		}
	}
	return true
}

func (svc *SVC) CalRadius(dataSet *Matrix) *Vector {
	// calculate radius
	r, _ := dataSet.Dims()
	r1, _ := svc.sv.Dims()
	clss := make(Vector, r)
	for i := 0; i < r; i++ {
		clss[i] += svc.kernel(dataSet.Row(i), dataSet.Row(i), svc.args)
		for j := 0; j < r1; j++ {
			clss[i] -= 2 * svc.a.At(j) * svc.kernel(svc.sv.Row(j), dataSet.Row(i), svc.args)
		}
	}
	clss = *(clss.AddNum(svc.bOffset))
	for i := range clss {
		clss[i] = math.Sqrt(clss[i])
	}
	return &clss
}

func (svc *SVC) Fit(dataSet *Matrix) {
	// construct Q matrix for solving
	r, c := dataSet.Dims()
	Q := ZeroMatrix(r, r)
	for i := 0; i < r; i++ {
		for j := i; j < r; j++ {
			Qval := 1.
			Qval *= svc.kernel(dataSet.Row(i), dataSet.Row(j), svc.args)
			Q._array[i][j], Q._array[j][i] = Qval, Qval
		}
	}

	// solve for a and w simultaneously by coordinate descent
	w := make(Vector, c)
	a := make(Vector, r)
	svc.w, svc.a = &w, &a
	delta := 10000000000.0
	for delta > svc.tolerance {
		delta = 0
		for i := 0; i < r; i++ {
			g := Q.Row(i).Dot(svc.a) - Q.At(i, i)
			adelta := svc.a.At(i) - math.Min(math.Max(svc.a.At(i)-g/Q.At(i, i), 0.), svc.C)
			svc.w = svc.w.Add(dataSet.Row(i).MulNum(adelta))
			delta += math.Abs(adelta)
			(*(svc.a))[i] -= adelta
		}
	}

	// get support vector
	na := make(Vector, 0, svc.a.Length())
	idx := make([]int, 0, svc.a.Length())
	for i := range *(svc.a) {
		if svc.a.At(i) >= svc.C/100. {
			na = append(na, svc.a.At(i))
			idx = append(idx, i)
		}
	}
	svc.a = &na

	// fmt.Printf("%+v\n", idx)

	Qshrunk := ZeroMatrix(len(idx), len(idx))
	for i := range idx {
		for j := range idx {
			Qshrunk.Set(i, j, Q.At(i, j))
		}
	}
	// qsr, qsc := Qshrunk.Dims()
	// fmt.Println(qsr, qsc)

	svc.sv = new(Matrix)

	svc.sv._array = make(Data, svc.a.Length())
	for i := range *(svc.a) {
		svc.sv._array[i] = dataSet._array[int(svc.a.At(i))]
	}

	// calculate contribution of all SVs
	for i := range *(svc.a) {
		for j := range *(svc.a) {
			Qshrunk.Set(i, j, Qshrunk.At(i, j)*svc.a.At(i)*svc.a.At(j))
		}
	}

	svc.bOffset = Qshrunk.Sum(0).Sum()

	// select support vectors and solve for b to get the final classifier
	svc.b = svc.CalRadius(svc.sv).Mean()
}

func (svc *SVC) Classify(dataSet *Matrix) *Vector {
	// assign class labels to each vector based on connected graph components

	// build connected clusters
	r, _ := dataSet.Dims()
	unvisited := make([]int, r)
	for i := range unvisited {
		unvisited[i] = i
	}
	var clusters [][]int
	for {
		i := len(unvisited)
		if i <= 0 {
			break
		}
		// create a new cluster with the first unvisited node
		c := []int{unvisited[0]}
		unvisited = unvisited[1:]
		for j := 0; j < len(c) && len(unvisited) > 0; {
			// for all nodes in the cluster, add all connected unvisited nodes and remove them from the unvisited list
			var unvisitedNew []int
			for _, k := range unvisited {
				if svc.CheckConnected(dataSet.Row(c[j]), dataSet.Row(k), 5) {
					c = append(c, k)
				} else {
					unvisitedNew = append(unvisitedNew, k)
				}
				unvisited = unvisitedNew
				j += 1
			}
		}
		clusters = append(clusters, c)
	}

	// group components by classification
	classifications := make(Vector, r)
	classifications = *(classifications.SubNum(1))
	for i := 0; i < len(clusters); i++ {
		for _, c := range clusters[i] {
			classifications[c] = float64(i)
		}
	}
	return &classifications
}

func (svc *SVC) Predict(dataSet *Matrix) *Vector {
	// Predict classes for data X
	// 	NOTE: this should really be done with either the fitting data or a superset of the fitting data
	return svc.Classify(dataSet)
}
