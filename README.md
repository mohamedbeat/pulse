# Pulse
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![WIP](https://img.shields.io/badge/status-under_development-yellow)

A robust, extensible health monitoring tool written in Go.  
> **Status**: In Development  
> **Target**: Production-ready monitoring solution for services/endpoints.

> 🚧Note:🚧 
 The project is **Under Active Development**. We're refining features and adding new capabilities. Your feedback helps shape the final product!

![WIP](https://img.shields.io/badge/status-under_development-yellow)
![Not Production Ready](https://img.shields.io/badge/production-not_ready-red)
![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen) 

## 📋 What is Pulse?
Pulse is a high-concurrency health monitoring daemon in Go that actively probes HTTP endpoints with precise timing, latency thresholds, and status validation—delivering real-time, structured observability. Designed for reliability and extension, it’s the foundation for a full-stack monitoring system: from CLI checks today, to alerts, dashboards, and API-driven ops tomorrow.


## ✨ Key Features (Current)
👉 Check out our detailed [Roadmap](ROADMAP.md)
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
- Retry mechanism 


### 📋 Planned Features
- Response body matching (fields defined, implementation pending)
- TCP/DNS checkers (stubs exist, implementation pending)
- Persistence layer
- Metrics and uptime tracking
- Alerting system
- Web dashboard
- API layer
- Additional notification channels
- Deployment tooling
- Comprehensive testing

 👉 Check out our detailed [Roadmap](ROADMAP.md)

## 🤝 Contributing
We welcome contributors, especially during this development phase! This is the perfect time to help shape pulse.


## Development Setup
```bash
# 1. Fork and clone
git clone https://github.com/mohamedbeat/pulse.git
cd pulse

# 2. Install dev dependencies
go mod tidy
```

⚠️ This is beta software - please be aware
 Report issues on our GitHub Issues page.


## 📄 License
Pulse is released under the [MIT License](LICENSE).

## 📬 Stay Updated
Star the project to show your support.

Watch this repo for releases.

C
