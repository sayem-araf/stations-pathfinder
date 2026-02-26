package algorithm

// FindShortestPath returns a single shortest path between start and end.
// If no path exists, it returns nil.
func (g *Graph) FindShortestPath(start, end string, blocked map[string]bool) Path {
	if start == end {
		return Path{start}
	}

	visited := make(map[string]bool)
	prev := make(map[string]string)
	queue := []string{start}
	visited[start] = true

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

			queue = append(queue, neighbor)
		}
	}

	return nil
}

// FindPaths finds multiple disjoint paths between start and end
func (g *Graph) FindPaths(start, end string) []Path {
	neighbors := g.Adj[start]
	
	// STEP 1: For each neighbor, find the shortest path
	initialPaths := make([]Path, 0)
	for _, neighbor := range neighbors {
		path := g.findPathThroughNeighbor(start, end, neighbor)
		if path != nil {
			initialPaths = append(initialPaths, path)
		}
	}
	
	if len(initialPaths) == 0 {
		return nil
	}
	
	// STEP 2: Try to resolve conflicts by finding alternative paths
	// Keep trying until we have maximum disjoint paths
	bestPaths := g.findMaximumDisjointPaths(start, end, neighbors, initialPaths)
	
	// STEP 3: Return best set of paths
	if len(bestPaths) >= 2 {
		return bestPaths
	}
	
	// For maps where paths must overlap, return all paths
	if len(initialPaths) > 0 {
		return initialPaths
	}
	
	return nil
}

// findMaximumDisjointPaths tries to find the maximum number of disjoint paths
func (g *Graph) findMaximumDisjointPaths(start, end string, neighbors []string, initialPaths []Path) []Path {
	// Check if initial paths are disjoint
	disjoint := g.selectDisjointPaths(initialPaths)
	
	// If we already have as many disjoint paths as neighbors, we're done
	if len(disjoint) == len(neighbors) {
		return disjoint
	}
	
	// If we're missing disjoint paths, try to find alternative routes
	// that avoid conflicts
	if len(disjoint) < len(neighbors) {
		// Build a conflict map showing which stations are shared
		conflicts := g.findConflicts(initialPaths)
		
		// For each path that has conflicts, try to find an alternative
		improvedPaths := make([]Path, len(initialPaths))
		copy(improvedPaths, initialPaths)
		
		for i, path := range initialPaths {
			neighbor := path[1] // First hop
			
			// Check if this path has conflicts
			hasConflict := false
			for j := 1; j < len(path)-1; j++ {
				if len(conflicts[path[j]]) > 1 {
					hasConflict = true
					break
				}
			}
			
			if hasConflict {
				// Build blocked set from OTHER paths
				blocked := make(map[string]bool)
				for k, otherPath := range improvedPaths {
					if k != i {
						for j := 1; j < len(otherPath)-1; j++ {
							blocked[otherPath[j]] = true
						}
					}
				}
				
				// Try to find alternative path avoiding blocked stations
				altPath := g.findPathThroughNeighborAvoiding(start, end, neighbor, blocked)
				if altPath != nil {
					improvedPaths[i] = altPath
				}
			}
		}
		
		// Check if improved paths are better
		improvedDisjoint := g.selectDisjointPaths(improvedPaths)
		if len(improvedDisjoint) > len(disjoint) {
			return improvedDisjoint
		}
	}
	
	return disjoint
}

// findConflicts returns a map of station -> list of path indices using it
func (g *Graph) findConflicts(paths []Path) map[string][]int {
	conflicts := make(map[string][]int)
	for i, path := range paths {
		for j := 1; j < len(path)-1; j++ {
			station := path[j]
			conflicts[station] = append(conflicts[station], i)
		}
	}
	return conflicts
}

// selectDisjointPaths greedily selects the maximum set of disjoint paths
func (g *Graph) selectDisjointPaths(paths []Path) []Path {
	if len(paths) == 0 {
		return nil
	}
	
	var selected []Path
	blocked := make(map[string]bool)
	
	// Try each starting order to find the best combination
	var bestSelection []Path
	
	for startIdx := 0; startIdx < len(paths); startIdx++ {
		selected = []Path{}
		blocked = make(map[string]bool)
		
		for i := 0; i < len(paths); i++ {
			idx := (startIdx + i) % len(paths)
			path := paths[idx]
			
			// Check if path conflicts with already selected paths
			hasConflict := false
			for j := 1; j < len(path)-1; j++ {
				if blocked[path[j]] {
					hasConflict = true
					break
				}
			}
			
			if !hasConflict {
				selected = append(selected, path)
				// Block intermediate stations
				for j := 1; j < len(path)-1; j++ {
					blocked[path[j]] = true
				}
			}
		}
		
		if len(selected) > len(bestSelection) {
			bestSelection = make([]Path, len(selected))
			copy(bestSelection, selected)
		}
	}
	
	return bestSelection
}

// findPathThroughNeighbor finds a path that uses the specified neighbor as the first hop
func (g *Graph) findPathThroughNeighbor(start, end, firstHop string) Path {
	if firstHop == end {
		return Path{start, end}
	}
	
	visited := make(map[string]bool)
	prev := make(map[string]string)
	queue := []string{firstHop}
	visited[start] = true
	visited[firstHop] = true
	prev[firstHop] = start
	
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		
		for _, neighbor := range g.Adj[current] {
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

// findPathThroughNeighborAvoiding finds a path through neighbor while avoiding blocked stations
func (g *Graph) findPathThroughNeighborAvoiding(start, end, firstHop string, blocked map[string]bool) Path {
	if firstHop == end {
		return Path{start, end}
	}
	
	visited := make(map[string]bool)
	prev := make(map[string]string)
	queue := []string{firstHop}
	visited[start] = true
	visited[firstHop] = true
	prev[firstHop] = start
	
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		
		for _, neighbor := range g.Adj[current] {
			// Skip blocked stations (unless it's the destination)
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

// pathsEqual checks if two paths are identical
func pathsEqual(p1, p2 Path) bool {
	if len(p1) != len(p2) {
		return false
	}
	for i := range p1 {
		if p1[i] != p2[i] {
			return false
		}
	}
	return true
}
