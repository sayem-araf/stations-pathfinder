# Pathfinder - Testing Guide

<<<<<<< HEAD
Quick reference for testing all functionality.
=======
This guide provides step-by-step instructions to test and verify all functionality of the Stations Pathfinder project.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Testing the CLI Version](#testing-the-cli-version)
- [Mandatory Test Cases](#mandatory-test-cases)
- [Testing the Web UI](#testing-the-web-ui)
- [Testing Error Handling](#testing-error-handling)
- [Performance Testing](#performance-testing)
- [Running Unit Tests](#running-unit-tests)
>>>>>>> 624c052 (adding github as second repo)

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

<<<<<<< HEAD
## Quick Start
=======
## Project Structure

```
pathfinder/
├── main.go        # Main entry point (CLI and web server)
├── parser/        # Input parsing and validation logic
├── algorithm/     # Pathfinding algorithms and scheduler
├── web/           # Web UI server and assets
├── maps/          # Network map files
├── go.mod         # Go module definition
└── README.md      # Project documentation
```

---

## Testing the CLI Version

The CLI version is the core pathfinder that runs from the command line.

### 1. Navigate to Project Directory
>>>>>>> 624c052 (adding github as second repo)

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

<<<<<<< HEAD
Build program:
=======
**Verification:**
- ✅ Exactly 3 lines of output (3 turns)
- ✅ All trains (T1-T4) move from waterloo to st_pancras
- ✅ Trains use 2 different paths (via victoria and via euston)

### 3. Count Number of Turns

```bash
go run . maps/london.map waterloo st_pancras 4 | wc -l
```

**Expected Output:** `3`

This confirms the solution is optimal with minimal turns.

### 4. Test with Different Train Counts

**2 trains:**
```bash
go run . maps/london.map waterloo st_pancras 2
```

**Expected:** 2 lines (2 turns)

**10 trains:**
```bash
go run . maps/london.map waterloo st_pancras 10
```

**Expected:** More turns as more trains need to traverse

### 5. Test with Different Maps

**Small map:**
```bash
go run . maps/small.map A D 5
```

**Beginning map:**
```bash
go run . maps/beginning.map a e 3
```

**Bond map (more complex):**
```bash
go run . maps/bond.map bond eon 6
```

### 6. Build Executable (Optional)

Build a standalone executable:

```bash
go build -o pathfinder .
```

Then run it:

>>>>>>> 624c052 (adding github as second repo)
```bash
go build -o pathfinder .
./pathfinder maps/london.map waterloo st_pancras 4
```

---

<<<<<<< HEAD
## Web Interface
=======
## Mandatory Test Cases

These test cases are **REQUIRED** for the project audit. All must pass for the project to be accepted.

---

### Test 1: Valid Train Movement Rules (CRITICAL)

**Objective:** Verify trains move correctly according to the rules.

**For EVERY test, verify that:**
1. ✅ All trains successfully reach the end station
2. ✅ No more than one train is in a station at any time (except start/end)
3. ✅ The same track is NOT used more than once in a single turn
4. ✅ Each train moves ONLY ONCE per turn

**How to verify:**
- Manually review the output line by line
- Check that train movements don't violate any rules
- Confirm all trains eventually reach the destination

**Example validation for London Map (4 trains):**
```
Turn 1: T1-victoria T2-euston
  ✓ T1 moves waterloo→victoria (track: waterloo-victoria)
  ✓ T2 moves waterloo→euston (track: waterloo-euston)
  ✓ No track used twice
  ✓ No station collision

Turn 2: T1-st_pancras T2-st_pancras T3-victoria T4-euston
  ✓ T1 moves victoria→st_pancras (track: victoria-st_pancras)
  ✓ T2 moves euston→st_pancras (track: euston-st_pancras)
  ✓ T3 moves waterloo→victoria (track: waterloo-victoria)
  ✓ T4 moves waterloo→euston (track: waterloo-euston)
  ✓ No track used twice
  ✓ Victoria is now empty (T1 left)
  ✓ Euston is now empty (T2 left)

Turn 3: T3-st_pancras T4-st_pancras
  ✓ All trains reached destination
```

---

### Test 2: London Map - 2 Trains (Multiple Routes)

**Command:**
```bash
go run . maps/london.map waterloo st_pancras 2
```

**Expected:**
- ✅ Completes in 2 turns
- ✅ Uses MORE THAN ONE route (not just one path repeatedly)
- ✅ Output format: `T1-victoria` and `T2-euston` (or similar)

**Example Output:**
```
T1-victoria T2-euston
T1-st_pancras T2-st_pancras
```

**Verification:**
- [ ] 2 lines (2 turns)
- [ ] Trains use different paths (one via victoria, one via euston)
- [ ] All trains reach st_pancras

---

### Test 3: London Map - 3 Trains (Multiple Routes)

**Command:**
```bash
go run . maps/london.map waterloo st_pancras 3
```

**Expected:**
- ✅ Completes in 2-3 turns
- ✅ Uses MORE THAN ONE route
- ✅ All 3 trains reach destination

**Example Output:**
```
T1-victoria T2-euston
T1-st_pancras T2-st_pancras T3-victoria
T3-st_pancras
```

**Verification:**
- [ ] 3 lines or less
- [ ] Multiple paths utilized
- [ ] All trains reach st_pancras

---

### Test 4: London Map - 4 Trains (Multiple Routes)

**Command:**
```bash
go run . maps/london.map waterloo st_pancras 4
```

**Expected:**
- ✅ Completes in 3 turns
- ✅ Uses MORE THAN ONE route
- ✅ All 4 trains reach destination

**Example Output:**
```
T1-victoria T2-euston
T1-st_pancras T2-st_pancras T3-victoria T4-euston
T3-st_pancras T4-st_pancras
```

**Verification:**
- [ ] Exactly 3 lines (3 turns)
- [ ] Uses both victoria and euston paths
- [ ] All trains reach st_pancras

---

### Test 5: London Map - 100 Trains (Multiple Routes)

**Command:**
```bash
go run . maps/london.map waterloo st_pancras 100
```

**Expected:**
- ✅ Uses MORE THAN ONE route efficiently
- ✅ All 100 trains reach destination
- ✅ Completes in reasonable time (< 5 seconds)

**Verification:**
```bash
# Count total turns
go run . maps/london.map waterloo st_pancras 100 | wc -l

# Verify all trains reached destination
go run . maps/london.map waterloo st_pancras 100 | tail -1
```

**Check:**
- [ ] Completes successfully
- [ ] Uses parallel paths (not sequential)
- [ ] All trains eventually show "T100-st_pancras"

---

### Test 6: London Map - 1 Train (Single Route ONLY)

**Command:**
```bash
go run . maps/london.map waterloo st_pancras 1
```

**Expected:**
- ✅ Uses ONLY ONE route (not multiple)
- ✅ Completes in 2 turns
- ✅ Train reaches destination

**Example Output (via victoria):**
```
T1-victoria
T1-st_pancras
```

**OR (via euston):**
```
T1-euston
T1-st_pancras
```

**Verification:**
- [ ] Exactly 2 lines (2 turns)
- [ ] Only one path used (either victoria OR euston)
- [ ] Train reaches st_pancras

---

### Test 7: Output Format Validation

**Command:**
```bash
go run . maps/london.map waterloo st_pancras 4
```

**Expected Output Format:**
```
T1-victoria T2-euston
T1-st_pancras T2-st_pancras T3-victoria T4-euston
T3-st_pancras T4-st_pancras
```

**Verify:**
- [ ] Format is exactly: `TN-station` (e.g., `T1-victoria`)
- [ ] Space-separated on same line for simultaneous moves
- [ ] Each line represents one turn
- [ ] Train IDs are sequential (T1, T2, T3, etc.)
- [ ] Station names match the map file exactly

---

### Test 8: Bond Map - 4 Trains (Max 6 Turns)

**Command:**
```bash
go run . maps/bond.map bond_square space_port 4
```

**Expected:**
- ✅ Completes in NO MORE THAN 6 TURNS
- ✅ All 4 trains reach space_port

**Example Output:**
```
T1-apple_avenue
T1-orange_junction T2-apple_avenue
T1-space_port T2-orange_junction T3-apple_avenue
T2-space_port T3-orange_junction T4-apple_avenue
T3-space_port T4-orange_junction
T4-space_port
```

**Verification:**
```bash
# Count turns
go run . maps/bond.map bond_square space_port 4 | wc -l
```

**Check:**
- [ ] Output is 6 lines or less
- [ ] All trains reach space_port
- [ ] Train movements are valid

---

### Test 9: Jungle Map - 10 Trains (Max 8 Turns)

**Command:**
```bash
go run . maps/jungle.map jungle desert 10
```

**Expected:**
- ✅ Completes in NO MORE THAN 8 TURNS
- ✅ All 10 trains reach desert
- ✅ Uses multiple paths efficiently

**Verification:**
```bash
# Count turns
go run . maps/jungle.map jungle desert 10 | wc -l

# Check last line to see if T10 reaches desert
go run . maps/jungle.map jungle desert 10 | tail -1
```

**Check:**
- [ ] Output is 8 lines or less
- [ ] All trains T1-T10 reach desert
- [ ] Multiple paths utilized

---

### Test 10: Beginning Map - 20 Trains (Max 11 Turns)

**Command:**
```bash
go run . maps/beginning.map beginning terminus 20
```

**Expected:**
- ✅ Completes in NO MORE THAN 11 TURNS
- ✅ All 20 trains reach terminus

**Example Output:**
```
T1-terminus T2-near
T2-far T3-terminus T4-near
T2-terminus T4-far T5-terminus T6-near
T4-terminus T6-far T7-terminus T8-near
T6-terminus T8-far T9-terminus T10-near
T8-terminus T10-far T11-terminus T12-near
T10-terminus T12-far T13-terminus T14-near
T12-terminus T14-far T15-terminus T16-near
T14-terminus T16-far T17-terminus T18-near
T16-terminus T18-far T19-terminus
T18-terminus T20-terminus
```

**Verification:**
```bash
# Count turns
go run . maps/beginning.map beginning terminus 20 | wc -l
```

**Check:**
- [ ] Output is 11 lines or less
- [ ] All trains T1-T20 reach terminus

---

### Test 11: One Map - 4 Trains (Max 6 Turns)

**Command:**
```bash
go run . maps/one.map two four 4
```

**Expected:**
- ✅ Completes in NO MORE THAN 6 TURNS
- ✅ All 4 trains reach station 'four'

**Example Output:**
```
T1-three
T1-one T2-three
T1-four T2-one T3-three
T2-four T3-one T4-three
T3-four T4-one
T4-four
```

**Verification:**
```bash
# Count turns
go run . maps/one.map two four 4 | wc -l
```

**Check:**
- [ ] Output is 6 lines or less
- [ ] All trains reach 'four'

---

### Test 12: Beethoven Map - 9 Trains (Max 6 Turns)

**Command:**
```bash
go run . maps/beethoven.map beethoven part 9
```

**Expected:**
- ✅ Completes in NO MORE THAN 6 TURNS
- ✅ All 9 trains reach 'part'
- ✅ Uses multiple paths efficiently

**Example Output:**
```
T1-verdi T3-handel
T1-part T2-verdi T3-mozart T5-handel
T2-part T3-part T4-verdi T5-mozart T7-handel
T4-part T5-part T6-verdi T7-mozart T9-handel
T6-part T7-part T8-verdi T9-mozart
T8-part T9-part
```

**Verification:**
```bash
# Count turns
go run . maps/beethoven.map beethoven part 9 | wc -l
```

**Check:**
- [ ] Output is 6 lines or less
- [ ] All trains T1-T9 reach 'part'
- [ ] Multiple routes used

---

### Test 13: Small Map - 9 Trains (Max 8 Turns)

**Command:**
```bash
go run . maps/small.map small large 9
```

**Expected:**
- ✅ Completes in NO MORE THAN 8 TURNS
- ✅ All 9 trains reach 'large'
- ✅ Uses multiple paths efficiently

**Verification:**
```bash
# Count turns
go run . maps/small.map small large 9 | wc -l

# Verify last train reaches destination
go run . maps/small.map small large 9 | grep "T9-large"
```

**Check:**
- [ ] Output is 8 lines or less
- [ ] All trains T1-T9 reach 'large'
- [ ] Multiple routes utilized

---

### Test 14: Additional Error Handling Tests

**Connection to Non-Existent Station:**

Create a test map `testghost.map`:
```
stations:
A,0,0
B,1,1

connections:
A-B
B-C
```

```bash
go run . maps/testghost.map A B 2
```

**Expected:** `Error: Unknown station in connection: C`

---

**Map with More Than 10,000 Stations:**

```bash
go run . maps/testbigmap.map start end 1
```

**Expected:** 
- If map has > 10,000 stations: `Error: Map contains more than 10000 stations`
- OR: Completes successfully if ≤ 10,000 stations

---

### Test 15: Verify No Rule Violations

**For ANY test output, check:**

1. **No station collisions:**
   ```bash
   # Example manual check for turn 2 of london.map:
   # T1-st_pancras T2-st_pancras T3-victoria T4-euston
   # At this point: victoria has T3, euston has T4, st_pancras has T1 and T2 (END station - allowed)
   ```

2. **No track reuse in same turn:**
   ```bash
   # If turn shows: T1-victoria T2-victoria
   # Both use waterloo-victoria track = INVALID!
   ```

3. **Each train moves once per turn:**
   ```bash
   # If turn shows: T1-victoria T1-euston
   # T1 appears twice = INVALID!
   ```

4. **All trains finish:**
   ```bash
   # Check last turn contains all trains at destination
   go run . maps/london.map waterloo st_pancras 4 | grep -o "T[0-9]*-st_pancras" | sort -u
   # Should show: T1-st_pancras, T2-st_pancras, T3-st_pancras, T4-st_pancras
   ```

---

## Testing the Web UI

The web UI provides an interactive visualization of train movements.

### 1. Start the Web Server
>>>>>>> 624c052 (adding github as second repo)

Start server:
```bash
go run . -w
```

<<<<<<< HEAD
Open browser: `http://localhost:8080`
=======
**Expected Output:**
```
Pathfinder Web UI starting on http://localhost:8080
Open your browser and navigate to the URL above
```
>>>>>>> 624c052 (adding github as second repo)

Test: Select map → Enter trains → Click "Run Pathfinder"

<<<<<<< HEAD
Stop server: Press `Ctrl+C`
=======
### 2. Open in Browser

Open your web browser and navigate to:
```
http://localhost:8080
```

**Expected Result:**
- Dark themed interface loads
- Form with dropdown and input fields on the left
- Large canvas visualization area on the right

### 3. Test Web UI Functionality

**Step-by-step test:**

1. **Select Map:** Choose `london.map` from dropdown
2. **Number of Trains:** Enter `4`
3. **Click:** "Run Pathfinder" button

**Note:** The web UI automatically selects start and end stations based on the map file.

**Expected Result:**
- Success message: "Found 2 path(s) in 3 turn(s) from waterloo to st_pancras!"
- Network graph displays with:
  - White station (waterloo - start)
  - Red station (st_pancras - end)
  - Gray stations (victoria, euston - intermediate)
  - White path lines connecting stations
  - Green trains with IDs (T1, T2, T3, T4)
- Statistics show: 2 paths, 3 turns, 4 trains
- Movement log shows turn-by-turn actions
- Animation automatically starts

### 4. Test with Different Maps

Try other maps in the web UI:
- `beethoven.map` - automatically uses beethoven → part
- `bond.map` - automatically uses bond_square → space_port
- `jungle.map` - automatically uses jungle → desert
- `beginning.map` - automatically uses beginning → terminus

**Verify:**
- Trains move smoothly along paths
- Current turn is highlighted in log
- Trains don't overlap at intermediate stations
- Animation progresses automatically

### 5. Stop the Web Server

In the terminal where the server is running:
- Press **Ctrl + C** to stop the server
>>>>>>> 624c052 (adding github as second repo)

---

## Mandatory Tests

### London Map Tests
```bash
<<<<<<< HEAD
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
=======
go run .
>>>>>>> 624c052 (adding github as second repo)
```

### Other Maps (Turn Limits)
```bash
<<<<<<< HEAD
# Bond - 4 trains in 6 turns or less
go run . maps/bond.map bond_square space_port 4
=======
go run . maps/london.map waterloo st_pancras 4 extra
```
>>>>>>> 624c052 (adding github as second repo)

# Jungle - 10 trains in 8 turns or less
go run . maps/jungle.map jungle desert 10

# Beginning - 20 trains in 11 turns or less
go run . maps/beginning.map beginning terminus 20

<<<<<<< HEAD
# One - 4 trains in 6 turns or less
go run . maps/one.map two four 4
=======
**Non-existent start station:**
```bash
go run . maps/london.map invalid_station st_pancras 4
```
>>>>>>> 624c052 (adding github as second repo)

# Beethoven - 9 trains in 6 turns or less
go run . maps/beethoven.map beethoven part 9

<<<<<<< HEAD
# Small - 9 trains in 8 turns or less
go run . maps/small.map small large 9
=======
**Non-existent end station:**
```bash
go run . maps/london.map waterloo invalid_station 4
```

**Expected:** `Error: Invalid end station`

### 3. Same Start and End

```bash
go run . maps/london.map waterloo waterloo 4
```

**Expected:** `Error: Start station can't be the same as end station`

### 4. Invalid Train Count

**Non-numeric value:**
```bash
go run . maps/london.map waterloo st_pancras abc
```

**Expected:** `Error: Number of trains must be at least 1`

**Zero trains:**
```bash
go run . maps/london.map waterloo st_pancras 0
```

**Expected:** `Error: Number of trains must be at least 1`

**Negative trains:**
```bash
go run . maps/london.map waterloo st_pancras -5
```

**Expected:** `Error: Number of trains must be at least 1`

### 5. Invalid Map File

**Non-existent file:**
```bash
go run . maps/nonexistent.map waterloo st_pancras 4
```

**Expected:** `Error: Failed to read file: ...`

### 6. Test Error Maps

The project includes test maps for specific error scenarios:

**No path exists:**
```bash
go run . maps/testnopath.map A C 2
```

**Expected:** `Error: No path exists between A and C`

**Duplicate station names:**
```bash
go run . maps/testdupnames.map A B 2
```

**Expected:** `Error: Duplicate station name: ...`

**Invalid coordinates:**
```bash
go run . maps/testbadcoords.map A B 2
```

**Expected:** `Error: Station ... has invalid ... coordinate: ...`

**Duplicate coordinates:**
```bash
go run . maps/testsamecoords.map A B 2
```

**Expected:** `Error: Stations ... and ... share the same coordinates ...`

**Invalid station names:**
```bash
go run . maps/testinvalidstations.map A B 2
```

**Expected:** `Error: Invalid station name: ...`

**No stations section:**
```bash
go run . maps/testnostations.map A B 2
```

**Expected:** `Error: Map does not contain a stations section`

**No connections section:**
```bash
go run . maps/testnoconnections.map A B 2
```

**Expected:** `Error: Map does not contain a connections section`

**Duplicate connections:**
```bash
go run . maps/testduplicateroutes.map A B 2
```

**Expected:** `Error: Duplicate connection between ... and ...`

---

## Performance Testing

Test the pathfinder with complex networks to verify performance.

### 1. Large Map Test

```bash
time go run . maps/testbigmap.map start end 100
```

**Verify:**
- Completes in reasonable time (< 5 seconds)
- No memory issues
- Produces valid output

### 2. Many Trains Test

```bash
time go run . maps/london.map waterloo st_pancras 100
```

**Verify:**
- Handles 100+ trains efficiently
- Turn count scales appropriately
- No crashes or hangs

### 3. Complex Network Test

```bash
time go run . maps/bond.map bond_square space_port 50
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
>>>>>>> 624c052 (adding github as second repo)
```

---

<<<<<<< HEAD
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
=======
## Quick Verification Checklist

Use this checklist to verify all functionality:

### Mandatory Test Cases (MUST ALL PASS)
- [ ] **Valid Movement Rules:** All trains reach destination, no collisions, no rule violations
- [ ] **London 2 trains:** Uses multiple routes
- [ ] **London 3 trains:** Uses multiple routes
- [ ] **London 4 trains:** Completes in 3 turns, uses multiple routes
- [ ] **London 100 trains:** Uses multiple routes efficiently
- [ ] **London 1 train:** Uses ONLY ONE route
- [ ] **Output Format:** Correct `TN-station` format
- [ ] **Bond 4 trains:** Completes in ≤6 turns
- [ ] **Jungle 10 trains:** Completes in ≤8 turns
- [ ] **Beginning 20 trains:** Completes in ≤11 turns
- [ ] **One 4 trains:** Completes in ≤6 turns
- [ ] **Beethoven 9 trains:** Completes in ≤6 turns
- [ ] **Small 9 trains:** Completes in ≤8 turns

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
- [ ] Movement log updates correctly

### Error Handling (All Required)
- [ ] **Too few arguments:** Error displayed
- [ ] **Too many arguments:** Error displayed (or used for bonuses)
- [ ] **Invalid start station:** Error displayed
- [ ] **Invalid end station:** Error displayed
- [ ] **Same start and end:** Error displayed
- [ ] **No path exists:** Error displayed
- [ ] **Duplicate routes:** Error displayed
- [ ] **Invalid train count:** Error displayed
- [ ] **Invalid coordinates:** Error displayed
- [ ] **Duplicate coordinates:** Error displayed
- [ ] **Connection to non-existent station:** Error displayed
- [ ] **Duplicate station names:** Error displayed
- [ ] **Invalid station names:** Error displayed
- [ ] **No stations section:** Error displayed
- [ ] **No connections section:** Error displayed
- [ ] **More than 10,000 stations:** Error displayed

### Performance
- [ ] Large maps complete in reasonable time
- [ ] Many trains are handled efficiently
- [ ] No memory leaks or crashes
- [ ] Complex networks work correctly

### Additional Tricky Cases
- [ ] Test with complex networks
- [ ] Test with edge cases (single path, multiple equal-length paths)
- [ ] Test with various train counts (1, 2, 10, 100, 1000)
- [ ] Verify optimal turn counts

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
cd pathfinder

# 2. Test basic functionality
go run . maps/london.map waterloo st_pancras 4

# 3. Verify turn count
go run . maps/london.map waterloo st_pancras 4 | wc -l

# 4. Test error handling
go run . maps/london.map invalid waterloo 4

# 5. Build executable
go build -o pathfinder .
>>>>>>> 624c052 (adding github as second repo)

# No path exists
go run . maps/testnopath.map A C 2

<<<<<<< HEAD
# Invalid train count (non-numeric)
go run . maps/london.map waterloo st_pancras abc
=======
# 7. Start web server (in new terminal)
go run . -w
>>>>>>> 624c052 (adding github as second repo)

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

<<<<<<< HEAD
=======
## Contact & Support

For issues or questions:
- Review the main README.md
- Check error messages carefully
- Verify map file format
- Ensure Go version compatibility

**Project Status:** ✅ Complete and Functional
>>>>>>> 624c052 (adding github as second repo)
**Last Updated:** March 1, 2026
