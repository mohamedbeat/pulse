# Health Check Monitor — Go Project Roadmap

> A robust, extensible health monitoring tool written in Go.  
> **Status**: Planning  
> **Target**: Production-ready monitoring solution for services/endpoints.

---

## Phase 1: Foundation (Week 1–2)

**Goal**: ✅ Basic working health checker with config and concurrency.

### Week 1
- [ ] **Project scaffolding**
  - Define project layout (`cmd/`, `internal/`, `pkg/`, `configs/`)
  - Set up `main.go` with entrypoint
- [ ] **Core health check logic**
  - Single HTTP(S) GET check with timeout
  - Parse response status code + latency
  - Basic stdout logging
- [ ] **Configuration**
  - YAML config schema (e.g., `config.yaml`)
  - Load endpoints, intervals, timeouts
  - Validate required fields

### Week 2
- [ ] **Concurrency & scheduling**
  - Goroutine-based concurrent checks
  - `time.Ticker` for repeating checks (e.g., `10m`)
  - Channel-based result aggregation
- [ ] **Result handling**
  - Structured `Result` type (status, latency, timestamp, error)
  - In-memory storage (slice/map)
  - Basic error recovery (retry once on failure)

✅ **Phase 1 Deliverable**: CLI tool that monitors 1+ endpoints indefinitely and logs status.

---

## Phase 2: Core Features (Week 3–4)

**Goal**: 🚀 Production-viable with persistence, alerts, and UI.

### Week 3
- [ ] **Multiple endpoint types**
  - HTTP/HTTPS (GET, POST, HEAD)
  - Custom headers & auth (Bearer, Basic)
  - TCP port checks
  - Response body matching (regex/string)
- [ ] **Persistence**
  - SQLite backend (`/data/monitor.db`)
  - Schema: `checks`, `endpoints`, `alerts`
  - Automatic cleanup (retain last 30 days)
- [ ] **Metrics**
  - Uptime % calculation
  - Avg/min/max response time
  - Downtime duration tracking

### Week 4
- [ ] **Alerting**
  - Thresholds (e.g., `3 consecutive failures`)
  - Recovery detection
  - Console + email (SMTP) alerts
- [ ] **Web dashboard**
  - Embedded HTTP server (`:8080`)
  - Real-time status page (HTML + minimal JS)
  - Endpoint list with status badges

✅ **Phase 2 Deliverable**: Self-contained binary with SQLite, email alerts, and live dashboard.

---

## Phase 3: Advanced Features (Week 5–6)

**Goal**: 🌐 Enterprise-grade extensibility and integrations.

### Week 5
- [ ] **Notification integrations**
  - Slack/Discord webhooks
  - PagerDuty (v2 Events API)
  - Custom webhook support
- [ ] **Configuration enhancements**
  - Hot-reload on config change (fsnotify)
  - Env var overrides (`HC_INTERVAL=300s`)
  - Validate config on load
- [ ] **Advanced checks**
  - SSL certificate expiry (< 30 days)
  - DNS lookup checks
  - Status code ranges (`2xx`, `3xx`)

### Week 6
- [ ] **API layer**
  - REST API (`/api/v1/status`, `/api/v1/history`)
  - CRUD for endpoints (add/remove/update)
  - JWT-based authentication
- [ ] **Performance & resilience**
  - HTTP client pooling & keep-alive
  - Circuit breaker for failing endpoints
  - CPU/memory profiling

✅ **Phase 3 Deliverable**: API-first monitor with multi-channel alerts and config reload.

---

## Phase 4: Production Readiness (Week 7–8)

**Goal**: 🛡️ Secure, observable, and deployable.

### Week 7
- [ ] **Observability**
  - Prometheus metrics (`/metrics`)
    - `healthcheck_up`, `healthcheck_latency_seconds`
  - Structured JSON logging (Zap/Slog)
  - Tracing (OpenTelemetry)
- [ ] **Security**
  - TLS for web/API (`--tls-cert`, `--tls-key`)
  - RBAC for API/dashboard
  - Rate limiting (5 req/s per IP)

### Week 8
- [ ] **Deployment**
  - Dockerfile (multi-stage, scratch base)
  - Helm chart (for Kubernetes)
  - systemd service template
- [ ] **Quality & docs**
  - Test coverage ≥ 80% (unit + integration)
  - CLI help (`--help`, `--version`)
  - `README.md` with examples
  - `docs/` (config reference, API spec)

✅ **Phase 4 Deliverable**: Production-ready, containerized, observable, and documented.

---

## Milestones & Success Criteria

| Phase | Deadline | Success Criteria |
|-------|----------|------------------|
| **1** | Week 2   | ✅ `./health-monitor --config config.yaml` runs indefinitely, checks endpoints, logs to stdout |
| **2** | Week 4   | ✅ SQLite storage + email alerts + live dashboard on `:8080` |
| **3** | Week 6   | ✅ API + webhook alerts + config hot-reload |
| **4** | Week 8   | ✅ Docker image + Prometheus metrics + ≥80% test coverage |

---

## Daily Workflow

| Day       | Focus                           |
|-----------|---------------------------------|
| **Mon**   | Plan features, review PRs       |
| **Tue**   | Core implementation             |
| **Wed**   | Edge cases & error handling     |
| **Thu**   | Testing + bug fixes             |
| **Fri**   | Docs + next week prep           |

---

## Tech Stack

| Component       | Choice                          |
|-----------------|---------------------------------|
| Language        | Go                              |
| Config          | YAML (Viper)                    |
| Storage         | SQLite (embedded)               |
| Web Server      | `net/http`                      |
| CLI             | `urfave/cli/v2`                 |
| Logging         | `log/slog`                      |
| Testing         | `testify`, `httpexpect`         |
| Build           | Go modules + `goreleaser` (v2)  |

---

> 📝 **Tip**: Start small—even Phase 1 is valuable! Iterate fast, ship often.
