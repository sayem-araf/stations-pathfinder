package algorithm
// This file defines the Scheduler struct and its methods for simulating train movements along paths.
import (
	"fmt"
	"sort"
	"strings"
)
// Train represents a train with its ID, assigned path, current position on the path, and whether it has finished its journey.
type Train struct {
	ID       int
	Path     Path
	Position int
	Done     bool
}
// Scheduler manages the scheduling of trains along given paths from a start station to an end station. 
// It keeps track of the paths, number of trains, and the start and end stations.
type Scheduler struct {
	paths     []Path
	numTrains int
	start     string
	end       string
}
// NewScheduler creates a new Scheduler instance with the provided paths, number of trains, start station, and end station.
func NewScheduler(paths []Path, numTrains int, start, end string) *Scheduler {
	return &Scheduler{
		paths:     paths,
		numTrains: numTrains,
		start:     start,
		end:       end,
	}
}
// Run executes the scheduling simulation. It creates the trains, simulates their movements turn by turn, 
// and prints the moves until all trains have finished or a maximum number of turns is reached.
func (s *Scheduler) Run() {
	if len(s.paths) == 0 {
		return
	}
	trains := s.createTrains()
	occupied := make(map[string]bool)
	for turn := 0; turn < 10000; turn++ {
		if s.allFinished(trains) {
			break
		}
		moves := s.simulateTurn(trains, occupied)
		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		}
	}
}
// createTrains initializes the Train instances based on the available paths and the number of trains. 
// It assigns trains to paths in a round-robin manner, 
// and then tries to optimize the assignment by moving trains from longer paths to shorter ones if it reduces the total number of turns needed for all trains to finish.
func (s *Scheduler) createTrains() []*Train {
	trains := make([]*Train, s.numTrains)
// Sort paths by length (shortest first)
	sorted := make([]Path, len(s.paths))
	copy(sorted, s.paths)
	sort.Slice(sorted, func(i, j int) bool {
		return len(sorted[i]) < len(sorted[j])
	})

	// Simple round-robin assignment
	assignment := make([]int, s.numTrains)
	for i := range assignment {
		assignment[i] = i % len(sorted)
	}

	// Fix: if last assigned train is on a longer path than shortest,
	// and moving it to shortest path reduces total turns, do it
	if len(sorted) > 1 {
		improved := true
		for improved {
			improved = false
			counts := make([]int, len(sorted))
			for _, p := range assignment {
				counts[p]++
			}
			finishTurn := make([]int, len(sorted))
			for i, p := range sorted {
				if counts[i] > 0 {
					finishTurn[i] = (len(p) - 1) + counts[i] - 1
				}
			}
			// Find max finish turn
			maxFinish := 0
			for _, f := range finishTurn {
				if f > maxFinish {
					maxFinish = f
				}
			}
			// Try moving last train from slowest path to any faster path
			for slowest := len(sorted) - 1; slowest > 0; slowest-- {
				if finishTurn[slowest] != maxFinish {
					continue
				}
				// Try moving last train on this path to path 0 (shortest)
				newCounts := make([]int, len(sorted))
				copy(newCounts, counts)
				newCounts[slowest]--
				newCounts[0]++
				newFinish := 0
				for i, p := range sorted {
					if newCounts[i] > 0 {
						f := (len(p) - 1) + newCounts[i] - 1
						if f > newFinish {
							newFinish = f
						}
					}
				}
				if newFinish < maxFinish {
					// Find last train on slowest path and reassign
					for i := s.numTrains - 1; i >= 0; i-- {
						if assignment[i] == slowest {
							assignment[i] = 0
							improved = true
							break
						}
					}
					break
				}
			}
		}
	}
// Create Train instances based on final assignment
	for i := 0; i < s.numTrains; i++ {
		trains[i] = &Train{
			ID:       i + 1,
			Path:     sorted[assignment[i]],
			Position: 0,
			Done:     false,
		}
	}
	return trains
}
// simulateTurn simulates one turn of train movements. It checks each train in order,
func (s *Scheduler) simulateTurn(trains []*Train, occupied map[string]bool) []string {
	var moves []string
	usedEdges := make(map[string]bool)
	targetedStations := make(map[string]bool)
// Process trains in order of their IDs to ensure deterministic behavior
	for _, train := range trains {
		if train.Done {
			continue
		}
		if train.Position >= len(train.Path)-1 {
			train.Done = true
			continue
		}
//		Check if the train can move to the next station without conflicts
		currStation := train.Path[train.Position]
		nextPos := train.Position + 1
		nextStation := train.Path[nextPos]
		edgeKey := edgeID(currStation, nextStation)

		if usedEdges[edgeKey] {
			continue
		}
//	The train can move if the next station is not occupied or targeted by another train, unless it's the end station which can be shared
		isEnd := nextStation == s.end
		if !isEnd {
			if occupied[nextStation] || targetedStations[nextStation] {
				continue
			}
		}
//	Move the train: mark current station as unoccupied (if not start/end), mark next station as occupied, and update train position
		if currStation != s.start && currStation != s.end {
			occupied[currStation] = false
		}
//	Mark the edge as used for this turn to prevent other trains from using it
		usedEdges[edgeKey] = true
		if !isEnd {
			occupied[nextStation] = true
			targetedStations[nextStation] = true
		}
//	Update train position and check if it has reached the end
		train.Position = nextPos
		if train.Position >= len(train.Path)-1 {
			train.Done = true
		}
		moves = append(moves, fmt.Sprintf("T%d-%s", train.ID, nextStation))
	}
	return moves
}
// edgeID generates a unique identifier for an edge between two stations, regardless of the order of the stations.
func edgeID(a, b string) string {
	if a < b {
		return a + "|" + b
	}
	return b + "|" + a
}
//	allFinished checks if all trains have completed their journeys by verifying that each train is either marked as done or has reached the end of its path.
func (s *Scheduler) allFinished(trains []*Train) bool {
	for _, t := range trains {
		if !t.Done && t.Position < len(t.Path)-1 {
			return false
		}
	}
	return true
}
// The following wrapper methods are provided for testing purposes, allowing external code to call the internal methods of the Scheduler struct.
func (s *Scheduler) CreateAndDistributeTrains() []*Train {
	return s.createTrains()
}

func (s *Scheduler) SimulateTurnWrapper(trains []*Train, occupied map[string]bool) []string {
	return s.simulateTurn(trains, occupied)
}

func (s *Scheduler) AllTrainsFinished(trains []*Train) bool {
	return s.allFinished(trains)
}