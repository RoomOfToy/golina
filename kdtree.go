package golina

import (
	"fmt"
	"math"
	"strconv"
)

type KDTree struct {
	root  *Node
	count int
}

type Node struct {
	point       *Vector
	left, right *Node
}

func NewNode(p *Vector) *Node {
	return &Node{
		point: p,
		left:  nil,
		right: nil,
	}
}

// recursively
func insert(n *Node, p *Vector, depth int) *Node {
	if n == nil {
		return NewNode(p)
	}
	currDim := depth % p.Length()
	if compare(p, n.point, currDim) < 0 {
		n.left = insert(n.left, p, depth+1)
	} else {
		n.right = insert(n.right, p, depth+1)
	}
	return n
}

func compare(p, q *Vector, dim int) float64 {
	return p.At(dim) - q.At(dim)
}

func (n *Node) Insert(p *Vector) *Node { // after insert, return node self
	return insert(n, p, 0)
}

func (t *KDTree) Insert(p *Vector) bool {
	if t.root == nil {
		t.root = NewNode(p)
		t.count++
		return true
	}
	if t.root.Insert(p) != nil {
		t.count++
		return true
	}
	return false
}

func isPointSame(p1, p2 *Vector) bool {
	for i := range *p1 {
		if p1.At(i) != p2.At(i) {
			return false
		}
	}
	return true
}

// recursively
func search(root *Node, p *Vector, depth int) (*Node, bool) {
	if root == nil {
		return nil, false
	}
	if isPointSame(root.point, p) {
		return root, true
	}
	currDim := depth % p.Length()
	if compare(p, root.point, currDim) < 0 {
		return search(root.left, p, depth+1)
	}
	return search(root.right, p, depth+1)
}

func (t *KDTree) Search(p *Vector) (*Node, bool) {
	return search(t.root, p, 0)
}

func ternaryMinValue(x, y, z float64) float64 {
	return math.Min(x, math.Min(y, z))
}

func findMinValue(root *Node, dim, depth int) float64 {
	if root == nil {
		return math.MaxFloat64
	}
	currDim := depth % root.point.Length()
	if currDim == dim {
		if root.left == nil {
			return root.point.At(dim)
		}
		return math.Min(root.point.At(dim), findMinValue(root.left, dim, depth+1))
	}
	return ternaryMinValue(root.point.At(dim), findMinValue(root.left, dim, depth+1), findMinValue(root.right, dim, depth+1))
}

func (t *KDTree) FindMinValue(dim int) float64 {
	return findMinValue(t.root, dim, 0)
}

func ternaryMinNode(x, y, z *Node, dim int) *Node {
	res := x
	if y != nil && y.point.At(dim) < res.point.At(dim) {
		res = y
	}
	if z != nil && z.point.At(dim) < res.point.At(dim) {
		res = z
	}
	return res
}

func findMinNode(root *Node, dim, depth int) *Node {
	if root == nil {
		return nil
	}
	currDim := depth % root.point.Length()
	if currDim == dim {
		if root.left == nil {
			return root
		}
		return findMinNode(root.left, dim, depth+1)
	}
	return ternaryMinNode(root, findMinNode(root.left, dim, depth+1), findMinNode(root.right, dim, depth+1), dim)
}

func (t *KDTree) FindMinNode(dim int) *Node {
	return findMinNode(t.root, dim, 0)
}

func deleteNode(root *Node, p *Vector, depth int) *Node { // return root after modification
	if root == nil {
		return nil
	}
	currDim := depth % p.Length()
	if isPointSame(root.point, p) {
		if root.right != nil {
			rminNode := findMinNode(root.right, currDim, 0)
			root.point = rminNode.point
			root.right = deleteNode(root.right, rminNode.point, depth+1)
		} else if root.left != nil {
			lminNode := findMinNode(root.left, currDim, 0)
			root.point = lminNode.point
			root.right = deleteNode(root.left, lminNode.point, depth+1)
		} else {
			root = nil
		}
		return root
	}
	if compare(p, root.point, currDim) < 0 {
		root.left = deleteNode(root.left, p, depth+1)
	} else {
		root.right = deleteNode(root.right, p, depth+1)
	}
	return root
}

func (t *KDTree) DeleteNode(p *Vector) *Node { // return root of modified tree
	return deleteNode(t.root, p, 0)
}

func (n *Node) String() string {
	return fmt.Sprintf("Node->point: %s", n.point.String())
}

// TODO: find a better way for pretty print
func printPreOrder(n *Node, depth int) string {
	if n == nil {
		return "<nil>\n"
	}
	res := "depth: " + strconv.Itoa(depth) + " root: " + n.String()
	depth++
	res += "Left: " + printPreOrder(n.left, depth)
	res += "Right: " + printPreOrder(n.right, depth)
	return res
}

func (t *KDTree) String() string {
	root := t.root
	if root == nil {
		return "<nil>"
	}
	return printPreOrder(root, 0)
}
