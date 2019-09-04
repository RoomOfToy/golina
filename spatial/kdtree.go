package spatial

import (
	"fmt"
	"golina/matrix"
	"math"
	"strconv"
)

type KDTree struct {
	Root  *Node
	Count int
}

type Node struct {
	Point       *matrix.Vector
	Left, Right *Node
}

func NewNode(p *matrix.Vector) *Node {
	return &Node{
		Point: p,
		Left:  nil,
		Right: nil,
	}
}

// recursively
func insert(n *Node, p *matrix.Vector, depth int) *Node {
	if n == nil {
		return NewNode(p)
	}
	currDim := depth % p.Length()
	if compare(p, n.Point, currDim) < 0 {
		n.Left = insert(n.Left, p, depth+1)
	} else {
		n.Right = insert(n.Right, p, depth+1)
	}
	return n
}

func compare(p, q *matrix.Vector, dim int) float64 {
	return p.At(dim) - q.At(dim)
}

func (n *Node) Insert(p *matrix.Vector) *Node { // after insert, return node self
	return insert(n, p, 0)
}

func (t *KDTree) Insert(p *matrix.Vector) bool {
	if t.Root == nil {
		t.Root = NewNode(p)
		t.Count++
		return true
	}
	if t.Root.Insert(p) != nil {
		t.Count++
		return true
	}
	return false
}

func isPointSame(p1, p2 *matrix.Vector) bool {
	for i := range *p1 {
		if p1.At(i) != p2.At(i) {
			return false
		}
	}
	return true
}

// recursively
func search(root *Node, p *matrix.Vector, depth int) (*Node, bool) {
	if root == nil {
		return nil, false
	}
	if isPointSame(root.Point, p) {
		return root, true
	}
	currDim := depth % p.Length()
	if compare(p, root.Point, currDim) < 0 {
		return search(root.Left, p, depth+1)
	}
	return search(root.Right, p, depth+1)
}

func (t *KDTree) Search(p *matrix.Vector) (*Node, bool) {
	return search(t.Root, p, 0)
}

func ternaryMinValue(x, y, z float64) float64 {
	return math.Min(x, math.Min(y, z))
}

func findMinValue(root *Node, dim, depth int) float64 {
	if root == nil {
		return math.MaxFloat64
	}
	currDim := depth % root.Point.Length()
	if currDim == dim {
		if root.Left == nil {
			return root.Point.At(dim)
		}
		return math.Min(root.Point.At(dim), findMinValue(root.Left, dim, depth+1))
	}
	return ternaryMinValue(root.Point.At(dim), findMinValue(root.Left, dim, depth+1), findMinValue(root.Right, dim, depth+1))
}

func (t *KDTree) FindMinValue(dim int) float64 {
	return findMinValue(t.Root, dim, 0)
}

func ternaryMinNode(x, y, z *Node, dim int) *Node {
	res := x
	if y != nil && y.Point.At(dim) < res.Point.At(dim) {
		res = y
	}
	if z != nil && z.Point.At(dim) < res.Point.At(dim) {
		res = z
	}
	return res
}

func findMinNode(root *Node, dim, depth int) *Node {
	if root == nil {
		return nil
	}
	currDim := depth % root.Point.Length()
	if currDim == dim {
		if root.Left == nil {
			return root
		}
		return findMinNode(root.Left, dim, depth+1)
	}
	return ternaryMinNode(root, findMinNode(root.Left, dim, depth+1), findMinNode(root.Right, dim, depth+1), dim)
}

func (t *KDTree) FindMinNode(dim int) *Node {
	return findMinNode(t.Root, dim, 0)
}

func deleteNode(root *Node, p *matrix.Vector, depth int) *Node { // return root after modification
	if root == nil {
		return nil
	}
	currDim := depth % p.Length()
	if isPointSame(root.Point, p) {
		if root.Right != nil {
			rminNode := findMinNode(root.Right, currDim, 0)
			root.Point = rminNode.Point
			root.Right = deleteNode(root.Right, rminNode.Point, depth+1)
		} else if root.Left != nil {
			lminNode := findMinNode(root.Left, currDim, 0)
			root.Point = lminNode.Point
			root.Right = deleteNode(root.Left, lminNode.Point, depth+1)
		} else {
			root = nil
		}
		return root
	}
	if compare(p, root.Point, currDim) < 0 {
		root.Left = deleteNode(root.Left, p, depth+1)
	} else {
		root.Right = deleteNode(root.Right, p, depth+1)
	}
	return root
}

func (t *KDTree) DeleteNode(p *matrix.Vector) *Node { // return root of modified tree
	return deleteNode(t.Root, p, 0)
}

func (n *Node) String() string {
	return fmt.Sprintf("Node->point: %s", n.Point.String())
}

// TODO: find a better way for pretty print
func printPreOrder(n *Node, depth int) string {
	if n == nil {
		return "<nil>\n"
	}
	res := "depth: " + strconv.Itoa(depth) + " root: " + n.String()
	depth++
	res += "Left: " + printPreOrder(n.Left, depth)
	res += "Right: " + printPreOrder(n.Right, depth)
	return res
}

func (t *KDTree) String() string {
	root := t.Root
	if root == nil {
		return "<nil>"
	}
	return printPreOrder(root, 0)
}
