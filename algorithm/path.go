package algorithm

// FindPaths finds one or more shortest paths between start and end.
//
// Guarantees:
// - Returned paths are valid in the graph
// - All paths start at `start` and end at `end`
// - Paths are ordered from shortest to longest
// - Paths attempt to minimize shared edges/stations
//
// If no path exists, an empty slice is returned.

func reconstructPath(prev map[string]string, start, end string) Path {
	var path []string
	for at := end; at != ""; at = prev[at] {
		path = append([]string{at}, path...)
		if at == start {
			break
		}
	}
	return path
}
