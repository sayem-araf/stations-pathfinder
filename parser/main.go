package main

import (
	"fmt"
	"os"
	"strconv"

	"gitea.kood.tech/sayemaraf/pathfinder/algorithm"
)

func main() {

	args := os.Args[1:]

	if len(args) != 4 {
		PrintError("Invalid arguments")
		os.Exit(1)
	}

	filePath := args[0]
	startStation := args[1]
	endStation := args[2]
	trains := args[3]

	fileContent, err := os.ReadFile(filePath)

	if err != nil {
		PrintError("Could not read file: " + filePath)
		os.Exit(1)
	}

	lines := NormalizeInput(fileContent)
	
	stations, connections := ParseMap(lines)

	stationExists := make(map[string]bool)
	for _, s := range stations {
    	stationExists[s.Name] = true
	}

	if !stationExists[startStation] {
    PrintError("Start station does not exist")
	}

	if !stationExists[endStation] {
    PrintError("End station does not exist")
	}

	if startStation == endStation {
		PrintError("Start station can't be the same as end station")
	}

	graph := algorithm.NewGraph(stations, connections)
	paths := graph.FindPaths(startStation, endStation)

	for i, path := range paths {
		fmt.Printf("Path %d: %v\n", i+1, path)
	}

	numTrains, err := strconv.Atoi(trains)
	if err != nil || numTrains <= 0 {
    	PrintError("Invalid number of trains")
	}
}

func PrintError(msg string) {
	fmt.Fprintln(os.Stderr, "Error:", msg)
	os.Exit(1)
}