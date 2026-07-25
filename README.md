# HLFA - Heuristic Logic & Flow Analyzer

HLFA adalah tools Bug Bounty & Vulnerability Scanner yang TIDAK menggunakan AI/LLM.
HLFA menggunakan algoritma HEURISTIK, DIFFERENTIAL ANALYSIS, dan STATE TRACKING
untuk menemukan kerentanan yang TIDAK DIKETAHUI (Zero-day, Logic Flaws, IDOR, SSRF).

## Arsitektur

- **Core Engine**: Go (Golang) - concurrency, network speed, worker pool
- **Heuristic Logic**: Python (FastAPI) - parsing kompleks, graph logic
- **Database**: PostgreSQL + Redis + Neo4j
- **Frontend**: React.js + TypeScript + Vite + TailwindCSS
- **Real-time**: Socket.io (WebSocket)
- **Browser Automation**: Playwright
- **Containerization**: Docker + Docker Compose + Helm

## Prinsip Kerja HLFA

1. **Amati Struktur** → Parse HTML/JSON, klasifikasikan endpoint
2. **Pahami Alur Bisnis** → Bangun Graph alur aplikasi
3. **Manipulasi State** → Kirim request yang dimutasi, bandingkan response
4. **Rantai Temuan** → Gabungkan temuan kecil menjadi exploit chain
5. **Buktikan** → Generate PoC otomatis

## Quick Start

```bash
# Development
make dev

# Production
make build
docker-compose -f deploy/docker-compose.prod.yml up
```

## Project Structure

```
hlfa/
├── cmd/                    # CLI & Server entry points
├── internal/               # Core packages
├── pkg/                    # Shared packages
├── python/                 # Heuristic engine (Python)
├── frontend/               # React.js UI
├── plugins/                # External plugins
├── deploy/                 # Docker & K8s configs
├── migrations/             # Database migrations
└── docs/                   # Documentation
```

## Documentation

See `docs/` folder for detailed documentation.
