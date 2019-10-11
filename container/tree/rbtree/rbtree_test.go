package rbtree

import (
	"fmt"
	"golina/container"
	"testing"
)

func TestRBTree(t *testing.T) {
	root := NewNode(0)
	fmt.Println(root)
	rbTree := NewRBTree(root, container.IntComparator)
	fmt.Println(rbTree)
	/*
		a := []int{13, 21, 5, 2, 1, 86, 33, 79, 25, 18}
		for _, i := range a {
			fmt.Println(rbTree)
			rbTree.Insert(i)
		}
		fmt.Println(rbTree)
	*/
}
