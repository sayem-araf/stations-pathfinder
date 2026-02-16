package parser

import (
	"fmt"
	"strings"

	"gitea.kood.tech/sayemaraf/pathfinder/algorithm"
)

// Removes comments, trims spaces and removes empty lines
func NormalizeInput(data []byte) []string {

	// Converts data into a slice split by newLines
	raw := strings.Split(string(data), "\n")
	var lines []string

	// Loops through lines in raw
	for _, line := range raw {
		// Checks line for # comments and then removes comment
		if idx := strings.Index(line, "#"); idx != -1 {
			line = line[:idx]
		}

		line = strings.TrimSpace(line)

		// Adds non-empty lines to slice
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines
}

// Builds maps
func BuildStationMaps(stations []*algorithm.Station) (map[string]*algorithm.Station, map[string]bool) {
	stationsMap := make(map[string]*algorithm.Station)
	stationExists := make(map[string]bool)

	for _, s := range stations {
		stationsMap[s.Name] = s
		stationExists[s.Name] = true
	}

	return stationsMap, stationExists
}

// Takes normalized lines and returns a slice of stations and a slice of connections
func ParseMap(lines []string) ([]*algorithm.Station, [][2]string) {

	section := ""                                   // tracks whether we are in "stations" or "connections"
	stations := make(map[string]*algorithm.Station) // temporary map to hold stations by name
	connections := [][2]string{}                    // slice of station pairs representing connections

	for _, line := range lines {
		switch line {
		case "stations:":
			section = "stations"
			continue
		case "connections:":
			section = "connections"
			continue
		}

		// Parses station lines
		if section == "stations" {
			parts := strings.Split(line, ",") // expected name,X,Y

			name := strings.TrimSpace(parts[0])
			X := strings.TrimSpace(parts[1])
			Y := strings.TrimSpace(parts[2])

			// Converts coordinates from string to int
			xi := MustParseInt(X, name, "X")
			yi := MustParseInt(Y, name, "Y")

			ValidateNewStation(stations, name)

			// Add the station to the map so it can be looked up by name
			stations[name] = &algorithm.Station{Name: name, X: xi, Y: yi}
			continue
		}

		// Parses connection lines
		if section == "connections" {
			parts := strings.Split(line, "-") // expected: stationA-stationB

			a := strings.TrimSpace(parts[0])
			b := strings.TrimSpace(parts[1])

			// Adds connection as a pair
			connections = append(connections, [2]string{a, b})
		}
	}

	// Creates an empty slice with capacity equal to the number of stations
	stationsSlice := make([]*algorithm.Station, 0, len(stations))

	// Adds stations from the map into the slice
	for _, s := range stations {
		stationsSlice = append(stationsSlice, s)
	}

	return stationsSlice, connections
}

// SimulateTrainMovement simulates and outputs train movements turn by turn
func SimulateTrainMovement(paths []algorithm.Path, numTrains int, start, end string) {
	type Train struct {
		ID       int
		Path     algorithm.Path
		Position int
	}

	// Create trains and assign them to paths
	trains := make([]*Train, numTrains)
	for i := 0; i < numTrains; i++ {
		pathIndex := i % len(paths)
		trains[i] = &Train{
			ID:       i + 1,
			Path:     paths[pathIndex],
			Position: 0,
		}
	}

	// Track which stations are currently occupied (excluding start and end)
	occupied := make(map[string]bool)

	// Simulate movement turn by turn
	for {
		var moves []string
		allFinished := true

		for _, train := range trains {
			// Check if train has reached the end
			if train.Position >= len(train.Path)-1 {
				continue
			}

			allFinished = false

			// Try to move train to next station
			nextPos := train.Position + 1
			nextStation := train.Path[nextPos]

			// Check if next station is available (not occupied)
			// Start and end stations can have unlimited trains
			canMove := nextStation == end || nextStation == start || !occupied[nextStation]

			if canMove {
				// Mark current station as free (if not start)
				if train.Position > 0 {
					currentStation := train.Path[train.Position]
					if currentStation != start {
						occupied[currentStation] = false
					}
				}

				// Move train
				train.Position = nextPos

				// Mark new station as occupied (if not end)
				if nextStation != end {
					occupied[nextStation] = true
				}

				// Record the move
				moves = append(moves, fmt.Sprintf("T%d-%s", train.ID, nextStation))
			}
		}

		// Print all moves for this turn
		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		}

		// Exit when all trains have finished
		if allFinished {
			break
		}
	}
}
