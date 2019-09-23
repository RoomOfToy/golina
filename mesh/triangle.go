package mesh

import (
	"fmt"
	"golina/matrix"
)

type Triangle struct {
	Id        uint
	Vertexes  []Point
	Neighbors []*Triangle
}

var cnt uint

func NewTriangle(v0, v1, v2 Point) *Triangle {
	cnt++
	return &Triangle{
		Id:        cnt,
		Vertexes:  []Point{v0, v1, v2},
		Neighbors: make([]*Triangle, 3),
	}
}

func (tr *Triangle) IsCoincidentWith(pt Point) bool {
	return PEqual(tr.Vertexes[0], pt) || PEqual(tr.Vertexes[1], pt) || PEqual(tr.Vertexes[2], pt)
}

func (tr *Triangle) AssignNeighbors(tr0, tr1, tr2 *Triangle) {
	tr.Neighbors = []*Triangle{tr0, tr1, tr2}
}

func (tr *Triangle) String() string {
	return fmt.Sprintf("Triangle ID: %d\n", tr.Id) + fmt.Sprintf("	Vertex[0]: %s	Vertex[1]: %s	Vertex[2]: %s", tr.Vertexes[0].Vector, tr.Vertexes[1].Vector, tr.Vertexes[2].Vector) + fmt.Sprintf("	Neighbors[0] ID: %d	Neighbors[0] ID: %d	Neighbors[0] ID: %d", tr.Neighbors[0].Id, tr.Neighbors[1].Id, tr.Neighbors[2].Id)
}

// Incremental Insertion Algorithm
type DelaunayTriangle struct {
	iv              *initVarTriangulation
	Points          Points
	AuxiliaryPoints []Point
	ProjectedPoints []Point
	Triangles       []*Triangle
}

func NewDelaunayTriangle(points Points, iv *initVarTriangulation) *DelaunayTriangle {
	ap := make([]Point, iv.initHullVerticesCnt)
	for i := 0; i < iv.initHullVerticesCnt; i++ {
		switch i / 2 {
		case 0, 1, 2:
			if i%2 == 0 {
				ap[i] = Point{NewPoint(1, 1, 1).MulNum(iv.unitSphereRadius), false}
			} else {
				ap[i] = Point{NewPoint(-1, -1, -1).MulNum(iv.unitSphereRadius), false}
			}
		default:
			ap[i] = NewPoint(0, 0, 0)
		}
	}
	// N points -> 8 + (N - 6) * 2 triangles
	return &DelaunayTriangle{
		Points:          points,
		AuxiliaryPoints: ap,
		ProjectedPoints: make([]Point, len(points.Data)),
		Triangles:       make([]*Triangle, 0, 8+(len(points.Data)-6)*2),
	}
}

func (dt *DelaunayTriangle) BuildInitialHull(points *Points) {
	initialVertices := make([]Point, dt.iv.initHullVerticesCnt)
	initialHullFaces := make([]*Triangle, dt.iv.initHullFacesCnt)
	dist := make(matrix.Vector, dt.iv.initHullVerticesCnt)
	minDist := make(matrix.Vector, dt.iv.initHullVerticesCnt) // {0, 0, 0, 0, 0, 0}
	for _, pt := range dt.Points.Data {
		for i := range dist {
			dist[i] = dt.GetDistance(dt.AuxiliaryPoints[i], Point{&pt, false})
			if minDist[i] == 0 || dist[i] < minDist[i] {
				minDist[i] = dist[i]
			}
		}

		for i := range dist {
			idx, val := dist.Min()
			if minDist[i] == dist[i] && idx == i && dist[i] == val {
				initialVertices[i] = Point{
					Vector:    &pt,
					IsVisited: false,
				}
			}
		}
	}

	vertex0Idx := []int{0, 0, 0, 0, 1, 1, 1, 1}
	vertex1Idx := []int{4, 3, 5, 2, 2, 4, 3, 5}
	vertex2Idx := []int{2, 4, 3, 5, 4, 3, 5, 2}

	for i := 0; i < dt.iv.initHullFacesCnt; i++ {
		v0 := initialVertices[vertex0Idx[i]]
		v1 := initialVertices[vertex1Idx[i]]
		v2 := initialVertices[vertex2Idx[i]]

		tri := NewTriangle(v0, v1, v2)
		initialHullFaces[i] = tri
		dt.Triangles = append(dt.Triangles, tri)
	}

	neighbor0Idx := []int{1, 2, 3, 0, 7, 4, 5, 6}
	neighbor1Idx := []int{4, 5, 6, 7, 0, 1, 2, 3}
	neighbor2Idx := []int{3, 0, 1, 2, 5, 6, 7, 4}

	for i := 0; i < dt.iv.initHullFacesCnt; i++ {
		n0 := initialHullFaces[neighbor0Idx[i]]
		n1 := initialHullFaces[neighbor1Idx[i]]
		n2 := initialHullFaces[neighbor2Idx[i]]
		initialHullFaces[i].Neighbors = []*Triangle{n0, n1, n2}
	}

	for i := 0; i < dt.iv.initHullVerticesCnt; i++ {
		initialVertices[i].IsVisited = true
	}
}

func (dt *DelaunayTriangle) InsertPoint(point Point) {}

func (dt *DelaunayTriangle) RemoveExtraTriangles() {}

func (dt *DelaunayTriangle) SplitTriangle(tri *Triangle, pt *Point) {}

func (dt *DelaunayTriangle) FixNeighborhood(target, oldNeighbor, newNeighbor *Triangle) {}

func (dt *DelaunayTriangle) DoLocalOptimization(tri0, tri1 *Triangle) {}

func (dt *DelaunayTriangle) TrySwapDiagonal(tri0, tri1 *Triangle) bool {
	return false
}

func (dt *DelaunayTriangle) GetDistance(pt0, pt1 Point) float64 {
	return pt0.Sub(pt1.Vector).Norm()
}

func (dt *DelaunayTriangle) GetDetPoints(pt0, pt1, pt2 Point) float64 {
	return new(matrix.Matrix).Init(matrix.Data{*(pt0.Vector), *(pt1.Vector), *(pt2.Vector)}).Det()
}

func (dt *DelaunayTriangle) GetDetMatrix(pts *Points) float64 {
	return pts.Det()
}
