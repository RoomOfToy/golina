package mesh

import (
	"golina/matrix"
	"log"
)

type Point struct {
	*matrix.Vector
}

func NewPoint(x, y, z float64) Point {
	return Point{&(matrix.Vector{x, y, z})}
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
