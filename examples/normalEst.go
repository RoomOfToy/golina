package main

import (
	"fmt"
	"golina/examples/pcl2mesh"
	"time"
)

func main() {
	start := time.Now()
	iv := pcl2mesh.GetInitVar()
	ne := pcl2mesh.NewNormalEst("ism_train_cat.txt", iv)
	fmt.Printf("generate grid system time consumption: %f\n", time.Now().Sub(start).Seconds())
	// fmt.Println(ne.PointsWithVoxelID.VoxelID)

	start = time.Now()
	ne.Voxelization()
	fmt.Printf("voxelization time consumption: %f\n", time.Now().Sub(start).Seconds())

	start = time.Now()
	ne.FindValidVoxel()
	fmt.Printf("voxel validation time consumption: %f\n", time.Now().Sub(start).Seconds())
}
