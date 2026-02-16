package main

import (
	"flag"
	"fmt"
	"os"

	"gitea.kood.tech/sayemaraf/pathfinder/algorithm"
	"gitea.kood.tech/sayemaraf/pathfinder/parser"
)

func main() {
	webMode := flag.Bool("w", false, "starts web server")
	flag.Parse()

	if *webMode {
		// Run web server
		fmt.Println("Starting web server...")
		fmt.Println("Run: go run web/server.go")
		fmt.Println("Or build: go build -o pathfinder-web web/server.go && ./pathfinder-web")
		os.Exit(0)
	}

	// CLI Mode - same as parser/main.go
	args := flag.Args() // Get remaining args after flags

	parser.ValidateArgs(args)

	filePath := args[0]
	startStation := args[1]
	endStation := args[2]
	trains := args[3]

	fileContent := parser.MustReadFile(filePath)
	lines := parser.NormalizeInput(fileContent)
	stations, connections := parser.ParseMap(lines)
	stationsMap, stationExists := parser.BuildStationMaps(stations)

	parser.ValidateSections(stationsMap, connections)
	parser.ValidateStationMap(stationsMap)
	parser.ValidateCoordinates(stationsMap)
	parser.ValidateStations(startStation, endStation, stationExists)
	parser.ValidateConnections(connections, stationExists)

	graph := algorithm.NewGraph(stations, connections)
	parser.ValidatePathExists(startStation, endStation, graph)

	paths := graph.FindPaths(startStation, endStation)
	numTrains := parser.ValidateTrains(trains)

	// Simulate train movements (THE MISSING PIECE!)
	parser.SimulateTrainMovement(paths, numTrains, startStation, endStation)
}
