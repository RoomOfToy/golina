package matrix

import (
	"math"
)

var cos = math.Cos
var sin = math.Sin

// TransformOnRow transforms original matrix with input matrix and returns a new matrix
//	notice: affine transformations
// https://en.wikipedia.org/wiki/Affine_transformation
func TransformOnRow(t, transMat *Matrix) *Matrix {
	row, col := t.Dims()
	resMat := ZeroMatrix(row, col)
	newRowMat := ZeroMatrix(1, col)
	for i := range t.Data {
		newRowMat.Data[0] = append(t.Data[i], 1.)                      // 1 x (col+1)
		resMat.Data[i] = transMat.Mul(newRowMat.T()).T().Data[0][:col] // (col+1) x (col+1) x (col+1) x 1 -> (col+1) x 1 -> 1 x (col+1) -> 1 x col
	}
	return resMat
}

// Stretch stretches matrix along input coordinates
// if not use transform matrix, scale can be simply done by t.T().Row(i).MulNum(coordinates[i])
// 	squeeze: just x = k, y = 1 / k for stretch
func Stretch(t *Matrix, coordinates ...float64) *Matrix {
	cl := len(coordinates)
	if cl == 0 {
		return t
	}
	_, col := t.Dims()
	cl = MinInt(cl, col) // discard unused coordinates
	transMat := IdentityMatrix(col + 1)
	for i := 0; i < cl; i++ {
		transMat.Set(i, i, coordinates[i])
	}
	return TransformOnRow(t, transMat)
}

// Rotate2D rotates matrix in 2D
//	notice: rotation angle θ in counter-clockwise
func Rotate2D(t *Matrix, angle float64) *Matrix {
	_, col := t.Dims()
	if col != 2 {
		panic("this is 2D rotation function, 3D please use Rotate3D")
	}
	angle = math.Pi * angle / 180.
	rotMat := new(Matrix).Init(Data{{cos(angle), -sin(angle), 0}, {sin(angle), cos(angle), 0}, {0, 0, 1}})
	return TransformOnRow(t, rotMat)
}

// Rotate3D rotates matrix in 3D
//	notice: rotation angle θ in counter-clockwise with vector axis
func Rotate3D(t *Matrix, angle float64, axis *Vector) *Matrix {
	row, col := t.Dims()
	if col != 3 {
		panic("this is 3D rotation function, 2D please use Rotate2D")
	}
	angle = math.Pi * angle / 180.
	resMat := ZeroMatrix(row, col)
	for i := range t.Data {
		p := &t.Data[i]
		// Prot = Pcos(θ) + (n cross P)sin(θ) + n(n dot P)(1 - cos(θ)) -> θ is clockwise
		resMat.Data[i] = *(p.MulNum(cos(angle)).Add(axis.Cross(p).MulNum(sin(angle))).Add(axis.MulNum(axis.Dot(p) * (1 - cos(angle)))))
	}
	return resMat
}

// Translate translates matrix along input coordinates
func Translate(t *Matrix, coordinates ...float64) *Matrix {
	cl := len(coordinates)
	if cl == 0 {
		return t
	}
	_, col := t.Dims()
	cl = MinInt(cl, col)
	transMat := IdentityMatrix(col + 1)
	for i := 0; i < cl; i++ {
		transMat.Set(i, col, coordinates[i])
	}
	return TransformOnRow(t, transMat)
}

// Shear2D 2D shear function
// 	hx: parallel to x, hy: parallel to y
// 	x1 = x + hx * y
// 	y1 = y + hy * x
func Shear2D(t *Matrix, coordinates ...float64) *Matrix {
	cl := len(coordinates)
	if cl == 0 {
		return t
	}
	_, col := t.Dims()
	if col != 2 || cl > 2 {
		panic("this is 2D shear function, 3D please use Shear3D")
	}
	transMat := IdentityMatrix(3)
	switch cl {
	case 1:
		transMat.Set(0, 1, coordinates[0])
	case 2:
		transMat.Set(0, 1, coordinates[0])
		transMat.Set(1, 0, coordinates[1])
	}
	return TransformOnRow(t, transMat)
}

// Shear3D 3D shear function
// 	hxy, hxz, hyx, hyz, hzx, hzy
// 	x1 = x + hxy * y + hxz * z
// 	y1 = hyx * x + y + hyz * z
// 	z1 = hzx * x + hzy * y + z
func Shear3D(t *Matrix, coordinates ...float64) *Matrix {
	cl := len(coordinates)
	if cl == 0 {
		return t
	}
	_, col := t.Dims()
	if col != 3 || cl < 3 {
		panic("this is 3D shear function, 2D please use Shear2D")
	}
	transMat := IdentityMatrix(4)
	switch cl {
	case 1:
		transMat.Set(0, 1, coordinates[0])
	case 2:
		transMat.Set(0, 1, coordinates[0])
		transMat.Set(0, 2, coordinates[1])
	case 3:
		transMat.Set(0, 1, coordinates[0])
		transMat.Set(0, 2, coordinates[1])
		transMat.Set(1, 0, coordinates[2])
	case 4:
		transMat.Set(0, 1, coordinates[0])
		transMat.Set(0, 2, coordinates[1])
		transMat.Set(1, 0, coordinates[2])
		transMat.Set(1, 2, coordinates[3])
	case 5:
		transMat.Set(0, 1, coordinates[0])
		transMat.Set(0, 2, coordinates[1])
		transMat.Set(1, 0, coordinates[2])
		transMat.Set(1, 2, coordinates[3])
		transMat.Set(2, 0, coordinates[4])
	case 6:
		transMat.Set(0, 1, coordinates[0])
		transMat.Set(0, 2, coordinates[1])
		transMat.Set(1, 0, coordinates[2])
		transMat.Set(1, 2, coordinates[3])
		transMat.Set(2, 0, coordinates[4])
		transMat.Set(2, 1, coordinates[5])
	}
	return TransformOnRow(t, transMat)
}

// Kabsch calculates superimpose rotation matrix (Kabsch Algorithm) and returns two matrices,
// one represents linear transformation, and the other represents translation
//	https://en.wikipedia.org/wiki/Kabsch_algorithm
func Kabsch(P, Q *Matrix) (linear *Matrix, translation *Vector) { // X -> AX + B, A: linear transformation, B: translation
	rp, cp := P.Dims()
	rq, cq := Q.Dims()
	if cp != cq || rp != rq || cp != 3 {
		panic("dimension mismatch")
	}
	// find scale, it's distances sum ratio
	distP, distQ := 0., 0.
	for i := 0; i < rp-1; i++ {
		distP += P.Row(i + 1).Sub(P.Row(i)).Norm()
		distQ += Q.Row(i + 1).Sub(Q.Row(i)).Norm()
	}
	distP += P.Row(0).Sub(P.Row(rp - 1)).Norm()
	distQ += Q.Row(0).Sub(Q.Row(rp - 1)).Norm()
	if FloatEqual(distQ, 0.) {
		panic("invalid scale")
	}
	scale := distQ / distP
	Q = Q.MulNum(1. / scale)

	// move to centroid
	centeredP := P.Sub(P.Mean(0).Tile(0, rp))
	centeredQ := Q.Sub(Q.Mean(0).Tile(0, rq))
	// SVD
	U, _, V := SVD(centeredP.T().Mul(centeredQ))
	// Rotation
	d := V.Mul(U.T()).Det()
	if d > 0 {
		d = 1.
	} else {
		d = -1.
	}
	I := IdentityMatrix(3)
	I.Set(2, 2, d)
	rotMatrix := V.Mul(I).Mul(U.T())
	linear = rotMatrix.MulNum(scale)
	translation = Q.Mean(0).Sub(rotMatrix.MulVec(P.Mean(0))).MulNum(scale)
	return
}

// ToAffineMatrix transforms matrix into affine matrix
func ToAffineMatrix(t *Matrix) *Matrix {
	row, col := t.Dims()
	nt := ZeroMatrix(row+1, col+1)
	nt.SetSubMatrix(0, 0, t)
	nt.Set(row, col, 1)
	return nt
}
