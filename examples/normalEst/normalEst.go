package main

import (
	"fmt"
	"golina/mesh"
	"time"
)

func main() {
	start := time.Now()
	iv := mesh.GetInitVar()
	ne := mesh.NewNormalEst("ism_train_cat.txt", iv)
	fmt.Printf("Generate grid system time consumption: %fs\n", time.Now().Sub(start).Seconds())
	ne.Process("ism_train_cat_normal.txt")
}
