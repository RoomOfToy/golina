package tree

import (
	"fmt"
	"golina/container"
)

// Binary Search Tree
//	http://cslibrary.stanford.edu/110/BinaryTrees.html
type BSTree struct {
	Root       *Node
	Comparator container.Comparator
	Diffidence func(a, b interface{}) interface{}
}

func treeInterfaceAssertion() {
	var _ Tree = (*BSTree)(nil)
}

func NewBSTree() *BSTree {
	return &BSTree{
		Root:       nil,
		Comparator: nil,
	}
}

type Node struct {
	data        interface{}
	left, right *Node
}

func NewNode(data interface{}) *Node {
	return &Node{
		data:  data,
		left:  nil,
		right: nil,
	}
}

func (bst *BSTree) lookup(node *Node, target interface{}) bool {
	if node == nil {
		return false
	} else {
		switch bst.Comparator(target, node.data) {
		case 0:
			return true
		case -1:
			return bst.lookup(node.left, target)
		case 1:
			return bst.lookup(node.right, target)
		default:
			return false
		}
	}
}

func (bst *BSTree) Lookup(target interface{}) bool {
	return bst.lookup(bst.Root, target)
}

func (bst *BSTree) insert(node *Node, data interface{}) *Node {
	if node == nil {
		return NewNode(data)
	} else {
		switch bst.Comparator(data, node.data) {
		case -1, 0:
			node.left = bst.insert(node.left, data)
		case 1:
			node.right = bst.insert(node.right, data)
		}
		return node
	}
}

func (bst *BSTree) Insert(data interface{}) {
	bst.Root = bst.insert(bst.Root, data)
}

func (bst *BSTree) print(node *Node) {
	if node == nil {
		return
	}
	bst.print(node.left)
	fmt.Printf("%+v ", node.data)
	bst.print(node.right)
}

func (bst *BSTree) Print() {
	bst.print(bst.Root)
	fmt.Println()
}

func (bst *BSTree) Empty() bool {
	return bst.Root == nil
}

func (bst *BSTree) size(node *Node) int {
	if node == nil {
		return 0
	} else {
		return bst.size(node.left) + 1 + bst.size(node.right)
	}
}

func (bst *BSTree) Size() int {
	return bst.size(bst.Root)
}

func (bst *BSTree) value(node *Node, dataSlice *[]interface{}) {
	if node == nil {
		return
	}
	bst.value(node.left, dataSlice)
	*dataSlice = append(*dataSlice, node.data)
	bst.value(node.right, dataSlice)
}

func (bst *BSTree) Values() []interface{} {
	var dataSlice []interface{}
	bst.value(bst.Root, &dataSlice)
	return dataSlice
}

func (bst *BSTree) Clear() {
	bst.Root = nil
}

func (bst *BSTree) maxDepth(node *Node) int {
	if node == nil {
		return 0
	} else {
		lDepth := bst.maxDepth(node.left)
		rDepth := bst.maxDepth(node.right)
		if lDepth > rDepth {
			return lDepth + 1
		} else {
			return rDepth + 1
		}
	}
}

func (bst *BSTree) MaxDepth() int {
	return bst.maxDepth(bst.Root)
}

func (bst *BSTree) minValue(node *Node) interface{} {
	currentNode := node
	for currentNode.left != nil {
		currentNode = currentNode.left
	}
	return currentNode.data
}

func (bst *BSTree) MinValue() interface{} {
	return bst.minValue(bst.Root)
}

func (bst *BSTree) hasPathSum(node *Node, sum interface{}) bool {
	if node == nil {
		return sum == 0
	} else {
		sum = bst.Diffidence(sum, node.data)
		return bst.hasPathSum(node.left, sum) || bst.hasPathSum(node.right, sum)
	}
}

func (bst *BSTree) HasPathSum(sum interface{}) bool {
	if bst.Diffidence == nil {
		panic("no Difference function for node data")
	}
	return bst.hasPathSum(bst.Root, sum)
}

func (bst *BSTree) printPaths(node *Node, path []interface{}) {
	if node == nil {
		return
	}
	path = append(path, node.data)
	if node.left == nil && node.right == nil {
		for _, p := range path {
			fmt.Printf("%+v ", p)
		}
		fmt.Println()
	} else {
		bst.printPaths(node.left, path)
		bst.printPaths(node.right, path)
	}
}

// all root to leaf paths
func (bst *BSTree) PrintPaths() {
	var paths []interface{}
	bst.printPaths(bst.Root, paths)
}

func (bst *BSTree) mirror(node *Node) {
	if node == nil {
		return
	} else {
		bst.mirror(node.left)
		bst.mirror(node.right)
		// swap
		node.left, node.right = node.right, node.left
	}
}

// in-place change
func (bst *BSTree) Mirror() {
	bst.mirror(bst.Root)
}

func (bst *BSTree) doubleTree(node *Node) {
	if node == nil {
		return
	}

	bst.doubleTree(node.left)
	bst.doubleTree(node.right)

	oldLeft := node.left
	node.left = NewNode(node.data)
	node.left.left = oldLeft
}

func (bst *BSTree) DoubleTree() {
	bst.doubleTree(bst.Root)
}

func (bst *BSTree) sameTree(nodeA, nodeB *Node) bool {
	if nodeA == nil && nodeB == nil {
		return true
	} else if nodeA != nil && nodeB != nil {
		return nodeA.data == nodeB.data && bst.sameTree(nodeA.left, nodeB.left) && bst.sameTree(nodeA.right, nodeB.right)
	} else {
		return false
	}
}

func (bst *BSTree) SameTree(bstB *BSTree) bool {
	return bst.sameTree(bst.Root, bstB.Root)
}

func CountTrees(numKeys int) int {
	if numKeys <= 1 {
		return 1
	} else {
		sum, left, right := 0, 0, 0
		for root := 0; root < numKeys; root++ {
			// left tree node num
			left = CountTrees(root)
			// right tree node num
			right = CountTrees(numKeys - 1 - root)
			sum += left * right
		}
		return sum
	}
}

func (bst *BSTree) isBST(node *Node, min, max interface{}) bool {
	if node == nil {
		return true
	}

	if bst.Comparator(node.data, min) < 0 || bst.Comparator(node.data, max) > 0 {
		return false
	}

	return bst.isBST(node.left, min, node.data) && bst.isBST(node.right, node.data, max)
}

func (bst *BSTree) IsBST(min, max interface{}) bool {
	return bst.isBST(bst.Root, min, max)
}
