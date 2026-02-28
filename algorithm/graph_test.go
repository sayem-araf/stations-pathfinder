package algorithm
// This file contains tests for the graph and pathfinding algorithm. 
// It uses the testing package to define test cases that verify the correctness of the FindPaths function. 
// The test sets up a simple graph with stations and connections, then calls FindPaths to find paths between two stations and prints the results for verification.
import (
	"fmt"
	"testing"
)
// This test verifies that the FindPaths function correctly identifies the shortest path between two stations in a simple graph. 
// It sets up a graph with four stations (A, B, C, D) and connections between them, then calls FindPaths to find paths from station A to station C. 
// Finally, it prints the found paths for verification.
func TestFindShortestPath(t *testing.T) {
	stations := []*Station{
		{Name: "A"},
		{Name: "B"},
		{Name: "C"},
		{Name: "D"},
	}
// The connections between the stations are defined as follows:
// - A is connected to B and D
// - B is connected to A and C
// - C is connected to B and D
// - D is connected to A and C	
	connections := [][2]string{
		{"A", "B"},
		{"B", "C"},
		{"A", "D"},
		{"D", "C"},
	}
// Create a new graph with the defined stations and connections
	g := NewGraph(stations, connections)
// Call FindPaths to find paths from station A to station C
	path := g.FindPaths("A", "C")
// Print the found paths for verification
	fmt.Println(path)
}
