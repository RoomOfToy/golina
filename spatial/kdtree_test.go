package spatial

import (
	"fmt"
	"golina/matrix"
	"testing"
)

func TestNewNode(t *testing.T) {
	p := &matrix.Vector{1, 2, 3}
	n := NewNode(p)
	if n.Point != p || n.Left != nil || n.Right != nil {
		t.Fail()
	}
}

func TestNode_String(t *testing.T) {
	p := &matrix.Vector{1, 2, 3}
	n := NewNode(p)
	if n.String() != "Node->point: {1.000000, 2.000000, 3.000000}\n" {
		t.Fail()
	}
}

func TestKDTree_String(t *testing.T) {
	m := matrix.GenerateRandomMatrix(10, 3)
	fmt.Println(m)
	tree := KDTree{}
	for i := range m.Data {
		tree.Insert(&(m.Data[i]))
	}
	fmt.Println(tree.String())
}

func TestNode_Insert(t *testing.T) {
	p, q := &matrix.Vector{1, 2, 3}, &matrix.Vector{4, 5, 6}
	n := NewNode(p)
	n.Insert(q)
	if n.Right.Point != q {
		t.Fail()
	}
}

func TestKDTree_Insert(t *testing.T) {
	p := &matrix.Vector{1, 2, 3}
	tree := KDTree{}
	tree.Insert(p)
	q := &matrix.Vector{4, 5, 6}
	tree.Insert(q)
	if tree.Count != 2 || tree.Root.Point != p || tree.Root.Right.Point != q {
		t.Fail()
	}
}

func TestKDTree_Search(t *testing.T) {
	m := matrix.GenerateRandomMatrix(10, 3)
	tree := KDTree{}
	for i := range m.Data {
		tree.Insert(&(m.Data[i]))
	}
	node, res := tree.Search(&(m.Data[5]))
	if !res || node.Point != &(m.Data[5]) {
		t.Fail()
	}
}

func TestKDTree_FindMinValue(t *testing.T) {
	m := matrix.GenerateRandomMatrix(10, 3)
	tree := KDTree{}
	for i := range m.Data {
		tree.Insert(&(m.Data[i]))
	}
	_, value := m.Col(0).Min()
	if tree.FindMinValue(0) != value {
		t.Fail()
	}
}

func TestKDTree_FindMinNode(t *testing.T) {
	m := matrix.GenerateRandomMatrix(10, 3)
	tree := KDTree{}
	for i := range m.Data {
		tree.Insert(&(m.Data[i]))
	}
	idx, _ := m.Col(0).Min()
	if tree.FindMinNode(0).Point != m.Row(idx) {
		t.Fail()
	}
}

func TestKDTree_DeleteNode(t *testing.T) {
	m := matrix.GenerateRandomMatrix(10, 3)
	tree := KDTree{}
	for i := range m.Data {
		tree.Insert(&(m.Data[i]))
	}
	nt := KDTree{tree.DeleteNode(&(m.Data[5])), tree.Count - 1}
	_, res := nt.Search(&(m.Data[5]))
	if res {
		t.Fail()
	}
}
