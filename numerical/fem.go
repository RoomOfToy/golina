package numerical

// reference: http://www.math.purdue.edu/~zcai/math615/matlab_fem.pdf

// Nodal Coordinate Matrix
//	dimension of this matrix is nodeNum x spacialDim
//	2D array: [][]float64
//	e.g. nodes = [[0., 0.], [2., 0.], [0., 3.], [2., 3.], [0., 6.], [2., 6.]]
type Nodes [][]float64

// Element Connectivity Matrix
//	a matrix of node numbers where each row of the matrix contains the connectivity of an element
//	elements connectivity are all ordered in a counter-clockwise fashion to keep a positive Jacobian's
//	https://en.wikipedia.org/wiki/Jacobian_matrix_and_determinant
//	e.g. elements := [[1, 2, 3], [2, 4, 3], [4, 5, 2], [6, 5, 4]]
type Elements [][]int

// Boundary Connectivity Matrix
//	boundary elements always 1D lower than spacial dimension of the problem: 3D -> 2D boundary, 2D -> 1D boundary
//	rightEdge = [[2, 4], [4, 6]]
//	nodesOnBoundary = unique(rightEdge) = [2, 4, 6]
type Boundaries [][]int

// Polynomial coefficients of func f
//	0 -> n, store in Vector
// 	f = &Vector{0, 1, 2, 3}
