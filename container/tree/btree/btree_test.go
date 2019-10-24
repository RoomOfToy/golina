package btree

import (
	"fmt"
	"golina/container"
	"golina/container/tree"
	"math"
	"strconv"
	"testing"
)

func TestBTree(t *testing.T) {
	var _ tree.Tree = (*BTree)(nil)

	bTree := NewBTree(3, container.IntComparator)
	if !bTree.Empty() || bTree.Size() != 0 {
		t.Fail()
	}

	a := []int{1, 2, 3, 4, 5, 6, 7}
	for _, v := range a {
		bTree.Insert(&Item{
			Key:   v,
			Value: fmt.Sprintf("[%s]", strconv.Itoa(v)),
		})
	}
	if bTree.Empty() || bTree.Size() != 7 {
		t.Fail()
	}

	for i, v := range a {
		value, found := bTree.Lookup(v)
		if !found || value.(string) != fmt.Sprintf("[%s]", strconv.Itoa(v)) {
			fmt.Printf("At index %d, value %d\n", i, v)
			t.Fail()
		}
	}

	for i, v := range bTree.Values() {
		if v.(string) != fmt.Sprintf("[%s]", strconv.Itoa(a[i])) {
			fmt.Printf("At index %d, value %s\n", i, v)
			t.Fail()
		}
	}

	for _, v := range a[:3] {
		bTree.Delete(v)
	}

	for i, v := range bTree.Values() {
		if v.(string) != fmt.Sprintf("[%s]", strconv.Itoa(a[3:][i])) {
			fmt.Printf("At index %d, value %s\n", i, v)
			t.Fail()
		}
	}

	if bTree.Size() != 4 {
		t.Fail()
	}

	bTree.Clear()
	if !bTree.Empty() || bTree.Size() != 0 {
		t.Fail()
	}
}

func BenchmarkBTree_Insert(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))

		bTree := NewBTree(10, container.IntComparator)

		rn := 0
		for i := 0; i < n; i++ {
			rn = container.GenerateRandomInt()
			bTree.Insert(&Item{
				Key:   rn,
				Value: rn,
			})
		}

		num := container.GenerateRandomInt()
		b.ResetTimer()

		b.Run("size-"+strconv.Itoa(n), func(b *testing.B) {
			for i := 1; i < b.N; i++ {
				bTree.Insert(&Item{
					Key:   num,
					Value: num,
				})
			}
		})
	}
}

/*
BenchmarkBTree_Insert/size-10-8         	20000000	        99.9 ns/op
BenchmarkBTree_Insert/size-100-8        	20000000	       110 ns/op
BenchmarkBTree_Insert/size-1000-8       	10000000	       146 ns/op
 */
