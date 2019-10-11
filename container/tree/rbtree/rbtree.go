package rbtree

import (
	"fmt"
	"golina/container"
)

// https://en.wikipedia.org/wiki/Red%E2%80%93black_tree

const (
	Red   = false
	Black = true
)

type Node struct {
	value                       interface{}
	color                       bool
	leftTree, rightTree, parent *Node
}

func NewNode(value interface{}) *Node {
	return &Node{
		value:     value,
		color:     Red,
		leftTree:  nil,
		rightTree: nil,
		parent:    nil,
	}
}

func (node *Node) grandparent() *Node {
	if node.parent == nil {
		return nil
	}
	return node.parent.parent
}

func (node *Node) uncle() *Node {
	if node.grandparent() == nil {
		return nil
	}
	if node.parent == node.grandparent().rightTree {
		return node.grandparent().leftTree
	} else {
		return node.grandparent().rightTree
	}
}

func (node *Node) sibling() *Node {
	if node.parent.leftTree == node {
		return node.parent.rightTree
	} else {
		return node.parent.leftTree
	}
}

// in-order
func (node *Node) String() string {
	if node == nil {
		return "[nil]"
	}
	s := ""
	if node.leftTree != nil {
		s += node.leftTree.String() + " "
	}
	c := "Red"
	if node.color {
		c = "Black"
	}
	s += fmt.Sprintf("%v(%s)", node.value, c)
	if node.rightTree != nil {
		s += " " + node.rightTree.String()
	}
	return "[" + s + "]\n"
}

type RBTree struct {
	Root       *Node
	Comparator container.Comparator
}

func NewRBTree(root *Node, comparator container.Comparator) *RBTree {
	rbTree := &RBTree{Root: root, Comparator: comparator}
	rbTree.Root.color = Black
	return rbTree
}

/*
Rotate left on Y:
     gp                gp
     /                 /
    X                 Y
   / \               / \
  a   Y    ----->   X   c
     / \           / \
    b   c         a   b
*/
func (rbTree *RBTree) rotateLeft(Y *Node) {
	if Y.parent == nil {
		rbTree.Root = Y
		return
	}

	gp := Y.grandparent()
	X := Y.parent
	b := Y.leftTree

	X.rightTree = b
	if b != nil {
		b.parent = X
	}
	Y.leftTree = X
	X.parent = Y

	if rbTree.Root == X {
		rbTree.Root = Y
	}
	Y.parent = gp

	if gp != nil {
		if gp.leftTree == X {
			gp.leftTree = Y
		} else {
			gp.rightTree = Y
		}
	}
}

/*
Rotate right on Y:
     gp                gp
     /                 /
    X                 Y
   / \               / \
  Y   a    ----->   b   X
 / \                   / \
b   c                 c   a
*/
func (rbTree *RBTree) rotateRight(Y *Node) {
	gp := Y.grandparent()
	X := Y.parent
	c := Y.rightTree

	X.leftTree = c

	if c != nil {
		c.parent = X
	}
	Y.rightTree = X
	X.parent = Y

	if rbTree.Root == X {
		rbTree.Root = Y
	}
	Y.parent = gp

	if gp != nil {
		if gp.leftTree == X {
			gp.leftTree = Y
		} else {
			gp.rightTree = Y
		}
	}
}

var NIL = &Node{
	value:     nil,
	color:     Black,
	leftTree:  nil,
	rightTree: nil,
	parent:    nil,
}

/*
Y color: Red
     gp                gp
     /                 /
    X                 X
   / \               / \
  a   Y    <---->   Y   a
     / \           / \
    b   c         b   c
*/
func (rbTree *RBTree) insertCase(Y *Node) {
	if Y.parent == nil {
		rbTree.Root = Y
		Y.color = Black
		return
	}

	gp := Y.grandparent()
	X := Y.parent
	a := Y.uncle()
	b := Y.leftTree
	c := Y.rightTree

	if X.color == Red {
		if Y.uncle().color == Red {
			X.color, a.color = Black, Black
			gp.color = Red
			rbTree.insertCase(gp)
		} else {
			// rotate to left
			if X.rightTree == Y && gp.leftTree == X {
				rbTree.rotateLeft(Y)
				Y.color = Black
				b.color, c.color = Red, Red
			} else if X.leftTree == Y && gp.rightTree == X {
				// rotate to right
				rbTree.rotateRight(Y)
				Y.color = Black
				b.color, c.color = Red, Red
			} else if X.leftTree == Y && gp.leftTree == X {
				X.color = Black
				gp.color = Red
				rbTree.rotateRight(X)
			} else if X.rightTree == Y && gp.rightTree == X {
				X.color = Black
				gp.color = Red
				rbTree.rotateLeft(X)
			}
		}
	}
}

func (rbTree *RBTree) insert(node *Node, value interface{}) {
	if rbTree.Comparator(node.value, value) >= 0 {
		if node.leftTree != nil {
			rbTree.insert(node.leftTree, value)
		} else {
			tmp := NewNode(value)
			tmp.leftTree, tmp.rightTree = NIL, NIL
			tmp.parent = node
			node.leftTree = tmp
			rbTree.insertCase(tmp)
		}
	} else {
		if node.rightTree != nil {
			rbTree.insert(node.rightTree, value)
		} else {
			tmp := NewNode(value)
			tmp.leftTree, tmp.rightTree = NIL, NIL
			tmp.parent = node
			node.rightTree = tmp
			rbTree.insertCase(tmp)
		}
	}
}

func (rbTree *RBTree) Insert(value interface{}) {
	if rbTree.Root == nil {
		rbTree.Root = NewNode(value)
		rbTree.Root.color = Black
	}
	rbTree.insert(rbTree.Root, value)
}

func (rbTree *RBTree) minValue(node *Node) interface{} {
	if node.leftTree == nil {
		return node.value
	}
	return rbTree.minValue(node.leftTree)
}

func (rbTree *RBTree) MinValue() interface{} {
	return rbTree.minValue(rbTree.Root)
}

func (rbTree *RBTree) deleteOneChild(X *Node) {
	Y, a := NIL, NIL
	if X.leftTree != nil {
		Y = X.leftTree
		a = X.rightTree
	} else {
		Y = X.rightTree
		a = X.leftTree
	}

	gp := X.parent

	if gp == nil && Y == NIL && a == nil {
		X = nil
		rbTree.Root = X
		return
	}

	if gp == nil {
		X = nil
		Y.parent = X
		rbTree.Root = Y
		rbTree.Root.color = Black
		return
	}

	if gp.leftTree == X {
		gp.leftTree = Y
	} else {
		gp.rightTree = Y
	}
	Y.parent = X.parent

	if X.color == Black {
		if Y.color == Red {
			Y.color = Black
		} else {
			rbTree.deleteCase(Y)
		}
	}
	X = nil
}

func (rbTree *RBTree) deleteCase(Y *Node) {
	X := Y.parent
	if X == nil {
		Y.color = Black
		return
	}

	S := Y.sibling()
	if S.color == Red {
		Y.parent.color = Red
		S.color = Black
		if Y == X.leftTree {
			rbTree.rotateLeft(X)
		} else {
			rbTree.rotateRight(X)
		}
	}

	if X.color == Black && S.color == Black && S.leftTree.color == Black && S.rightTree.color == Black {
		S.color = Red
		rbTree.deleteCase(X)
	} else if X.color == Red && S.color == Black && S.leftTree.color == Black && S.rightTree.color == Black {
		S.color = Red
		X.color = Black
	} else {
		if S.color == Black {
			if Y == X.leftTree && S.leftTree.color == Red && S.rightTree.color == Black {
				S.color = Red
				S.leftTree.color = Black
				rbTree.rotateRight(S.leftTree)
			} else if Y == X.rightTree && S.leftTree.color == Black && S.rightTree.color == Red {
				S.color = Red
				S.rightTree.color = Black
				rbTree.rotateLeft(S.rightTree)
			}
		}

		S.color = X.color
		X.color = Black
		if Y == X.leftTree {
			S.rightTree.color = Black
			rbTree.rotateLeft(S)
		} else {
			S.leftTree.color = Black
			rbTree.rotateRight(S)
		}
	}
}

func (rbTree *RBTree) deleteTree(node *Node) {
	if node == nil {
		return
	}
	rbTree.deleteTree(node.leftTree)
	rbTree.deleteTree(node.rightTree)
	node = nil
}

func (rbTree *RBTree) Clear() {
	rbTree.deleteTree(rbTree.Root)
}
