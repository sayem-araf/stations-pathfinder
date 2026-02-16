package main

import (
	"fmt"
	"os"

	"gitea.kood.tech/sayemaraf/pathfinder/algorithm"
)

func main() {

	// Reads CLI arguments and skips the program name
	args := os.Args[1:]

	// Checks that four arguments are used
	ValidateArgs(args)

	filePath := args[0]
	startStation := args[1]
	endStation := args[2]
	trains := args[3]

	// Stores the file
	fileContent := MustReadFile(filePath)

	lines := NormalizeInput(fileContent)
	
	// Stores the data from stations and connections sections
	stations, connections := ParseMap(lines)

	// stationsMap stores stations by name and stationExist is used for booelan checking
	stationsMap, stationExists := BuildStationMaps(stations)

	ValidateSections(stationsMap, connections)
	ValidateStationNames(stationsMap)
	ValidateCoordinates(stationsMap)

	seenRoutes := make(map[string]bool)
	
	for _, c := range connections {
   		a, b := c[0], c[1]
    	ValidateConnection(stationsMap, a, b)
    	ValidateRoute(stationsMap, a, b, seenRoutes)
	}

	ValidateStations(startStation, endStation, stationExists)

	// Builds the graph from stations and connections
	graph := algorithm.NewGraph(stations, connections)

	ValidatePathExists(startStation, endStation, graph)

	// Finds all valid paths
	paths := graph.FindPaths(startStation, endStation)

	// Converts and validates number of trains
	numTrains := ValidateTrains(trains)

	fmt.Println(numTrains)
	fmt.Println(paths)
}