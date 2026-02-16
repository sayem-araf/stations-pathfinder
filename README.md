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

## Error handler



## Usage

### Bonus Function

