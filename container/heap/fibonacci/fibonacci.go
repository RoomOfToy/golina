package fibonacci

import (
	"fmt"
	"golina/container/heap"
	"math"
)

// FibonacciHeap
//	https://en.wikipedia.org/wiki/Fibonacci_heap
//	https://www.cnblogs.com/skywang12345/p/3659060.html
//	minimum heap (consist of a series of minimum ordered tree)
type Heap struct {
	// maintain a pointer to the root containing the minimum key
	root    *node
	itemNum int
}

// node consists of Fibonacci Heap, used internally
type node struct {
	item heap.Item
	// children are linked using a circular doubly linked list (parent / child)
	parent, child *node
	// each child has two siblings (left / right)
	left, right *node
	// isMarked indicates whether current node lost first child from last time when it became another node's child
	isMarked bool
	// number of children
	degree int
}

// NewHeap creates a new empty heap
func NewHeap() *Heap {
	return &Heap{
		root:    nil,
		itemNum: 0,
	}
}

// FindMin returns the minimum item inside the heap (root)
func (h *Heap) FindMin() heap.Item {
	if h.root == nil {
		return nil
	}
	return h.root.item
}

// DeleteMin (Extract Min) delete the root node from doubly linked list and returns its item
//	then find the new root
func (h *Heap) DeleteMin() heap.Item {
	if h.root == nil {
		return nil
	}
	r := h.root
	// add root's children (child and its left / right siblings) to heap's root list
	for {
		if x := r.child; x != nil {
			x.parent = nil
			// sibling
			if x.right != x {
				r.child = x.right
				x.right.left = x.left
				x.left.right = x.right
			} else {
				r.child = nil
			}
			// add x to root doubly linked list (h.insert)
			x.left = r.left
			x.right = r
			r.left.right = x
			r.left = x
		} else {
			break
		}
	}

	// remove r from heap's root doubly linked list
	r.left.right = r.right
	r.right.left = r.left

	// check roots after removal
	if r == r.right {
		h.root = nil
	} else {
		h.root = r.right
		// consolidate heap
		h.consolidate()
	}
	h.itemNum--
	return r.item
}

func (h *Heap) consolidate() {
	// use a map to tell whether there exists two sub-trees (roots) which have the same degree
	degreeMap := make(map[int]*node)
	head := h.root
	tail := head.left
	for {
		x := head
		next := head.right
		deg := head.degree
		for {
			if y, found := degreeMap[deg]; !found {
				// there's no roots with the same degree
				break
			} else {
				// the less one is the root (x)
				if y.item.Compare(x.item) < 0 {
					x, y = y, x
				}
				h.link(y, x)
				delete(degreeMap, deg)
				deg++
			}
		}
		degreeMap[deg] = x
		// loop over the whole list
		if head == tail {
			break
		}
		// move to next
		head = next
	}
	// reconstruct the heap from degreeMap
	h.root = nil
	for _, n := range degreeMap {
		h.insertNode(n)
	}
}

// link node to root
func (h *Heap) link(n, r *node) {
	// remove node n from heap's root list
	n.right.left = n.left
	n.left.right = n.right
	// link n to r (n as r's child and increase r's degree)
	n.parent = r
	r.degree++

	if r.child == nil {
		// r has no child, so n is the only child in the children doubly linked list
		r.child = n
		n.left = n
		n.right = n
	} else {
		// insert to r's children doubly linked list
		h.insert(n, r.child)
	}
	// isMarked indicates whether current node lost first child from last time when it became another node's child
	// here n just becomes r's child, so set n.isMarked to false
	n.isMarked = false
}

// Insert inserts item into heap, just insert it into heap's roots doubly linked list
//	if the item less than min, replace the min
func (h *Heap) Insert(item heap.Item) {
	nNode := &node{
		item:     item,
		parent:   nil,
		child:    nil,
		left:     nil,
		right:    nil,
		isMarked: false,
		degree:   0,
	}
	h.insertNode(nNode)
	h.itemNum++
}

func (h *Heap) insertNode(nNode *node) {
	if h.root == nil {
		nNode.left = nNode
		nNode.right = nNode
		h.root = nNode
		return
	}
	// insert node before root
	h.insert(nNode, h.root)
	// if the item less than min, replace the min (root)
	if nNode.item.Compare(h.root.item) < 0 {
		h.root = nNode
	}
}

// insert node into heap's root list
//	insert node before root, which means the `tail` of doubly linked list
func (h *Heap) insert(nNode, root *node) {
	nNode.left = root.left
	root.left.right = nNode
	nNode.right = root
	root.left = nNode
}

// deleteNode deletes node from heap and returns its item
//	1. decrease node's item to a value less than min (root's item)
//	2. call DeleteMin
// this function need to define how to decrease item, so ignored here

// DecreaseKey decrease input node's item to nItem
//	1. cut decreased node from its heap and add this node (single node or root of a sub-tree) to roots list
//	2. cascading cut on decreased node's parent node to ensure the min heap property
//	3. update heap root (min)
func (h *Heap) DecreaseKey(n *node, nItem heap.Item) {
	if nItem.Compare(n.item) >= 0 {
		fmt.Printf("decrease failed: the new item(%+v) is no smaller than current item(%+v)\n", nItem, n.item)
		return
	}

	n.item = nItem
	if p := n.parent; p != nil && p.item.Compare(n.item) > 0 {
		// cut node from parent node and add node to roots list
		h.cut(n, p)
		h.cascadingCut(p)
	}
	// update heap root (min)
	if n.item.Compare(h.root.item) < 0 {
		h.root = n
	}
}

// cut node from its parent
func (h *Heap) cut(n, p *node) {
	// remove node from its parent's children list and decrease its parent's degree
	if n.right == n {
		// no sibling
		p.child = nil
	} else {
		p.child = n.right
		n.right.left = n.left
		n.left.right = n.right
	}
	p.degree--
	// add n to roots list
	h.insert(n, h.root)

	n.parent = nil
	n.isMarked = false
}

// cascadingCut recursively cut nodes starting from the root of tree whose child has been cut
func (h *Heap) cascadingCut(n *node) {
	if p := n.parent; p != nil {
		return
	} else {
		if n.isMarked == false {
			// n has been cut a child
			n.isMarked = true
		} else {
			h.cut(n, p)
			h.cascadingCut(p)
		}
	}
}

// IncreaseKey decrease input node's item to nItem
//	1. add increased node's children (child and child's siblings) into roots list
//	2. cut (cut and cascadingCut) increased node and add it into roots list
//	3. update heap root (min) if n is root
func (h *Heap) IncreaseKey(n *node, nItem heap.Item) {
	if nItem.Compare(n.item) <= 0 {
		fmt.Printf("decrease failed: the new item(%+v) is no larger than current item(%+v)\n", nItem, n.item)
		return
	}

	for {
		if child := n.child; child == nil {
			break
		} else {
			// remove child from children list
			child.left.right = child.right
			child.right.left = child.left

			// update n.child
			if child.right == child {
				n.child = nil
			} else {
				n.child = child.right
			}

			// add child into roots list
			h.insert(child, h.root)
			child.parent = nil
		}
	}

	// add node into roots list
	n.degree = 0
	n.item = nItem
	if p := n.parent; p != nil {
		h.cut(n, p)
		h.cascadingCut(p)
	} else {
		// update heap root (min) if n is root
		if h.root == n {
			right := n.right
			for right != n {
				if h.root.item.Compare(right.item) > 0 {
					h.root = right
				}
				right = right.right
			}
		}
	}
}

// Meld returns union of two heaps (notice: in place change)
//	for efficiency consideration, add to heap which has larger maxDegree to achieve less operations
func (h *Heap) Meld(ah *Heap) *Heap {
	if h.root == nil {
		return ah
	}
	if ah.root == nil {
		return h
	}

	h1, h2 := h, ah

	if ah.maxDegree() > h.maxDegree() {
		h1, h2 = h2, h1
	}

	h.cat(h2.root, h1.root)
	if h2.root.item.Compare(h1.root.item) < 0 {
		h1.root = h2.root
	}
	h1.itemNum += h2.itemNum
	return h1
}

// cat append n2 to n1
//	notice: n1 and n2 are doubly linked list node, this method is different from insert
func (h *Heap) cat(n1, n2 *node) {
	var tmp *node
	tmp = n1.right
	n1.right = n2.right
	n2.right.left = n1
	n2.right = tmp
	tmp.left = n2
}

// maxDegree estimates max degree of heap
func (h *Heap) maxDegree() int {
	return int(math.Log2(float64(h.itemNum))) + 1 // ceil
}

// Size returns item number inside the heap
func (h *Heap) Size() int {
	return h.itemNum
}

// Empty returns true if no item inside heap
func (h *Heap) Empty() bool {
	return h.root == nil
}

// Clear clears heap by setting its root to nil
func (h *Heap) Clear() {
	h.root = nil
	h.itemNum = 0
}

// Values returns values inside the heap by recursively call DeleteMin
//	Warning: after call this method, the heap will be cleared!
//	TODO: other way to traverse
func (h *Heap) Values() []interface{} {
	num := h.itemNum
	var res []interface{}
	for i := 0; i < num; i++ {
		res = append(res, h.DeleteMin())
	}
	return res
}

// TODO: Search, Update, Delete
