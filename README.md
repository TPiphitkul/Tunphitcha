# go-adaptive-gw

Starter GitHub-style project for **A Go-based Adaptive Security Model for Microservices Gateways**.

This repository contains:
- A minimal Gateway (Go) with an adaptive layer (profiler, heuristic risk scorer, policy decision, enforcement).
- Three mock microservices (`user`, `catalog`, `order`).
- Docker & Docker Compose files to run everything locally.
- k6 scripts for baseline and attack scenarios.

## Quick start (local, requires Docker & Docker Compose)

```bash
# build and start containers
docker compose up --build

# gateway listens on http://localhost:8080
# health: http://localhost:8080/health
```

Run baseline load:
```bash
# install k6 and run
k6 run deploy/k6/baseline.js
```

Run attack load:
```bash
k6 run deploy/k6/attack.js
```

## Project layout

```
go-adaptive-gw/
├─ cmd/gateway/
│  └─ main.go
├─ internal/
│  ├─ profiler/
│  │  └─ profiler.go
│  ├─ risk/
│  │  └─ heuristic.go
│  ├─ policy/
│  │  └─ policy.go
│  └─ enforcement/
│     └─ enforce.go
├─ services/
│  ├─ user/
│  │  └─ main.go
│  ├─ catalog/
│  │  └─ main.go
│  └─ order/
│     └─ main.go
├─ deploy/
│  ├─ docker-compose.yml
│  └─ k6/
│     ├─ baseline.js
│     └─ attack.js
├─ Dockerfile.gateway
├─ Dockerfile.service
└─ README.md
```

## Notes
- The adaptive risk scorer is heuristic by default. You can replace `internal/risk` with an ML-based scorer (exposed via HTTP) later.
- Redis is included in docker-compose as a placeholder for counters and shared state (e.g., request-per-minute).
