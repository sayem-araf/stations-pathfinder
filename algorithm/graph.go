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
	Stations map[string]*Station
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
	g := &Graph{
		Stations: make(map[string]*Station),
		Adj:      make(map[string][]string),
	}

	for _, s := range Stations {
		g.Stations[s.Name] = s
		g.Adj[s.Name] = []string{}
	}

	for _, c := range connections {
		a, b := c[0], c[1]
		g.Adj[a] = append(g.Adj[a], b)
		g.Adj[b] = append(g.Adj[b], a)
	}

	return g
}
