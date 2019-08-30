package golina

import (
	"fmt"
	"testing"
)

func TestNewNode(t *testing.T) {
	p := &Vector{1, 2, 3}
	n := NewNode(p)
	if n.point != p || n.left != nil || n.right != nil {
		t.Fail()
	}
}

func TestNode_String(t *testing.T) {
	p := &Vector{1, 2, 3}
	n := NewNode(p)
	if n.String() != "Node->point: {1.000000, 2.000000, 3.000000}\n" {
		t.Fail()
	}
}

func TestKDTree_String(t *testing.T) {
	m := GenerateRandomMatrix(10, 3)
	fmt.Println(m)
	tree := KDTree{}
	for i := range m._array {
		tree.Insert(&(m._array[i]))
	}
	fmt.Println(tree.String())
}

func TestNode_Insert(t *testing.T) {
	p, q := &Vector{1, 2, 3}, &Vector{4, 5, 6}
	n := NewNode(p)
	n.Insert(q)
	if n.right.point != q {
		t.Fail()
	}
}

func TestKDTree_Insert(t *testing.T) {
	p := &Vector{1, 2, 3}
	tree := KDTree{}
	tree.Insert(p)
	q := &Vector{4, 5, 6}
	tree.Insert(q)
	if tree.count != 2 || tree.root.point != p || tree.root.right.point != q {
		t.Fail()
	}
}

func TestKDTree_Search(t *testing.T) {
	m := GenerateRandomMatrix(10, 3)
	tree := KDTree{}
	for i := range m._array {
		tree.Insert(&(m._array[i]))
	}
	node, res := tree.Search(&(m._array[5]))
	if !res || node.point != &(m._array[5]) {
		t.Fail()
	}
}

func TestKDTree_FindMinValue(t *testing.T) {
	m := GenerateRandomMatrix(10, 3)
	tree := KDTree{}
	for i := range m._array {
		tree.Insert(&(m._array[i]))
	}
	_, value := m.Col(0).Min()
	if tree.FindMinValue(0) != value {
		t.Fail()
	}
}

func TestKDTree_FindMinNode(t *testing.T) {
	m := GenerateRandomMatrix(10, 3)
	tree := KDTree{}
	for i := range m._array {
		tree.Insert(&(m._array[i]))
	}
	idx, _ := m.Col(0).Min()
	if tree.FindMinNode(0).point != m.Row(idx) {
		t.Fail()
	}
}

func TestKDTree_DeleteNode(t *testing.T) {
	m := GenerateRandomMatrix(10, 3)
	tree := KDTree{}
	for i := range m._array {
		tree.Insert(&(m._array[i]))
	}
	nt := KDTree{tree.DeleteNode(&(m._array[5])), tree.count - 1}
	_, res := nt.Search(&(m._array[5]))
	if res {
		t.Fail()
	}
}
