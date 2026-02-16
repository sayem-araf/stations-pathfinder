# Stations Pathfinder - Testing & Usage Instructions

This guide provides step-by-step instructions to test and verify all functionality of the Stations Pathfinder project.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Testing the CLI Version](#testing-the-cli-version)
- [Testing the Web UI](#testing-the-web-ui)
- [Testing Error Handling](#testing-error-handling)
- [Performance Testing](#performance-testing)
- [Running Unit Tests](#running-unit-tests)

---

## Prerequisites

Before testing, ensure you have:
- **Go 1.22.2 or higher** installed
- A terminal (Terminal.app on macOS, or VS Code integrated terminal)
- A web browser (Safari, Chrome, Firefox)

### Verify Go Installation

```bash
go version
```

Expected output: `go version go1.22.2 darwin/arm64` (or similar)

---

## Project Structure

```
pathfinder/
├── parser/         # CLI entry point and parsing logic
├── algorithm/      # Pathfinding algorithms
├── web/           # Web UI server
├── maps/          # Network map files
├── go.mod         # Go module definition
└── README.md      # Project documentation
```

---

## Testing the CLI Version

The CLI version is the core pathfinder that runs from the command line.

### 1. Navigate to Project Directory

```bash
cd /Users/sayemaraf/Desktop/VSCodePorject/pathfinder/pathfinder
```

### 2. Basic Test - London Map

Test with the standard example from requirements:

```bash
go run ./parser maps/london.map waterloo st_pancras 4
```

**Expected Output:**
```
T1-victoria T2-euston
T1-st_pancras T2-st_pancras T3-victoria T4-euston
T3-st_pancras T4-st_pancras
```

**Verification:**
- ✅ Exactly 3 lines of output (3 turns)
- ✅ All trains (T1-T4) move from waterloo to st_pancras
- ✅ Trains use 2 different paths (via victoria and via euston)

### 3. Count Number of Turns

```bash
go run ./parser maps/london.map waterloo st_pancras 4 | wc -l
```

**Expected Output:** `3`

This confirms the solution is optimal with minimal turns.

### 4. Test with Different Train Counts

**2 trains:**
```bash
go run ./parser maps/london.map waterloo st_pancras 2
```

**Expected:** 2 lines (2 turns)

**10 trains:**
```bash
go run ./parser maps/london.map waterloo st_pancras 10
```

**Expected:** More turns as more trains need to traverse

### 5. Test with Different Maps

**Small map:**
```bash
go run ./parser maps/small.map A D 5
```

**Beginning map:**
```bash
go run ./parser maps/beginning.map a e 3
```

**Bond map (more complex):**
```bash
go run ./parser maps/bond.map bond eon 6
```

### 6. Build Executable (Optional)

Build a standalone executable:

```bash
go build -o pathfinder ./parser
```

Then run it:

```bash
./pathfinder maps/london.map waterloo st_pancras 4
```

---

## Testing the Web UI

The web UI provides an interactive visualization of train movements.

### 1. Start the Web Server

```bash
go run web/server.go
```

**Expected Output:**
```
🚂 Pathfinder Web UI starting on http://localhost:8080
Open your browser and navigate to the URL above
```

**Note:** The terminal will stay here - this is normal. The server is running.

### 2. Open in Browser

Open your web browser and navigate to:
```
http://localhost:8080
```

**Expected Result:**
- Purple gradient interface loads
- Form with dropdown and input fields on the left
- Large canvas visualization area on the right
- Legend showing station/path/train colors

### 3. Test Web UI Functionality

**Step-by-step test:**

1. **Select Map:** Choose `london.map` from dropdown
2. **Start Station:** Enter `waterloo`
3. **End Station:** Enter `st_pancras`
4. **Number of Trains:** Enter `4`
5. **Click:** "🚀 Run Pathfinder" button

**Expected Result:**
- Success message: "Found 2 path(s) in 3 turn(s)!"
- Network graph displays with:
  - Green station (waterloo - start)
  - Red station (st_pancras - end)
  - Gray-blue stations (victoria, euston)
  - Colored path lines with numbers
  - Orange trains with IDs (T1, T2, T3, T4)
- Statistics show: 2 paths, 3 turns, 4 trains
- Movement log shows turn-by-turn actions
- Animation automatically starts

### 4. Test Animation Controls

- **Play:** Click ▶️ to start/resume animation
- **Pause:** Click ⏸️ to pause animation
- **Reset:** Click 🔄 to restart from beginning

**Verify:**
- Trains move one station per turn
- Current turn is highlighted in log
- Trains don't overlap at intermediate stations

### 5. Test with Different Maps

Try other maps in the web UI:
- `beethoven.map` with stations: beethoven → part
- `bond.map` with stations: bond → eon
- `jungle.map` with custom start/end

### 6. Stop the Web Server

In the terminal where the server is running:
- Press **Ctrl + C** to stop the server

---

## Testing Error Handling

The program must handle various error conditions gracefully.

### 1. Invalid Arguments

**Too few arguments:**
```bash
go run ./parser
```

**Expected:** `Error: Too few arguments`

**Too many arguments:**
```bash
go run ./parser maps/london.map waterloo st_pancras 4 extra
```

**Expected:** `Error: Too many arguments`

### 2. Invalid Station Names

**Non-existent start station:**
```bash
go run ./parser maps/london.map invalid_station st_pancras 4
```

**Expected:** `Error: Invalid start station`

**Non-existent end station:**
```bash
go run ./parser maps/london.map waterloo invalid_station 4
```

**Expected:** `Error: Invalid end station`

### 3. Same Start and End

```bash
go run ./parser maps/london.map waterloo waterloo 4
```

**Expected:** `Error: Start station can't be the same as end station`

### 4. Invalid Train Count

**Non-numeric value:**
```bash
go run ./parser maps/london.map waterloo st_pancras abc
```

**Expected:** `Error: Number of trains must be at least 1`

**Zero trains:**
```bash
go run ./parser maps/london.map waterloo st_pancras 0
```

**Expected:** `Error: Number of trains must be at least 1`

**Negative trains:**
```bash
go run ./parser maps/london.map waterloo st_pancras -5
```

**Expected:** `Error: Number of trains must be at least 1`

### 5. Invalid Map File

**Non-existent file:**
```bash
go run ./parser maps/nonexistent.map waterloo st_pancras 4
```

**Expected:** `Error: Failed to read file: ...`

### 6. Test Error Maps

The project includes test maps for specific error scenarios:

**No path exists:**
```bash
go run ./parser maps/testnopath.map A C 2
```

**Expected:** `Error: No path exists between A and C`

**Duplicate station names:**
```bash
go run ./parser maps/testdupnames.map A B 2
```

**Expected:** `Error: Duplicate station name: ...`

**Invalid coordinates:**
```bash
go run ./parser maps/testbadcoords.map A B 2
```

**Expected:** `Error: Station ... has invalid ... coordinate: ...`

**Duplicate coordinates:**
```bash
go run ./parser maps/testsamecoords.map A B 2
```

**Expected:** `Error: Stations ... and ... share the same coordinates ...`

**Invalid station names:**
```bash
go run ./parser maps/testinvalidstations.map A B 2
```

**Expected:** `Error: Invalid station name: ...`

**No stations section:**
```bash
go run ./parser maps/testnostations.map A B 2
```

**Expected:** `Error: Map does not contain a stations section`

**No connections section:**
```bash
go run ./parser maps/testnoconnections.map A B 2
```

**Expected:** `Error: Map does not contain a connections section`

**Duplicate connections:**
```bash
go run ./parser maps/testduplicateroutes.map A B 2
```

**Expected:** `Error: Duplicate connection between ... and ...`

---

## Performance Testing

Test the pathfinder with complex networks to verify performance.

### 1. Large Map Test

```bash
time go run ./parser maps/testbigmap.map start end 100
```

**Verify:**
- Completes in reasonable time (< 5 seconds)
- No memory issues
- Produces valid output

### 2. Many Trains Test

```bash
time go run ./parser maps/london.map waterloo st_pancras 100
```

**Verify:**
- Handles 100+ trains efficiently
- Turn count scales appropriately
- No crashes or hangs

### 3. Complex Network Test

```bash
time go run ./parser maps/bond.map bond eon 50
```

**Verify:**
- Finds optimal paths in complex network
- Multiple paths are utilized
- Performance is acceptable

---

## Running Unit Tests

If unit tests are available in the `algorithm/` directory:

### 1. Run All Tests

```bash
go test ./...
```

**Expected:** All tests pass

### 2. Run Tests with Verbose Output

```bash
go test -v ./algorithm
```

**Expected:** Detailed test results for each test case

### 3. Run Tests with Coverage

```bash
go test -cover ./algorithm
```

**Expected:** Code coverage percentage

### 4. Run Specific Test

```bash
go test -run TestGraphCreation ./algorithm
```

---

## Quick Verification Checklist

Use this checklist to verify all functionality:

### CLI Functionality
- [ ] Basic pathfinding works (london.map example)
- [ ] Output format is correct (TN-station format)
- [ ] Turn count is optimal
- [ ] Multiple paths are found
- [ ] Works with different train counts
- [ ] Executable can be built

### Web UI Functionality
- [ ] Server starts successfully
- [ ] Web page loads in browser
- [ ] Map dropdown populates
- [ ] Pathfinding runs from UI
- [ ] Network graph displays correctly
- [ ] Trains animate properly
- [ ] Play/Pause/Reset controls work
- [ ] Movement log updates correctly

### Error Handling
- [ ] Invalid arguments are caught
- [ ] Invalid stations are detected
- [ ] Same start/end is rejected
- [ ] Invalid train count is caught
- [ ] File errors are handled
- [ ] Map validation errors work
- [ ] All error test maps produce correct errors

### Performance
- [ ] Large maps complete in reasonable time
- [ ] Many trains are handled efficiently
- [ ] No memory leaks or crashes
- [ ] Complex networks work correctly

---

## Troubleshooting

### "command not found: go"
- Install Go from https://golang.org/dl/
- Add Go to your PATH

### "cannot find package"
```bash
go mod tidy
```

### Web server shows blank page
- Ensure you're in the project root directory
- Check that `web/templates/index.html` and `web/static/app.js` exist
- Look for errors in browser console (F12 → Console)

### "port already in use"
```bash
lsof -i :8080
kill -9 [PID]
```

### Pathfinder hangs
- Press Ctrl+C to cancel
- Check if the map has valid paths between start and end
- Verify the map format is correct

---

## Example Test Session

Complete test session from start to finish:

```bash
# 1. Navigate to project
cd /Users/sayemaraf/Desktop/VSCodePorject/pathfinder/pathfinder

# 2. Test basic functionality
go run ./parser maps/london.map waterloo st_pancras 4

# 3. Verify turn count
go run ./parser maps/london.map waterloo st_pancras 4 | wc -l

# 4. Test error handling
go run ./parser maps/london.map invalid waterloo 4

# 5. Build executable
go build -o pathfinder ./parser

# 6. Test executable
./pathfinder maps/london.map waterloo st_pancras 4

# 7. Start web server (in new terminal)
go run web/server.go

# 8. Test in browser at http://localhost:8080

# 9. Stop server with Ctrl+C
```

---

## Expected Results Summary

### CLI Output Format
```
TN-station [TM-station ...]
```
- Each line = one turn
- Space-separated train movements
- Train IDs: T1, T2, T3, etc.
- Station names as defined in map

### Web UI Output
- Visual network graph
- Animated train movements
- Turn-by-turn movement log
- Statistics display
- Interactive controls

### Error Output Format
```
Error: [descriptive error message]
```
- Printed to stderr
- Exit status 1
- Clear description of the problem

---

## Contact & Support

For issues or questions:
- Review the main README.md
- Check error messages carefully
- Verify map file format
- Ensure Go version compatibility

**Project Status:** ✅ Complete and Functional
**Last Updated:** February 16, 2026
