package btree

import (
	"fmt"
	"golina/container"
)

// https://en.wikipedia.org/wiki/B-tree

// Item stored int Node
type Item struct {
	Key   interface{}
	Value interface{}
}

func (item *Item) String() string {
	return fmt.Sprintf("{%+v: %+v}", item.Key, item.Value)
}

// Node in BTree
type Node struct {
	Items    []*Item
	Children []*Node
	Parent   *Node
}

func NewNode(item *Item) *Node {
	return &Node{
		Items:    []*Item{item},
		Children: []*Node{}, // value ascending from left to right
		Parent:   nil,
	}
}

type BTree struct {
	Root       *Node
	Comparator container.Comparator
	M          int // MaxNumOfChildren on each node, MinNumOfChildren on each non-leaf node = ceil(M / 2), Middle for split = (M - 1) / 2
	ItemsNum   int // Items number of the tree (include all items on all nodes)
}

func NewBTree(M int, comparator container.Comparator) *BTree {
	if M < 3 {
		panic("Maximum number of children on each node should be at least 3")
	}
	return &BTree{
		Root:       nil,
		Comparator: comparator,
		M:          M,
	}
}

func (bTree *BTree) Empty() bool {
	return bTree.ItemsNum == 0
}

func (bTree *BTree) Size() int {
	return bTree.ItemsNum
}

func (bTree *BTree) isLeaf(node *Node) bool {
	return len(node.Children) == 0 // single root is also treated as leaf
}

func (bTree *BTree) isFull(node *Node) bool {
	return len(node.Items) == bTree.M-1
}

func (bTree *BTree) requireSplit(node *Node) bool {
	return len(node.Items) > bTree.M-1
}

func (bTree *BTree) Insert(item *Item) {
	if bTree.Root == nil {
		bTree.Root = NewNode(item)
		bTree.ItemsNum++
	} else {
		if bTree.insert(bTree.Root, item) {
			bTree.ItemsNum++
		}
	}
}

func (bTree *BTree) insert(node *Node, item *Item) (inserted bool) {
	if bTree.isLeaf(node) {
		return bTree.insertLeaf(node, item)
	}
	return bTree.insertInternal(node, item)
}

func (bTree *BTree) findInsertPosition(node *Node, item *Item) (index int, found bool) {
	// binary search, fit for larger M
	// for small M, directly traversal may be more efficient
	head, tail := 0, len(node.Items)-1
	mid := 0
	for head <= tail {
		mid = (head + tail) / 2
		switch bTree.Comparator(node.Items[mid].Key, item.Key) {
		case 1:
			tail = mid - 1
		case 0:
			return mid, true // find same item (key) in child, should replace the old item value
		case -1:
			head = mid + 1
		}
	}
	return head, false
}

func (bTree *BTree) insertItemAtIdxIntoNode(node *Node, index int, item *Item) {
	node.Items = append(node.Items, nil)
	if index < len(node.Items) {
		copy(node.Items[index+1:], node.Items[index:])
	}
	node.Items[index] = item
}

func (bTree *BTree) insertNodeAtIdxIntoChildren(parent *Node, index int, node *Node) {
	parent.Children = append(parent.Children, nil)
	if index < len(parent.Children) {
		copy(parent.Children[index+1:], parent.Children[index:])
	}
	parent.Children[index] = node
}

func (bTree *BTree) insertLeaf(node *Node, item *Item) (inserted bool) {
	insertPos, found := bTree.findInsertPosition(node, item)
	// find item with the same key, just update item value, elements number remains the same
	if found {
		node.Items[insertPos] = item
		return false
	}
	bTree.insertItemAtIdxIntoNode(node, insertPos, item)
	// check split
	bTree.split(node)
	return true
}

func (bTree *BTree) insertInternal(node *Node, item *Item) (inserted bool) {
	insertPos, found := bTree.findInsertPosition(node, item)
	// find item with the same key, just update item value, elements number remains the same
	if found {
		node.Items[insertPos] = item
		return false
	}
	// insert leaf first then split
	return bTree.insert(node.Children[insertPos], item)
}

func (bTree *BTree) split(node *Node) {
	// no need to split
	if !bTree.requireSplit(node) {
		return
	}

	// split parent
	// if parent is root, will create a new root and tree height += 1
	// non-root parent has constrains on its item num:
	//	if M is odd -> mid = (M - 1) / 2
	//	if M is even -> mid = M / 2 - 1 (right sub-tree items number >= left sub-tree items number) = (M - 1) / 2
	//	mid = (M - 1) / 2
	if node == bTree.Root {
		bTree.splitRoot()
	} else {
		bTree.splitNonRoot(node)
	}
}

func (bTree *BTree) splitRoot() {
	mid := (bTree.M - 1) / 2
	// split into left, right sub-trees
	leftSubTree, rightSubTree := bTree.splitIntoLR(bTree.Root, mid)
	// create new root
	newRoot := &Node{
		Items:    []*Item{bTree.Root.Items[mid]},
		Children: []*Node{leftSubTree, rightSubTree},
		Parent:   nil,
	}
	// set sub-trees' parent to new root
	leftSubTree.Parent, rightSubTree.Parent = newRoot, newRoot
	// set new root
	bTree.Root = newRoot
}

func (bTree *BTree) splitIntoLR(node *Node, mid int) (leftSubTree, rightSubTree *Node) {
	// create left, right sub-trees
	leftSubTree = &Node{
		Items:    node.Items[:mid],
		Children: nil,
		Parent:   nil,
	}
	rightSubTree = &Node{
		Items:    node.Items[mid+1:],
		Children: nil,
		Parent:   nil,
	}
	if len(node.Children) != 0 {
		leftSubTree.Children = node.Children[:mid+1]
		rightSubTree.Children = node.Children[mid+1:]
		// set subtrees' parent
		setParent(leftSubTree.Children, leftSubTree)
		setParent(rightSubTree.Children, rightSubTree)
	}
	return
}

func setParent(nodes []*Node, parent *Node) {
	for _, node := range nodes {
		node.Parent = parent
	}
}

func (bTree *BTree) splitNonRoot(node *Node) {
	mid := (bTree.M - 1) / 2
	parent := node.Parent

	// split into left, right sub-trees
	leftSubTree, rightSubTree := bTree.splitIntoLR(node, mid)
	// set sub-trees' parent
	leftSubTree.Parent, rightSubTree.Parent = parent, parent

	// insert node's middle item into parent
	item := node.Items[mid]
	insertPos, _ := bTree.findInsertPosition(parent, item)
	bTree.insertItemAtIdxIntoNode(parent, insertPos, item)

	// set parent's newly inserted item's corresponding node to leftSubTree
	parent.Children[insertPos] = leftSubTree

	// set parent's newly inserted item's next node to rightSubTree
	bTree.insertNodeAtIdxIntoChildren(parent, insertPos+1, rightSubTree)

	// check split
	bTree.split(parent)
}

func (bTree *BTree) Lookup(key interface{}) (value interface{}, found bool) {
	node, index, found := bTree.lookupRec(bTree.Root, key)
	if found {
		return node.Items[index].Value, true
	}
	return nil, false
}

func (bTree *BTree) lookupRec(startNode *Node, key interface{}) (node *Node, index int, found bool) {
	if bTree.Empty() {
		return nil, -1, false
	}
	node = startNode
	for {
		index, found = bTree.findInsertPosition(node, &Item{
			Key:   key,
			Value: nil,
		})
		if found {
			return node, index, true
		}
		if bTree.isLeaf(node) {
			return nil, -1, false
		}
		node = node.Children[index]
	}
}

func (bTree *BTree) Values() []interface{} {
	values := make([]interface{}, bTree.Size())
	var items []*Item
	bTree.items(bTree.Root, &items)
	for idx, i := range items {
		values[idx] = i.Value
	}
	return values
}

func (bTree *BTree) items(node *Node, items *[]*Item) {
	if len(node.Children) != 0 {
		mid := (len(node.Children)-1)/2 + 1
		// left sub-tree
		for _, c := range node.Children[:mid] {
			bTree.items(c, items)
		}
		// mid
		*items = append(*items, node.Items...)
		// right sub-tree
		for _, c := range node.Children[mid:] {
			bTree.items(c, items)
		}
	} else {
		*items = append(*items, node.Items...)
	}
}

func (bTree *BTree) Clear() {
	*bTree = *NewBTree(bTree.M, bTree.Comparator)
}

// TODO: delete
