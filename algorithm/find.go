package algorithm

// FindShortestPath returns a single shortest path between start and end.
// If no path exists, it returns nil.

// func (g is a receiver variable *Graph point to Graph and means that the function can read the graph, modify it and avoid copying large data.)

func (g *Graph) FindShortestPath(start, end string, blocked map[string]bool) Path {
	if start == end { // if start is end
		return Path{start} // returns path
	}

	visited := make(map[string]bool) // Makes a map of the visited stations
	prev := make(map[string]string)  // makes a map of the previous station
	queue := []string{start}         // makes a queue
	visited[start] = true            // adds start into the visited stations

	for len(queue) > 0 { // BFS loop
		current := queue[0] // first in first out queue, queue can be A, B, C and the current is A
		queue = queue[1:]   // while the queue is B, C

		for _, neighbor := range g.Adj[current] { // checks every station connected to current
			if blocked[neighbor] && neighbor != end { // if neigbor is blocked and its not the final station then skip
				continue
			}
			if visited[neighbor] { // skips already visited
				continue
			}

			visited[neighbor] = true // marks as visited
			prev[neighbor] = current // records how we got there prev["B"] = "A"

			if neighbor == end { // Checks if we found the end
				return reconstructPath(prev, start, end)
			}

			queue = append(queue, neighbor) // adds the neigbor to the end of the queue
		}
	}

	return nil
}

func (g *Graph) FindPaths(start, end string) []Path { // belong to graph takes stat and end and returns a slice
	var paths []Path // creates a container for all paths

	blocked := make(map[string]bool) // creates a blocked map

	// force blockking blocked["B"] = true

	for {
		path := g.FindShortestPath(start, end, blocked) // Finds shortes path while respecting blocked paths
		if path == nil {                                // exit if loop if no paths
			break
		}

		paths = append(paths, path) // save the found path

		// Block intermediate stations (not start/end)
		for i := 1; i < len(path)-1; i++ {
			blocked[path[i]] = true
		}
	}

	return paths // retuns path
}
