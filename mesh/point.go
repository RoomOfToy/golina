package mesh

import (
	"fmt"
	"golina/matrix"
	"log"
)

type Point struct {
	*matrix.Vector
}

func NewPoint(x, y, z float64) Point {
	return Point{&(matrix.Vector{x, y, z})}
}

func PEqual(pt1, pt2 *Point) bool {
	return pt1.At(0) == pt2.At(0) && pt1.At(1) == pt2.At(1) && pt1.At(2) == pt2.At(2)
}

func (point *Point) String() string {
	return fmt.Sprintf("%s", point.Vector)
}

type Points struct {
	*matrix.Matrix
}

func NewPoints(path string) Points {
	mat, err := matrix.Load3DToMatrix(path)
	if err != nil {
		log.Fatal(err)
	}
	return Points{mat}
}

func (points *Points) MaxXYZ() Point {
	_, x := points.Col(0).Max()
	_, y := points.Col(1).Max()
	_, z := points.Col(2).Max()
	return NewPoint(x, y, z)
}

func (points *Points) MinXYZ() Point {
	_, x := points.Col(0).Min()
	_, y := points.Col(1).Min()
	_, z := points.Col(2).Min()
	return NewPoint(x, y, z)
}

func (points *Points) PointsNum() int {
	r, _ := points.Dims()
	return r
}

func (points *Points) String() string {
	return fmt.Sprintf("%s", points.Matrix)
}
