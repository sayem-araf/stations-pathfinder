package algorithm

import (
	"fmt"
	"strings"
)

type Train struct {
	ID       int
	Path     Path
	Position int
	Done     bool
}

type Scheduler struct {
	paths     []Path
	numTrains int
	start     string
	end       string
}

func NewScheduler(paths []Path, numTrains int, start, end string) *Scheduler {
	return &Scheduler{
		paths:     paths,
		numTrains: numTrains,
		start:     start,
		end:       end,
	}
}

// Run executes the full simulation for CLI usage
func (s *Scheduler) Run() {
	if len(s.paths) == 0 {
		return
	}

	trains := s.createTrains()
	occupied := make(map[string]bool) // stations currently holding a train

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

// createTrains distributes trains across paths in round-robin order
func (s *Scheduler) createTrains() []*Train {
	trains := make([]*Train, s.numTrains)
	for i := 0; i < s.numTrains; i++ {
		pathIdx := i % len(s.paths)
		trains[i] = &Train{
			ID:       i + 1,
			Path:     s.paths[pathIdx],
			Position: 0,
			Done:     false,
		}
	}
	return trains
}

// simulateTurn moves each train that can legally advance this turn.
// Fixed-block rules enforced:
//  1. Each track (edge) can only be used once per turn.
//  2. Only one train may occupy an intermediate station at a time.
//  3. Start and end stations have unlimited capacity.
//  4. A train cannot enter a station already targeted by another train this turn.
func (s *Scheduler) simulateTurn(trains []*Train, occupied map[string]bool) []string {
	var moves []string

	// Track edges used this turn: "stationA|stationB" (always sorted alphabetically)
	usedEdges := make(map[string]bool)

	// Track stations being moved INTO this turn (prevents two trains targeting same station)
	targetedStations := make(map[string]bool)

	for _, train := range trains {
		if train.Done {
			continue
		}
		if train.Position >= len(train.Path)-1 {
			train.Done = true
			continue
		}

		currStation := train.Path[train.Position]
		nextPos := train.Position + 1
		nextStation := train.Path[nextPos]

		edgeKey := edgeID(currStation, nextStation)

		// Rule 1: each track can only be used once per turn
		if usedEdges[edgeKey] {
			continue
		}

		isEnd := nextStation == s.end

		// Rule 2 & 4: next station must be free and not already targeted this turn
		if !isEnd {
			if occupied[nextStation] || targetedStations[nextStation] {
				continue
			}
		}

		// --- Train can legally move ---

		// Free the current station if it's intermediate
		if currStation != s.start && currStation != s.end {
			occupied[currStation] = false
		}

		// Claim the edge and the next station
		usedEdges[edgeKey] = true

		if !isEnd {
			occupied[nextStation] = true
			targetedStations[nextStation] = true
		}

		train.Position = nextPos
		if train.Position >= len(train.Path)-1 {
			train.Done = true
		}

		moves = append(moves, fmt.Sprintf("T%d-%s", train.ID, nextStation))
	}

	return moves
}

// edgeID returns a canonical sorted key for an undirected edge
func edgeID(a, b string) string {
	if a < b {
		return a + "|" + b
	}
	return b + "|" + a
}

func (s *Scheduler) allFinished(trains []*Train) bool {
	for _, t := range trains {
		if !t.Done && t.Position < len(t.Path)-1 {
			return false
		}
	}
	return true
}

// --- Exported wrappers for web integration ---

func (s *Scheduler) CreateAndDistributeTrains() []*Train {
	return s.createTrains()
}

func (s *Scheduler) SimulateTurnWrapper(trains []*Train, occupied map[string]bool) []string {
	return s.simulateTurn(trains, occupied)
}

func (s *Scheduler) AllTrainsFinished(trains []*Train) bool {
	return s.allFinished(trains)
}
