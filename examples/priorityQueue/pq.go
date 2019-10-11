package main

import (
	"fmt"
	"golina/container/tree"
)

type Item struct {
	priority int
	content  interface{}
}

type PriorityQueue struct {
	items *tree.MaxHeap
}

func ItemComparator(a, b interface{}) int {
	itemA, itemB := a.(*Item), b.(*Item)
	A, B := itemA.priority, itemB.priority
	if A > B {
		return 1
	} else if A == B {
		return 0
	} else {
		return -1
	}
}

func NewPQ() *PriorityQueue {
	maxH := tree.MaxHeap{
		Comparator: ItemComparator,
	}
	maxH.Init()
	return &PriorityQueue{items: &maxH}
}

func (pq *PriorityQueue) push(item *Item) {
	pq.items.Push(item)
}

func (pq *PriorityQueue) pop() *Item {
	return pq.items.Pop().(*Item)
}

func (pq *PriorityQueue) len() int {
	return pq.items.Size()
}

func (pq *PriorityQueue) top() *Item {
	return pq.items.Values()[0].(*Item)
}

func (pq *PriorityQueue) values() []Item {
	var items []Item
	for _, v := range pq.items.Values() {
		items = append(items, *(v.(*Item)))
	}
	return items
}

func main() {
	items := []Item{
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

	pq := NewPQ()
	for i := range items {
		pq.push(&items[i])
	}

	fmt.Printf("Top: %+v\n", pq.top())

	fmt.Printf("Items in pq: %+v\n", pq.values())

	for pq.len() > 0 {
		fmt.Printf("Sorted items: %+v\n", pq.pop())
	}
}
