package pcl2mesh

import "golina/matrix"

type NormalEst struct {
	InitVar *initVar
	Grid    *Grid
	// encode like lattice, current voxel as center, from left-bottom to right-top
	RelativeNeighborPosition [][]float64
	PointNormals             Points
	PointsWithVoxelID        PointsWithVoxelID
	Voxels                   map[int]Voxel // map[voxel_id]Voxel
	ValidVoxelIDs            []int
}

func NewNormalEst(path string, iv *initVar) *NormalEst {
	ne := new(NormalEst)
	ne.InitVar = iv
	ne.Grid = NewGrid(path, iv.voxelSize[0], iv.voxelSize[1], iv.voxelSize[2])
	ne.RelativeNeighborPosition = [][]float64{
		{-1, -1, -1}, {0, -1, -1}, {1, -1, -1}, {-1, 0, -1}, {0, 0, -1}, {1, 0, -1}, {-1, 1, -1}, {0, 1, -1}, {1, 1, -1},
		{-1, -1, 0}, {0, -1, 0}, {1, -1, 0}, {-1, 0, 0}, {0, 0, 0}, {1, 0, 0}, {-1, 1, 0}, {0, 1, 0}, {1, 1, 0},
		{-1, -1, 1}, {0, -1, 1}, {1, -1, 1}, {-1, 0, 1}, {0, 0, 1}, {1, 0, 1}, {-1, 1, 1}, {0, 1, 1}, {1, 1, 1},
	}
	ne.PointsWithVoxelID = ne.Grid.ConvertXYZToVoxelID()
	return ne
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
		for i := range idx_array {
			points.Data[i] = ne.PointsWithVoxelID.Points.Data[i]
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
			ne.ValidVoxelIDs = append(ne.ValidVoxelIDs, id)
		}
	}
}

// ComputePlane for all valid voxels
func (ne *NormalEst) ComputeVoxelPlaneInfo() {
	var v Voxel
	for _, id := range ne.ValidVoxelIDs {
		v = ne.Voxels[id]
		v.ComputePlane()
	}
}

// Find good voxels according to conditions in initVar
func (ne *NormalEst) FindGoodVoxel() {}

// Compute and assign point normals
func (ne *NormalEst) AssignPointNormal() {}

// Call methods above in right order
func (ne *NormalEst) Process() {}
