package bheap

import (
	"container/heap"
	"golina/container"
)

// use heap interface from "container/heap"
type Heap struct {
	values     []interface{}
	comparator container.Comparator
}

func (h Heap) Len() int {
	return len(h.values)
}

func (h Heap) Less(i, j int) bool {
	return h.comparator(h.values[i], h.values[j]) < 0
}

func (h Heap) Swap(i, j int) {
	h.values[i], h.values[j] = h.values[j], h.values[i]
}

func (h *Heap) Push(v interface{}) {
	h.values = append(h.values, v)
}

func (h *Heap) Pop() interface{} {
	size := len(h.values)
	lastNode := h.values[size-1]
	h.values = h.values[:size-1]
	return lastNode
}

func (h *Heap) Size() int {
	return len(h.values)
}

func (h *Heap) Empty() bool {
	return len(h.values) == 0
}

func (h *Heap) Clear() {
	h.values = ([]interface{})(nil)
}

func (h *Heap) Values() []interface{} {
	return h.values
}

func (h *Heap) set(i int, v interface{}) {
	if i > len(h.values)-1 || i < 0 {
		panic("invalid index")
	}
	h.values[i] = v
}

type MinHeap struct {
	heap       *Heap
	Comparator container.Comparator
}

func (h *MinHeap) Init() {
	h.heap = &Heap{
		values:     nil,
		comparator: h.Comparator,
	}
	heap.Init(h.heap)
}

func (h *MinHeap) Push(v interface{}) {
	heap.Push(h.heap, v)
}

func (h *MinHeap) Pop() interface{} {
	return heap.Pop(h.heap)
}

func (h *MinHeap) Remove(i int) interface{} {
	return heap.Remove(h.heap, i)
}

func (h *MinHeap) Set(i int, v interface{}) {
	h.heap.set(i, v)
	heap.Fix(h.heap, i)
}

func (h *MinHeap) Size() int {
	return len(h.heap.values)
}

func (h *MinHeap) Empty() bool {
	return len(h.heap.values) == 0
}

func (h *MinHeap) Clear() {
	h.heap.values = ([]interface{})(nil)
}

func (h *MinHeap) Values() []interface{} {
	return h.heap.values
}

type MaxHeap struct {
	heap       *Heap
	Comparator container.Comparator
}

func (h *MaxHeap) Init() {
	h.heap = &Heap{
		values: nil,
		comparator: func(a, b interface{}) int {
			return -h.Comparator(a, b)
		},
	}
	heap.Init(h.heap)
}

func (h *MaxHeap) Push(v interface{}) {
	heap.Push(h.heap, v)
}

func (h *MaxHeap) Pop() interface{} {
	return heap.Pop(h.heap)
}

func (h *MaxHeap) Remove(i int) interface{} {
	return heap.Remove(h.heap, i)
}

func (h *MaxHeap) Set(i int, v interface{}) {
	h.heap.set(i, v)
	heap.Fix(h.heap, i)
}

func (h *MaxHeap) Size() int {
	return len(h.heap.values)
}

func (h *MaxHeap) Empty() bool {
	return len(h.heap.values) == 0
}

func (h *MaxHeap) Clear() {
	h.heap.values = ([]interface{})(nil)
}

func (h *MaxHeap) Values() []interface{} {
	return h.heap.values
}
