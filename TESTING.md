# Pathfinder - Testing Guide

Quick reference for testing all functionality.

---

## Prerequisites

- Go 1.22.2 or higher installed
- Terminal application
- Web browser

Check Go installation:
```bash
go version
```

---

## Quick Start

Navigate to project:
```bash
cd pathfinder
```

Basic test:
```bash
go run . maps/london.map waterloo st_pancras 4
```

Expected output (3 lines):
```
T1-victoria T2-euston
T1-st_pancras T2-st_pancras T3-victoria T4-euston
T3-st_pancras T4-st_pancras
```

Build program:
```bash
go build -o pathfinder .
./pathfinder maps/london.map waterloo st_pancras 4
```

---

## Web Interface

Start server:
```bash
go run . -w
```

Open browser: `http://localhost:8080`

Test: Select map → Enter trains → Click "Run Pathfinder"

Stop server: Press `Ctrl+C`

---

## Mandatory Tests

### London Map Tests
```bash
# 2 trains - should complete in 2 turns
go run . maps/london.map waterloo st_pancras 2

# 3 trains - should complete in 3 turns or less
go run . maps/london.map waterloo st_pancras 3

# 4 trains - should complete in exactly 3 turns
go run . maps/london.map waterloo st_pancras 4

# 100 trains - should find multiple routes
go run . maps/london.map waterloo st_pancras 100

# 1 train - should complete in 2 turns
go run . maps/london.map waterloo st_pancras 1
```

### Other Maps (Turn Limits)
```bash
# Bond - 4 trains in 6 turns or less
go run . maps/bond.map bond_square space_port 4

# Jungle - 10 trains in 8 turns or less
go run . maps/jungle.map jungle desert 10

# Beginning - 20 trains in 11 turns or less
go run . maps/beginning.map beginning terminus 20

# One - 4 trains in 6 turns or less
go run . maps/one.map two four 4

# Beethoven - 9 trains in 6 turns or less
go run . maps/beethoven.map beethoven part 9

# Small - 9 trains in 8 turns or less
go run . maps/small.map small large 9
```

---

## Error Handling Tests

```bash
# Too few arguments
go run .

# Too many arguments
go run . maps/london.map waterloo st_pancras 4 extra

# Invalid start station
go run . maps/london.map invalid st_pancras 4

# Invalid end station
go run . maps/london.map waterloo invalid 4

# Same start and end
go run . maps/london.map waterloo waterloo 4

# No path exists
go run . maps/testnopath.map A C 2

# Invalid train count (non-numeric)
go run . maps/london.map waterloo st_pancras abc

# Zero trains
go run . maps/london.map waterloo st_pancras 0

# Negative trains
go run . maps/london.map waterloo st_pancras -5

# Duplicate connections
go run . maps/testduplicateroutes.map A B 2

# Invalid coordinates
go run . maps/testbadcoords.map A B 2

# Duplicate coordinates
go run . maps/testsamecoords.map A B 2

# Invalid station names
go run . maps/testinvalidstations.map A B 2

# Duplicate station names
go run . maps/testdupnames.map A B 2

# No stations section
go run . maps/testnostations.map A B 2

# No connections section
go run . maps/testnoconnections.map A B 2

# Ghost station in connections
go run . maps/testghost.map A B 2
```

All should display appropriate error messages.

---

## Performance Tests

```bash
# Large map
go run . maps/testbigmap.map start end 100

# Many trains
go run . maps/london.map waterloo st_pancras 100

# Complex network
go run . maps/jungle.map jungle desert 50
```

Should complete quickly without hanging.

---

**Last Updated:** March 1, 2026
