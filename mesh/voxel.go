package mesh

import (
	"golina/container"
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
	NeighborIDs    container.IntSet
	IsValid        bool
	IsGood         bool
}

func NewVoxel(id int, points Points) *Voxel {
	return &Voxel{
		Id:          id,
		Points:      points,
		NumOfPoints: points.PointsNum(),
	}
}

// in place
func (v *Voxel) AddPoints(points Points) {
	v.Points = Points{v.Points.Concatenate(points.Matrix, 0)}
	v.NumOfPoints = v.Points.PointsNum()
}

// compute plane on valid voxels, call on demand
func (v *Voxel) ComputePlane() {
	cov := v.Points.CovMatrix()
	eigVec, eigVal := matrix.EigenDecompose(cov)
	v.PlaneCenter = Point{v.Points.Mean(0)}
	v.PlaneNormal = Point{eigVec.Col(0)}
	v.PlaneMSE = eigVal.At(0, 0)
	v.PlaneCurvature = eigVal.At(0, 0) / eigVal.Sum(-1).At(0)
}

func (v *Voxel) ComputeNormalSimilarity(v1 *Voxel) float64 {
	return math.Abs(v.PlaneNormal.Dot(v1.PlaneNormal.Vector))
}
