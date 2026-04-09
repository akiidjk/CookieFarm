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
    CLI_API["client/api"]
    CLI_WS["client/websockets"]
    CLI_EXPLOIT["client/exploit"]
    CLI_CONFIG["client/config"]
    CLI_TUI["client/tui"]

    PKG_LOGGER["pkg/logger"]
    PKG_MODELS["pkg/models"]
    PKG_SYSTEM["pkg/system"]

    CLI_MAIN --> CLI_CMD
    CLI_MAIN --> CLI_CONFIG
    CLI_MAIN --> CLI_TUI
    CLI_MAIN --> PKG_LOGGER

    CLI_CMD --> CLI_EXPLOIT
    CLI_CMD --> CLI_CONFIG
    CLI_CMD --> CLI_API

    CLI_TUI --> CLI_CONFIG
    CLI_TUI --> PKG_LOGGER

    CLI_API --> CLI_CONFIG
    CLI_API --> PKG_LOGGER
    CLI_API --> PKG_MODELS

    CLI_WS --> CLI_CONFIG
    CLI_WS --> PKG_LOGGER
    CLI_WS --> PKG_MODELS

    CLI_EXPLOIT --> CLI_CONFIG
    CLI_EXPLOIT --> CLI_WS
    CLI_EXPLOIT --> PKG_LOGGER
    CLI_EXPLOIT --> PKG_SYSTEM
```

- `client/main.go` checks for TUI mode and delegates to either `client/tui` or `client/cmd`.
- `client/exploit` is the core: it calls `client/websockets` to stream captured flags to the server, and uses `client/config` and `pkg/system` for process management.
- `client/api` provides HTTP calls (`GetConfig`, `Login`, `SubmitBatchDirect`) and imports `pkg/models` for shared data types. [17](https://www.notion.so/Dependency-Graph-31d5d8cb6b3a8055b91fc1683691ef96?pvs=21)
- `client/websockets` connects to the server WebSocket endpoint, handles the circuit breaker, and dispatches `ConfigEvent` messages by updating `client/config`.
- `client/tui` depends only on `client/config` and `pkg/logger` plus the Charmbracelet UI libraries.

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
