package rbtree

import (
	"fmt"
	"golina/container"
	"golina/container/tree"
	"math"
	"strconv"
	"testing"
)

func TestRBTree(t *testing.T) {
	root := NewNode(0)
	rbTree := NewRBTree(root, container.IntComparator)

	var _ tree.Tree = (*RBTree)(nil)

	if rbTree.Empty() {
		t.Fail()
	}

	a := []int{1, 2, 3, -5, -3, 5, -8, 8, 4, 6}
	for _, i := range a {
		rbTree.Insert(i)
	}

	if rbTree.Root.value.(int) != 2 {
		t.Fail()
	}
	if rbTree.MinValue().(int) != -8 {
		t.Fail()
	}

	fmt.Println(rbTree)

	ok := rbTree.Delete(0)
	if !ok {
		t.Fail()
	}

	aSorted := []int{-8, -5, -3, 1, 2, 3, 4, 5, 6, 8}
	for i, v := range rbTree.Values() {
		if v.(int) != aSorted[i] {
			t.Fail()
		}
	}
	if rbTree.Size() != len(aSorted) {
		t.Fail()
	}

	b := []int{-3, 1, 2, 3}
	for _, v := range b {
		if ok := rbTree.Delete(v); !ok {
			t.Fail()
		}
	}

	aSorted = []int{-8, -5, 4, 5, 6, 8}
	for i, v := range rbTree.Values() {
		if v.(int) != aSorted[i] {
			t.Fail()
		}
	}

	rbTree.Clear()
	if !rbTree.Empty() || rbTree.Size() != 0 || rbTree.Values() != nil {
		t.Fail()
	}
}

func BenchmarkRBTree_Insert(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))

		rbTree := new(RBTree)
		rbTree.Comparator = container.IntComparator

		rn := 0
		for i := 0; i < n; i++ {
			rn = container.GenerateRandomInt()
			rbTree.Insert(rn)
		}

		num := container.GenerateRandomInt()
		b.ResetTimer()

		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			for i := 1; i < b.N; i++ {
				rbTree.Insert(num)
			}
		})
	}
}

/*
BenchmarkRBTree_Insert/size-10-8         	 3000000	       450 ns/op
BenchmarkRBTree_Insert/size-100-8        	 3000000	       452 ns/op
BenchmarkRBTree_Insert/size-1000-8       	 3000000	       539 ns/op
 */
