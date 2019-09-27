package mesh

import (
	"container/list"
	"fmt"
	"golina/matrix"
	"os"
)

var vcnt uint

type Vertex struct {
	Id uint
	Point
	IsVisited bool
}

func NewVertex(x, y, z float64) Vertex {
	vcnt++
	return Vertex{
		Id:        vcnt,
		Point:     Point{&(matrix.Vector{x, y, z})},
		IsVisited: false,
	}
}

func VertexEqual(pt1, pt2 Vertex) bool {
	return pt1.At(0) == pt2.At(0) && pt1.At(1) == pt2.At(1) && pt1.At(2) == pt2.At(2)
}

type Triangle struct {
	Id        uint
	Vertexes  []Vertex
	Neighbors []*Triangle
}

var cnt uint

func NewTriangle(v0, v1, v2 Vertex) *Triangle {
	cnt++
	return &Triangle{
		Id:        cnt,
		Vertexes:  []Vertex{v0, v1, v2},
		Neighbors: make([]*Triangle, 3),
	}
}

func (tr *Triangle) IsCoincidentWith(pt Vertex) bool {
	return VertexEqual(tr.Vertexes[0], pt) || VertexEqual(tr.Vertexes[1], pt) || VertexEqual(tr.Vertexes[2], pt)
}

func (tr *Triangle) AssignNeighbors(tr0, tr1, tr2 *Triangle) {
	tr.Neighbors = []*Triangle{tr0, tr1, tr2}
}

func (tr *Triangle) String() string {
	return fmt.Sprintf("Triangle ID: %d\n", tr.Id) + fmt.Sprintf("	Vertex[0]: %s	Vertex[1]: %s	Vertex[2]: %s", tr.Vertexes[0].Vector, tr.Vertexes[1].Vector, tr.Vertexes[2].Vector) + fmt.Sprintf("	Neighbors[0] ID: %d	Neighbors[0] ID: %d	Neighbors[0] ID: %d\n", tr.Neighbors[0].Id, tr.Neighbors[1].Id, tr.Neighbors[2].Id)
}

// Incremental Insertion Algorithm
type DelaunayTriangle struct {
	iv              *initVarTriangulation
	Points          Points
	AuxiliaryPoints []Vertex
	ProjectedPoints []Vertex
	Triangles       *list.List
}

func NewDelaunayTriangle(points Points, iv *initVarTriangulation) *DelaunayTriangle {
	ap := make([]Vertex, iv.initHullVerticesCnt)
	for i := 0; i < iv.initHullVerticesCnt; i++ {
		ap[i] = NewVertex(
			matrix.Ternary(i%2 == 0, 1., -1.).(float64)*matrix.Ternary(i/2 == 0, iv.unitSphereRadius, 0.).(float64),
			matrix.Ternary(i%2 == 0, 1., -1.).(float64)*matrix.Ternary(i/2 == 1, iv.unitSphereRadius, 0.).(float64),
			matrix.Ternary(i%2 == 0, 1., -1.).(float64)*matrix.Ternary(i/2 == 2, iv.unitSphereRadius, 0.).(float64),
		)
		ap[i].IsVisited = true
	}
	// N points -> 8 + (N - 6) * 2 triangles
	return &DelaunayTriangle{
		iv:              iv,
		Points:          points,
		AuxiliaryPoints: ap,
		ProjectedPoints: make([]Vertex, len(points.Data)),
		Triangles:       list.New(),
		// Triangles:       make([]*Triangle, 0, 8+(len(points.Data)-6)*2),
	}
}

func (dt *DelaunayTriangle) ProjectPointsToUnitSphere() {
	for i := range dt.Points.Data {
		p := dt.Points.Data[i].MulNum(dt.iv.unitSphereRadius / dt.Points.Data[i].Norm())
		dt.ProjectedPoints[i] = NewVertex(p.At(0), p.At(1), p.At(2))
	}
}

// after projection
func (dt *DelaunayTriangle) BuildInitialHull() {
	initialVertices := make([]Vertex, dt.iv.initHullVerticesCnt)
	for i := range initialVertices {
		initialVertices[i] = Vertex{
			Point:     Point{dt.AuxiliaryPoints[i].Vector.Normalize()},
			IsVisited: true,
		}
	}
	initialHullFaces := make([]*Triangle, dt.iv.initHullFacesCnt)
	dist := make(matrix.Vector, dt.iv.initHullVerticesCnt)
	minDist := make(matrix.Vector, dt.iv.initHullVerticesCnt) // {0, 0, 0, 0, 0, 0}
	for _, pt := range dt.ProjectedPoints {
		for i := range dist {
			dist[i] = dt.GetDistance(dt.AuxiliaryPoints[i], pt)
			if minDist[i] == 0 || dist[i] < minDist[i] {
				minDist[i] = dist[i]
			}
		}

		for i := range dist {
			idx, val := dist.Min()
			if minDist[i] == dist[i] && idx == i && dist[i] == val {
				initialVertices[i] = pt
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
		dt.Triangles.PushBack(tri)
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

func (dt *DelaunayTriangle) InsertPoint(point Vertex) {
	det := []float64{0, 0, 0}
	for e := dt.Triangles.Front(); e != nil; e = e.Next() {
		tri := e.Value.(*Triangle)
		det[0] = dt.GetDetPoints(tri.Vertexes[0], tri.Vertexes[1], point)
		det[1] = dt.GetDetPoints(tri.Vertexes[1], tri.Vertexes[2], point)
		det[2] = dt.GetDetPoints(tri.Vertexes[2], tri.Vertexes[0], point)

		// if this point projected into an existing triangle, split the existing triangle to 3 new ones
		if det[0] >= 0 && det[1] >= 0 && det[2] >= 0 {
			// no common vertex -> point is projected into triangle interior
			if !tri.IsCoincidentWith(point) {
				dt.SplitTriangle(tri, point)
			}
			return
		} else if det[1] >= 0 && det[2] >= 0 {
			// projected on one side -> search neighbors
			tri = tri.Neighbors[0]
		} else if det[0] >= 0 && det[2] >= 0 {
			tri = tri.Neighbors[1]
		} else if det[0] >= 0 && det[1] >= 0 {
			tri = tri.Neighbors[2]
		} else if det[0] >= 0 {
			tri = tri.Neighbors[1]
		} else if det[1] >= 0 {
			tri = tri.Neighbors[2]
		} else if det[2] >= 0 {
			tri = tri.Neighbors[0]
		} else {
			continue
		}
	}
}

func (dt *DelaunayTriangle) RemoveExtraTriangles() {
	for e := dt.Triangles.Front(); e != nil; e = e.Next() {
		tri := e.Value.(*Triangle)
		isExtraTri := false
		for _, v := range tri.Vertexes {
			for _, p := range dt.AuxiliaryPoints {
				if VertexEqual(v, p) {
					isExtraTri = true
					break
				}
			}
			if isExtraTri == true {
				break
			}
		}

		if isExtraTri {
			dt.Triangles.Remove(e)
		}
	}
}

func (dt *DelaunayTriangle) SplitTriangle(tri *Triangle, point Vertex) {
	newTri1 := NewTriangle(point, tri.Vertexes[1], tri.Vertexes[2])
	newTri2 := NewTriangle(point, tri.Vertexes[2], tri.Vertexes[0])

	// flip
	tri.Vertexes[2] = tri.Vertexes[1]
	tri.Vertexes[1] = tri.Vertexes[0]
	tri.Vertexes[0] = point

	newTri1.AssignNeighbors(tri, tri.Neighbors[1], newTri2)
	newTri2.AssignNeighbors(newTri1, tri.Neighbors[2], tri)
	tri.AssignNeighbors(newTri2, tri.Neighbors[0], newTri1)

	dt.FixNeighborhood(newTri1.Neighbors[1], tri, newTri1)
	dt.FixNeighborhood(newTri2.Neighbors[1], tri, newTri2)

	dt.Triangles.PushBack(newTri1)
	dt.Triangles.PushBack(newTri2)

	// optimize triangles according to delaunay triangulation definition
	//dt.DoLocalOptimization(tri, tri.Neighbors[1])
	//dt.DoLocalOptimization(newTri1, newTri1.Neighbors[1])
	//dt.DoLocalOptimization(newTri2, newTri2.Neighbors[1])
}

func (dt *DelaunayTriangle) FixNeighborhood(target, oldNeighbor, newNeighbor *Triangle) {
	for i := 0; i < 3; i++ {
		if target.Neighbors[i] == oldNeighbor {
			target.Neighbors[i] = newNeighbor
			break
		}
	}
}

func (dt *DelaunayTriangle) DoLocalOptimization(tri0, tri1 *Triangle) {
	for i := 0; i < 3; i++ {
		if tri1.Vertexes[i] == tri0.Vertexes[0] || tri1.Vertexes[i] == tri0.Vertexes[1] || tri1.Vertexes[i] == tri0.Vertexes[2] {
			continue
		}
		pts := Points{new(matrix.Matrix).Init(matrix.Data{
			*(tri1.Vertexes[i].Sub(tri0.Vertexes[0].Vector)),
			*(tri1.Vertexes[i].Sub(tri0.Vertexes[1].Vector)),
			*(tri1.Vertexes[i].Sub(tri0.Vertexes[2].Vector)),
		})}
		if dt.GetDetMatrix(pts) <= 0 {
			break
		}
		if dt.TrySwapDiagonal(tri0, tri1) {
			return
		}
	}
}

func (dt *DelaunayTriangle) TrySwapDiagonal(tri0, tri1 *Triangle) bool {
	for j := 0; j < 3; j++ {
		for k := 0; k < 3; k++ {
			if tri0.Vertexes[j] != tri1.Vertexes[0] && tri0.Vertexes[j] != tri1.Vertexes[1] && tri0.Vertexes[j] != tri1.Vertexes[2] && tri0.Vertexes[k] != tri1.Vertexes[0] && tri0.Vertexes[k] != tri1.Vertexes[1] && tri0.Vertexes[k] != tri1.Vertexes[2] {
				tri0.Vertexes[(j+2)%3] = tri1.Vertexes[k]
				tri1.Vertexes[(k+2)%3] = tri0.Vertexes[j]

				tri0.Neighbors[(j+1)%3] = tri1.Neighbors[(k+2)%3]
				tri1.Neighbors[(k+1)%3] = tri0.Neighbors[(j+2)%3]
				tri0.Neighbors[(j+2)%3] = tri1
				tri1.Neighbors[(k+2)%3] = tri0

				dt.FixNeighborhood(tri0.Neighbors[(j+1)%3], tri1, tri0)
				dt.FixNeighborhood(tri1.Neighbors[(k+1)%3], tri0, tri1)

				dt.DoLocalOptimization(tri0, tri0.Neighbors[j])
				dt.DoLocalOptimization(tri0, tri0.Neighbors[(j+1)%3])
				dt.DoLocalOptimization(tri1, tri1.Neighbors[k])
				dt.DoLocalOptimization(tri1, tri1.Neighbors[(k+1)%3])

				return true
			}
		}
	}
	return false
}

func (dt *DelaunayTriangle) GetDistance(pt0, pt1 Vertex) float64 {
	return pt0.Sub(pt1.Vector).Norm()
}

func (dt *DelaunayTriangle) GetDetPoints(pt0, pt1, pt2 Vertex) float64 {
	return new(matrix.Matrix).Init(matrix.Data{*(pt0.Vector), *(pt1.Vector), *(pt2.Vector)}).Det()
}

func (dt *DelaunayTriangle) GetDetMatrix(pts Points) float64 {
	return pts.Det()
}

// 3 vertexes to form a triangle
func (dt *DelaunayTriangle) WriteVertexIDToTxt(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	for e := dt.Triangles.Front(); e != nil; e = e.Next() {
		_, err = fmt.Fprintf(file, "%d %d %d\n", e.Value.(*Triangle).Vertexes[0].Id, e.Value.(*Triangle).Vertexes[1].Id, e.Value.(*Triangle).Vertexes[2].Id)
	}

	return err
}
