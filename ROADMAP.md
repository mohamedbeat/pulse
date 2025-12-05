# Health Check Monitor — Go Project Roadmap

> A robust, extensible health monitoring tool written in Go.  
> **Status**: In Development  
> **Target**: Production-ready monitoring solution for services/endpoints.

---

## Phase 1: Foundation

**Goal**: ✅ Basic working health checker with config and concurrency.

### Project Structure
- [x] Set up `main.go` with entrypoint
- [ ] Define project layout (`cmd/`, `internal/`, `pkg/`, `configs/`)

### Core Health Check Logic
- [x] HTTP(S) GET/POST/PUT/PATCH/DELETE checks with timeout
- [x] Parse response status code + latency
- [x] Status code classification (2xx = up, 5xx = down, others = degraded)
- [x] Structured logging with JSON output

### Configuration
- [x] YAML config schema (`pulse.yml`)
- [x] Load endpoints, intervals, timeouts
- [x] Validate required fields
- [x] Global defaults with endpoint overrides
- [x] Config file path resolution (supports directories and files)

### Concurrency & Scheduling
- [x] Goroutine-based concurrent checks (one per endpoint)
- [x] `time.Ticker` for repeating checks per endpoint interval
- [x] Buffered channel-based result aggregation
- [ ] Graceful shutdown support (`stop` channel)

### Result Handling
- [x] Structured `Result` type (status, latency, timestamp, error, message)
- [x] Status types: `up`, `down`, `unreachable`, `degraded`
- [ ] In-memory storage (slice/map)
- [ ] Basic error recovery (retry once on failure)

✅ **Phase 1 Status**: Core functionality complete. CLI tool monitors endpoints indefinitely and logs status.

---

## Phase 2: Core Features

**Goal**: 🚀 Production-viable with persistence, alerts, and UI.

### Multiple Endpoint Types
- [x] HTTP/HTTPS (GET, POST, PUT, PATCH, DELETE)
- [x] Custom headers support
- [x] Expected status code matching (`must_match_status`)
- [x] Max latency threshold checking
- [ ] Response body matching (regex/string) - fields exist but not implemented
- [ ] TCP port checks - stub exists (`TCPChecker`)
- [ ] DNS lookup checks - stub exists (`DNSChecker`)

### Advanced HTTP Features
- [ ] HTTP client pooling & keep-alive (configured in `HTTPChecker`)
- [x] Per-endpoint timeout configuration
- [x] Per-endpoint interval configuration
- [ ] SSL certificate expiry checking (< 30 days)

### Persistence
- [ ] SQLite backend (`/data/monitor.db`)
- [ ] Schema: `checks`, `endpoints`, `alerts`
- [ ] Automatic cleanup (retain last 30 days)

### Metrics
- [ ] Uptime % calculation
- [ ] Avg/min/max response time
- [ ] Downtime duration tracking

### Alerting
- [ ] Thresholds (e.g., `3 consecutive failures`)
- [ ] Recovery detection
- [ ] Console alerts (structured logging exists)
- [ ] Email (SMTP) alerts

### Web Dashboard
- [ ] Embedded HTTP server (`:8080`)
- [ ] Real-time status page (HTML + minimal JS)
- [ ] Endpoint list with status badges

✅ **Phase 2 Status**: HTTP checks with advanced features (headers, latency, status matching) are working. Persistence, metrics, alerts, and dashboard pending.

---

## Phase 3: Advanced Features

**Goal**: 🌐 Enterprise-grade extensibility and integrations.

### Notification Integrations
- [ ] Slack/Discord webhooks
- [ ] PagerDuty (v2 Events API)
- [ ] Custom webhook support

### Configuration Enhancements
- [ ] Hot-reload on config change (fsnotify)
- [ ] Env var overrides (`HC_INTERVAL=300s`)
- [x] Validate config on load

### Advanced Checks
- [ ] SSL certificate expiry (< 30 days)
- [ ] DNS lookup checks (stub exists)
- [ ] Status code ranges (`2xx`, `3xx`) - partially supported via status code ranges

### API Layer
- [ ] REST API (`/api/v1/status`, `/api/v1/history`)
- [ ] CRUD for endpoints (add/remove/update)
- [ ] JWT-based authentication

### Performance & Resilience
- [ ] HTTP client pooling & keep-alive
- [ ] Circuit breaker for failing endpoints
- [ ] CPU/memory profiling

✅ **Phase 3 Status**: Not started. Foundation is in place for extensibility.

---

## Phase 4: Production Readiness

**Goal**: 🛡️ Secure, observable, and deployable.

### Observability
- [ ] Prometheus metrics (`/metrics`)
  - `healthcheck_up`, `healthcheck_latency_seconds`
- [x] Structured JSON logging (custom logger with JSON output)
- [ ] Tracing (OpenTelemetry)

### Security
- [ ] TLS for web/API (`--tls-cert`, `--tls-key`)
- [ ] RBAC for API/dashboard
- [ ] Rate limiting (5 req/s per IP)

### Deployment
- [ ] Dockerfile (multi-stage, scratch base)
- [ ] Helm chart (for Kubernetes)
- [ ] systemd service template

### Quality & Docs
- [ ] Test coverage ≥ 80% (unit + integration)
- [ ] CLI help (`--help`, `--version`)
- [ ] `README.md` with examples
- [ ] `docs/` (config reference, API spec)

✅ **Phase 4 Status**: Not started. Logging infrastructure exists.

---

## Milestones & Success Criteria

| Phase | Status | Success Criteria |
|-------|--------|------------------|
| **1** | ✅ Complete | CLI tool runs indefinitely, checks endpoints, logs to stdout |
| **2** | 🚧 In Progress | SQLite storage + email alerts + live dashboard on `:8080` |
| **3** | 📋 Planned | API + webhook alerts + config hot-reload |
| **4** | 📋 Planned | Docker image + Prometheus metrics + ≥80% test coverage |

---

## Tech Stack

| Component       | Choice                          | Status        |
|-----------------|---------------------------------|---------------|
| Language        | Go                              | ✅            |
| Config          | YAML (Viper)                    | ✅            |
| Storage         | No yet                          | 📋 Planned    |
| Web Server      | `net/http`                      | 📋 Planned    |
| CLI             | `flag` (standard library)       | ✅            |
| Logging         | Custom JSON logger              | ✅            |
| Testing         | `testify`, `httpexpect`         | 📋 Planned    |
| Build           | Go modules                      | ✅            |

---

## Current Implementation Status

### ✅ Completed Features
- HTTP health checks with multiple methods (GET, POST, PUT, PATCH, DELETE)
- YAML configuration with validation
- Concurrent endpoint monitoring with goroutines
- Per-endpoint intervals and timeouts
- Custom headers support
- Status code classification and strict matching
- Max latency threshold detection
- Structured JSON logging
- Buffered result channels

### 🚧 In Progress / Partially Complete
- Response body matching (fields defined, implementation pending)
- TCP/DNS checkers (stubs exist, implementation pending)

### 📋 Planned Features
- Persistence layer (SQLite)
- Metrics and uptime tracking
- Alerting system
- Web dashboard
- API layer
- Additional notification channels
- Deployment tooling
- Comprehensive testing

---

> 📝 **Note**: The project has a solid foundation with core HTTP monitoring capabilities. Focus areas for next steps: persistence, alerting, and dashboard.
