package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gitea.kood.tech/sayemaraf/pathfinder/algorithm"
)

type PathfinderRequest struct {
	MapFile      string `json:"mapFile"`
	StartStation string `json:"startStation"`
	EndStation   string `json:"endStation"`
	NumTrains    int    `json:"numTrains"`
}

type PathfinderResponse struct {
	Success   bool               `json:"success"`
	Error     string             `json:"error,omitempty"`
	Stations  map[string]Station `json:"stations"`
	Paths     [][]string         `json:"paths"`
	Movements [][]TrainMovement  `json:"movements"`
}

type Station struct {
	Name string `json:"name"`
	X    int    `json:"x"`
	Y    int    `json:"y"`
}

type TrainMovement struct {
	TrainID int    `json:"trainId"`
	Station string `json:"station"`
}

func Start() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve main page
	http.HandleFunc("/", serveHome)

	// API endpoints
	http.HandleFunc("/api/maps", listMaps)
	http.HandleFunc("/api/pathfind", handlePathfind)

	port := ":8080"
	fmt.Printf("Pathfinder Web UI starting on http://localhost%s\n", port)
	fmt.Println("Open your browser and navigate to the URL above")
	log.Fatal(http.ListenAndServe(port, nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func listMaps(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("../maps")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var mapFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".map") {
			mapFiles = append(mapFiles, file.Name())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapFiles)
}

func handlePathfind(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PathfinderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request: "+err.Error())
		return
	}

	// Normalize station names to lowercase for case-insensitive matching
	req.StartStation = strings.ToLower(strings.TrimSpace(req.StartStation))
	req.EndStation = strings.ToLower(strings.TrimSpace(req.EndStation))

	// Read and parse the map file
	mapPath := filepath.Join("../maps", req.MapFile)
	fileContent, err := os.ReadFile(mapPath)
	if err != nil {
		sendError(w, "Failed to read map file: "+err.Error())
		return
	}

	lines := normalizeInput(fileContent)
	stations, connections := parseMap(lines)

	// Build station map with lowercase keys
	stationsMap := make(map[string]*algorithm.Station)
	for _, s := range stations {
		stationsMap[strings.ToLower(s.Name)] = s
	}

	// Validate inputs
	if err := validate(req, stationsMap); err != nil {
		sendError(w, err.Error())
		return
	}

	// Build graph and find paths
	graph := algorithm.NewGraph(stations, connections)
	paths := graph.FindPaths(req.StartStation, req.EndStation)

	if len(paths) == 0 {
		sendError(w, "No path exists between stations")
		return
	}

	// Simulate train movements
	movements := simulateMovements(paths, req.NumTrains, req.StartStation, req.EndStation)

	// Convert stations for response
	stationsResponse := make(map[string]Station)
	for name, s := range stationsMap {
		stationsResponse[name] = Station{
			Name: s.Name,
			X:    s.X,
			Y:    s.Y,
		}
	}

	// Convert paths
	pathsResponse := make([][]string, len(paths))
	for i, path := range paths {
		pathsResponse[i] = []string(path)
	}

	response := PathfinderResponse{
		Success:   true,
		Stations:  stationsResponse,
		Paths:     pathsResponse,
		Movements: movements,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func simulateMovements(paths []algorithm.Path, numTrains int, start, end string) [][]TrainMovement {
	type Train struct {
		ID       int
		Path     algorithm.Path
		Position int
	}

	trains := make([]*Train, numTrains)
	for i := 0; i < numTrains; i++ {
		pathIndex := i % len(paths)
		trains[i] = &Train{
			ID:       i + 1,
			Path:     paths[pathIndex],
			Position: 0,
		}
	}

	occupied := make(map[string]bool)
	var allMovements [][]TrainMovement

	for {
		var turnMovements []TrainMovement
		allFinished := true

		for _, train := range trains {
			if train.Position >= len(train.Path)-1 {
				continue
			}

			allFinished = false
			nextPos := train.Position + 1
			nextStation := train.Path[nextPos]

			canMove := nextStation == end || nextStation == start || !occupied[nextStation]

			if canMove {
				if train.Position > 0 {
					currentStation := train.Path[train.Position]
					if currentStation != start {
						occupied[currentStation] = false
					}
				}

				train.Position = nextPos

				if nextStation != end {
					occupied[nextStation] = true
				}

				turnMovements = append(turnMovements, TrainMovement{
					TrainID: train.ID,
					Station: nextStation,
				})
			}
		}

		if len(turnMovements) > 0 {
			allMovements = append(allMovements, turnMovements)
		}

		if allFinished {
			break
		}
	}

	return allMovements
}

func validate(req PathfinderRequest, stations map[string]*algorithm.Station) error {
	if req.StartStation == req.EndStation {
		return fmt.Errorf("start and end stations cannot be the same")
	}
	if _, exists := stations[req.StartStation]; !exists {
		return fmt.Errorf("start station does not exist")
	}
	if _, exists := stations[req.EndStation]; !exists {
		return fmt.Errorf("end station does not exist")
	}
	if req.NumTrains < 1 {
		return fmt.Errorf("number of trains must be at least 1")
	}
	return nil
}

func sendError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PathfinderResponse{
		Success: false,
		Error:   message,
	})
}

func normalizeInput(data []byte) []string {
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

func parseMap(lines []string) ([]*algorithm.Station, [][2]string) {
	section := ""
	stations := make(map[string]*algorithm.Station)
	var connections [][2]string

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
			if len(parts) != 3 {
				continue
			}

			name := strings.TrimSpace(parts[0])
			x, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
			y, _ := strconv.Atoi(strings.TrimSpace(parts[2]))

			stations[name] = &algorithm.Station{Name: name, X: x, Y: y}
			continue
		}

		if section == "connections" {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				continue
			}

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
