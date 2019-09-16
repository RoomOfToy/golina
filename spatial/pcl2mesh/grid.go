package pcl2mesh

type Grid struct {
	NumOfValidVoxel    int
	Cols, Rows, Depths int
	MinXYZ, MaxXYZ     Point
	Points             Points
}

func NewGrid() *Grid {
	return &Grid{
		NumOfValidVoxel: 0,
		Cols:            0,
		Rows:            0,
		Depths:          0,
		MinXYZ:          Point{},
		MaxXYZ:          Point{},
		Points:          Points{},
	}
}
