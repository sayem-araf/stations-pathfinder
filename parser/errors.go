package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"gitea.kood.tech/sayemaraf/pathfinder/algorithm"
)

func PrintError(msg string) {
	fmt.Fprintln(os.Stderr, "Error:", msg)
	os.Exit(1)
}

func BuildStationMaps(stations []*algorithm.Station) (map[string]*algorithm.Station, map[string]bool) {
    stationsMap := make(map[string]*algorithm.Station)
    stationExists := make(map[string]bool)

    for _, s := range stations {
        if _, exists := stationsMap[s.Name]; exists {
            PrintError("Duplicate station name: " + s.Name)
        }
        stationsMap[s.Name] = s
        stationExists[s.Name] = true
    }

    return stationsMap, stationExists
}

func MustReadFile(filePath string) []byte {
	data, err := os.ReadFile(filePath)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to read file: %s", err))
	}
	return data
}

func ValidateArgs(args []string) {
	if len(args) < 4 {
		PrintError("Too few arguments")
	}

	if len(args) > 4 {
		PrintError("Too many arguments")
	}
}

func ValidateStations(start string, end string, stations map[string]bool) {
	if start == end {
		PrintError("Start station can't be the same as end station")
	}
	
	if !stations[start]{
		PrintError("Invalid start station")
	}

	if !stations[end]{
		PrintError("Invalid end station")
	}
}

func ValidateConnections(connections [][2]string, stations map[string]bool) {
    seen := make(map[string]bool)

    for _, c := range connections {
        a, b := c[0], c[1]

        if a == b {
            PrintError("Duplicate connection between " + a + " and " + b) // catch self-loop
        }

        if !stations[a] {
            PrintError("Unknown station in connection: " + a)
        }
        if !stations[b] {
            PrintError("Unknown station in connection: " + b)
        }

        // Generate a canonical key (sorted) so "A-B" and "B-A" are treated the same
        key := a
        if a > b {
            key = b + "-" + a
        } else {
            key = a + "-" + b
        }

        if seen[key] {
            PrintError("Duplicate connection between " + a + " and " + b)
        }

        seen[key] = true
    }
}

func ValidateTrains(trains string) int {
	numTrains, err := strconv.Atoi(trains)
	if err != nil || numTrains < 1 {
		PrintError("Number of trains must be at least 1")
	}
	return numTrains
}

func ValidateStationMap(stations map[string]*algorithm.Station) {
	validName := regexp.MustCompile(`^[A-Za-z0-9_]+$`)
	seenNames := make(map[string]bool)

	for name := range stations {
		if !validName.MatchString(name) {
			PrintError("Invalid station name: " + name)
		}

		if seenNames[name] {
			PrintError("Duplicate station name: " + name)
		}
		seenNames[name] = true
	}

}

func ValidateSections(stations map[string]*algorithm.Station, connections [][2]string) {
	if len(stations) == 0 {
		PrintError("Map does not contain a stations section")
	}
	if len(connections) == 0 {
		PrintError("Map does not contain a connections section")
	}
	if len(stations) > 10000 {
		PrintError("Map contains more than 10000 stations")
	}
}

func ValidateCoordinates(stations map[string]*algorithm.Station) {
	coords := make(map[string]string)

	for name, s := range stations {
		if s.X < 0 || s.Y < 0 {
			PrintError(fmt.Sprintf("Error: Station %s has invalid coordinates (%d, %d)", name, s.X, s.Y))
		}

		key := fmt.Sprintf("%d_%d", s.X, s.Y)
		if existing, exists := coords[key]; exists {
			PrintError(fmt.Sprintf("Error: Stations %s and %s share the same coordinates (%d, %d)", existing, name, s.X, s.Y))
		}
		coords[key] = name
	}
}

func ValidatePathExists(start, end string, g *algorithm.Graph) {
    if g.FindShortestPath(start, end, map[string]bool{}) == nil {
        PrintError(fmt.Sprintf("No path exists between %s and %s", start, end))
    }
}

func MustParseInt(value string, stationName string, coord string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		PrintError(fmt.Sprintf("Station %s has invalid %s coordinate: %q", stationName, coord, value))
	}
	return i
}