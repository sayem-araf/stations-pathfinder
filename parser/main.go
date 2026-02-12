package main

import (
	"fmt"
	"os"

	"gitea.kood.tech/sayemaraf/pathfinder/algorithm"
)

func main() {

	args := os.Args[1:]

	ValidateArgs(args)

	filePath := args[0]
	startStation := args[1]
	endStation := args[2]
	trains := args[3]

	fileContent := MustReadFile(filePath)

	lines := NormalizeInput(fileContent)
	
	stations, connections := ParseMap(lines)

	stationsMap, stationExists := BuildStationMaps(stations)

	ValidateSections(stationsMap, connections)
	ValidateStationMap(stationsMap)
	ValidateCoordinates(stationsMap)
	ValidateStations(startStation, endStation, stationExists)
	ValidateConnections(connections, stationExists)

	graph := algorithm.NewGraph(stations, connections)

	ValidatePathExists(startStation, endStation, graph)

	paths := graph.FindPaths(startStation, endStation)

	numTrains := ValidateTrains(trains)

	fmt.Println(numTrains)
	fmt.Println(paths)
}