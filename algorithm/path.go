package algorithm

/* FindPaths finds one or more shortest paths between start and end.
 Guarantees:
- Returned paths are valid in the graph
- All paths start at `start` and end at `end`
- Paths are ordered from shortest to longest
- Paths attempt to minimize shared edges/stations

If no path exists, an empty slice is returned.
*/
//
func reconstructPath(prev map[string]string, start, end string) Path { // prev map is a map of the breadcrumb trail
	var path []string                        // creates a empty slice
	for at := end; at != ""; at = prev[at] { // start from the end station, keep looping while at is not empty, move one step backward each iteration
		path = append([]string{at}, path...) // inserts at at the front of the path C -> B, C -> A, B, C
		if at == start {                     // once reached start stops walking backwards
			break
		}
	}
	return path // returns the path taken
}
