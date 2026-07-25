# ARCHITECTURE.md

## System Architecture

### High-Level Overview

HLFA consists of three main components:

1. **Frontend** - React/TypeScript SPA with real-time WebSocket updates
2. **Backend** - Go HTTP API server with scanning orchestration
3. **Python Engine** - FastAPI-based heuristic analysis service

### Data Flow

```
1. User starts scan via Frontend
   ↓
2. Frontend sends POST /api/v1/scans to Backend
   ↓
3. Backend validates target and creates Scan record in PostgreSQL
   ↓
4. Backend submits tasks to worker pool
   ↓
5. Workers execute scanners (Web, Network, etc.)
   ↓
6. Scanners collect baseline and mutated responses
   ↓
7. Responses sent to Python Engine for differential analysis
   ↓
8. Analysis results returned as Findings
   ↓
9. Findings stored in PostgreSQL and cached in Redis
   ↓
10. Frontend receives real-time updates via WebSocket
```

### Component Details

#### Go Backend

**Responsibilities:**
- REST API handling (scan creation, status, findings)
- Worker pool management (concurrent scanning)
- Scanner orchestration
- Database operations (PostgreSQL)
- Redis queue management
- WebSocket server for real-time updates

**Key Packages:**
- `internal/api` - HTTP handlers and routing
- `internal/engine` - Scan orchestration and worker pool
- `internal/scanner` - Web and network scanners
- `internal/repository` - Database access layer

#### Python Engine

**Responsibilities:**
- Differential response analysis
- Advanced state tracking
- Context profiling
- Vulnerability detection heuristics

**Key Modules:**
- `differential_analyzer.py` - Response comparison
- `advanced_state_tracker.py` - Session state monitoring
- `context_profiler.py` - Endpoint analysis

#### Frontend

**Responsibilities:**
- User interface and interactions
- Real-time scan monitoring
- Flow graph visualization
- Finding management and review

**Key Pages:**
- Dashboard - Quick scan and stats
- Findings - Finding list and detail view
- Flow Graph - Application flow visualization
- Workspace - Team management
- Settings - Configuration

### Database Schema

#### PostgreSQL Tables

**scans**
- id (UUID)
- workspace_id (UUID)
- status (enum: pending, running, completed, failed)
- progress (int 0-100)
- findings_count (int)
- started_at, completed_at (timestamp)

**findings**
- id (UUID)
- scan_id (UUID)
- target (string)
- severity (enum: CRITICAL, HIGH, MEDIUM, LOW, INFO)
- category (enum: IDOR, SQLi, XSS, SSRF, LogicFlaw, etc.)
- title, description (text)
- evidence (JSONB)
- reproduction_steps (text[])
- poc_curl (text)
- status (enum: OPEN, CONFIRMED, FALSE_POSITIVE, FIXED)
- assigned_to (UUID)

**workspaces**
- id (UUID)
- name (string)
- owner_id (UUID)
- description, settings (JSONB)

**users**
- id (UUID)
- username, email (string)
- password_hash (string)
- role (enum: admin, user, analyst)

#### Redis Keys

- `scan:{scanId}` - Scan metadata and status
- `queue:recon_tasks` - Reconnaissance task queue
- `queue:web_scan_tasks` - Web scanning task queue
- `queue:network_scan_tasks` - Network scanning task queue
- `rate:{userId}` - Rate limit counters

#### Neo4j Nodes and Relationships

**Nodes:**
- Endpoint (url, method, params)
- Session (token, user_id)
- Parameter (name, type, sensitivity)
- Flow (name, steps)

**Relationships:**
- CONNECTS_TO - Endpoint connections
- REQUIRES_AUTH - Authentication dependencies
- PRODUCES - Output relationships
- CONSUMES - Input relationships

### Scanning Pipeline

```
┌─────────────────────────┐
│   1. Task Creation      │
│  - Parse target         │
│  - Validate URL         │
│  - Create Task object   │
└───────────┬─────────────┘
            ↓
┌─────────────────────────┐
│  2. Queue Task          │
│  - Submit to worker     │
│  - Track in Redis       │
└───────────┬─────────────┘
            ↓
┌─────────────────────────┐
│  3. Scanner Selection   │
│  - Match target type    │
│  - Load appropriate     │
│    scanner module       │
└───────────┬─────────────┘
            ↓
┌─────────────────────────┐
│  4. Baseline Request    │
│  - Send GET request     │
│  - Record response      │
│  - Extract metrics      │
└───────────┬─────────────┘
            ↓
┌─────────────────────────┐
│  5. Mutation Testing    │
│  - Modify parameters    │
│  - Inject payloads      │
│  - Track state changes  │
└───────────┬─────────────┘
            ↓
┌─────────────────────────┐
│  6. Differential        │
│     Analysis            │
│  - Compare responses    │
│  - Detect anomalies     │
│  - Generate findings    │
└───────────┬─────────────┘
            ↓
┌─────────────────────────┐
│  7. Finding Storage     │
│  - Save to PostgreSQL   │
│  - Cache in Redis       │
│  - Notify subscribers   │
└─────────────────────────┘
```

### Deployment Architecture

#### Docker Compose
- All services in one Docker network
- Single-host deployment suitable for development
- Services auto-start and health-checked

#### Kubernetes
- Microservices-oriented
- Horizontal pod autoscaling
- Load balanced with ingress
- Persistent volumes for databases
- Secrets management
- Resource quotas per pod

### Performance Considerations

1. **Concurrency** - Worker pool defaults to 50 concurrent scanners
2. **Caching** - Redis cache layer reduces database hits
3. **Connection Pooling** - PostgreSQL and Redis use connection pools
4. **Request Timeouts** - Configurable per target
5. **Rate Limiting** - Per-user and per-IP limits

### Security Architecture

1. **Authentication** - JWT tokens with 24-hour expiry
2. **Authorization** - Role-based access control (RBAC)
3. **Encryption** - TLS for in-transit, AES-256 for at-rest
4. **Audit Logging** - All API calls and findings changes logged
5. **Input Validation** - All inputs validated on backend
6. **CORS** - Configurable cross-origin policies

### Extensibility

1. **Plugin System** - Load custom scanners dynamically
2. **Webhook Notifications** - Slack, Discord, Jira integration
3. **Custom Heuristics** - Python-based rule engine
4. **API-First** - All functionality available via REST API
