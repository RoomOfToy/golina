package tree

import (
	"golina/container"
	"golina/container/tree/bheap"
	"golina/container/tree/rbtree"
	"math"
	"strconv"
	"testing"
)

func BenchmarkTree_Insert(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))
		nums := make([]int, n)
		for i := range nums {
			nums[i] = container.GenerateRandomInt()
		}
		b.ResetTimer()

		/*
			b.Run("BST: size-"+strconv.Itoa(n), func(b *testing.B) {
				rbTree := new(bstree.BSTree)
				rbTree.Comparator = container.IntComparator
				b.ResetTimer()
				for i := 1; i < b.N; i++ {
					for _, num := range nums {
						rbTree.Insert(num)
					}
				}
			})
		*/

		b.Run("Red-Black Tree: size-"+strconv.Itoa(n), func(b *testing.B) {
			rbTree := new(rbtree.RBTree)
			rbTree.Comparator = container.IntComparator
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				for _, num := range nums {
					rbTree.Insert(num)
				}
			}
		})

		b.Run("Binary-Heap: size-"+strconv.Itoa(n), func(b *testing.B) {
			minH := new(bheap.MinHeap)
			minH.Comparator = container.IntComparator
			minH.Init()
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				for _, num := range nums {
					minH.Push(num)
				}
			}
		})
	}
}
