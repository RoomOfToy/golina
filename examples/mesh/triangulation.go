package main

import (
	"fmt"
	"golina/mesh"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	// init
	start := time.Now()
	points := mesh.NewPoints("ism_train_cat.txt")
	iv := mesh.GetInitVarTriangulation()
	dt := mesh.NewDelaunayTriangle(points, iv)
	dt.ProjectPointsToUnitSphere()
	dt.BuildInitialHull()
	fmt.Println("Init time consumption: ", time.Now().Sub(start))
	fmt.Println(dt.AuxiliaryPoints)
	fmt.Println(dt.Triangles)
	fmt.Println(len(dt.ProjectedPoints))

	// triangulation
	// goroutine stack exceeds 1000000000-byte limit
	start = time.Now()
	for _, p := range dt.ProjectedPoints {
		if !p.IsVisited {
			dt.InsertPoint(p)
		}
	}
	dt.RemoveExtraTriangles()
	fmt.Println("triangulation time consumption: ", time.Now().Sub(start))
	fmt.Println(dt.Triangles.Len())
}
