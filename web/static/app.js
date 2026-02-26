// Global state
let canvas, ctx;
let stations = {};
let connections = [];
let paths = [];
let movements = [];
let currentTurn = 0;
let animationInterval = null;
let positionMap = {};
let startStation = '';
let endStation = '';
let trainPaths = {};
let trainPositions = {};
let animationProgress = 0;

// Colors
const BG_COLOR = '#0d0e12';
const PATH_COLOR = '#ffffff'; // Single white color for all routes
const STATION_COLOR = '#888888'; // Gray for intermediate stations
const TRAIN_COLOR = '#00ff88'; // Green for all trains
const START_COLOR = '#ffffff'; // White for start
const END_COLOR = '#ff3366'; // Red for end

// Initialize
window.addEventListener('load', () => {
    canvas = document.getElementById('canvas');
    ctx = canvas.getContext('2d');
    
    resizeCanvas();
    window.addEventListener('resize', resizeCanvas);
    
    loadMaps();
});

function resizeCanvas() {
    canvas.width = canvas.offsetWidth;
    canvas.height = canvas.offsetHeight;
    if (Object.keys(stations).length > 0) {
        drawNetwork();
    }
}

// Load available maps
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

// Run pathfinder
async function runPathfinder() {
    const mapFile = document.getElementById('mapSelect').value;
    const numTrains = parseInt(document.getElementById('numTrains').value);

    if (!mapFile) {
        showMessage('Please select a map', 'error');
        return;
    }
    if (!numTrains || numTrains < 1) {
        showMessage('Please enter a valid number of trains', 'error');
        return;
    }

    // Stop any existing animation
    if (animationInterval) {
        clearInterval(animationInterval);
        animationInterval = null;
    }

    showMessage('Running pathfinder...', 'info');

    try {
        const response = await fetch('/api/pathfind', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                mapFile,
                startStation: '',
                endStation: '',
                numTrains
            })
        });

        const data = await response.json();

        if (!data.success) {
            showMessage('Error: ' + data.error, 'error');
            return;
        }

        // Store data
        stations = data.stations;
        paths = data.paths;
        movements = data.movements;
        startStation = data.startStation;
        endStation = data.endStation;
        currentTurn = 0;
        animationProgress = 0;

        // Assign each train to a path
        trainPaths = {};
        for (let i = 0; i < numTrains; i++) {
            const pathIndex = i % paths.length;
            trainPaths[i + 1] = { pathIndex: pathIndex };
        }

        // Build all connections from paths
        connections = [];
        const seenConnections = new Set();
        
        paths.forEach(path => {
            for (let i = 0; i < path.length - 1; i++) {
                const conn = [path[i], path[i + 1]].sort().join('|');
                if (!seenConnections.has(conn)) {
                    seenConnections.add(conn);
                    connections.push([path[i], path[i + 1]]);
                }
            }
        });

        // Calculate positions
        calculatePositions();

        // Draw initial state
        drawNetwork();

        // Show stats
        showStats(data);

        // Show controls
        document.getElementById('controls').style.display = 'block';
        document.getElementById('movementLog').style.display = 'block';
        updateMovementLog();
        updateTurnInfo();

        showMessage(`Found ${paths.length} path(s) in ${movements.length} turn(s) from ${startStation} to ${endStation}`, 'success');

        // Start animation with delay
        setTimeout(() => startAnimation(), 500);

    } catch (error) {
        showMessage('Request failed: ' + error.message, 'error');
    }
}

// Calculate station positions on canvas
function calculatePositions() {
    const stationArray = Object.values(stations);
    
    if (stationArray.length === 0) return;

    // Find min/max coordinates
    let minX = Infinity, maxX = -Infinity;
    let minY = Infinity, maxY = -Infinity;

    stationArray.forEach(s => {
        if (s.x < minX) minX = s.x;
        if (s.x > maxX) maxX = s.x;
        if (s.y < minY) minY = s.y;
        if (s.y > maxY) maxY = s.y;
    });

    const rangeX = maxX - minX || 1;
    const rangeY = maxY - minY || 1;

    const padding = 80;
    const width = canvas.width - 2 * padding;
    const height = canvas.height - 2 * padding;

    positionMap = {};
    stationArray.forEach(s => {
        const x = padding + ((s.x - minX) / rangeX) * width;
        const y = padding + ((s.y - minY) / rangeY) * height;
        positionMap[s.name.toLowerCase()] = { x, y };
    });
}

// Draw the network
function drawNetwork() {
    ctx.fillStyle = BG_COLOR;
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    if (Object.keys(stations).length === 0) return;

    // Draw all connections in white
    ctx.strokeStyle = PATH_COLOR;
    ctx.lineWidth = 3;
    connections.forEach(([from, to]) => {
        const p1 = positionMap[from.toLowerCase()];
        const p2 = positionMap[to.toLowerCase()];
        if (p1 && p2) {
            ctx.beginPath();
            ctx.moveTo(p1.x, p1.y);
            ctx.lineTo(p2.x, p2.y);
            ctx.stroke();
        }
    });

    // Draw all stations
    Object.entries(stations).forEach(([name, station]) => {
        const pos = positionMap[name.toLowerCase()];
        if (!pos) return;

        // Determine color and size
        let color = STATION_COLOR;
        let radius = 10;
        if (name.toLowerCase() === startStation) {
            color = START_COLOR;
            radius = 16;
        } else if (name.toLowerCase() === endStation) {
            color = END_COLOR;
            radius = 16;
        }

        // Draw station circle
        ctx.fillStyle = color;
        ctx.beginPath();
        ctx.arc(pos.x, pos.y, radius, 0, Math.PI * 2);
        ctx.fill();
        
        // Draw station border
        ctx.strokeStyle = BG_COLOR;
        ctx.lineWidth = 2;
        ctx.beginPath();
        ctx.arc(pos.x, pos.y, radius, 0, Math.PI * 2);
        ctx.stroke();

        // Draw station name
        ctx.fillStyle = '#ffffff';
        ctx.font = 'bold 13px monospace';
        ctx.textAlign = 'center';
        
        const yOffset = (name.toLowerCase() === startStation || name.toLowerCase() === endStation) ? -25 : -18;
        ctx.fillText(station.name, pos.x, pos.y + yOffset);
    });

    // Draw trains
    drawTrainsAnimated();
}

// Draw trains with smooth animation along tracks
function drawTrainsAnimated() {
    if (movements.length === 0) return;

    // Get current turn movements
    const turnMovements = currentTurn < movements.length ? movements[currentTurn] : [];
    
    // Group trains by their target station for this turn
    const movingTrains = new Map();
    turnMovements.forEach(move => {
        movingTrains.set(move.trainId, move.station);
    });

    // Get all train IDs
    const allTrains = Object.keys(trainPaths).map(id => parseInt(id));

    allTrains.forEach(trainId => {
        const targetStation = movingTrains.get(trainId);

        // Determine train's current and next position
        let currentPos, nextPos;
        
        if (targetStation) {
            // Train is moving this turn
            const prevStation = getPreviousStation(trainId, currentTurn);
            currentPos = positionMap[prevStation.toLowerCase()];
            nextPos = positionMap[targetStation.toLowerCase()];
            
            if (currentPos && nextPos) {
                // Interpolate position based on animation progress
                const x = currentPos.x + (nextPos.x - currentPos.x) * animationProgress;
                const y = currentPos.y + (nextPos.y - currentPos.y) * animationProgress;
                
                drawTrain(x, y, trainId);
            }
        } else {
            // Train is waiting at a station
            const station = getCurrentStation(trainId, currentTurn);
            const pos = positionMap[station.toLowerCase()];
            
            if (pos) {
                drawTrain(pos.x, pos.y, trainId);
            }
        }
    });
}

// Draw a single train
function drawTrain(x, y, trainId) {
    // Draw train circle with glow effect
    ctx.shadowBlur = 15;
    ctx.shadowColor = TRAIN_COLOR;
    ctx.fillStyle = TRAIN_COLOR;
    ctx.beginPath();
    ctx.arc(x, y, 12, 0, Math.PI * 2);
    ctx.fill();
    
    // Reset shadow
    ctx.shadowBlur = 0;
    
    // Draw train border
    ctx.strokeStyle = '#000';
    ctx.lineWidth = 2;
    ctx.beginPath();
    ctx.arc(x, y, 12, 0, Math.PI * 2);
    ctx.stroke();

    // Draw train ID
    ctx.fillStyle = '#000';
    ctx.font = 'bold 11px monospace';
    ctx.textAlign = 'center';
    ctx.fillText('T' + trainId, x, y + 4);
}

// Get the station where a train was before this turn
function getPreviousStation(trainId, turn) {
    if (turn === 0) return startStation;
    
    // Look through previous movements to find where this train was
    for (let t = turn - 1; t >= 0; t--) {
        const move = movements[t].find(m => m.trainId === trainId);
        if (move) return move.station;
    }
    
    return startStation;
}

// Get the current station of a train
function getCurrentStation(trainId, turn) {
    // Look through movements up to current turn
    for (let t = turn; t >= 0; t--) {
        const move = movements[t].find(m => m.trainId === trainId);
        if (move) return move.station;
    }
    
    return startStation;
}

// Animation with slower speed for clarity
function startAnimation() {
    if (animationInterval) return;
    
    const fps = 60;
    const turnDuration = 2000; // 2 seconds per turn
    const frameTime = 1000 / fps;
    const progressPerFrame = frameTime / turnDuration;
    
    animationInterval = setInterval(() => {
        animationProgress += progressPerFrame;
        
        if (animationProgress >= 1) {
            animationProgress = 0;
            currentTurn++;
            
            if (currentTurn >= movements.length) {
                clearInterval(animationInterval);
                animationInterval = null;
                currentTurn = movements.length - 1;
                animationProgress = 1;
                showMessage('Animation complete! All trains reached destination.', 'success');
            }
            
            updateTurnInfo();
            updateMovementLog();
        }
        
        drawNetwork();
    }, frameTime);
}

function updateTurnInfo() {
    const info = document.getElementById('turnInfo');
    info.textContent = `Turn ${currentTurn + 1} of ${movements.length}`;
}

// Show stats
function showStats(data) {
    const stats = document.getElementById('stats');
    const content = document.getElementById('statsContent');
    
    content.innerHTML = `
        <div style="display: flex; flex-direction: column; gap: 12px; font-size: 13px;">
            <div><strong>Start:</strong> ${startStation}</div>
            <div><strong>End:</strong> ${endStation}</div>
            <div><strong>Paths:</strong> ${data.paths.length}</div>
            <div><strong>Turns:</strong> ${data.movements.length}</div>
            <div><strong>Trains:</strong> ${document.getElementById('numTrains').value}</div>
            <div><strong>Stations:</strong> ${Object.keys(data.stations).length}</div>
        </div>
    `;
    
    stats.style.display = 'block';
}

// Update movement log
function updateMovementLog() {
    const logContent = document.getElementById('logContent');
    logContent.innerHTML = '';

    movements.forEach((turnMovements, idx) => {
        const turnDiv = document.createElement('div');
        turnDiv.style.padding = '8px';
        turnDiv.style.marginBottom = '4px';
        turnDiv.style.background = idx === currentTurn ? '#00ff8822' : '#16181d';
        turnDiv.style.borderRadius = '4px';
        turnDiv.style.borderLeft = idx === currentTurn ? '3px solid #00ff88' : '3px solid transparent';

        const movementsText = turnMovements
            .map(m => `T${m.trainId}-${m.station}`)
            .join(' ');

        turnDiv.textContent = `Turn ${idx + 1}: ${movementsText}`;
        logContent.appendChild(turnDiv);

        if (idx === currentTurn) {
            turnDiv.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
        }
    });
}

// Show message
function showMessage(text, type) {
    const messageDiv = document.getElementById('message');
    messageDiv.textContent = text;
    
    const colors = {
        error: '#ff3366',
        success: '#00ff88',
        info: '#ffffff'
    };
    
    messageDiv.style.background = colors[type] + '22';
    messageDiv.style.color = colors[type];
    messageDiv.style.border = `1px solid ${colors[type]}`;
}
