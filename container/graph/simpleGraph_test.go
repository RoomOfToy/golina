package graph

import (
	"fmt"
	"testing"
)

func TestGraph(t *testing.T) {
	graph, err := NewGraphFromJson("test.json", "graph")
	if err != nil {
		panic(err)
	}
	fmt.Println(graph)
	if graph.Empty() {
		t.Fail()
	}
	if graph.NodeNum() != 8 || graph.Size() != 8 {
		t.Fail()
	}
	for _, v := range graph.Values() {
		if _, existed := graph.GetNodes()[v.(StringID)]; !existed {
			t.Fail()
		}
	}
	if node, existed := graph.GetNode(StringID("A")); !existed {
		t.Fail()
	} else {
		if graph.EdgeNum() != 8 {
			t.Fail()
		}
		in, out, err := graph.GetEdges(node.ID())
		if err != nil {
			t.Fail()
		}
		if len(in) != 4 || len(out) != 4 {
			t.Fail()
		}

		es, err := graph.GetOutEdges(node.ID())
		if err != nil {
			t.Fail()
		}
		for _, e := range es {
			if e.Source() != node {
				t.Fail()
			}
		}
		targets, err := graph.GetTargets(node.ID())
		if err != nil {
			t.Fail()
		}
		if len(targets) != 4 {
			t.Fail()
		}
		keys := []string{"B", "D", "G", "H"}
		for _, k := range keys {
			if n, existed := targets[StringID(k)]; !existed {
				t.Fail()
			} else {
				if n.ID() != StringID(k) {
					t.Fail()
				}
			}
		}

		es, err = graph.GetInEdges(node.ID())
		if err != nil {
			t.Fail()
		}
		for _, e := range es {
			if e.Target() != node {
				t.Fail()
			}
		}
		sources, err := graph.GetSources(node.ID())
		if err != nil {
			t.Fail()
		}
		if len(sources) != 4 {
			t.Fail()
		}
		for _, k := range keys {
			if n, existed := sources[StringID(k)]; !existed {
				t.Fail()
			} else {
				if n.ID() != StringID(k) {
					t.Fail()
				}
			}
		}
	}

	// can not replace node with new id, so replace node function is useless here...
	err = graph.ReplaceNode(StringID("A"), NewNode("X"))
	if err == nil {
		t.Fail()
	}

	err = graph.ReplaceEdge(StringID("A"), StringID("B"), 10)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	e, err := graph.GetEdge(StringID("A"), StringID("B"))
	if err != nil {
		t.Fail()
	}
	if e.Weight() != 10 {
		t.Fail()
	}
	fmt.Println(e)

	err = graph.DeleteEdge(StringID("A"), StringID("B"))
	if err != nil {
		t.Fail()
	}
	_, err = graph.GetEdge(StringID("A"), StringID("B"))
	if err == nil {
		t.Fail()
	}

	if deleted := graph.DeleteNode(StringID("A")); !deleted {
		t.Fail()
	}
	fmt.Println(graph)
	if _, existed := graph.GetNode(StringID("A")); existed {
		t.Fail()
	}

	graph.Clear()
	if !graph.Empty() {
		t.Fail()
	}
}
