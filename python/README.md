# HLFA Python Heuristic Engine

This directory contains the Python-based heuristic analysis engine for HLFA.

## Components

### 1. Differential Analyzer (`differential_analyzer.py`)
- Compares baseline vs mutated responses
- Detects response anomalies
- Identifies timing attacks
- Calculates similarity scores
- Generates unified diffs

### 2. Advanced State Tracker (`advanced_state_tracker.py`)
- Tracks session state changes over time
- Records token/cookie modifications
- Detects state anomalies
- Maintains state timeline
- Identifies unusual patterns

### 3. Context Profiler (`context_profiler.py`)
- Extracts endpoint context
- Identifies sensitive parameters
- Detects authentication requirements
- Classifies input types
- Profiles endpoint characteristics

### 4. FastAPI Server (`api/main.py`)
- RESTful API for heuristic analysis
- Real-time analysis endpoints
- State tracking endpoints
- Endpoint profiling service
- Health check monitoring

## Installation

```bash
pip install -r requirements.txt
```

## Running

```bash
uvicorn api.main:app --host 0.0.0.0 --port 5000
```

## API Endpoints

### Differential Analysis
- `POST /api/v1/differential/baseline` - Set baseline response
- `POST /api/v1/differential/analyze` - Analyze against baseline

### State Tracking
- `POST /api/v1/state/capture` - Capture state snapshot
- `GET /api/v1/state/timeline/{user_id}` - Get state timeline

### Context Profiling
- `POST /api/v1/profile/endpoint` - Profile an endpoint
- `GET /api/v1/profile/endpoints` - Get all profiled endpoints

### Health
- `GET /health` - Health check

## Example Usage

```bash
# Set baseline
curl -X POST http://localhost:5000/api/v1/differential/baseline \
  -H "Content-Type: application/json" \
  -d '{"status": 200, "headers": {}, "body": "response", "response_time_ms": 150}'

# Analyze response
curl -X POST http://localhost:5000/api/v1/differential/analyze \
  -H "Content-Type: application/json" \
  -d '{"status": 200, "headers": {}, "body": "modified", "response_time_ms": 200}'
```
