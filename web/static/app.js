// Global state
let canvas, ctx;
let stations = {};
let connections = [];
let paths = [];
let movements = [];
let currentTurn = 0;
let animationInterval = null;
let isPlaying = false;
let startStation = '';
let endStation = '';

// Initialize when page loads
document.addEventListener('DOMContentLoaded', () => {
    canvas = document.getElementById('canvas');
    ctx = canvas.getContext('2d');
    loadMaps();
});

// Load available map files
async function loadMaps() {
    try {
        const response = await fetch('/api/maps');
        const maps = await response.json();
        
        const select = document.getElementById('mapSelect');
        select.innerHTML = '<option value="">-- Select a map --</option>';
        
        maps.forEach(map => {
            const option = document.createElement('option');
            option.value = map;
            option.textContent = map.replace('.map', '');
            select.appendChild(option);
        });
    } catch (error) {
        showMessage('Failed to load maps: ' + error.message, 'error');
    }
}

// Run the pathfinder
async function runPathfinder() {
    const mapFile = document.getElementById('mapSelect').value;
    const start = document.getElementById('startStation').value.trim();
    const end = document.getElementById('endStation').value.trim();
    const numTrains = parseInt(document.getElementById('numTrains').value);

    if (!mapFile) {
        showMessage('Please select a map file', 'error');
        return;
    }

    if (!start || !end) {
        showMessage('Please enter start and end stations', 'error');
        return;
    }

    if (isNaN(numTrains) || numTrains < 1) {
        showMessage('Please enter a valid number of trains', 'error');
        return;
    }

    // Show loading
    document.getElementById('runButton').disabled = true;
    document.getElementById('message').innerHTML = '<div class="loading"><div class="spinner"></div>Running pathfinder...</div>';

    try {
        const response = await fetch('/api/pathfind', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                mapFile: mapFile,
                startStation: start,
                endStation: end,
                numTrains: numTrains
            })
        });

        const data = await response.json();

        if (!data.success) {
            showMessage(data.error, 'error');
            return;
        }

        // Store results
        stations = data.stations;
        paths = data.paths;
        movements = data.movements;
        startStation = start.toLowerCase();
        endStation = end.toLowerCase();
        currentTurn = 0;

        // Extract all connections from paths
        connections = [];
        const connSet = new Set();
        paths.forEach(path => {
            for (let i = 0; i < path.length - 1; i++) {
                const key = [path[i], path[i+1]].sort().join('-');
                if (!connSet.has(key)) {
                    connSet.add(key);
                    connections.push([path[i], path[i+1]]);
                }
            }
        });

        // Update stats
        document.getElementById('pathCount').textContent = paths.length;
        document.getElementById('turnCount').textContent = movements.length;
        document.getElementById('trainCount').textContent = numTrains;
        document.getElementById('stats').style.display = 'block';
        document.getElementById('log').style.display = 'block';

        // Draw initial state with trains at start
        drawNetwork();
        updateLog();
        showMessage(`Found ${paths.length} path(s) in ${movements.length} turn(s)!`, 'success');

        // Auto-start animation
        setTimeout(startAnimation, 500);

    } catch (error) {
        showMessage('Error: ' + error.message, 'error');
    } finally {
        document.getElementById('runButton').disabled = false;
    }
}

// Draw trains with better visibility
function drawTrains(positionMap) {
    let currentMovements;
    
    // If currentTurn is 0, show all trains at start station
    if (currentTurn === 0) {
        const trainIds = [];
        movements.forEach(turn => {
            turn.forEach(move => {
                if (!trainIds.includes(move.trainId)) {
                    trainIds.push(move.trainId);
                }
            });
        });
        
        currentMovements = trainIds.map(id => ({ trainId: id, station: startStation }));
    } else {
        currentMovements = movements[currentTurn - 1];
    }
    
    // Group trains by station
    const trainsByStation = {};
    currentMovements.forEach(move => {
        if (!trainsByStation[move.station]) {
            trainsByStation[move.station] = [];
        }
        trainsByStation[move.station].push(move.trainId);
    });
    
    // Draw trains around stations
    Object.entries(trainsByStation).forEach(([stationName, trainIds]) => {
        const pos = positionMap[stationName];
        if (!pos) return;
        
        const numTrains = trainIds.length;
        const radius = numTrains > 1 ? 35 : 30;
        
        trainIds.forEach((trainId, idx) => {
            let trainX, trainY;
            
            if (numTrains === 1) {
                // Single train: position above station
                trainX = pos.x;
                trainY = pos.y + 28;
            } else {
                // Multiple trains: arrange in circle
                const angle = (idx / numTrains) * Math.PI * 2 - Math.PI / 2;
                trainX = pos.x + Math.cos(angle) * radius;
                trainY = pos.y + Math.sin(angle) * radius;
            }
            
            // Train shadow
            ctx.fillStyle = 'rgba(0, 0, 0, 0.3)';
            ctx.beginPath();
            ctx.arc(trainX + 2, trainY + 2, 18, 0, Math.PI * 2);
            ctx.fill();
            
            // Train body (gradient)
            const gradient = ctx.createRadialGradient(trainX - 5, trainY - 5, 2, trainX, trainY, 18);
            gradient.addColorStop(0, '#FFB74D');
            gradient.addColorStop(1, '#FF6F00');
            ctx.fillStyle = gradient;
            ctx.beginPath();
            ctx.arc(trainX, trainY, 18, 0, Math.PI * 2);
            ctx.fill();
            
            // Train border
            ctx.strokeStyle = 'white';
            ctx.lineWidth = 3;
            ctx.stroke();
            
            // Train ID
            ctx.fillStyle = 'white';
            ctx.font = 'bold 13px Arial';
            ctx.textAlign = 'center';
            ctx.textBaseline = 'middle';
            ctx.shadowColor = 'rgba(0,0,0,0.5)';
            ctx.shadowBlur = 3;
            ctx.fillText('T' + trainId, trainX, trainY);
            ctx.shadowBlur = 0;
        });
    });
}

// Draw the network graph with force-directed layout
function drawNetwork() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    if (Object.keys(stations).length === 0) {
        ctx.fillStyle = '#999';
        ctx.font = '20px Arial';
        ctx.textAlign = 'center';
        ctx.fillText('Select a map and run pathfinder to see visualization', canvas.width/2, canvas.height/2);
        return;
    }

    // Create position map
    const positionMap = {};
    const coords = Object.values(stations);
    const minX = Math.min(...coords.map(s => s.x));
    const maxX = Math.max(...coords.map(s => s.x));
    const minY = Math.min(...coords.map(s => s.y));
    const maxY = Math.max(...coords.map(s => s.y));
    
    const padding = 60;
    const rangeX = maxX - minX || 1;
    const rangeY = maxY - minY || 1;
    
    // First pass: calculate normalized positions
    Object.entries(stations).forEach(([name, station]) => {
        const normalizedX = (station.x - minX) / rangeX;
        const normalizedY = (station.y - minY) / rangeY;
        
        positionMap[name] = {
            x: padding + normalizedX * (canvas.width - 2 * padding),
            y: padding + normalizedY * (canvas.height - 2 * padding)
        };
    });
    
    // Second pass: Force start and end stations to opposite corners for maximum distance
    if (startStation && endStation && positionMap[startStation] && positionMap[endStation]) {
        // Place start station at bottom-left corner
        positionMap[startStation] = {
            x: padding + 20,
            y: canvas.height - padding - 20
        };
        
        // Place end station at top-right corner
        positionMap[endStation] = {
            x: canvas.width - padding - 20,
            y: padding + 20
        };
        
        // Redistribute other stations to avoid overlap with start/end
        Object.entries(stations).forEach(([name, station]) => {
            if (name !== startStation && name !== endStation) {
                const normalizedX = (station.x - minX) / rangeX;
                const normalizedY = (station.y - minY) / rangeY;
                
                // Adjust positions to use middle area
                positionMap[name] = {
                    x: padding + 80 + normalizedX * (canvas.width - 2 * padding - 160),
                    y: padding + 80 + normalizedY * (canvas.height - 2 * padding - 160)
                };
            }
        });
    }

    // Draw connections in background (very light)
    ctx.strokeStyle = '#e0e0e0';
    ctx.lineWidth = 2;
    connections.forEach(([from, to]) => {
        const p1 = positionMap[from];
        const p2 = positionMap[to];
        if (p1 && p2) {
            ctx.beginPath();
            ctx.moveTo(p1.x, p1.y);
            ctx.lineTo(p2.x, p2.y);
            ctx.stroke();
        }
    });

    // Draw paths with distinct colors and labels
    const pathColors = ['#2196F3', '#FF9800', '#9C27B0', '#4CAF50', '#F44336', '#00BCD4', '#FF5722'];
    paths.forEach((path, idx) => {
        const color = pathColors[idx % pathColors.length];
        ctx.strokeStyle = color;
        ctx.lineWidth = 5;
        
        for (let i = 0; i < path.length - 1; i++) {
            const p1 = positionMap[path[i]];
            const p2 = positionMap[path[i + 1]];
            if (p1 && p2) {
                // Draw thick colored line for path
                ctx.beginPath();
                ctx.moveTo(p1.x, p1.y);
                ctx.lineTo(p2.x, p2.y);
                ctx.stroke();
                
                // Draw path label
                const midX = (p1.x + p2.x) / 2;
                const midY = (p1.y + p2.y) / 2;
                
                // Background circle for path number
                ctx.fillStyle = 'white';
                ctx.beginPath();
                ctx.arc(midX, midY, 15, 0, Math.PI * 2);
                ctx.fill();
                
                // Path number
                ctx.fillStyle = color;
                ctx.font = 'bold 14px Arial';
                ctx.textAlign = 'center';
                ctx.textBaseline = 'middle';
                ctx.fillText(`${idx + 1}`, midX, midY);
            }
        }
    });

    // Draw stations
    Object.entries(stations).forEach(([name, station]) => {
        const pos = positionMap[name];
        if (!pos) return;
        
        // Station appearance based on type
        let color = '#607D8B';  // Default gray-blue
        let size = 12;
        let strokeColor = '#fff';
        let strokeWidth = 3;
        
        if (name === startStation) {
            color = '#4CAF50';  // Green
            size = 18;
            strokeWidth = 4;
        } else if (name === endStation) {
            color = '#F44336';  // Red
            size = 18;
            strokeWidth = 4;
        }
        
        // Draw station shadow
        ctx.fillStyle = 'rgba(0, 0, 0, 0.15)';
        ctx.beginPath();
        ctx.arc(pos.x + 2, pos.y + 2, size, 0, Math.PI * 2);
        ctx.fill();
        
        // Draw station outline
        ctx.fillStyle = strokeColor;
        ctx.beginPath();
        ctx.arc(pos.x, pos.y, size + strokeWidth, 0, Math.PI * 2);
        ctx.fill();
        
        // Draw station
        ctx.fillStyle = color;
        ctx.beginPath();
        ctx.arc(pos.x, pos.y, size, 0, Math.PI * 2);
        ctx.fill();
        
        // Draw station name with better styling
        ctx.font = 'bold 13px Arial';
        ctx.textAlign = 'center';
        const textWidth = ctx.measureText(name).width;
        
        // Text background
        ctx.fillStyle = 'rgba(255, 255, 255, 0.95)';
        ctx.fillRect(pos.x - textWidth/2 - 6, pos.y - 35, textWidth + 12, 22);
        
        // Text border
        ctx.strokeStyle = color;
        ctx.lineWidth = 2;
        ctx.strokeRect(pos.x - textWidth/2 - 6, pos.y - 35, textWidth + 12, 22);
        
        // Station name
        ctx.fillStyle = '#333';
        ctx.fillText(name, pos.x, pos.y - 24);
    });

    // Draw trains
    drawTrains(positionMap);
}

// Animation controls
function startAnimation() {
    if (movements.length === 0) return;
    
    if (currentTurn >= movements.length) {
        currentTurn = 0;
    }
    
    isPlaying = true;
    updateLog();
    
    animationInterval = setInterval(() => {
        currentTurn++;
        
        if (currentTurn > movements.length) {
            pauseAnimation();
            return;
        }
        
        drawNetwork();
        updateLog();
    }, 1500);
}

function pauseAnimation() {
    isPlaying = false;
    if (animationInterval) {
        clearInterval(animationInterval);
        animationInterval = null;
    }
}

function resetAnimation() {
    pauseAnimation();
    currentTurn = 0;
    drawNetwork();
    updateLog();
}

// Update movement log
function updateLog() {
    const logEntries = document.getElementById('logEntries');
    logEntries.innerHTML = '';
    
    // Add initial state showing trains at start
    const trainIds = [];
    movements.forEach(turn => {
        turn.forEach(move => {
            if (!trainIds.includes(move.trainId)) {
                trainIds.push(move.trainId);
            }
        });
    });
    
    if (trainIds.length > 0) {
        const startText = trainIds.map(id => `T${id}-${startStation}`).join(' ');
        const startEntry = document.createElement('div');
        startEntry.className = 'log-entry';
        startEntry.style.background = '#e8f5e9';
        startEntry.style.borderLeft = '4px solid #4CAF50';
        
        if (currentTurn === 0) {
            startEntry.style.background = '#e3f2fd';
            startEntry.style.fontWeight = 'bold';
            startEntry.style.borderLeft = '4px solid #2196F3';
        }
        
        startEntry.innerHTML = `<strong>Initial:</strong> ${startText}`;
        logEntries.appendChild(startEntry);
    }
    
    movements.forEach((turn, i) => {
        const turnText = turn.map(m => `T${m.trainId}-${m.station}`).join(' ');
        
        const entry = document.createElement('div');
        entry.className = 'log-entry';
        entry.innerHTML = `<strong>Turn ${i + 1}:</strong> ${turnText}`;
        
        if (i === currentTurn - 1) {
            entry.style.background = '#e3f2fd';
            entry.style.fontWeight = 'bold';
            entry.style.borderLeft = '4px solid #2196F3';
        }
        
        logEntries.appendChild(entry);
    });
    
    // Auto-scroll to current turn
    const currentIndex = currentTurn === 0 ? 0 : currentTurn;
    const currentEntry = logEntries.children[currentIndex];
    if (currentEntry) {
        currentEntry.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }
}

// Show message to user
function showMessage(message, type) {
    const messageDiv = document.getElementById('message');
    messageDiv.innerHTML = `<div class="${type}">${message}</div>`;
    
    setTimeout(() => {
        messageDiv.innerHTML = '';
    }, 5000);
}
