package fibonacci

import (
	"golina/container/heap"
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

	for i, v := range h.Values() {
		if v.(item) != aSorted[i] {
			t.Fail()
		}
	}

	if !h.Empty() {
		t.Fail()
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

	if h3.itemNum != len(c) || h3.FindMin().(item) != 1 {
		t.Fail()
	}

	for i, v := range h3.Values() {
		if v.(item) != c[i] {
			t.Fail()
		}
	}

	for i := range a {
		h1.Insert(a[i])
	}
	h1.Clear()

	if !h1.Empty() || h1.itemNum != 0 {
		t.Fail()
	}

	// a := []item{12, 7, 25, 15, 28, 33, 41, 1}
	for i := range a {
		h1.Insert(a[i])
	}

	aRevised := []item{7, 12, 15, 25, 28, 33, 36, 41}
	h1.IncreaseKey(h1.root, item(36))
	if h1.FindMin().(item) != 7 {
		t.Fail()
	}
	for i, v := range h1.Values() {
		if v.(item) != aRevised[i] {
			t.Fail()
		}
	}

	for i := range a {
		h1.Insert(a[i])
	}
	aRevised = []item{-1, 1, 7, 15, 25, 28, 33, 41}
	h1.DecreaseKey(h1.root.right.right, item(-1))
	if h1.FindMin().(item) != -1 {
		t.Fail()
	}
	for i, v := range h1.Values() {
		if v.(item) != aRevised[i] {
			t.Fail()
		}
	}
}
