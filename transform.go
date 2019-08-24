package golina

import (
	"math"
)

var cos = math.Cos
var sin = math.Sin

// Affine transformations
// https://en.wikipedia.org/wiki/Affine_transformation
func TransformOnRow(t, transMat *Matrix) *Matrix {
	row, col := t.Dims()
	resMat := ZeroMatrix(row, col)
	newRowMat := ZeroMatrix(1, col)
	for i := range t._array {
		newRowMat._array[0] = append(t._array[i], 1.)                      // 1 x (col+1)
		resMat._array[i] = transMat.Mul(newRowMat.T()).T()._array[0][:col] // (col+1) x (col+1) x (col+1) x 1 -> (col+1) x 1 -> 1 x (col+1) -> 1 x col
	}
	return resMat
}

func Stretch(t *Matrix, coordinates ...float64) *Matrix {
	cl := len(coordinates)
	if cl == 0 {
		return t
	}
	_, col := t.Dims()
	cl = min(cl, col) // discard unused coordinates
	transMat := IdentityMatrix(col + 1)
	for i := 0; i < cl; i++ {
		transMat.Set(i, i, coordinates[i])
	}
	return TransformOnRow(t, transMat)
}

// if not use transform matrix, scale can be simply done by t.T().Row(i).MulNum(coordinates[i])

// squeeze: just x = k, y = 1 / k for stretch

// 2D, rotation angle θ in counter-clockwise
func Rotate2D(t *Matrix, angle float64) *Matrix {
	_, col := t.Dims()
	if col != 2 {
		panic("this is 2D rotation function, 3D please use Rotate3D")
	}
	angle = math.Pi * angle / 180.
	rotMat := new(Matrix).Init(Data{{cos(angle), -sin(angle), 0}, {sin(angle), cos(angle), 0}, {0, 0, 1}})
	return TransformOnRow(t, rotMat)
}

// 3D, rotation angle θ in counter-clockwise with vector axis
func Rotate3D(t *Matrix, angle float64, axis *Vector) *Matrix {
	row, col := t.Dims()
	if col != 3 {
		panic("this is 3D rotation function, 2D please use Rotate2D")
	}
	angle = math.Pi * angle / 180.
	resMat := ZeroMatrix(row, col)
	for i := range t._array {
		p := &t._array[i]
		// Prot = Pcos(θ) + (n cross P)sin(-θ) + n(n dot P)(1 - cos(θ)) -> θ is clockwise, so here -θ
		resMat._array[i] = *(p.MulNum(cos(angle)).Add(axis.Cross(p).MulNum(sin(-angle))).Add(axis.MulNum(axis.Dot(p) * (1 - cos(angle)))))
	}
	return resMat
}

func Translate(t *Matrix, coordinates ...float64) *Matrix {
	cl := len(coordinates)
	if cl == 0 {
		return t
	}
	_, col := t.Dims()
	cl = min(cl, col)
	transMat := IdentityMatrix(col + 1)
	for i := 0; i < cl; i++ {
		transMat.Set(i, col, coordinates[i])
	}
	return TransformOnRow(t, transMat)
}

// hx: parallel to x, hy: parallel to y
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
		transMat.Set(0, 1, coordinates[0]) // x1 = x + hx * y
		transMat.Set(1, 0, coordinates[1]) // y1 = y + hy * x
	}
	return TransformOnRow(t, transMat)
}

func Shear3D(t *Matrix, coordinates ...float64) *Matrix { // hxy, hxz, hyx, hyz, hzx, hzy
	cl := len(coordinates) // x1 = x + hxy * y + hxz * z
	if cl == 0 {           // y1 = hyx * x + y + hyz * z
		return t // z1 = hzx * x + hzy * y + z
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
