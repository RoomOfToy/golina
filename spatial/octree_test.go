package spatial

import (
	"testing"
)

func TestOctree_GetNodeTreeDepth(t *testing.T) {
	node := OctreeNode{
		Code:     127,
		HasChild: 3,
	}
	tree := Octree{Nodes: map[uint32]*OctreeNode{}}
	tree.Nodes[node.Code] = &node
	if tree.GetNodeTreeDepth(&node) != 2 {
		t.Fail()
	}
}

func TestOctree_GetParentNode(t *testing.T) {
	pnode := OctreeNode{
		Code:     15,
		HasChild: 7,
	}
	node := OctreeNode{
		Code:     127,
		HasChild: 3,
	}
	tree := Octree{Nodes: map[uint32]*OctreeNode{}}
	tree.Nodes[pnode.Code] = &pnode
	tree.Nodes[node.Code] = &node
	if tree.GetParentNode(&node) != &pnode {
		t.Fail()
	}
}

func TestOctree_LookupNode(t *testing.T) {
	node := OctreeNode{
		Code:     127,
		HasChild: 3,
	}
	tree := Octree{Nodes: map[uint32]*OctreeNode{}}
	tree.Nodes[node.Code] = &node
	if tree.LookupNode(127) != &node {
		t.Fail()
	}
}
