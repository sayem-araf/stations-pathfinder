package main
// This file contains the main function that serves as the entry point for the pathfinding application.
import (
	"flag"
	"os"

	"gitea.kood.tech/sayemaraf/pathfinder/algorithm"
	"gitea.kood.tech/sayemaraf/pathfinder/parser"
	"gitea.kood.tech/sayemaraf/pathfinder/web"
)
// The main function parses command-line flags to determine whether to start the web server or run the pathfinding algorithm in command-line mode. 
// If the -w flag is provided, it starts the web server. Otherwise, it expects command-line arguments for the file path, start station, end station, and number of trains. 
// It reads and processes the input file to build the graph, validates the input data, finds paths between the specified stations, and runs the scheduler to simulate train movements along those paths.
func main() {
	webMode := flag.Bool("w", false, "starts web server")
	flag.Parse()

	if *webMode {
		web.Start()
		return
	}
// Validate that the correct number of command-line arguments are provided (file path, start station, end station, number of trains).
	args := flag.Args()
	parser.ValidateArgs(args)

	filePath := args[0]
	startStation := args[1]
	endStation := args[2]
	trains := args[3]
// Read the input file, normalize the content, and parse the stations and connections to build the graph.
	fileContent := parser.MustReadFile(filePath)
	lines := parser.NormalizeInput(fileContent)

	stations, connections := parser.ParseMap(lines)

	stationsMap, stationExists := parser.BuildStationMaps(stations)
	parser.ValidateSections(stationsMap, connections)
	parser.ValidateStationNames(stationsMap)
	parser.ValidateCoordinates(stationsMap)
// Validate that all connections reference valid stations and that there are no duplicate routes.
	seenRoutes := make(map[string]bool)
	for _, c := range connections {
		a, b := c[0], c[1]
		parser.ValidateRoute(stationsMap, a, b, seenRoutes)
	}
// Validate that the start and end stations exist in the station map and are not the same, and that the number of trains is a valid integer.
	parser.ValidateStations(startStation, endStation, stationExists)

	graph := algorithm.NewGraph(stations, connections)
	parser.ValidatePathExists(startStation, endStation, graph)

	numTrains := parser.ValidateTrains(trains)

	paths := graph.FindPaths(startStation, endStation)

// Create a new Scheduler instance with the found paths, number of trains, start station, and end station, and run the scheduling simulation.
	scheduler := algorithm.NewScheduler(paths, numTrains, startStation, endStation)
	scheduler.Run()

	os.Exit(0)
}
