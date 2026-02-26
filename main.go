package main

import (
	"flag"
	"os"

	"gitea.kood.tech/sayemaraf/pathfinder/algorithm"
	"gitea.kood.tech/sayemaraf/pathfinder/parser"
	"gitea.kood.tech/sayemaraf/pathfinder/web"
)

func main() {
	webMode := flag.Bool("w", false, "starts web server")
	flag.Parse()

	if *webMode {
		web.Start()
		return
	}

	// Use flag.Args() so it works correctly whether or not -w is passed
	args := flag.Args()

	// Checks that four arguments are used
	parser.ValidateArgs(args)

	filePath := args[0]
	startStation := args[1]
	endStation := args[2]
	trains := args[3]

	// Stores the file
	fileContent := parser.MustReadFile(filePath)
	lines := parser.NormalizeInput(fileContent)

	// Stores the data from stations and connections sections
	stations, connections := parser.ParseMap(lines)

	// stationsMap stores stations by name and stationExists is used for boolean checking
	stationsMap, stationExists := parser.BuildStationMaps(stations)

	parser.ValidateSections(stationsMap, connections)
	parser.ValidateStationNames(stationsMap)
	parser.ValidateCoordinates(stationsMap)

	seenRoutes := make(map[string]bool)
	for _, c := range connections {
		a, b := c[0], c[1]
		parser.ValidateRoute(stationsMap, a, b, seenRoutes)
	}

	parser.ValidateStations(startStation, endStation, stationExists)

	// Builds the graph from stations and connections
	graph := algorithm.NewGraph(stations, connections)

	parser.ValidatePathExists(startStation, endStation, graph)

	// Converts and validates number of trains
	numTrains := parser.ValidateTrains(trains)

	// Finds all valid paths
	paths := graph.FindPaths(startStation, endStation)

	// Run the scheduler and print train movements
	scheduler := algorithm.NewScheduler(paths, numTrains, startStation, endStation)
	scheduler.Run()

	os.Exit(0)
}