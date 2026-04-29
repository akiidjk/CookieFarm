# CookieFarm - Dependency Graph

CookieFarm is a monorepo composed of a Go server, a Go client, shared Go packages, a Python exploiter library, a React/Vite server frontend, and a monitoring stack.

## 1. Top-Level Component Architecture

```mermaid
graph TD
    subgraph "cookiefarm/client"
        CLI["ckc Go Client"]
        CLI_CKP["client/internal/ckp"]
        CLI_API["client/internal/api"]
        CLI_EXP["client/internal/exploit"]
    end

    subgraph "exploiter/python"
        PY["cookiefarm Python Exploiter"]
    end

    subgraph "cookiefarm/server"
        SRV["cks Go Server"]
        API["server/api"]
        SRV_CKP["server/internal/ckp"]
        DB[("SQLite\ncookiefarm.db")]
        FRONT["React/Vite Frontend"]
        PROTOS["Protocol Adapters\npkg/protocols"]
    end

    subgraph "monitoring"
        GRAFANA["Grafana"]
        PROM["Prometheus"]
        NODE["node-exporter"]
        PROC["process-exporter"]
    end

    CLI --> CLI_EXP
    CLI_EXP --> CLI_CKP
    CLI --> CLI_API
    PY --> CLI_EXP
    CLI_CKP -- "CKP TCP :7777" --> SRV_CKP
    CLI_API -- "HTTP REST /api/v1" --> API
    SRV --> API
    SRV --> SRV_CKP
    API --> DB
    SRV_CKP --> DB
    API --> FRONT
    SRV --> PROTOS
    PROTOS -- "submit flags" --> EXT["External Flag Checker"]
    GRAFANA --> PROM
    PROM --> NODE
    PROM --> PROC
```

## 2. Server Package Dependencies

```mermaid
graph TD
    SRV_MAIN["server/main.go"]
    SRV_CMD["server/cmd"]
    SRV_API["server/api"]
    SRV_CKP["server/internal/ckp"]
    SRV_CORE["server/internal/core"]
    SRV_DB["server/internal/database"]
    SRV_CONFIG["server/pkg/config"]
    SRV_POOL["server/pkg/pool"]

    PKG_SHARED["pkg/config\nsharedconfig"]
    PKG_LOGGER["pkg/logger"]
    PKG_MODELS["pkg/models"]
    PKG_PROTOCOLS["pkg/protocols"]

    SRV_MAIN --> SRV_CMD
    SRV_CMD --> SRV_API
    SRV_CMD --> SRV_CKP
    SRV_CMD --> SRV_CORE
    SRV_CMD --> SRV_DB
    SRV_CMD --> SRV_CONFIG
    SRV_CMD --> PKG_LOGGER

    SRV_API --> SRV_CKP
    SRV_API --> SRV_CORE
    SRV_API --> SRV_DB
    SRV_API --> SRV_CONFIG
    SRV_API --> PKG_LOGGER
    SRV_API --> PKG_MODELS

    SRV_CKP --> SRV_DB
    SRV_CKP --> SRV_CONFIG
    SRV_CKP --> SRV_POOL
    SRV_CKP --> PKG_LOGGER

    SRV_CORE --> SRV_DB
    SRV_CORE --> SRV_CONFIG
    SRV_CORE --> PKG_LOGGER
    SRV_CORE --> PKG_MODELS
    SRV_CORE --> PKG_PROTOCOLS

    SRV_DB --> PKG_LOGGER
    SRV_DB --> PKG_PROTOCOLS

    SRV_CONFIG --> PKG_SHARED
    SRV_CONFIG --> PKG_PROTOCOLS
```

- `server/main.go` imports `server/cmd`.
- `server/cmd` wires configuration, database, core runner, CKP server, and Fiber API startup.
- `server/api` owns HTTP routing, auth middleware, Swagger routes, frontend fallback, and config broadcasts to connected CKP clients.
- `server/internal/ckp` owns the raw TCP listener, accepted connection registry, frame parsing, and config writes.
- `server/internal/core` owns the flag checker submission loop and TTL cleanup loop.
- `server/internal/database` owns SQLite schema access, sqlc generated queries, mapping helpers, and `FlagCollector`.
- `server/pkg/config` owns server runtime configuration, environment variables, JWT secret, and active checker submit function.
- `server/pkg/pool` is used by the CKP server worker pool.

## 3. Client Package Dependencies

```mermaid
graph TD
    CLI_MAIN["client/main.go"]
    CLI_CMD["client/cmd"]
    CLI_API["client/internal/api"]
    CLI_CKP["client/internal/ckp"]
    CLI_EXPLOIT["client/internal/exploit"]
    CLI_SUBMITTER["client/internal/submitter"]
    CLI_TEMPLATE["client/internal/template"]
    CLI_TUI["client/internal/tui"]
    CLI_CONFIG["client/pkg/config"]
    CLI_PROCESS["client/pkg/process"]

    SRV_DB["server/internal/database"]
    SHAREDCFG["pkg/config\nsharedconfig"]
    PKG_LOGGER["pkg/logger"]
    PKG_MODELS["pkg/models"]
    PKG_SYSTEM["pkg/system"]

    CLI_MAIN --> CLI_CMD
    CLI_MAIN --> CLI_CONFIG
    CLI_MAIN --> CLI_TUI
    CLI_MAIN --> PKG_LOGGER

    CLI_CMD --> CLI_API
    CLI_CMD --> CLI_CKP
    CLI_CMD --> CLI_EXPLOIT
    CLI_CMD --> CLI_SUBMITTER
    CLI_CMD --> CLI_TEMPLATE
    CLI_CMD --> CLI_CONFIG

    CLI_CKP --> CLI_CONFIG
    CLI_CKP --> SRV_DB
    CLI_CKP --> SHAREDCFG
    CLI_CKP --> PKG_LOGGER

    CLI_API --> CLI_CONFIG
    CLI_API --> SRV_DB
    CLI_API --> SHAREDCFG
    CLI_API --> PKG_LOGGER
    CLI_API --> PKG_MODELS

    CLI_EXPLOIT --> CLI_CONFIG
    CLI_EXPLOIT --> CLI_PROCESS
    CLI_EXPLOIT --> SRV_DB
    CLI_EXPLOIT --> PKG_LOGGER

    CLI_SUBMITTER --> CLI_API
    CLI_SUBMITTER --> SRV_DB
    CLI_SUBMITTER --> PKG_LOGGER

    CLI_TEMPLATE --> CLI_CONFIG
    CLI_TEMPLATE --> PKG_LOGGER
    CLI_TEMPLATE --> PKG_SYSTEM

    CLI_TUI --> CLI_CMD
    CLI_TUI --> CLI_CONFIG
    CLI_TUI --> CLI_EXPLOIT
    CLI_TUI --> CLI_TEMPLATE
    CLI_TUI --> PKG_LOGGER

    CLI_CONFIG --> SHAREDCFG
    CLI_CONFIG --> PKG_LOGGER
    CLI_CONFIG --> PKG_SYSTEM
```

- `client/main.go` initializes config and delegates to the TUI or Cobra command tree.
- `client/cmd` owns CLI commands and chooses CKP transport by default for exploit runs.
- `client/internal/ckp` depends on config, shared config, logger, and `server/database.Flag` for the binary flag payload model.
- `client/internal/api` wraps HTTP calls for auth, config, direct submission, and exploit upload.
- `client/internal/exploit` starts Python subprocesses, parses stdout, and emits flag channels.
- `client/internal/submitter` provides the HTTP fallback path used by `--submit`.
- `client/internal/template` manages generated exploit templates.
- `client/internal/tui` bridges Bubble Tea UI actions to command/package operations.
- `client/pkg/config` is the atomic runtime config singleton.
- `client/pkg/process` is a leaf package for cross-platform process management.

## 4. Shared Package Dependencies

```mermaid
graph TD
    PKG_LOGGER["pkg/logger"] --> ZEROLOG["github.com/rs/zerolog"]
    PKG_LOGGER --> LIPGLOSS["charm.land/lipgloss"]
    PKG_LOGGER --> FANG["github.com/charmbracelet/fang"]

    PKG_PROTOCOLS["pkg/protocols"] --> PKG_LOGGER
    PKG_PROTOCOLS --> PLUGIN["Go plugin package"]

    PKG_MODELS["pkg/models"] --> SRV_DB["server/internal/database"]
    PKG_CONFIG["pkg/config\nsharedconfig"]
    PKG_SYSTEM["pkg/system"]
```

- `pkg/config` defines the shared CTF metadata model imported as `sharedconfig`.
- `pkg/logger` wraps zerolog and Charm tooling for logs and CLI presentation.
- `pkg/models` defines status constants and HTTP request envelopes.
- `pkg/protocols` defines checker response types, built-in protocols, and dynamic protocol loading.
- `pkg/system` provides filesystem helpers such as tilde expansion.

## 5. CKP Dependency Position

CKP sits between the exploit parser and the server collector:

```mermaid
graph LR
    PY["Python exploit stdout"] --> PARSER["client/internal/exploit Parser"]
    PARSER --> FLAGS["chan database.Flag"]
    FLAGS --> CKPC["client/internal/ckp"]
    CKPC -->|"TCP :7777"| CKPS["server/internal/ckp"]
    CKPS --> COLL["server/internal/database FlagCollector"]
    COLL --> DB["SQLite"]
```

The HTTP submitter remains available but is not the default exploit-run path:

```mermaid
graph LR
    FLAGS["chan database.Flag"] --> SUB["client/internal/submitter"]
    SUB --> API["server/api"]
    API --> DB["SQLite"]
    API --> PROTO["pkg/protocols"]
```
