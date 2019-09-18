package mesh

import (
	"fmt"
	"math"
)

type Grid struct {
	VoxelSize                                  []float64
	GridSize                                   []float64
	NumOfVoxel, NumOfValidVoxel                int
	ColsOfVoxels, RowsOfVoxels, DepthsOfVoxels int // X: col, Y: row, Z: depth
	MinXYZ, MaxXYZ                             Point
	Points                                     Points
	NumOfPoints                                int
}

func NewGrid(path string, voxelX, VoxelY, voxelZ float64) *Grid {
	points := NewPoints(path)
	numOfPoints, _ := points.Dims()
	minXYZ := points.MinXYZ()
	maxXYZ := points.MaxXYZ()
	grideSize := maxXYZ.Sub(minXYZ.Vector)
	voxelSize := []float64{voxelX, VoxelY, voxelZ}
	colRowDepth := make([]int, 3)
	for i, v := range *grideSize {
		// v :=0 -> zero division panic
		if math.Mod(v, voxelSize[i]) == 0 {
			// since this conversion from float to int, it may out of int range
			// better to use int64 if voxel num too large
			colRowDepth[i] = int(math.Floor(v / voxelSize[i]))
		} else {
			colRowDepth[i] = int(math.Floor(v/voxelSize[i])) + 1
		}
	}
	return &Grid{
		VoxelSize:       voxelSize,
		GridSize:        *grideSize,
		NumOfVoxel:      colRowDepth[0] * colRowDepth[1] * colRowDepth[2],
		NumOfValidVoxel: 0,
		ColsOfVoxels:    colRowDepth[0],
		RowsOfVoxels:    colRowDepth[1],
		DepthsOfVoxels:  colRowDepth[2],
		MinXYZ:          minXYZ,
		MaxXYZ:          maxXYZ,
		Points:          points,
		NumOfPoints:     numOfPoints,
	}
}

type PointsWithVoxelID struct {
	Points  Points
	VoxelID *[]int
}

func (pwv PointsWithVoxelID) String() string {
	return fmt.Sprintf("Points: \n%+v\nVoxelID: \n%+v\n", pwv.Points, pwv.VoxelID)
}

// voxel_id = depth_idx * (cols * rows) + row_idx * cols + col_idx
// []int
func (g *Grid) ConvertXYZToVoxelID() PointsWithVoxelID {
	idx := g.Points.Sub(g.MinXYZ.Tile(0, g.NumOfPoints))
	idxX := idx.Col(0).MulNum(1. / g.VoxelSize[0]).MapFloat(math.Floor)
	idxY := idx.Col(1).MulNum(1. / g.VoxelSize[1]).MapFloat(math.Floor)
	idxZ := idx.Col(2).MulNum(1. / g.VoxelSize[2]).MapFloat(math.Floor)
	voxelID := idxZ.MulNum(g.ColsOfVoxels * g.RowsOfVoxels).Add(idxY.MulNum(g.ColsOfVoxels)).Add(idxX).MapInt(func(x float64) int {
		return int(x)
	})
	return PointsWithVoxelID{
		Points:  g.Points,
		VoxelID: voxelID,
	}
}

func (g *Grid) ConvertVoxelIDToGridCoord(voxelID int) (gridCoord []int) {
	depthIdx := voxelID / (g.ColsOfVoxels * g.RowsOfVoxels)
	rowIdx := (voxelID - depthIdx*(g.ColsOfVoxels*g.RowsOfVoxels)) / g.ColsOfVoxels
	colIdx := voxelID - depthIdx*(g.ColsOfVoxels*g.RowsOfVoxels) - rowIdx*g.ColsOfVoxels
	return []int{colIdx, rowIdx, depthIdx}
}
