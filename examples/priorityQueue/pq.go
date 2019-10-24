package main

import (
	"fmt"
	"golina/container/heap/bheap"
)

type item struct {
	priority int
	content  interface{}
}

type priorityQueue struct {
	items *bheap.MaxHeap
}

func itemComparator(a, b interface{}) int {
	itemA, itemB := a.(*item), b.(*item)
	A, B := itemA.priority, itemB.priority
	if A > B {
		return 1
	} else if A == B {
		return 0
	} else {
		return -1
	}
}

func newPQ() *priorityQueue {
	maxH := bheap.MaxHeap{
		Comparator: itemComparator,
	}
	maxH.Init()
	return &priorityQueue{items: &maxH}
}

func (pq *priorityQueue) push(item *item) {
	pq.items.Push(item)
}

func (pq *priorityQueue) pop() *item {
	return pq.items.Pop().(*item)
}

func (pq *priorityQueue) len() int {
	return pq.items.Size()
}

func (pq *priorityQueue) top() *item {
	return pq.items.Values()[0].(*item)
}

func (pq *priorityQueue) values() []item {
	var items []item
	for _, v := range pq.items.Values() {
		items = append(items, *(v.(*item)))
	}
	return items
}

func main() {
	items := []item{
		{
			priority: 3,
			content:  "hello",
		},
		{
			priority: 6,
			content:  "world",
		},
		{
			priority: 1,
			content:  "yo",
		},
	}

	pq := newPQ()
	for i := range items {
		pq.push(&items[i])
	}

	fmt.Printf("Top: %+v\n", pq.top())

	fmt.Printf("Items in pq: %+v\n", pq.values())

	for pq.len() > 0 {
		fmt.Printf("Sorted items: %+v\n", pq.pop())
	}
}
