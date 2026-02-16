package parser

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"gitea.kood.tech/sayemaraf/pathfinder/algorithm"
)

// Helper function for printing error messages
func PrintError(msg string) {
	fmt.Fprintln(os.Stderr, "Error:", msg)
	os.Exit(1)
}

// Reads the input file and prints error if fails
func MustReadFile(filePath string) []byte {
	data, err := os.ReadFile(filePath)
	if err != nil {
		PrintError(fmt.Sprintf("Failed to read file: %s", err))
	}
	return data
}

// Checks if too many or too few arguments are used
func ValidateArgs(args []string) {
	if len(args) < 4 {
		PrintError("Too few arguments")
	}

	if len(args) > 4 {
		PrintError("Too many arguments")
	}
}

func ValidateNewStation(stations map[string]*algorithm.Station, name string) {
	if _, exists := stations[name]; exists {
		PrintError("Duplicate station name: " + name)
	}
}

func ValidateConnection(stations map[string]*algorithm.Station, a, b string) {
	if _, ok := stations[a]; !ok {
		PrintError("Unknown station in connection: " + a)
	}
	if _, ok := stations[b]; !ok {
		PrintError("Unknown station in connection: " + b)
	}
}

func ValidateRoute(stations map[string]*algorithm.Station, a, b string, seenRoutes map[string]bool) {
	if _, ok := stations[a]; !ok {
		PrintError("Unknown station in connection: " + a)
	}
	if _, ok := stations[b]; !ok {
		PrintError("Unknown station in connection: " + b)
	}

	if a == b {
		PrintError("Connection cannot link a station to itself: " + a)
	}

	key := a
	if a > b {
		key = b + "-" + a
	} else {
		key = a + "-" + b
	}

	if seenRoutes[key] {
		PrintError("Duplicate connection between " + a + " and " + b)
	}

	seenRoutes[key] = true
}

// Checks if stations in the input arguments are valid
func ValidateStations(start string, end string, stations map[string]bool) {
	if start == end {
		PrintError("Start station can't be the same as end station")
	}

	if !stations[start] {
		PrintError("Invalid start station")
	}

	if !stations[end] {
		PrintError("Invalid end station")
	}
}

// Checks if number of trains is valid
func ValidateTrains(trains string) int {
	numTrains, err := strconv.Atoi(trains)
	if err != nil || numTrains < 1 {
		PrintError("Number of trains must be at least 1")
	}
	return numTrains
}

// Checks map for invalid or duplicate station names
func ValidateStationNames(stations map[string]*algorithm.Station) {
	validName := regexp.MustCompile(`^[A-Za-z0-9_]+$`)

	for name := range stations {
		if !validName.MatchString(name) {
			PrintError("Invalid station name: " + name)
		}
	}

}

// Checks if map has the required sections and valid number of stations
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

// Checks if station coordinates are valid positive integers and also prints error if stations share the same coordinates
func ValidateCoordinates(stations map[string]*algorithm.Station) {
	coords := make(map[string]string)

	for name, s := range stations {
		if s.X < 0 || s.Y < 0 {
			PrintError(fmt.Sprintf("Station %s has invalid coordinates (%d, %d)", name, s.X, s.Y))
		}

		key := fmt.Sprintf("%d_%d", s.X, s.Y)
		if existing, exists := coords[key]; exists {
			PrintError(fmt.Sprintf("Stations %s and %s share the same coordinates (%d, %d)", existing, name, s.X, s.Y))
		}
		coords[key] = name
	}
}

// Checks if the path exists between stations
func ValidatePathExists(start, end string, g *algorithm.Graph) {
	if g.FindShortestPath(start, end, map[string]bool{}) == nil {
		PrintError(fmt.Sprintf("No path exists between %s and %s", start, end))
	}
}

// Converts coordinate values from string to integer
func MustParseInt(value string, stationName string, coord string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		PrintError(fmt.Sprintf("Station %s has invalid %s coordinate: %q", stationName, coord, value))
	}
	return i
}
