package main

import (
	"fmt"
	"golina/matrix"
	"golina/mesh"
	"time"
)

func main() {
	start := time.Now()
	iv := mesh.GetInitVar()
	ne := mesh.NewNormalEst("ism_train_cat.txt", iv)
	fmt.Printf("generate grid system time consumption: %f\n", time.Now().Sub(start).Seconds())
	fmt.Println(ne.Grid.NumOfPoints)
	fmt.Println(ne.Grid.MinXYZ, ne.Grid.MaxXYZ)

	start = time.Now()
	ne.Voxelization()
	fmt.Printf("voxelization time consumption: %f\n", time.Now().Sub(start).Seconds())
	fmt.Println(ne.Grid.NumOfVoxel)

	start = time.Now()
	ne.FindValidVoxel()
	fmt.Printf("voxel validation time consumption: %f\n", time.Now().Sub(start).Seconds())
	fmt.Println(ne.Grid.NumOfValidVoxel)

	start = time.Now()
	ne.ComputeVoxelPlaneInfo()
	fmt.Printf("compute plane info on valid voxel time consumption: %f\n", time.Now().Sub(start).Seconds())

	start = time.Now()
	ne.AlignVoxelNormal()
	fmt.Printf("align voxel normal on valid voxel time consumption: %f\n", time.Now().Sub(start).Seconds())

	start = time.Now()
	ne.GetPointNormals()
	fmt.Printf("get point normals time consumption: %f\n", time.Now().Sub(start).Seconds())

	start = time.Now()
	err := matrix.WriteMatrixToTxt("ism_train_cat_normal.txt", ne.PointNormals.Matrix)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("write point normals to file time consumption: %f\n", time.Now().Sub(start).Seconds())
}
