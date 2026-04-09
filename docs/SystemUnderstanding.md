# System Understanding Document — CookieFarm

---

**Purpose:** This document outlines the architecture, components, interfaces, and data flows of the **CookieFarm** system, to ensure it meets the functional and non-functional requirements specified. It serves as a guide for developers, testers, and stakeholders during all stages of the software development lifecycle.

---

## 1. Overview

| Field | Value |
| --- | --- |
| **System Name** | CookieFarm |
| **Prepared By** | ByteTheCookies |
| **Date** | 2025 |
| **Version** | v2.0.0-rc |

CookieFarm is an Attack/Defense CTF framework inspired by DestructiveFarm, developed by **ByteTheCookies**. 

---

## 2. System Objectives

CookieFarm automates the entire flag lifecycle in Attack/Defense CTF competitions:

- **Automate exploit distribution** across all opponent teams, removing manual targeting overhead.
- **Collect flags** captured by exploit scripts and submit them to the official flag checker automatically.
- **Monitor flag statuses** (`UNSUBMITTED`, `ACCEPTED`, `DENIED`, `ERROR`) in real-time through a web interface.
- **Allow participants to focus exclusively on writing exploits** — the system handles the rest.

---

## 3. Scope

The system covers two major runtime boundaries:

**Server-side (Go):**

- Receives and stores flags submitted by clients into a SQLite database.
- Periodically submits batches of unsubmitted flags to the external flag checker service via pluggable protocol adapters.
- Exposes a web dashboard (SSR + HTMX) for monitoring.
- Exposes a RESTful API and a WebSocket endpoint for client communication.

**Client-side (Go + Python):**

- Provides a CLI (`ckc`) and TUI interface for managing exploits.
- Launches Python exploit scripts as subprocesses, captures their stdout as structured JSON flag data, and forwards them to the server via WebSocket.
- Supports exploit lifecycle management: create, run, test, list, stop, remove.

**Out of scope:**

- The exploit scripts themselves (user-provided Python).
- The external CTF flag checker service.
- The CTF game infrastructure (team IPs, scoring).

---

## 4. Assumptions and Constraints

- **Go 1.26.0+** is required for the server and client binaries.
- **Python 3+** is required to execute exploit scripts.
- **Docker** is required to deploy the server via `docker compose up --build`.
- The server is designed for **Linux/amd64** production deployments (Alpine Docker image).
- Protocol plugins (`.so` files) are compiled separately via `go build -buildmode=plugin`. They must match the Go runtime version of the server binary.
- The default password (`"password"`) **must** be changed before production deployment. Leaving it as default is a critical security risk.
- The server configuration can be provided either via a **YAML file** (`config.yml`) or via the **web form** at runtime.
- The web UI is currently based on HTMX/JavaScript and is **not fully up-to-date**.
- Flag TTL is measured in **ticks** (game rounds), not wall-clock time. Expired flags are automatically deleted.

---

## 5. System Architecture

**Overview:**
CookieFarm is a **distributed, event-driven, hybrid Go + Python** framework. The architecture separates responsibilities cleanly between a **central Go server** (flag storage, checker integration, web dashboard) and a **lightweight Go client** (exploit orchestration, flag capture, flag forwarding). The primary communication channel is **WebSocket**, with a fallback to **direct HTTP REST**. 

### 5.1 Components

| Component | Description |
| --- | --- |
| **`cmd/server` — Server Entrypoint** | Main entrypoint for the Go server binary (`cks`). Initializes the database, loads configuration, starts the Fiber HTTP server and background loops. |
| **`cmd/client` — Client Entrypoint** | Main entrypoint for the Go client binary (`ckc`). Bootstraps the CLI (Cobra) and TUI (Bubble Tea) interfaces. |
| **`internal/server/server` — HTTP Server** | Built on Fiber v2. Registers all routes (view, public API, private API protected by JWT, WebSocket). Handles CORS, rate limiting, static file serving, and compression.  |
| **`internal/server/core` — Flag Processing Engine** | Contains `StartFlagProcessingLoop` (periodic batch submission to flag checker) and `ValidateFlagTTL` (expired flag cleanup). Both run as independent goroutines with cancellable contexts. |
| **`internal/server/database` — Database Layer** | SQLite-backed persistence using a connection pool of size 20 with WAL mode. Provides CRUD operations for flags. Includes a `FlagCollector` singleton with an in-memory buffer (100 flags) and periodic flush (10s). |
| **`internal/server/websockets` — WebSocket Manager** | Manages connected clients, routes events (`flag`, `config`) to handlers. The `FlagHandler` routes incoming flags to the `FlagCollector`.  |
| **`internal/server/controllers` — Stats Controller** | Provides flag collector statistics: total received, flushed, buffer size, efficiency ratio.  |
| **`pkg/protocols` — Protocol Plugin System** | Dynamically loads `.so` plugins at runtime via Go's `plugin` package. Each plugin exposes a `Submit(host, token, flags)` function. The built-in `cc_http` protocol submits flags via HTTP PUT with `X-Team-Token`.  |
| **`internal/client/exploit` — Exploit Executor** | Launches Python exploit scripts as subprocesses. Reads stdout line-by-line, parsing structured JSON into `ClientData` flag objects. Manages exploit lifecycle (run, stop, restart on config change).  |
| **`internal/client/websockets` — Client WebSocket Submitter** | Maintains a WebSocket connection to the server and forwards parsed flags in real-time. Supports reconnection on failure.  |
| **`internal/client/api` — REST API Client** | HTTP client for server interactions: login (JWT cookie retrieval), config fetch, direct flag/batch submission.  |
| **`internal/server/ui` — Template Engine** | Server-side rendering using Go HTML templates (Fiber template engine). Serves dashboard and login views. |
| **`frontend/` — Next.js Frontend (WIP)** | A new React/Next.js-based frontend under development, currently disabled in `docker-compose.yml`.  |
| **`pkg/models` — Shared Data Models** | Defines all shared structs: `ClientData`, `ConfigShared`, `ConfigServer`, `ConfigClient`, `Service`, `SubmitFlagsRequest`.  |

### 5.2 System Diagram

The following diagram illustrates the overall system topology and data flow:

```mermaid
graph TD
    subgraph "CTF Competition Network"
        FC["Flag Checker Service\\n(External)"]
        TEAMS["Opponent Team VMs\\n(10.10.X.1)"]
    end

    subgraph "CookieFarm Client\\n(ckc binary - Go)"
        CLI["CLI / TUI Interface\\n(Cobra + Bubble Tea)"]
        EE["Exploit Executor\\n(internal/client/exploit)"]
        PY["Python Exploit Script\\n(.py - user provided)"]
        PARSER["stdout Parser\\n(ParseLine - JSON)"]
        WS_CLI["WebSocket Submitter\\n(internal/client/websockets)"]
        HTTP_CLI["HTTP API Client\\n(internal/client/api)"]
    end

    subgraph "CookieFarm Server\\n(cks binary - Go + Docker)"
        FIBER["Fiber HTTP Server\\n(internal/server/server)"]
        AUTH["JWT Auth\\n(bcrypt + HS256)"]
        WS_SRV["WebSocket Manager\\n(internal/server/websockets)"]
        FC_COLL["Flag Collector\\n(buffer + periodic flush)"]
        DB["SQLite Database\\n(WAL mode, pool=20)"]
        FLAG_LOOP["Flag Processing Loop\\n(internal/server/core)"]
        TTL["Flag TTL Validator\\n(internal/server/core)"]
        PROTO["Protocol Plugin\\n(cc_http.so / custom.so)"]
        UI["Web Dashboard\\n(HTMX / HTML templates)"]
    end

    CLI --> EE
    EE --> PY
    PY --> PARSER
    PARSER --> WS_CLI
    PARSER --> HTTP_CLI
    WS_CLI -->|"WebSocket /ws"| WS_SRV
    HTTP_CLI -->|"POST /api/v1/submit-flags-standalone"| FIBER
    WS_SRV --> FC_COLL
    FC_COLL --> DB
    FIBER --> AUTH
    AUTH --> DB
    FLAG_LOOP --> DB
    FLAG_LOOP --> PROTO
    PROTO -->|"HTTP PUT + X-Team-Token"| FC
    TTL --> DB
    FIBER --> UI
    EE -->|"exec subprocess"| TEAMS
```

---

## 6. Data Design

**Data Flow Description:**
Flags flow in one direction: from Python exploit stdout → Go client parser → WebSocket (or HTTP) → Server → SQLite. The server's flag processing loop independently reads unsubmitted flags from SQLite on a configurable timer, submits them to the external flag checker via a protocol plugin, and updates their status (`ACCEPTED`, `DENIED`, `ERROR`) back in SQLite. 

### 6.1 Data Entities

| Entity Name | Description |
| --- | --- |
| **`ClientData`** | Core flag entity. Contains `flag_code`, `service_name`, `port_service`, `team_id`, `status`, `username`, `exploit_name`, `submit_time`, `response_time`, `msg`. This is the primary record persisted in SQLite and exchanged between client and server.  |
| **`ConfigShared`** | Aggregated configuration struct shared between server and client contexts. Contains `ConfigServer`, `ConfigClient`, and a `Configured` boolean flag.  |
| **`ConfigServer`** | Server-side configuration: `url_flag_checker`, `team_token`, `protocol`, `tick_time`, `submit_flag_checker_time`, `max_flag_batch_size`, `flag_ttl`, `start_time`, `end_time`.  |
| **`ConfigClient`** | Client-side configuration: `services` list, `regex_flag`, `format_ip_teams`, `range_ip_teams`, `my_team_id`, `nop_team`, `url_flag_ids`.  |
| **`Service`** | Represents a single exploitable CTF service, with a `name` and `port`.  |
| **`ResponseProtocol`** | Response from the flag checker per flag: `status` (`ACCEPTED`/`DENIED`/`RESUBMIT`/`ERROR`), `flag`, `msg`.  |
| **`ParsedFlagOutput`** | JSON structure emitted by Python exploit scripts to stdout. Contains `status`, `flag_code`, `name_service`, `message`, `team_id`, `port_service`.  |
| **`StatusBatchOutput`** | JSON structure for batch statistics emitted by exploit scripts: `total_flag`, `success_team`, `failed_team`.  |
| **`FlagCollector`** | Server-side in-memory buffer (singleton). Holds up to 100 flags before flushing to SQLite. Tracks `CollectorStats` (total received, flushed, flush errors).  |

### 6.2 Data Flow Diagrams

**Flag Capture & Submission Flow:**

```mermaid
sequenceDiagram
    participant PY as "Python Exploit Script"
    participant PARSER as "stdout Parser (Go Client)"
    participant WS as "WebSocket Submitter"
    participant SRV as "CookieFarm Server"
    participant COLL as "Flag Collector (Buffer)"
    participant DB as "SQLite DB"
    participant FLAG_LOOP as "Flag Processing Loop"
    participant PROTO as "Protocol Plugin (.so)"
    participant FC as "External Flag Checker"

    PY->>PARSER: "JSON line to stdout"
    PARSER->>PARSER: "ParseLine() -> ClientData"
    PARSER->>WS: "ClientData via flagsChan"
    WS->>SRV: "WebSocket event {type: flag}"
    SRV->>COLL: "AddFlag(ClientData)"
    COLL->>COLL: "Buffer (max 100 flags)"
    COLL->>DB: "Flush (every 10s or on full buffer)"
    FLAG_LOOP->>DB: "GetUnsubmittedFlagCodeList(batchSize)"
    DB->>FLAG_LOOP: "[]string flag codes"
    FLAG_LOOP->>PROTO: "Submit(host, token, flags)"
    PROTO->>FC: "HTTP PUT + X-Team-Token"
    FC->>PROTO: "[]ResponseProtocol"
    PROTO->>FLAG_LOOP: "[]ResponseProtocol"
    FLAG_LOOP->>DB: "UpdateFlagsStatus(ACCEPTED/DENIED/ERROR)"
```

**Flag TTL Cleanup Flow:**

- `ValidateFlagTTL` runs on a ticker equal to `flagTTL × tickTime` seconds.
- On each tick, it calls `DeleteTTLFlag` to purge flags older than `flagTTL` ticks.

**Configuration Update Flow:**

- `POST /api/v1/config` → updates `SharedConfig` → cancels existing background goroutines → restarts `StartFlagProcessingLoop` and `ValidateFlagTTL` with new config → broadcasts config update to all WebSocket clients.

---

## 7. Interfaces

### 7.1 External Interfaces

| Name | Type | Description |
| --- | --- | --- |
| **Flag Checker Service** | HTTP (via Protocol Plugin) | The server submits captured flags via a dynamically loaded `.so` plugin. The default `cc_http` plugin sends an HTTP PUT request with flags as a JSON array and sets the `X-Team-Token` header. Accepts responses as `[]ResponseProtocol`.  |
| **Flag IDs Service** | HTTP (GET) | Optional external service for CyberChallenge-AD competitions that provides flag IDs. Configured via `url_flag_ids` in `ConfigClient`.  |
| **Opponent Team VMs** | TCP (via Python exploit) | Exploit scripts connect directly to opponent VMs using the IP pattern `format_ip_teams` (e.g., `10.10.{}.1`) across `range_ip_teams` team IDs, targeting the service port configured in `ConfigClient.Services`.  |
| **Web Browser (Dashboard)** | HTTP/HTTPS | End users access the server web interface at `http://<server_ip>:<port>`. Provides flag monitoring, status display, and server configuration form.  |

### 7.2 Internal Interfaces

| Name | Description |
| --- | --- |
| **`GET /api/v1/` (Public)** | Health check endpoint. Returns server status and current UTC time.  |
| **`POST /api/v1/auth/login` (Public, rate-limited)** | Authenticates with password, returns a JWT as an `HttpOnly` cookie (`token`) valid for 48 hours.  |
| **`GET /api/v1/protocols` (Public)** | Lists available protocol plugins (`.so` files) on the server filesystem.  |
| **`GET /api/v1/flags` (Private - JWT)** | Returns all flags stored in the database, ordered by submit time descending.  |
| **`GET /api/v1/flags/:limit` (Private - JWT)** | Returns paginated flags. Supports `?offset=N` query param.  |
| **`GET /api/v1/stats` (Private - JWT)** | Returns flag collector statistics: buffer size, total received/flushed, flush counts, efficiency ratio.  |
| **`GET /api/v1/config` (Private - JWT)** | Returns the current `ConfigShared` as JSON.  |
| **`POST /api/v1/submit-flags` (Private - JWT)** | Batch insert flags into the database (deferred checker submission via periodic loop).  |
| **`POST /api/v1/submit-flag` (Private - JWT)** | Insert a single flag AND immediately submit it to the flag checker.  |
| **`POST /api/v1/submit-flags-standalone` (Private - JWT)** | Batch insert flags AND immediately submit them all to the flag checker. Used by the client HTTP fallback.  |
| **`POST /api/v1/config` (Private - JWT)** | Updates `SharedConfig`, restarts background goroutines, broadcasts config update to all WS clients.  |
| **`DELETE /api/v1/delete-flag` (Private - JWT)** | Deletes a single flag by its `flag_code` query param.  |
| **`GET /ws` — WebSocket** | Real-time bidirectional channel. Client sends `{type: "flag", payload: ClientData}`. Server acknowledges with `{type: "flag_response"}`. Config push uses `{type: "config"}`.  |
| **`FlagCollector.AddFlag()` (Internal)** | Thread-safe in-process interface for adding flags to the buffer. Triggers immediate flush if buffer reaches 100 entries.  |
| **`protocols.LoadProtocol()` (Internal)** | Dynamically loads a protocol `.so` plugin by name and returns the `Submit` function pointer for use by the flag processing loop.  |

---

## 8. Security Considerations

- **Password Authentication:** The server requires a password set via the `PASSWORD` environment variable. Passwords are hashed at startup using **bcrypt** (`DefaultCost`). Login attempts compare against the stored bcrypt hash.
- **JWT Tokens:** After successful login, a **JWT (HS256)** is issued with a 24-hour expiry and embedded in an `HttpOnly`, `SameSite=Strict` cookie named `token`. The signing secret is a 32-byte cryptographically random value generated at each server startup.

Through the UI you can:

- View all received flags.
- Check the submission and acceptance status of each flag.
- Configure the server (unless you setup the configuration from YAML file).

---

## **9. Performance Requirements**

Still to be determined with numeric/statistic value, but as little CPU and memory as possible.

---

## 🗒️ Important Notes

- Never push directly to `dev` branch!!
- NEVER PUSH DIRECTLY TO `main` BRANCH!!
- Test your code before pushing (test environment in `/tests`)
- Make sure your branch is up to date with `dev` before creating a PR
- Delete your branch after it has been merged
