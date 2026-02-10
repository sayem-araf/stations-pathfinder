package main

import(
	"strconv"
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
			if len(parts) < 3 {
				PrintError("Invalid station line: " + line)
			}

			name := strings.TrimSpace(parts[0])
			X := strings.TrimSpace(parts[1])
			Y := strings.TrimSpace(parts[2])

			xi, err := strconv.Atoi(X)
			if err != nil {
				PrintError("Invalid x coordinate for station: " + name)
			}

			yi, err := strconv.Atoi(Y)
				if err != nil {
				PrintError("Invalid y coordinate for station: " + name)
			}

			if xi < 0 || yi < 0 {
    			PrintError("Coordinates must be positive for station: " + name)
			}

			if _, exists := stations[name]; exists {
				PrintError("Duplicate station: " + name)
			}

			stations[name] = &algorithm.Station{Name: name, X: xi, Y: yi}
    		continue
		}

		if section == "connections" {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				PrintError("Invalid connection line: " + line)
			}

			a := strings.TrimSpace(parts[0])
			b := strings.TrimSpace(parts[1])

			if _, exists := stations[a]; !exists {
				PrintError("Unknown station: " + a)
			}
			if _, exists := stations[b]; !exists {
				PrintError("Unknown station: " + b)
			}

			for _, c := range connections {
				if (c[0] == a && c[1] == b) || (c[0] == b && c[1] == a) {
					PrintError("Duplicate connection between " + a + " and " + b)
				}
			}
			connections = append(connections, [2]string{a, b})
		}
	}
	if len(stations) > 10000 {
		PrintError("Map contains too many stations")
	}
	stationsSlice := make([]*algorithm.Station, 0, len(stations))
	for _, s := range stations {
		stationsSlice = append(stationsSlice, s)
	}

	return stationsSlice, connections
}