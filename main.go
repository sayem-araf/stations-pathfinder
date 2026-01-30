package main

import (
	"fmt"
	"os"
	"strings"
)

type station struct {
	Name		string
	Connections []string
}

type graph struct {
	Stations map[string]*station
}

func main() {

	args := os.Args[1:]

	if len(args) != 4 {
		fmt.Println("Error: Invalid arguments")
		os.Exit(1)
	}

	filePath := args[0]
	startStation := args[1]
	endStation := args[2]
	trains := args[3]

	fileContent, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Println("Error: Could not read file:", err)
		os.Exit(1)
	}

	lines := normalizeInput(fileContent)

	_ = startStation
	_ = endStation
	_ = trains

	stations := make(map[string]*station)
	section := ""

	for _, line := range lines {
		if line == "stations:" {
			section = "stations"
			continue
		}

		if line == "connections:" {
			section = "connections"
			continue
		}

		if section == "stations" {
			parts := strings.Split(line, ",")
			name := strings.TrimSpace(parts[0])

			if _, exists := stations[name]; exists {
				fmt.Fprintln(os.Stderr, "Error: Duplicate station name", name)
        	os.Exit(1)
    		}
			
			
			stations[name] = &station{Name: name, Connections: []string{}}
			fmt.Println("Added station:", name)
		}
	}

}