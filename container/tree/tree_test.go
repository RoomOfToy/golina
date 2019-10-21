package tree

import (
	"golina/container"
	"golina/container/tree/bheap"
	"golina/container/tree/btree"
	"golina/container/tree/rbtree"
	"math"
	"strconv"
	"testing"
)

func BenchmarkTree_InsertOne(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))

		rbTree := new(rbtree.RBTree)
		rbTree.Comparator = container.IntComparator

		minH := new(bheap.MinHeap)
		minH.Comparator = container.IntComparator
		minH.Init()

		bTree := btree.NewBTree(10, container.IntComparator)

		rn := 0
		for i := 0; i < n; i++ {
			rn = container.GenerateRandomInt()

			rbTree.Insert(rn)

			minH.Push(rn)

			bTree.Insert(&btree.Item{
				Key:   rn,
				Value: rn,
			})
		}

		num := container.GenerateRandomInt()
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

			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				rbTree.Insert(num)
			}
		})

		b.Run("Binary-Heap: size-"+strconv.Itoa(n), func(b *testing.B) {

			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				minH.Push(num)
			}
		})

		b.Run("B Tree: size-"+strconv.Itoa(n), func(b *testing.B) {

			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				bTree.Insert(&btree.Item{
					Key:   num,
					Value: num,
				})
			}
		})
	}
}

/*
BenchmarkTree_Insert/Red-Black_Tree:_size-10-8         	  500000	      4970 ns/op
BenchmarkTree_Insert/Binary-Heap:_size-10-8            	 1000000	      1608 ns/op
BenchmarkTree_Insert/B_Tree:_size-10-8                 	 1000000	      1093 ns/op
BenchmarkTree_Insert/Red-Black_Tree:_size-100-8        	   30000	     50811 ns/op
BenchmarkTree_Insert/Binary-Heap:_size-100-8           	  100000	     15996 ns/op
BenchmarkTree_Insert/B_Tree:_size-100-8                	  100000	     15300 ns/op
BenchmarkTree_Insert/Red-Black_Tree:_size-1000-8       	    3000	    604496 ns/op
BenchmarkTree_Insert/Binary-Heap:_size-1000-8          	   10000	    164097 ns/op
BenchmarkTree_Insert/B_Tree:_size-1000-8               	   10000	    210689 ns/op
*/
