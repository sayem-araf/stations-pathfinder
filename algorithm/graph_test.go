package algorithm

import (
	"fmt"
	"testing"
)

func TestFindShortestPath(t *testing.T) {
	stations := []*Station{
		{Name: "A"},
		{Name: "B"},
		{Name: "C"},
		{Name: "D"},
	}

	connections := [][2]string{
		{"A", "B"},
		{"B", "C"},
		{"A", "D"},
		{"D", "C"},
	}

	g := NewGraph(stations, connections)

	path := g.FindPaths("A", "C")

	fmt.Println(path)
}
