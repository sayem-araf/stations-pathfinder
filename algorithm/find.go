package algorithm

// FindShortestPath returns a single shortest path between start and end.
// If no path exists, it returns nil.

// func (g is a receiver variable *Graph point to Graph and means that the function can read the graph, modify it and avoid copying large data.)

func (g *Graph) findShortestPath(start, end string, blocked map[string]bool) Path {
	if start == end {
		return Path{start}
	}

	visited := make(map[string]bool) // Makes a map of the visited stations
	prev := make(map[string]string)  // makes a map of the previous station
	queue := []string{start}         // makes a queue
	visited[start] = true            // adds start into the visited stations

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, neighbor := range g.Adj[current] {
			if blocked[neighbor] && neighbor != end {
				continue
			}
			if visited[neighbor] {
				continue
			}

			visited[neighbor] = true
			prev[neighbor] = current

			if neighbor == end {
				return reconstructPath(prev, start, end)
			}

			queue = append(queue, neighbor) // adds the neigbor to the end of the queue
		}
	}

	return nil
}

func (g *Graph) FindPaths(start, end string) []Path {
	var paths []Path

	blocked := make(map[string]bool)

	blocked["B"] = true

	for {
		path := g.findShortestPath(start, end, blocked)
		if path == nil {
			break
		}

		paths = append(paths, path)

		// Block intermediate stations (not start/end)
		for i := 1; i < len(path)-1; i++ {
			blocked[path[i]] = true
		}
	}

	return paths
}
