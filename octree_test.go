package golina

import (
	"testing"
)

func TestOctree_GetNodeTreeDepth(t *testing.T) {
	node := OctreeNode{
		code:     127,
		hasChild: 3,
	}
	tree := Octree{Nodes: map[uint32]*OctreeNode{}}
	tree.Nodes[node.code] = &node
	if tree.GetNodeTreeDepth(&node) != 2 {
		t.Fail()
	}
}

func TestOctree_GetParentNode(t *testing.T) {
	pnode := OctreeNode{
		code:     15,
		hasChild: 7,
	}
	node := OctreeNode{
		code:     127,
		hasChild: 3,
	}
	tree := Octree{Nodes: map[uint32]*OctreeNode{}}
	tree.Nodes[pnode.code] = &pnode
	tree.Nodes[node.code] = &node
	if tree.GetParentNode(&node) != &pnode {
		t.Fail()
	}
}

func TestOctree_LookupNode(t *testing.T) {
	node := OctreeNode{
		code:     127,
		hasChild: 3,
	}
	tree := Octree{Nodes: map[uint32]*OctreeNode{}}
	tree.Nodes[node.code] = &node
	if tree.LookupNode(127) != &node {
		t.Fail()
	}
}
