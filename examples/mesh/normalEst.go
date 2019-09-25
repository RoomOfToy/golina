package main

import (
	"golina/mesh"
)

func main() {
	mesh.NormalEstProcess("ism_train_cat.txt", "ism_train_cat_normal.txt")
}
