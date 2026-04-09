# CookieFarm — Dependency Graph

The project is a monorepo composed of four major areas: a **Go server**, a **Go client**, a **Python exploiter library**, and a **monitoring stack**. Below is a multi-level dependency breakdown.

---

## 1. Top-Level Component Architecture

```mermaid
graph TD
    subgraph "cookiefarm/client (Go)"
        CLI["Go Client Binary"]
    end

    subgraph "exploiter/python (Python)"
        PY["Python Exploiter\\n(cookiefarm pkg)"]
    end

    subgraph "cookiefarm/server (Go)"
        SRV["Go Server Binary"]
        DB[("SQLite DB\\ncookiefarm.db")]
        PROTOS["Protocol Plugins\\n(.so files)"]
        SRV --> DB
        SRV --> PROTOS
    end

    subgraph "monitoring/"
        GRAFANA["Grafana"]
        PROM["Prometheus"]
        NODE["node-exporter"]
        PROC["process-exporter"]
        GRAFANA --> PROM
        PROM --> NODE
        PROM --> PROC
    end

    CLI -- "WebSocket\\n(flags submission)" --> SRV
    CLI -- "HTTP REST API\\n(config, auth)" --> SRV
    PY -- "HTTP REST API\\n(submit-flags-standalone)" --> SRV
    SRV -- "Submits flags" --> EXT_CHK["External Flag Checker\\n(CTF Infrastructure)"]
```

## 2. Server Internal Package Dependencies

```mermaid
graph TD
    SRV_MAIN["server/main.go"]
    SRV_CMD["server/cmd"]
    SRV_API["server/api"]
    SRV_CORE["server/core"]
    SRV_SQLITE["server/sqlite"]
    SRV_WS["server/websockets"]
    SRV_CONFIG["server/config"]
    SRV_CTRL["server/controllers"]

    PKG_LOGGER["pkg/logger"]
    PKG_MODELS["pkg/models"]
    PKG_PROTOCOLS["pkg/protocols"]
    PKG_SYSTEM["pkg/system"]

    SRV_MAIN --> SRV_CMD
    SRV_MAIN --> PKG_LOGGER

    SRV_CMD --> SRV_API
    SRV_CMD --> SRV_CORE
    SRV_CMD --> SRV_SQLITE
    SRV_CMD --> SRV_CONFIG
    SRV_CMD --> PKG_LOGGER
    SRV_CMD --> PKG_MODELS

    SRV_API --> SRV_CONFIG
    SRV_API --> SRV_CORE
    SRV_API --> SRV_SQLITE
    SRV_API --> SRV_WS
    SRV_API --> SRV_CTRL
    SRV_API --> PKG_LOGGER
    SRV_API --> PKG_MODELS

    SRV_CORE --> SRV_CONFIG
    SRV_CORE --> SRV_SQLITE
    SRV_CORE --> PKG_LOGGER
    SRV_CORE --> PKG_MODELS
    SRV_CORE --> PKG_PROTOCOLS

    SRV_SQLITE --> SRV_CONFIG
    SRV_SQLITE --> PKG_LOGGER
    SRV_SQLITE --> PKG_MODELS
    SRV_SQLITE --> PKG_SYSTEM

    SRV_WS --> SRV_CONFIG
    SRV_WS --> SRV_SQLITE
    SRV_WS --> PKG_LOGGER
    SRV_WS --> PKG_MODELS

    SRV_CTRL --> SRV_SQLITE

    SRV_CONFIG --> PKG_MODELS
    SRV_CONFIG --> PKG_PROTOCOLS

    PKG_PROTOCOLS --> PKG_LOGGER
```

- `server/main.go` imports `server/cmd` and `pkg/logger`
- `server/cmd` (the Cobra root command) wires `server/api`, `server/core`, `server/sqlite`, and `server/config` together at startup.
- `server/api` handles all HTTP routing and calls into `server/sqlite`, `server/core`, `server/websockets`, `server/controllers`, and the shared `pkg/models`.
- `server/core` (flag processing loop) depends on `server/sqlite` to read unsubmitted flags and on `pkg/protocols` to dynamically load the protocol `.so` plugin.
- `server/sqlite` depends on `server/config` (for the DB path), `pkg/system`, `pkg/logger`, `pkg/models`, and the `crawshaw.io/sqlite` driver.
- `server/websockets` depends on `server/config` (for JWT secret), `server/sqlite` (via `FlagCollector`), `pkg/models`, and `pkg/logger`.
- `server/controllers` calls `server/sqlite` directly (via `FlagCollector`) to expose stats.
- `server/config` declares the shared config struct and the `Submit` function type, importing `pkg/models` and `pkg/protocols`.

---

## 3. Client Internal Package Dependencies

```mermaid
graph TD
    CLI_MAIN["client/main.go"]
    CLI_CMD["client/cmd"]
    CLI_API["client/internal/api"]
    CLI_WS["client/internal/websockets"]
    CLI_EXPLOIT["client/internal/exploit"]
    CLI_SUBMITTER["client/internal/submitter"]
    CLI_TEMPLATE["client/internal/template"]
    CLI_TUI["client/internal/tui"]
    CLI_CONFIG["client/pkg/config"]
    CLI_PROCESS["client/pkg/process"]

    PKG_LOGGER["pkg/logger"]
    PKG_MODELS["pkg/models"]
    PKG_SYSTEM["pkg/system"]
    SRV_DB["server/database"]
    SHAREDCFG["sharedconfig"]

    CLI_MAIN --> CLI_CMD
    CLI_MAIN --> CLI_CONFIG
    CLI_MAIN --> CLI_TUI
    CLI_MAIN --> PKG_LOGGER

    CLI_CMD --> CLI_EXPLOIT
    CLI_CMD --> CLI_CONFIG
    CLI_CMD --> CLI_API
    CLI_CMD --> CLI_SUBMITTER
    CLI_CMD --> CLI_TEMPLATE
    CLI_CMD --> CLI_WS

    CLI_TUI --> CLI_CONFIG
    CLI_TUI --> CLI_CMD
    CLI_TUI --> CLI_EXPLOIT
    CLI_TUI --> CLI_TEMPLATE
    CLI_TUI --> PKG_LOGGER

    CLI_API --> CLI_CONFIG
    CLI_API --> PKG_LOGGER
    CLI_API --> PKG_MODELS
    CLI_API --> SRV_DB
    CLI_API --> SHAREDCFG

    CLI_WS --> CLI_CONFIG
    CLI_WS --> PKG_LOGGER
    CLI_WS --> SRV_DB
    CLI_WS --> SHAREDCFG

    CLI_EXPLOIT --> CLI_CONFIG
    CLI_EXPLOIT --> CLI_PROCESS
    CLI_EXPLOIT --> PKG_LOGGER
    CLI_EXPLOIT --> SRV_DB

    CLI_SUBMITTER --> CLI_API
    CLI_SUBMITTER --> PKG_LOGGER
    CLI_SUBMITTER --> SRV_DB

    CLI_TEMPLATE --> CLI_CONFIG
    CLI_TEMPLATE --> PKG_LOGGER
    CLI_TEMPLATE --> PKG_SYSTEM

    CLI_CONFIG --> SHAREDCFG
    CLI_CONFIG --> PKG_LOGGER
```

The client is structured as a **monorepo-style multi-module Go project** with a clear separation between public packages (`pkg/`) and internal packages (`internal/`):

- `client/main.go` initialises the `ConfigManager` singleton, parses CLI arguments (via `Cobra` + `fang`), and delegates to either `client/internal/tui` or the CLI command tree.
- `client/pkg/config` is the foundational singleton shared by almost every package. It uses `sync/atomic` for lock-free reads and persists state to two YAML files (`client.yml`, `shared.yml`) and a plain `session` token file under `~/.config/cookiefarm/`.
- `client/pkg/process` is a leaf package with no internal dependencies. It provides cross-platform subprocess management (`StartWithContext`, `StartDetached`) with Unix process-group kill semantics (`SIGKILL` to `pgid`).
- `client/internal/exploit` is the execution core: it owns the `Exploits` singleton (thread-safe PID registry), calls `pkg/process` to launch Python scripts, pipes their output through a `Parser` (JSON → `database.Flag`), and exposes the flag stream as a Go channel (`<-chan database.Flag`).
- `client/internal/websockets` consumes the flag channel and streams flags to the server over a persistent WebSocket connection. It implements a **three-state circuit breaker** (Closed / HalfOpen / Open) and a `ConnectionMonitor` that tracks latency, message counts, and performs health-checks every 30 s. It also handles incoming `{type:"config"}` server pushes and updates `pkg/config` in-place.
- `client/internal/submitter` is the HTTP fallback path. It drains the flag channel in batches of 50 and calls `client/internal/api.SubmitBatchDirect`. It is activated when `exploit run/test --submit` is set or when the `exploit submit` command is used.
- `client/internal/api` is a thin HTTP singleton client (10 s timeout) covering `Login`, `GetConfig`, `SubmitBatchDirect`, and `SubmitFlag`. It uses `bytedance/sonic` for JSON and reads the JWT token from `pkg/config`.
- `client/internal/template` creates/removes Python exploit files in `~/.config/cookiefarm/exploits/` using an embedded `@exploit_manager`-decorated template.
- `client/internal/tui` implements the Bubble Tea TUI. It depends on `client/cmd` (to reuse the `LoginHandler`), `client/internal/exploit` and `client/internal/template` (for direct action execution), and `client/pkg/config` (for live config reads). The `CommandRunner` bridges TUI events to package calls; `CommandHandler` dispatches form submissions.

---

## 4. Shared Package (`pkg`) Dependencies

```mermaid
graph TD
    PKG_LOGGER["pkg/logger"] --> ZEROLOG["github.com/rs/zerolog"]
    PKG_LOGGER --> LIPGLOSS["github.com/charmbracelet/lipgloss"]
    PKG_LOGGER --> FANG["github.com/charmbracelet/fang"]

    PKG_PROTOCOLS["pkg/protocols"] --> PKG_LOGGER
    PKG_PROTOCOLS --> PLUGIN["Go built-in: plugin"]

    PKG_MODELS["pkg/models"]
    PKG_SYSTEM["pkg/system"]
```

- `pkg/logger` is a leaf shared package wrapping `zerolog`, `lipgloss`, and `fang` – consumed by both server and client.
- `pkg/protocols` dynamically loads `.so` protocol plugins using Go's `plugin` package and logs via `pkg/logger`.
- `pkg/models` defines all shared data structures (`ClientData`, `ConfigShared`, etc.) with no external Go dependencies.

---
