package mesh

import (
	"fmt"
	"golina/matrix"
	"golina/spatial"
	"math"
	"time"
)

type NormalEst struct {
	InitVar *initVar
	Grid    *Grid
	// encode like lattice, current voxel as center, from left-bottom to right-top
	RelativeNeighborPosition [][]int
	PointNormals             Points
	PointsWithVoxelID        PointsWithVoxelID
	Voxels                   map[int]Voxel // map[voxel_id]Voxel
	ValidVoxelIDs            []int
}

func NewNormalEst(path string, iv *initVar) *NormalEst {
	ne := new(NormalEst)
	ne.InitVar = iv
	ne.Grid = NewGrid(path, iv.voxelSize[0], iv.voxelSize[1], iv.voxelSize[2])
	// without {0, 0, 0}, which is the center and represents current voxel
	ne.RelativeNeighborPosition = [][]int{
		{-1, -1, -1}, {0, -1, -1}, {1, -1, -1}, {-1, 0, -1}, {0, 0, -1}, {1, 0, -1}, {-1, 1, -1}, {0, 1, -1}, {1, 1, -1},
		{-1, -1, 0}, {0, -1, 0}, {1, -1, 0}, {-1, 0, 0}, {1, 0, 0}, {-1, 1, 0}, {0, 1, 0}, {1, 1, 0},
		{-1, -1, 1}, {0, -1, 1}, {1, -1, 1}, {-1, 0, 1}, {0, 0, 1}, {1, 0, 1}, {-1, 1, 1}, {0, 1, 1}, {1, 1, 1},
	}
	ne.PointsWithVoxelID = ne.Grid.ConvertXYZToVoxelID()
	return ne
}

// Voxel has 26 neighbors under grid system
func (ne *NormalEst) GetNeighbors(v *Voxel) {
	gridCoord := ne.Grid.ConvertVoxelIDToGridCoord(v.Id) // X: Col, Y: Row, Z: Depth
	gridUpperEdge := []int{ne.Grid.ColsOfVoxels, ne.Grid.RowsOfVoxels, ne.Grid.DepthsOfVoxels}
	neighborVoxelID := 0
	// gridLowerEdge := []int{0, 0, 0}
	for _, move := range ne.RelativeNeighborPosition {
		for j, moveDir := range move {
			// neighbor on or out of grid border
			if gridCoord[j]+moveDir >= gridUpperEdge[j] || gridCoord[j]+moveDir < 0 {
				continue
			}
		}
		neighborVoxelID = (gridCoord[2]+move[2])*(gridUpperEdge[0]+gridUpperEdge[1]) + (gridCoord[1]+move[1])*gridUpperEdge[0] + (gridCoord[0] + move[0])
		// voxel position
		if neighborVoxelID < 0 || neighborVoxelID >= ne.Grid.NumOfVoxel {
			continue
		}
		neighborVoxel, exists := ne.Voxels[neighborVoxelID]
		// voxel existence
		if !exists {
			continue
		}
		// voxel validation
		if !neighborVoxel.IsValid {
			continue
		}
		v.NeighborIDs.Add(neighborVoxelID)
	}
}

// check whether voxel contains edge points or not based on its neighborhood
func (ne *NormalEst) IsVoxelContainsEdgePoints(v *Voxel) bool {
	// get neighbors
	if len(v.NeighborIDs) == 0 {
		ne.GetNeighbors(v)
	}
	// voxel has no neighbors -> not edge voxel
	if len(v.NeighborIDs) == 0 {
		v.IsGood = true
		return false
	}
	// check whether voxel has edge points
	badCnt := 0
	var normalAngleDiff float64
	var nv Voxel
	for nvID := range v.NeighborIDs {
		nv = ne.Voxels[nvID]
		normalAngleDiff = v.ComputeNormalSimilarity(&nv) * 180 / math.Pi
		if normalAngleDiff >= ne.InitVar.normalSimilarityThreshold {
			badCnt++
		}
		if badCnt >= ne.InitVar.badCnt {
			return true
		}
	}
	return false
}

// map[voxel_id][]point_idx
type set map[int][]int

// PointsWithVoxelID -> Voxels
func (ne *NormalEst) Voxelization() {
	voxelIDs := make(set, len(*(ne.PointsWithVoxelID.VoxelID)))
	for idx, voxel_id := range *(ne.PointsWithVoxelID.VoxelID) {
		if _, exists := voxelIDs[voxel_id]; !exists {
			voxelIDs[voxel_id] = []int{idx}
		} else {
			voxelIDs[voxel_id] = append(voxelIDs[voxel_id], idx)
		}
	}
	ne.Voxels = make(map[int]Voxel, len(voxelIDs))
	for voxel_id, idx_array := range voxelIDs {
		points := Points{matrix.ZeroMatrix(len(idx_array), 3)}
		for i, dix := range idx_array {
			points.Data[i] = ne.PointsWithVoxelID.Points.Data[dix]
		}
		ne.Voxels[voxel_id] = *(NewVoxel(voxel_id, points))
	}
}

// Generate ValidVoxelIDs according to initVar.planeMinPointNum
func (ne *NormalEst) FindValidVoxel() {
	if len(ne.Voxels) == 0 {
		ne.Voxelization()
	}
	ne.ValidVoxelIDs = make([]int, 0, len(ne.Voxels))
	for id, voxel := range ne.Voxels {
		if voxel.NumOfPoints >= ne.InitVar.planeMinPointNum {
			voxel.IsValid = true
			ne.ValidVoxelIDs = append(ne.ValidVoxelIDs, id)
		}
	}
	ne.Grid.NumOfValidVoxel = len(ne.ValidVoxelIDs)
}

// ComputePlane for all valid voxels
func (ne *NormalEst) ComputeVoxelPlaneInfo() {
	var v Voxel
	for _, id := range ne.ValidVoxelIDs {
		v = ne.Voxels[id]
		v.ComputePlane()
	}
}

// Align voxel normal with it neighbor
//	to avoid normal has 180 degree angle caused by Eigen Decomposition
//	TODO: randomly pick one now, need to see whether it is good enough
func (ne *NormalEst) AlignVoxelNormal() {
	var v, nv Voxel
	for _, id := range ne.ValidVoxelIDs {
		v = ne.Voxels[id]
		// get neighbors
		if len(v.NeighborIDs) == 0 {
			ne.GetNeighbors(&v)
		}
		// skip voxel which has no neighbors
		if len(v.NeighborIDs) == 0 {
			continue
		}
		// NeighborIDs is a set implemented by unordered map
		for nid := range v.NeighborIDs {
			nv = ne.Voxels[nid]
			// flip normal if vn dot nvn < 0
			if v.PlaneNormal.Dot(nv.PlaneNormal.Vector) < 0 {
				v.PlaneNormal = Point{v.PlaneNormal.MulNum(-1)}
			}
			// just randomly pick one
			break
		}
	}
}

// Find good voxels according to conditions in initVar
func (ne *NormalEst) FindGoodVoxel() {
	var v Voxel
	for _, id := range ne.ValidVoxelIDs {
		v = ne.Voxels[id]
		if v.PlaneMSE < ne.InitVar.planeMaxMSE {
			v.IsGood = true
		}
	}
}

// Compute and assign point normals
//	now only consider validation
//	TODO: use good instead of valid to get proper normal from voxel
func (ne *NormalEst) GetPointNormals() {
	// need ne.AlignVoxelNormal() first

	// point with its normal
	ne.PointNormals = Points{matrix.ZeroMatrix(ne.Grid.NumOfPoints, 6)}
	pointIdx := 0
	var v, nv Voxel
	var searchMat *matrix.Matrix
	for _, v = range ne.Voxels {
		if v.IsValid {
			for _, p := range v.Points.Data {
				ne.PointNormals.Data[pointIdx] = *(p.Concatenate(v.PlaneNormal.Vector))
				pointIdx++
			}
		} else {
			// not valid -> point normal compute from point to point normal calculation
			// kNN
			if len(v.NeighborIDs) == 0 {
				ne.GetNeighbors(&v)
			}
			points := v.Points
			for id := range v.NeighborIDs {
				nv = ne.Voxels[id]
				points.Concatenate(nv.Points.Matrix, 0)
			}
			// compute every point in voxel
			for _, p := range v.Points.Data {
				searchMat = spatial.KNearestNeighbors(points.Matrix, &p, ne.InitVar.searchRadius, spatial.EuclideanDistance)
				cov := searchMat.CovMatrix()
				// TODO: whether use MSE (eigVal) for further determination
				eigVec, _ := matrix.EigenDecompose(cov)
				ne.PointNormals.Data[pointIdx] = *(p.Concatenate(eigVec.Col(0)))
				pointIdx++
			}
		}
	}
}

// Call methods above in right order
func (ne *NormalEst) Process(OutPath string) {
	fmt.Println("	Number of Points", ne.Grid.NumOfPoints)
	fmt.Println("	Grid system minimum / maximum coordinates:")
	fmt.Printf("		Min: %s", ne.Grid.MinXYZ)
	fmt.Printf("		Max: %s", ne.Grid.MaxXYZ)

	start := time.Now()
	ne.Voxelization()
	fmt.Printf("Voxelization time consumption: %fs\n", time.Now().Sub(start).Seconds())
	fmt.Println("	Number of Voxels", ne.Grid.NumOfVoxel)

	start = time.Now()
	ne.FindValidVoxel()
	fmt.Printf("Voxel validation time consumption: %fs\n", time.Now().Sub(start).Seconds())
	fmt.Println("	Number of Valid Voxels", ne.Grid.NumOfValidVoxel)

	start = time.Now()
	ne.ComputeVoxelPlaneInfo()
	fmt.Printf("Compute plane info on valid voxel time consumption: %fs\n", time.Now().Sub(start).Seconds())

	start = time.Now()
	ne.AlignVoxelNormal()
	fmt.Printf("Align voxel normal on valid voxel time consumption: %fs\n", time.Now().Sub(start).Seconds())

	start = time.Now()
	ne.GetPointNormals()
	fmt.Printf("Get point normals time consumption: %fs\n", time.Now().Sub(start).Seconds())

	start = time.Now()
	err := matrix.WriteMatrixToTxt(OutPath, ne.PointNormals.Matrix)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Write point normals to %s time consumption: %fs\n", OutPath, time.Now().Sub(start).Seconds())
}
