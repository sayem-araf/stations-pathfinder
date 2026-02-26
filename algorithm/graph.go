package algorithm

type Station struct {
	Name string // Must be Capital
	X    int    // Must be Capital
	Y    int    // Must be Capital
}

type Path []string

type Graph struct {
	Stations map[string]*Station
	Adj      map[string][]string
}

func NewGraph(stations []*Station, connections [][2]string) *Graph {
	g := &Graph{
		Stations: make(map[string]*Station),
		Adj:      make(map[string][]string),
	}
	for _, s := range stations {
		g.Stations[s.Name] = s
		g.Adj[s.Name] = []string{}
	}
	for _, c := range connections {
		g.Adj[c[0]] = append(g.Adj[c[0]], c[1])
		g.Adj[c[1]] = append(g.Adj[c[1]], c[0])
	}
	return g
}
