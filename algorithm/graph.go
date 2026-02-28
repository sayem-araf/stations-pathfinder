// This file defines the graph structure and related types for the pathfinding algorithm.
package algorithm 
// Station represents a station in the graph with its name and coordinates.
type Station struct {
	Name string // Must be Capital
	X    int    // Must be Capital
	Y    int    // Must be Capital
}
// Path represents a sequence of station names from start to end.
type Path []string
// Graph represents the entire graph structure, including stations and their connections.
type Graph struct {
	Stations map[string]*Station
	Adj      map[string][]string
}
// NewGraph creates a new graph given a list of stations and their connections. 
// It initializes the Stations map and the adjacency list based on the provided data.
func NewGraph(stations []*Station, connections [][2]string) *Graph {
	g := &Graph{
		Stations: make(map[string]*Station),
		Adj:      make(map[string][]string),
	}// Populate the Stations map and initialize the adjacency list for each station
	for _, s := range stations {
		g.Stations[s.Name] = s
		g.Adj[s.Name] = []string{}
	}// Populate the adjacency list based on the connections between stations
	for _, c := range connections {
		g.Adj[c[0]] = append(g.Adj[c[0]], c[1])
		g.Adj[c[1]] = append(g.Adj[c[1]], c[0])
	}// Return the constructed graph
	return g
}
