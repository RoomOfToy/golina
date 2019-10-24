package tree

import (
	"golina/container"
	"golina/container/heap/bheap"
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
// Heap based on dynamic array
BenchmarkTree_InsertOne/Binary-Heap:_size-10-8            10000000       141 ns/op
BenchmarkTree_InsertOne/Binary-Heap:_size-100-8           10000000       127 ns/op
BenchmarkTree_InsertOne/Binary-Heap:_size-1000-8          20000000       118 ns/op
BenchmarkTree_InsertOne/Binary-Heap:_size-10000-8         10000000       102 ns/op
BenchmarkTree_InsertOne/Binary-Heap:_size-100000-8        20000000       113 ns/op

// Red-Black Tree based on double linked list
BenchmarkTree_InsertOne/Red-Black_Tree:_size-10-8          3000000       468 ns/op
BenchmarkTree_InsertOne/Red-Black_Tree:_size-100-8         3000000       501 ns/op
BenchmarkTree_InsertOne/Red-Black_Tree:_size-1000-8        3000000       525 ns/op
BenchmarkTree_InsertOne/Red-Black_Tree:_size-10000-8       3000000       637 ns/op
BenchmarkTree_InsertOne/Red-Black_Tree:_size-100000-8      2000000       946 ns/op

// B Tree with M = 10
BenchmarkTree_InsertOne/B_Tree:_size-10-8                 20000000        96 ns/op
BenchmarkTree_InsertOne/B_Tree:_size-100-8                10000000       125 ns/op
BenchmarkTree_InsertOne/B_Tree:_size-1000-8               10000000       152 ns/op
BenchmarkTree_InsertOne/B_Tree:_size-10000-8              10000000       185 ns/op
BenchmarkTree_InsertOne/B_Tree:_size-100000-8             10000000       223 ns/op
*/
