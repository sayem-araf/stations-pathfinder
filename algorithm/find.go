package algorithm
// This file implements the pathfinding algorithm to find one or more shortest paths between two stations in a graph, while minimizing shared edges/stations.
import "sort"
// FindShortestPath uses a breadth-first search (BFS) to find the shortest path from start to end, while respecting blocked stations. It returns a Path if found, or nil if no path exists.
func (g *Graph) FindShortestPath(start, end string, blocked map[string]bool) Path {
	if start == end {
		return Path{start}
	}
// BFS initialization
	visited := make(map[string]bool)
	prev := make(map[string]string)
	queue := []string{start}
	visited[start] = true
// BFS loop
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
// Explore neighbors
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
			queue = append(queue, neighbor)
		}
	}
	return nil
}
// FindPaths finds one or more shortest paths between start and end, while attempting to minimize shared edges/stations. It returns a slice of Paths ordered from shortest to longest. If no path exists, an empty slice is returned.
func (g *Graph) FindPaths(start, end string) []Path {
	seen := make(map[string]bool)
	var candidates []Path
// First, find all paths that start with each of the first hops from the start station
	firstHops := g.Adj[start]
// For each first hop, find a path and add it to candidates if it's unique
	for _, firstHop := range firstHops {
		path := g.bfsFrom(start, end, firstHop, nil)
		if path != nil {
			key := pathKey(path)
			if !seen[key] {
				seen[key] = true
				candidates = append(candidates, path)
			}
		}
	}
// Next, for each first hop, block the stations used by paths from other first hops and try to find a new path
	for i, firstHop := range firstHops {
		blocked := make(map[string]bool)
		for j, otherHop := range firstHops {
			if i == j {
				continue
			}// Find the path for this other first hop and block its stations
			path := g.bfsFrom(start, end, otherHop, nil)
			if path != nil {
				for k := 1; k < len(path)-1; k++ {
					blocked[path[k]] = true
				}
			}
		}// Now try to find a path from this first hop while blocking the stations used by other first hops
		path := g.bfsFrom(start, end, firstHop, blocked)
		if path != nil {
			key := pathKey(path)
			if !seen[key] {
				seen[key] = true
				candidates = append(candidates, path)
			}
		}
	}
// Finally, try all pairs of first hops to find disjoint paths
	for i := 0; i < len(firstHops); i++ {
		for j := i + 1; j < len(firstHops); j++ {
			p1 := g.bfsFrom(start, end, firstHops[i], nil)
			if p1 == nil {
				continue
			} // Block the stations used by p1 (except start and end) and try to find a path from the other first hop
			blocked := make(map[string]bool)
			for k := 1; k < len(p1)-1; k++ {
				blocked[p1[k]] = true
			} // Now try to find a path from the other first hop while blocking the stations used by p1
			p2 := g.bfsFrom(start, end, firstHops[j], blocked)
			if p2 != nil {
				key := pathKey(p2)
				if !seen[key] {
					seen[key] = true
					candidates = append(candidates, p2)
				}
			}
		}
	}
// Sort candidates by length and find the best disjoint set
	sort.Slice(candidates, func(i, j int) bool {
		return len(candidates[i]) < len(candidates[j])
	})
// Find the best set of disjoint paths from candidates
	best := g.findBestDisjointSet(candidates)
// If no disjoint paths found, return the single shortest path
	if len(best) == 0 {
		p := g.FindShortestPath(start, end, nil)
		if p != nil {
			best = []Path{p}
		}
	}
// Sort the best paths by length before returning
	sort.Slice(best, func(i, j int) bool {
		return len(best[i]) < len(best[j])
	})

	return best
}
// bfsFrom performs a BFS starting from the firstHop station, while respecting blocked stations. It returns a Path if found, or nil if no path exists.
func (g *Graph) bfsFrom(start, end, firstHop string, blocked map[string]bool) Path {
	if firstHop == end {
		return Path{start, end}
	}
// BFS initialization
	visited := make(map[string]bool)
	prev := make(map[string]string)
	visited[start] = true
	visited[firstHop] = true
	prev[firstHop] = start
	queue := []string{firstHop}
// BFS loop
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
// Explore neighbors
		for _, neighbor := range g.Adj[current] {
			if visited[neighbor] {
				continue
			}
			if blocked[neighbor] && neighbor != end {
				continue
			}
			visited[neighbor] = true
			prev[neighbor] = current
			if neighbor == end {
				return reconstructPath(prev, start, end)
			}
			queue = append(queue, neighbor)
		}
	}
	return nil
}
//	pathKey generates a unique string key for a given Path, used for deduplication in the seen map.		
func pathKey(p Path) string {
	result := ""
	for _, s := range p {
		result += s + "|"
	}
	return result
}
//	findBestDisjointSet takes a list of candidate paths and finds the largest subset of paths that are disjoint (i.e., do not share any stations except start and end). It uses a backtracking approach to explore all combinations of paths and keeps track of the best set found.
func (g *Graph) findBestDisjointSet(candidates []Path) []Path {
	var best []Path
// trySelect is a recursive function that tries to select paths from candidates while ensuring they are disjoint. It keeps track of the currently selected paths and the stations that are already used by those paths.
	var trySelect func(idx int, selected []Path, usedStations map[string]bool)
	trySelect = func(idx int, selected []Path, usedStations map[string]bool) {
		if len(selected) > len(best) {
			best = make([]Path, len(selected))
			copy(best, selected)
		}
		if idx >= len(candidates) {
			return
		}
		if len(selected)+len(candidates)-idx <= len(best) {
			return
		}
// Try to include the current candidate path if it doesn't conflict with already used stations
		path := candidates[idx]
		conflict := false
		for i := 1; i < len(path)-1; i++ {
			if usedStations[path[i]] {
				conflict = true
				break
			}
		}
// If no conflict, include this path and mark its stations as used
		if !conflict {
			newUsed := make(map[string]bool)
			for k, v := range usedStations {
				newUsed[k] = v
			}
			for i := 1; i < len(path)-1; i++ {
				newUsed[path[i]] = true
			}
			trySelect(idx+1, append(selected, path), newUsed)
		}
// Also try without including the current candidate path
		trySelect(idx+1, selected, usedStations)
	}
// Start the recursive search with an empty selection and no used stations
	trySelect(0, []Path{}, make(map[string]bool))
	return best
}
//	pathsEqual is a helper function that checks if two paths are equal by comparing their lengths and corresponding stations. It returns true if the paths are identical, and false otherwise.
func pathsEqual(p1, p2 Path) bool {
	if len(p1) != len(p2) {
		return false
	}
	for i := range p1 {// Compare each station in the paths
		if p1[i] != p2[i] {
			return false
		}
	}
	return true
}