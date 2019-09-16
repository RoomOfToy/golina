package pcl2mesh

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
