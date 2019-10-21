package deque

import (
	"fmt"
	"golina/container"
	"math"
	"strconv"
	"testing"
)

func TestDeque(t *testing.T) {

	var _ container.Container = (*Deque)(nil)

	dq := NewDeque(-10)

	if dq.Cap() != 8 || !dq.Empty() {
		t.Fail()
	}

	a := []int{1, 2, 3, 4, 5, 6, 7}
	for _, i := range a {
		dq.PushBack(i)
	}

	if dq.Size() != 7 || dq.PositionsCanPushBack() != 1 || dq.PositionsCanPopFront() != 7 {
		t.Fail()
	}

	for i, v := range dq.Values() {
		if v.(int) != a[i] {
			fmt.Printf("%d", v.(int))
			t.Fail()
		}
	}

	last, err := dq.PopBack()
	if err != nil || last != 7 || dq.Size() != 6 || dq.PositionsCanPushBack() != 2 || dq.PositionsCanPopFront() != 6 {
		t.Fail()
	}

	first, err := dq.PopFront()
	if err != nil || first != 1 || dq.Size() != 5 || dq.PositionsCanPushBack() != 3 || dq.PositionsCanPopFront() != 5 {
		t.Fail()
	}

	dq.PushFront(0)
	f, err := dq.Front()
	if err != nil || f != 0 || dq.Size() != 6 || dq.PositionsCanPushBack() != 2 || dq.PositionsCanPopFront() != 6 {
		t.Fail()
	}

	// [0, 2, 3, 4, 5, 6] + [1, 2, 3, 4, 5, 6, 7]
	for _, i := range a {
		dq.PushBack(i)
	}
	f, _ = dq.Front()
	b, err := dq.Back()
	if dq.Cap() != 16 || dq.Size() != 13 || err != nil || f != 0 || b != 7 || dq.PositionsCanPushBack() != 3 || dq.PositionsCanPopFront() != 13 {
		t.Fail()
	}

	// [7, 6, 5, 4, 3, 2, 1] + [0, 2, 3, 4, 5, 6] + [1, 2, 3, 4, 5, 6, 7]
	for _, i := range a {
		dq.PushFront(i)
	}
	f, _ = dq.Front()
	b, _ = dq.Back()
	if dq.Cap() != 32 || dq.Size() != 20 || f != 7 || b != 7 || dq.PositionsCanPushBack() != 12 || dq.PositionsCanPopFront() != 20 {
		fmt.Println(f, b)
		fmt.Println(dq)
		t.Fail()
	}

	dq.Clear()

	if !dq.Empty() || dq.Cap() != 32 || dq.Size() != 0 || dq.PositionsCanPushBack() != 32 || dq.PositionsCanPopFront() != 0 {
		t.Fail()
	}
}

func BenchmarkDeque(b *testing.B) {
	for k := 1.0; k <= 3; k++ {
		n := int(math.Pow(10, k))

		b.Run("Deque PushBack: size-"+strconv.Itoa(n), func(b *testing.B) {
			dq := NewDeque(n)
			rn := 0
			for i := 0; i < n; i++ {
				rn = container.GenerateRandomInt()
				dq.PushBack(rn)
			}
			num := container.GenerateRandomInt()
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				dq.PushBack(num)
			}
		})

		b.Run("Deque PushFront: size-"+strconv.Itoa(n), func(b *testing.B) {
			dq := NewDeque(n)
			rn := 0
			for i := 0; i < n; i++ {
				rn = container.GenerateRandomInt()
				dq.PushBack(rn)
			}
			num := container.GenerateRandomInt()
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				dq.PushFront(num)
			}
		})

		b.Run("Deque PopBack: size-"+strconv.Itoa(n), func(b *testing.B) {
			dq := NewDeque(n)
			rn := 0
			for i := 0; i < n; i++ {
				rn = container.GenerateRandomInt()
				dq.PushBack(rn)
			}
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				_, _ = dq.PopBack()
			}
		})

		b.Run("Deque PopFront: size-"+strconv.Itoa(n), func(b *testing.B) {
			dq := NewDeque(n)
			rn := 0
			for i := 0; i < n; i++ {
				rn = container.GenerateRandomInt()
				dq.PushBack(rn)
			}
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				_, _ = dq.PopFront()
			}
		})

		b.Run("Deque Front: size-"+strconv.Itoa(n), func(b *testing.B) {
			dq := NewDeque(n)
			rn := 0
			for i := 0; i < n; i++ {
				rn = container.GenerateRandomInt()
				dq.PushBack(rn)
			}
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				_, _ = dq.Front()
			}
		})

		b.Run("Deque Back: size-"+strconv.Itoa(n), func(b *testing.B) {
			dq := NewDeque(n)
			rn := 0
			for i := 0; i < n; i++ {
				rn = container.GenerateRandomInt()
				dq.PushBack(rn)
			}
			b.ResetTimer()
			for i := 1; i < b.N; i++ {
				_, _ = dq.Back()
			}
		})
	}
}

/*
BenchmarkDeque/Deque_PushBack:_size-10-8         	20000000	         51.4 ns/op
BenchmarkDeque/Deque_PushBack:_size-100-8        	30000000	         45.4 ns/op
BenchmarkDeque/Deque_PushBack:_size-1000-8       	30000000	         36.5 ns/op

BenchmarkDeque/Deque_PushFront:_size-10-8        	30000000	         55.5 ns/op
BenchmarkDeque/Deque_PushFront:_size-100-8       	30000000	         50.6 ns/op
BenchmarkDeque/Deque_PushFront:_size-1000-8      	30000000	         36.9 ns/op

BenchmarkDeque/Deque_PopBack:_size-10-8          	10000000	          118 ns/op
BenchmarkDeque/Deque_PopBack:_size-100-8         	10000000	          121 ns/op
BenchmarkDeque/Deque_PopBack:_size-1000-8        	10000000	          121 ns/op

BenchmarkDeque/Deque_PopFront:_size-10-8         	10000000	          122 ns/op
BenchmarkDeque/Deque_PopFront:_size-100-8        	10000000	          123 ns/op
BenchmarkDeque/Deque_PopFront:_size-1000-8       	10000000	          123 ns/op

BenchmarkDeque/Deque_Front:_size-10-8            	1000000000	         2.69 ns/op
BenchmarkDeque/Deque_Front:_size-100-8           	1000000000	         2.68 ns/op
BenchmarkDeque/Deque_Front:_size-1000-8          	1000000000	         2.71 ns/op

BenchmarkDeque/Deque_Back:_size-10-8             	500000000	         3.10 ns/op
BenchmarkDeque/Deque_Back:_size-100-8            	500000000	         3.11 ns/op
BenchmarkDeque/Deque_Back:_size-1000-8           	500000000	         3.08 ns/op
*/
