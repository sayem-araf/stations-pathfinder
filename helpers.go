package main

import(
	"strings"
)

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