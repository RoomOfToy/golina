package mesh

import (
	"fmt"
)

type Triangle struct {
	Id        uint
	Vertexes  []*Point
	Neighbors []*Triangle
}

var cnt uint

func NewTriangle(v0, v1, v2 *Point) *Triangle {
	cnt++
	return &Triangle{
		Id:        cnt,
		Vertexes:  []*Point{v0, v1, v2},
		Neighbors: make([]*Triangle, 3),
	}
}

func (tr *Triangle) IsCoincidentWith(pt *Point) bool {
	return PEqual(tr.Vertexes[0], pt) || PEqual(tr.Vertexes[1], pt) || PEqual(tr.Vertexes[2], pt)
}

func (tr *Triangle) AssignNeighbors(tr0, tr1, tr2 *Triangle) {
	tr.Neighbors = []*Triangle{tr0, tr1, tr2}
}

func (tr *Triangle) String() string {
	return fmt.Sprintf("Triangle ID: %d\n", tr.Id) + fmt.Sprintf("	Vertex[0]: %s	Vertex[1]: %s	Vertex[2]: %s", tr.Vertexes[0], tr.Vertexes[1], tr.Vertexes[2]) + fmt.Sprintf("	Neighbors[0] ID: %d	Neighbors[0] ID: %d	Neighbors[0] ID: %d", tr.Neighbors[0].Id, tr.Neighbors[1].Id, tr.Neighbors[2].Id)
}
