package pcl2mesh

import (
	"golina/matrix"
	"math"
)

type Voxel struct {
	Id             int
	Points         Points
	NumOfPoints    int
	PlaneCenter    Point
	PlaneNormal    Point
	PlaneMSE       float64
	PlaneCurvature float64
}

func NewVoxel(id int, points Points) *Voxel {
	return &Voxel{
		Id:          id,
		Points:      points,
		NumOfPoints: points.ElementsNum(),
		PlaneCenter: Point{points.Mean(0)},
		PlaneNormal: Point{},
	}
}

func (v *Voxel) AddPoints(points Points) {
	v.Points = Points{v.Points.Concatenate(points.Matrix, 0)}
}

func (v *Voxel) ComputePlane() {
	cov := v.Points.CovMatrix()
	eigVec, eigVal := matrix.EigenDecompose(cov)
	v.PlaneNormal = Point{eigVec.Col(0)}
	v.PlaneMSE = eigVal.At(0, 0)
	v.PlaneCurvature = eigVal.At(0, 0) / eigVal.Sum(-1).At(0)
}

func (v *Voxel) ComputeNormalSimilarity(v1 *Voxel) float64 {
	return math.Abs(v.PlaneNormal.Dot(v1.PlaneNormal.Vector))
}
