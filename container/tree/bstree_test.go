package tree

import (
	"golina/container"
	"testing"
)

func TestBSTree_Insert(t *testing.T) {
	bst := NewBTree()
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
	bst := NewBTree()
	if !bst.Empty() {
		t.Fail()
	}
	bst.Insert(1)
	if bst.Empty() {
		t.Fail()
	}
}

func TestBSTree_Size(t *testing.T) {
	bst := NewBTree()
	bst.Insert(1)
	if bst.Size() != 1 {
		t.Fail()
	}
}

func TestBSTree_Clear(t *testing.T) {
	bst := NewBTree()
	bst.Insert(1)
	bst.Clear()
	if !bst.Empty() {
		t.Fail()
	}
}

func TestBSTree_Lookup(t *testing.T) {
	bst := NewBTree()
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
	bst := NewBTree()
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
	bst := NewBTree()
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
	bst := NewBTree()
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
	bst := NewBTree()
	bst.Comparator = container.IntComparator
	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for _, v := range a {
		bst.Insert(v)
	}
	bst.PrintPaths()
}

func TestBSTree_Mirror(t *testing.T) {
	bst := NewBTree()
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
