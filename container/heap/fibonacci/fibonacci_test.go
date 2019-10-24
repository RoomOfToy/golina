package fibonacci

import (
	"golina/container"
	"golina/container/heap"
	"math"
	"strconv"
	"testing"
)

type item int

func (it item) Compare(ait heap.Item) int {
	if it > ait.(item) {
		return 1
	} else if it == ait.(item) {
		return 0
	} else {
		return -1
	}
}

func TestNewHeap(t *testing.T) {
	var _ heap.Heap = (*Heap)(nil)
	h := NewHeap()
	if h.root != nil || h.itemNum != 0 {
		t.Fail()
	}
}

func TestHeap(t *testing.T) {

	a := []item{12, 7, 25, 15, 28, 33, 41, 1}
	aSorted := []item{1, 7, 12, 15, 25, 28, 33, 41}
	b := []item{18, 35, 20, 42, 9, 31, 23, 6, 48, 11, 24, 52, 13, 2}
	c := []item{1, 2, 6, 7, 9, 11, 12, 13, 15, 18, 20, 23, 24, 25, 28, 31, 33, 35, 41, 42, 48, 52}

	h := NewHeap()
	if !h.Empty() || h.FindMin() != nil || h.DeleteMin() != nil {
		t.Fail()
	}

	for i := range a {
		h.Insert(a[i])
	}

	if h.Empty() {
		t.Fail()
	}

	if h.FindMin().(item) != 1 || h.Size() != 8 {
		t.Fail()
	}

	min := h.DeleteMin()
	if min.(item) != 1 || h.FindMin().(item) != 7 || h.Size() != 7 {
		t.Fail()
	}

	h.Insert(item(1))
	if h.FindMin().(item) != 1 || h.Size() != 8 {
		t.Fail()
	}

	for _, v := range h.Values() {
		if !h.Search(v.(item)) {
			t.Fail()
		}
	}

	for i, v := range h.PopAllItems() {
		if v.(item) != aSorted[i] {
			t.Fail()
		}
	}

	h = nil

	h1, h2 := NewHeap(), NewHeap()
	for i := range a {
		h1.Insert(a[i])
	}
	if h1.Meld(h2) != h1 || h2.Meld(h1) != h1 {
		t.Fail()
	}
	for i := range b {
		h2.Insert(b[i])
	}

	h3 := h1.Meld(h2)
	h3.Print()

	if h3.itemNum != len(c) || h3.FindMin().(item) != 1 {
		t.Fail()
	}

	for _, v := range c {
		if !h3.Search(v) {
			t.Fail()
		}
	}

	// c := []item{1, 2, 6, 7, 9, 11, 12, 13, 15, 18, 20, 23, 24, 25, 28, 31, 33, 35, 41, 42, 48, 52}
	err := h3.IncreaseKey(h3.root, item(36))
	if err != nil {
		t.Fail()
	}
	if h3.FindMin().(item) != 2 {
		t.Fail()
	}

	cRevised := []item{-1, 2, 6, 9, 11, 12, 13, 15, 18, 20, 23, 24, 25, 28, 31, 33, 35, 36, 41, 42, 48, 52}
	err = h3.DecreaseKey(h3.search(h3.root, item(7)), item(-1))
	if err != nil {
		t.Fail()
	}
	if h3.FindMin().(item) != -1 {
		t.Fail()
	}
	if h3.Search(item(100)) {
		t.Fail()
	}
	for _, v := range cRevised {
		if !h3.Search(v) {
			t.Fail()
		}
	}

	if it := h3.Delete(item(25), item(math.Inf(-1))); it.(item) != 25 {
		t.Fail()
	}

	n := h3.search(h3.root, item(41))
	if n == nil {
		t.Fail()
	}
	if !h3.Update(n, item(16)) {
		t.Fail()
	}

	n = h3.search(h3.root, item(24))
	if n == nil {
		t.Fail()
	}
	if !h3.Update(n, item(58)) {
		t.Fail()
	}

	h1.Print()

	h1.Clear()

	if !h1.Empty() || h1.itemNum != 0 {
		t.Fail()
	}
}

func BenchmarkHeap(b *testing.B) {
	for k := 1.0; k <= 5; k++ {
		n := int(math.Pow(10, k))

		h := NewHeap()

		rn := 0
		for i := 0; i < n; i++ {
			rn = container.GenerateRandomInt()
			h.Insert(item(rn))
		}

		num := item(container.GenerateRandomInt())
		b.ResetTimer()

		b.Run("Fibonacci Heap FindMin: size-"+strconv.Itoa(n), func(b *testing.B) {
			for i := 1; i < b.N; i++ {
				h.FindMin()
			}
		})

		b.Run("Fibonacci Heap DeleteMin: size-"+strconv.Itoa(n), func(b *testing.B) {
			for i := 1; i < b.N; i++ {
				h.DeleteMin()
			}
		})

		b.Run("Fibonacci Heap Insert: size-"+strconv.Itoa(n), func(b *testing.B) {
			for i := 1; i < b.N; i++ {
				h.Insert(num)
			}
		})

		b.Run("Fibonacci Heap Meld: size-"+strconv.Itoa(n), func(b *testing.B) {
			for i := 1; i < b.N; i++ {
				h.Meld(h)
			}
		})
	}
}

/*
BenchmarkHeap/Fibonacci_Heap_FindMin:_size-10-8         	2000000000	         0.28 ns/op
BenchmarkHeap/Fibonacci_Heap_FindMin:_size-100-8        	2000000000	         0.26 ns/op
BenchmarkHeap/Fibonacci_Heap_FindMin:_size-1000-8       	2000000000	         0.26 ns/op
BenchmarkHeap/Fibonacci_Heap_FindMin:_size-10000-8      	2000000000	         0.26 ns/op
BenchmarkHeap/Fibonacci_Heap_FindMin:_size-100000-8     	2000000000	         0.26 ns/op

BenchmarkHeap/Fibonacci_Heap_DeleteMin:_size-10-8       	1000000000	         2.38 ns/op
BenchmarkHeap/Fibonacci_Heap_DeleteMin:_size-100-8      	1000000000	         2.31 ns/op
BenchmarkHeap/Fibonacci_Heap_DeleteMin:_size-1000-8     	1000000000	         2.30 ns/op
BenchmarkHeap/Fibonacci_Heap_DeleteMin:_size-10000-8    	1000000000	         2.30 ns/op
BenchmarkHeap/Fibonacci_Heap_DeleteMin:_size-100000-8   	1000000000	         2.31 ns/op

BenchmarkHeap/Fibonacci_Heap_Insert:_size-10-8          	10000000	          123 ns/op
BenchmarkHeap/Fibonacci_Heap_Insert:_size-100-8         	10000000	          107 ns/op
BenchmarkHeap/Fibonacci_Heap_Insert:_size-1000-8        	20000000	          114 ns/op
BenchmarkHeap/Fibonacci_Heap_Insert:_size-10000-8       	20000000	          112 ns/op
BenchmarkHeap/Fibonacci_Heap_Insert:_size-100000-8      	10000000	          134 ns/op

BenchmarkHeap/Fibonacci_Heap_Meld:_size-10-8            	100000000	         19.5 ns/op
BenchmarkHeap/Fibonacci_Heap_Meld:_size-100-8           	100000000	         19.5 ns/op
BenchmarkHeap/Fibonacci_Heap_Meld:_size-1000-8          	100000000	         19.4 ns/op
BenchmarkHeap/Fibonacci_Heap_Meld:_size-10000-8         	100000000	         19.5 ns/op
BenchmarkHeap/Fibonacci_Heap_Meld:_size-100000-8        	100000000	         19.4 ns/op
 */
