# Stack Research

**Domain:** Zero-trust security gateway for AI agents  
**Researched:** Mon Jun 15 2026  
**Confidence:** HIGH  

## Recommended Technology Stack

### Core Language & Runtime
- **Go 1.22+** - Selected for performance, concurrency model, and extensive ecosystem for cloud-native services
- **Rationale:** Excellent for building high-performance network services, strong standard library for HTTP handling, good support for WASM if needed, and mature dependency management

### Containerization & Orchestration
- **Docker** - Container platform for consistent deployment
- **Rationale:** Matches deployment constraints, ensures environment consistency, and simplifies dependency management

### Communication & Transport
- **HTTP/1.1 & HTTP/2** - For MCP over HTTP transport
- **STDIO pipes** - For local MCP server communication (stdio transport)
- **Rationale:** Supports both standard MCP transports (HTTP and STDIO) as defined in MCP specification

### Authentication
- **JWT (JSON Web Tokens)** - For service-to-service authentication
- **Library:** github.com/golang-jwt/jwt/v5
- **Rationale:** Industry standard, stateless for scalability, integrates with existing identity systems

### Authorization & Policy Engine
- **Open Policy Agent (OPA)** - For policy evaluation and decision-making
- **Language:** Rego for policy definitions
- **Rationale:** Industry-standard for cloud-native authorization, strong ecosystem, separates policy from code

### Semantic Inspection (LLM-based)
- **Local LLM Inference** - For prompt safety and intent analysis
- **Options:** 
  - Llama.cpp (for Llama 3 8B or Mistral 7B models)
  - Ollama (for easy model management)
  - Hugging Face Transformers with ONNX Runtime
- **Rationale:** Ensures privacy (no prompt leakage), predictable latency/cost, and offline operation

### Observability
- **Structured Logging** - JSON-formatted logs for parsing and analysis
- **Library:** Uber Zap or Zerolog
- **Metrics:** Prometheus client library for Go
- **Tracing:** OpenTelemetry Go SDK
- **Rationale:** Production-ready observability stack with wide adoption and integration capabilities

### Storage & Persistence
- **Primary Database:** PostgreSQL - For audit logs, configuration, and relational data
- **Cache/Config Store:** Redis - For frequently accessed data, sessions, and pub/sub
- **Analytics Database:** ClickHouse - For queryable audit trail and usage analytics (optional for MVP)
- **Alternative for MVP:** BoltDB or SQLite - Embedded option for simpler deployment
- **Rationale:** Provides flexibility from simple embedded storage to scalable production architecture

### CI/CD & DevOps
- **GitHub Actions** - For automated testing, building, and deployment
- **Rationale:** Tight integration with repository hosting, supports Docker builds and deployments

### Security Headers & Middleware
- **Standard Security Headers:** X-Content-Type-Options, X-Frame-Options, X-XSS-Protection, Strict-Transport-Security, Content-Security-Policy
- **Rationale:** Basic web security hygiene; expected in any internet-facing service

### RPC/Protocol Handling
- **JSON-RPC 2.0** - For MCP protocol implementation
- **Rationale:** MCP specification uses JSON-RPC 2.0 over HTTP/STDIO

### Configuration Management
- **Library:** Viper or envconfig
- **Features:** Support for JSON/YAML files, environment variables, and dynamic updates
- **Rationale:** Flexible configuration loading with support for different sources

### Testing Framework
- **Go's built-in testing** - With testify for assertions
- **Rationale:** Standard, reliable, and well-integrated with Go toolchain

## Version Requirements & Constraints

### Language Versions
- Go: 1.22+ (for generics and improved performance)
- Docker: 20.10+ (for BuildKit and security features)

### Library Compatibility
- All libraries should be actively maintained with recent releases (<6 months old)
- Prefer libraries with good documentation and examples
- Avoid GPL-licensed libraries for core components to maintain permissive licensing

### Deployment Targets
- Primary: Linux/amd64 Docker containers
- Secondary: Linux/arm64 for edge/IoT deployments
- Local development: Supports macOS and Windows via Docker

## Rationale Summary

The selected stack balances:
1. **Performance:** Go's efficiency and concurrency for handling MCP traffic
2. **Security:** JWT + OPA + local LLMs for zero-trust principles
3. **Observability:** Full tracing, metrics, and logging for production operations
4. **Scalability:** Stateless design with externalized state (Redis/PostgreSQL)
5. **Operational Simplicity:** Docker packaging and GitHub Actions for CI/CD
6. **Compliance:** Audit logging and data protection features
7. **Community Adoption:** Open-source technologies with permissive licenses

## Sources

- Project Constraints (PROJECT.md): Deployment, Auth, OPA, LLM, Persistence, CI/CD, Licenseóta, Team, Timeline
- Architecture Research (ARCHITECTURE.md): Component Responsibilities table and Scaling Considerations
- Industry Standards: Zero Trust Architecture (NIST SP 800-207), MCP Specification, OPA Documentation
- Technology Selection: Based on evaluation of performance, security, ecosystem maturity, and alignment with project goals
