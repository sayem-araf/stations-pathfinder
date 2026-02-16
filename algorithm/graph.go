package algorithm

// Station represent a single node in the rail network
// It is assumed to be valdi and unique (validate upstream)
type Station struct {
	Name string
	X    int // x and y are optional if some kinda grid based maps is wanted.
	Y    int
}

// graph represents the rail network as an adjacency structure
// It is safe for traversal and may contain cycles

type Graph struct {
	Stations map[string]*Station // stations is the map of Station struct
	Adj      map[string][]string // adjacency list: station name -> neighbor names
}

// Path represents an ordered route from start to end
// start and end stations are included

type Path []string

// NewGraph constructs a Graph from validated stations and connections.
//
// Guarantees:
// - All stations exist
// - All connections are bidirectional
// - No duplicate edges are created

func NewGraph(Stations []*Station, connections [][2]string) *Graph {
	g := &Graph{ // receives data from Graph struct
		Stations: make(map[string]*Station), // makes a map of Stations from graph data
		Adj:      make(map[string][]string), // makes a Adj map pf the station graph data
	}

	for _, s := range Stations { // Loops trough all stations one by one, index is ignored
		g.Stations[s.Name] = s     // Stores station inside the graph station map using its name as the key
		g.Adj[s.Name] = []string{} // Creates an empty list of neigbors for this station
	}

	for _, c := range connections { // loops trough all of the connections
		a, b := c[0], c[1]             // extracts the two station in connections a and b and unpacks the connection into two endpoints
		g.Adj[a] = append(g.Adj[a], b) // adds station b as a neigbor of station a, a -> b
		g.Adj[b] = append(g.Adj[b], a) // adds station a as a neigbor of station b, a -> b
	}

	return g // returns graph
}
