package main

import(
	"strings"

	"gitea.kood.tech/sayemaraf/pathfinder/algorithm"
)

func NormalizeInput(data []byte) []string {
	raw := strings.Split(string(data), "\n")
	var lines []string

	for _, line := range raw {
		if idx := strings.Index(line, "#"); idx != -1 {
			line = line[:idx]
		}

		line = strings.TrimSpace(line)

		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines
}

func ParseMap(lines []string) ([]*algorithm.Station, [][2]string) {
	section := ""
	stations := make(map[string]*algorithm.Station)
	connections := [][2]string{}
	
	for _, line := range lines {
		switch line {
		case "stations:":
			section = "stations"
			continue
		case "connections:":
			section = "connections"
			continue
		}

		if section == "stations" {
			parts := strings.Split(line, ",")

			name := strings.TrimSpace(parts[0])
			X := strings.TrimSpace(parts[1])
			Y := strings.TrimSpace(parts[2])

			xi := MustParseInt(X, name, "X")
			yi := MustParseInt(Y, name, "Y")

			stations[name] = &algorithm.Station{Name: name, X: xi, Y: yi}
    		continue
		}

		if section == "connections" {
			parts := strings.Split(line, "-")

			a := strings.TrimSpace(parts[0])
			b := strings.TrimSpace(parts[1])

			connections = append(connections, [2]string{a, b})
		}
	}
	stationsSlice := make([]*algorithm.Station, 0, len(stations))
	for _, s := range stations {
		stationsSlice = append(stationsSlice, s)
	}

	return stationsSlice, connections
}