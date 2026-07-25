# HLFA Architecture

## Overview

HLFA is a modular, multi-language vulnerability scanner that doesn't rely on AI/LLM for detection.
Instead, it uses heuristic algorithms, differential analysis, and state tracking to identify vulnerabilities.

## Core Components

### 1. Go Backend (cmd/, internal/)
- **Engine**: Worker pool, task distribution, concurrent scanning
- **Scanner**: Plugin-based scanner system
- **Heuristic**: Core vulnerability detection algorithms
- **Stealth**: Evasion techniques, proxy rotation, header spoofing
- **API**: REST API with authentication and rate limiting

### 2. Python Heuristic Engine (python/)
- **Differential Analyzer**: Compare baseline vs mutated responses
- **State Tracker**: Track session state across requests
- **Flow Mapper**: Map application workflows
- **Context Profiler**: Analyze endpoint context

### 3. Frontend (frontend/)
- **React.js** with TypeScript
- **Real-time updates** via WebSocket
- **Flow visualization** with React Flow
- **Finding management** interface

### 4. Databases
- **PostgreSQL**: Findings, scans, users, workspaces
- **Redis**: Task queue, caching, rate limiting
- **Neo4j**: Application flow graphs, endpoint relationships

## Data Flow

```
User Input (URL)
    ↓
[Recon Engine] → Crawl, parse endpoints
    ↓
[Redis Queue] → Task distribution
    ↓
[Worker Pool] → Concurrent scanning (Go)
    ↓
[Scanner Plugins] → Web, Network, Logic scanners
    ↓
[Heuristic Engine (Python)] → Differential analysis, state tracking
    ↓
[PostgreSQL] → Store findings
    ↓
[WebSocket] → Live updates to Frontend
    ↓
[React Frontend] → Display findings, PoC reproduction
```

## Scanner Plugin System

Each scanner implements the `ScannerPlugin` interface:

```go
type ScannerPlugin interface {
    Name() string
    Category() string
    CanHandle(ctx context.Context, target Target) bool
    Scan(ctx context.Context, target Target) ([]UniversalFinding, error)
}
```

Scanner types:
- **Web Scanner**: HTML parsing, endpoint discovery, injection testing
- **Network Scanner**: Port scanning, service enumeration
- **Logic Scanner**: Business logic flaws, state manipulation

## Heuristic Principles

1. **Observe Structure**: Parse and classify endpoints by context
2. **Understand Flow**: Build application workflow graphs
3. **Manipulate State**: Send mutated requests, compare responses
4. **Chain Findings**: Combine findings into exploit chains
5. **Prove**: Generate PoC automatically

## Multi-Tenancy

- Workspace-based isolation
- RBAC (Role-Based Access Control)
- Rate limiting per tenant
- Quota management

## Security Features

- JWT-based authentication
- API key management
- Encrypted sensitive data (auth tokens, cookies)
- Audit logging
- CORS protection

## Scaling

- Horizontal scaling via Docker Compose / Kubernetes
- Redis-backed task queue
- Worker pool configuration
- Load balancing support
