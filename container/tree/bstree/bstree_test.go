package bstree

import (
	"golina/container"
	"golina/container/tree"
	"math"
	"strconv"
	"testing"
)

func TestBSTree_InterfaceAssertion(t *testing.T) {
	var _ tree.Tree = (*BSTree)(nil)
}

func TestBSTree_Insert(t *testing.T) {
	bst := NewBSTree()
	bst.Comparator = container.IntComparator
	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for _, v := range a {
		bst.Insert(v)
	}
	aSorted := []int{-8, -5, -3, 1, 2, 3, 4, 5, 6, 8}
	for i, v := range bst.Values() {
		if v.(int) != aSorted[i] {
			t.Fail()
		}
	}
}

func TestBSTree_Empty(t *testing.T) {
	bst := NewBSTree()
	if !bst.Empty() {
		t.Fail()
	}
	bst.Insert(1)
	if bst.Empty() {
		t.Fail()
	}
}

func TestBSTree_Size(t *testing.T) {
	bst := NewBSTree()
	bst.Insert(1)
	if bst.Size() != 1 {
		t.Fail()
	}
}

func TestBSTree_Clear(t *testing.T) {
	bst := NewBSTree()
	bst.Insert(1)
	bst.Clear()
	if !bst.Empty() {
		t.Fail()
	}
}

func TestBSTree_Lookup(t *testing.T) {
	bst := NewBSTree()
	bst.Comparator = container.IntComparator
	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for _, v := range a {
		bst.Insert(v)
	}
	for i := range a {
		if !bst.Lookup(a[i]) {
			t.Fail()
		}
	}
	if bst.Lookup(10) {
		t.Fail()
	}
}

func TestBSTree_MaxDepth(t *testing.T) {
	bst := NewBSTree()
	if bst.MaxDepth() != 0 {
		t.Fail()
	}
	bst.Comparator = container.IntComparator
	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for _, v := range a {
		bst.Insert(v)
	}
	if bst.MaxDepth() != 6 {
		t.Fail()
	}
}

func TestBSTree_MinValue(t *testing.T) {
	bst := NewBSTree()
	bst.Comparator = container.IntComparator
	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for _, v := range a {
		bst.Insert(v)
	}
	if bst.MinValue() != -8 {
		t.Fail()
	}
}

func TestBSTree_HasPathSum(t *testing.T) {
	bst := NewBSTree()
	bst.Comparator = container.IntComparator
	bst.Diffidence = func(sum, data interface{}) interface{} {
		return sum.(int) - data.(int)
	}
	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for _, v := range a {
		bst.Insert(v)
	}
	if !bst.HasPathSum(-12) || !bst.HasPathSum(-7) || !bst.HasPathSum(15) || !bst.HasPathSum(25) {
		t.Fail()
	}
}

func TestBSTree_PrintPaths(t *testing.T) {
	bst := NewBSTree()
	bst.Comparator = container.IntComparator
	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for _, v := range a {
		bst.Insert(v)
	}
	bst.PrintPaths()
}

func TestBSTree_Mirror(t *testing.T) {
	bst := NewBSTree()
	bst.Comparator = container.IntComparator
	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for _, v := range a {
		bst.Insert(v)
	}
	bst.Mirror()
	aReverseSorted := []int{8, 6, 5, 4, 3, 2, 1, -3, -5, -8}
	for i, v := range bst.Values() {
		if v.(int) != aReverseSorted[i] {
			t.Fail()
		}
	}
}

func TestBSTree_DoubleTree(t *testing.T) {
	bst := NewBSTree()
	bst.Comparator = container.IntComparator
	a := []int{2, 1, 3}
	for _, v := range a {
		bst.Insert(v)
	}
	bst.DoubleTree()
	aSorted := []int{1, 1, 2, 2, 3, 3}
	for i, v := range bst.Values() {
		if v.(int) != aSorted[i] {
			t.Fail()
		}
	}
}

func TestBSTree_SameTree(t *testing.T) {
	bstA, bstB := NewBSTree(), NewBSTree()
	bstA.Comparator, bstB.Comparator = container.IntComparator, container.IntComparator
	a := []int{2, 1, 3}
	for _, v := range a {
		bstA.Insert(v)
		bstB.Insert(v)
	}
	if !bstA.SameTree(bstB) {
		t.Fail()
	}
	bstB.Insert(5)
	if bstA.SameTree(bstB) {
		t.Fail()
	}
}

func TestCountTrees(t *testing.T) {
	if CountTrees(4) != 14 {
		t.Fail()
	}
}

func TestBSTree_IsBST(t *testing.T) {
	bst := NewBSTree()
	bst.Comparator = container.IntComparator
	a := []int{2, 1, 3}
	for _, v := range a {
		bst.Insert(v)
	}
	if !bst.IsBST(0, 5) {
		t.Fail()
	}
	bst.Clear()
	bst.Root = NewNode(1)
	bst.Root.left = NewNode(3)
	bst.Root.right = NewNode(2)
	if bst.IsBST(0, 5) {
		t.Fail()
	}
}

func BenchmarkBSTree_Insert(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))

		bst := new(BSTree)
		bst.Comparator = container.IntComparator
		rn := 0
		for i := 0; i < n; i++ {
			rn = container.GenerateRandomInt()
			bst.Insert(rn)
		}

		num := container.GenerateRandomInt()
		b.ResetTimer()

		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			for i := 1; i < b.N; i++ {
				bst.Insert(num)
			}
		})
	}
}

/*
BenchmarkBSTree_Insert/size-10-8         	   30000	    222062 ns/op
BenchmarkBSTree_Insert/size-100-8        	   30000	    221083 ns/op
BenchmarkBSTree_Insert/size-1000-8       	   30000	    216950 ns/op
 */
