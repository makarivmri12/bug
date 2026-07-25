# HLFA - Heuristic Logic & Flow Analyzer

## 🎯 Overview

HLFA is a production-ready Bug Bounty & Vulnerability Scanner that uses advanced heuristic analysis instead of traditional AI/LLM approaches. It detects zero-day vulnerabilities, logic flaws, IDOR issues, and SSRF attacks like a senior security engineer.

## 🚀 Features

### Core Scanning Engine
- **Differential Analysis** - Compare baseline vs mutated responses to detect anomalies
- **State Tracking** - Monitor session state changes and detect exploitation chains
- **Flow Mapping** - Build application flow graphs to identify attack vectors
- **Heuristic Detection** - Advanced pattern matching for vulnerability identification

### Vulnerability Detection
- IDOR (Insecure Direct Object Reference)
- SQL Injection
- Cross-Site Scripting (XSS)
- Server-Side Request Forgery (SSRF)
- Logic Flaws
- Authentication Bypass
- Broken Access Control
- Timing Attacks

### Architecture
- **Go Backend** - High-performance concurrent scanning with worker pools
- **Python Heuristics** - Advanced analysis engine with differential comparison
- **React Frontend** - Real-time dashboard with flow visualization
- **PostgreSQL** - Persistent finding storage
- **Redis** - Task queue and caching
- **Neo4j** - Application flow graph database

## 📦 Installation

### Docker Compose (Development)

```bash
git clone https://github.com/makarivmri12/bug.git
cd bug
make docker-up
```

Access the application:
- Frontend: http://localhost:3000
- API: http://localhost:8080
- Python Engine: http://localhost:5000

### Kubernetes Deployment

```bash
kubectl create namespace hlfa
helm install hlfa ./deploy/helm -n hlfa -f deploy/helm/values-prod.yaml
```

## 🛠️ Quick Start

### CLI Scanning

```bash
# Build the CLI
make build-cli

# Run a scan
./hlfa-cli --target https://example.com --method GET --output json
```

### API-based Scanning

```bash
# Start a scan
curl -X POST http://localhost:8080/api/v1/scans \
  -H "Content-Type: application/json" \
  -d '{
    "workspace_id": "default",
    "targets": [
      {"url": "https://example.com", "type": "web"}
    ],
    "scanner_types": ["web", "network"]
  }'

# Get findings
curl http://localhost:8080/api/v1/scans/{scanId}/findings
```

## 📊 Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (React + TypeScript)            │
│         (Dashboard, Flow Viz, Findings Explorer)           │
└────────────────────────┬────────────────────────────────────┘
                         │
                    HTTP/WebSocket
                         │
┌────────────────────────▼────────────────────────────────────┐
│              Backend API Server (Go)                        │
│  (REST Handlers, Authentication, Orchestration)            │
└────────────────────────┬────────────────────────────────────┘
                         │
        ┌────────────────┼────────────────┐
        │                │                │
   ┌────▼─────┐  ┌──────▼────┐  ┌───────▼─────┐
   │ PostgreSQL│  │   Redis   │  │   Neo4j     │
   │ (Storage) │  │  (Queue)  │  │ (Graphs)    │
   └───────────┘  └───────────┘  └─────────────┘
        │
   ┌────▼─────────────────────┐
   │  Scan Engine (Worker Pool)│
   │  (Go Concurrency)         │
   └────┬─────────────────────┐
        │                      │
    ┌───▼──────┐        ┌─────▼──────┐
    │ Web      │        │  Network   │
    │ Scanner  │        │  Scanner   │
    └──────────┘        └────────────┘
        │                      │
        │      HTTP Calls      │
        └──────────┬───────────┘
                   │
        ┌──────────▼──────────┐
        │ Python Heuristic    │
        │ Engine (FastAPI)    │
        │                     │
        │ - Differential      │
        │   Analysis          │
        │ - State Tracking    │
        │ - Context Profiling │
        └─────────────────────┘
```

## 🔌 Plugin System

Extend HLFA with custom scanners:

```go
// plugins/custom_scanner.go
package main

import (
    "github.com/makarivmri12/bug/internal/engine"
    "github.com/makarivmri12/bug/pkg/models"
)

type CustomScanner struct{}

func (cs *CustomScanner) Name() string { return "custom-scanner" }
func (cs *CustomScanner) Category() string { return "web" }
func (cs *CustomScanner) CanHandle(ctx context.Context, target models.Target) bool {
    return target.Type == models.TargetTypeWeb
}
func (cs *CustomScanner) Scan(ctx context.Context, target models.Target) ([]models.UniversalFinding, error) {
    // Your scanning logic
    return findings, nil
}

var Plugin engine.ScannerPlugin = &CustomScanner{}
```

## 🔔 Notifications

Integrate with external services:

```go
notificationService := webhook.NewNotificationService(logger)
notificationService.Configure(webhook.ChannelSlack, map[string]string{
    "webhook_url": "https://hooks.slack.com/...",
})
notificationService.NotifyFinding(finding)
```

## 📈 Performance

- **Concurrent Workers**: 50+ parallel scans
- **Response Time**: < 200ms for baseline requests
- **Throughput**: 1000+ requests/min per worker
- **Memory**: ~512MB base + 10MB per worker

## 🧪 Testing

```bash
# Run all tests
make test

# Run backend tests
make test-backend

# Run frontend tests
make test-frontend

# Run integration tests
make test-integration
```

## 📚 Documentation

- [API Documentation](./docs/API.md)
- [Deployment Guide](./deploy/README.md)
- [Architecture Guide](./docs/ARCHITECTURE.md)
- [Contributing Guide](./CONTRIBUTING.md)

## 🔐 Security

- All communication over HTTPS/TLS
- JWT-based authentication
- Role-based access control (RBAC)
- Data encryption at rest
- Audit logging for all actions

## 📝 License

MIT License - see LICENSE file for details

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for details.

## 📧 Support

For issues, feature requests, or questions:
- GitHub Issues: https://github.com/makarivmri12/bug/issues
- Email: support@hlfa.dev

## 🙏 Acknowledgments

Built with ❤️ by the HLFA Team
