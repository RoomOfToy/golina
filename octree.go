package golina

import (
	"fmt"
	"math/bits"
)

// base on hash map
type Octree struct {
	Nodes map[uint32]*OctreeNode // 2^32 nodes
}

func (t *Octree) GetParentNode(node *OctreeNode) *OctreeNode {
	return t.LookupNode(node.Code >> 3)
}

func (t *Octree) LookupNode(code uint32) *OctreeNode {
	return t.Nodes[code]
}

func (t *Octree) GetNodeTreeDepth(node *OctreeNode) int { // 0 - 10
	// 8 children -> 2^3 -> 3 bits (000: left-down - 111: right-top) -> one depth
	// golang x = 0 => 32 leading zeros, first node: 0000 0000 0000 0000 0000 0000 0000 0000 -> depth 0
	return (32 - bits.LeadingZeros32(node.Code)) / 3
}

// recursively: z-order curve
// 	000 -> 001 -> 010 -> 011 -> 100 -> 101 -> 110 -> 111 -> 000 000 ...
func (t *Octree) Transverse(node *OctreeNode) {
	var i uint32
	for i = 0; i < 8; i++ {
		if node.HasChild&(1<<i) != 0 {
			childCode := node.Code<<3 | i
			child := t.LookupNode(childCode)
			fmt.Println(child, child.Code)
			t.Transverse(child)
		}
	}
}

type OctreeNode struct {
	Code     uint32
	HasChild uint32
	Data     *Matrix
}
