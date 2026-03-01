# Pathfinder

---

## Table of Contents
- **[Description](#description)**
- **[Requirements](#requirements)**
- **[Project Structure](#project-structure)**
- **[Building](#building)**
    - [Linux & Visual Studio Code](#linux--visual-studio-code)
    - [Windows](#winbuild)
- **[Server setup](#server-setup)**
    - [Linux](#linux)
    - [Windows](#windows)
- **[Algorithm](#algorithm)**
- **[Error handler](#error-handler)**
- **[Usage](#usage)**
- **[Event Triggers](#event-triggers)**
- **[Bonus Function](#bonus-function)**

---

## Description

## Requirements

## Project Structure
pathfinder/                   
│                        
├─ algorithm/               
│  ├─ find.go #          
│  ├─ graph.go #            
│  └─ path.go #              
│                       
├─ maps/ # contains the network.map files                      
│                    
├─ parser/                       
│  ├─ errors.go #                          
│  └─ parser.go #                           
│                                 
├─ web/                              
│  ├─ static/                                   
│  │  └── app.js #                             
│  ├─ templates/                                
│  │  └── index.html #                                 
│  └─ server.go #                               
│                                             
├─ main.go # main for the program                                    
└─ README.md # Your reading this                                   

## Building

Before building **locate** the file with **CMD/Powershell**, **Terminal** or clone with **Visual Studio Code**: 

Windows ```cd c:\yourpath\to\pathfinder```

Linux ```cd /yourpath/to/pathfinder```

Visual Studio code ```git clone https://gitea.kood.tech/sayemaraf/pathfinder.git```

### Linux & Visual Studio Code 
If building on Linux **terminal** use:

For **Linux**: ```go build -o pathfinder```

For **Windows** ``` GOOS=windows GOARCH=amd64 go build -o pathfinder.exe ```

<h3 id="winbuild">Windows</h3>

If building on Windows **powershell/cmd** use:

For **Windows**: ```go build -o pathfinder.exe```

For **Linux CMD**: 
``` 
Set CGO_ENABLED=0
Set GOOS=linux
set GOARCH=amd64
go build -o pathfinder
```

For **Linux Powershell**: 
``` 
$env:CGO_ENABLED=0
$env:GOOS=linux
$env:GOARCH=amd64
go build -o pathfinder
```

## Server Setup
### Linux
After following the [building instructions](#building) **locate** the folder in terminal ```cd yourpath/to/pathfinder```,  
to run the code use ```./pathfinder```

If you haven’t built, you can run the code by **locating** the folder and using ```go run .```
### Windows
After following the [building instructions](#building) run the program with the pathfinder.exe in the folder or **locate** the folder on CMD/Powershell:  
```cd c:\yourpath\pathfinder```, if you are using CMD you can use ```pathfinder.exe``` and on Powershell ```.\pathfinder.exe```.  

If you haven’t built, you can run the code by **locating** the folder and using ```go run . ```

## Algorithm
The used Algorithm is The Breadth-First Search or BFS.

 BFS algorithm is a method for systematically exploring a graph or tree data structure,  
 visiting all nodes at the current depth level before moving on to nodes at the next depth level.

 For example:
 ```
GRAPH                         | WHAT BFS IS DOING
------------------------------|---------------------------------
      A                       | Lets say that Start is A and E is the End,    
     / \                      | in this case BFS works likes this.
    B   C                     | 
   / \   \                    | Start at A. Queue = [A]
  D   E   F                   |

Visit A                       | Dequeue A → visit it
                              | Enqueue neighbors B, C
                              | Queue = [B, C]

Visit B                       | Dequeue B → visit it
                              | Enqueue unvisited D, E
                              | Queue = [C, D, E]

Visit C                       | Dequeue C → visit it
                              | Enqueue unvisited F
                              | Queue = [D, E, F]

Visit D                       | Dequeue D → visit it
                              | No new neighbors
                              | Queue = [E, F]

Visit E (END)                 | Dequeue E → visit it
                              | Stops here because E is the End
                              | (F is never processed)
 ```


## Error handler

Here are all of the errors the program reports and why.

### 1. Invalid Arguments

**Too few arguments:**

If user runs program with under four arguments it prints out:

`Error: Too few arguments`

**Too many arguments:**

If user runs program with over four arguments it prints out:

`Error: Too many arguments`

### 2. Invalid Station Names

**Non-existent start station:**

If start station doesn't exist in the map it prints out:

`Error: Invalid start station`

**Non-existent end station:**

If end station doesn't exist in the map it prints out:

`Error: Invalid end station`

### 3. Same Start and End

If start and end station are the same it prints out:

`Error: Start station can't be the same as end station`

### 4. Invalid Train Count

If the number of trains is not valid it prints out:

`Error: Number of trains must be at least 1`

### 5. Invalid Map File

**Non-existent file:**

If the map file doesn't exist it prints out:

`Error: Failed to read file: ...`

**No path exists:**

If no path exists between stations it prints out:

`Error: No path exists between A and C`

**Duplicate station names:**

If station names are repeated in the map it prints out:

`Error: Duplicate station name: ...`

**Invalid coordinates:**

If the coordinates aren't valid positive integers it prints out:

`Error: Station ... has invalid ... coordinate: ...`

**Duplicate coordinates:**

If stations have duplicate coordinates it prints out:

`Error: Stations ... and ... share the same coordinates ...`

**Invalid station names:**

If the station names aren't valid it prints out:

`Error: Invalid station name: ...`

**No stations section:**

If the map file doesn't contain a section for stations it prints out:

`Error: Map does not contain a stations section`

**No connections section:**

If the map file doesn't contain a section for connections it prints out:

`Error: Map does not contain a connections section`

**Duplicate connections:**

If stations have duplicate connections it prints out:

`Error: Duplicate connection between ... and ...`

---
