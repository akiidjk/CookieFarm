# 🍪 CookieFarm vs DestructiveFarm — Benchmark Methodology

> **Scope:** A rigorous, reproducible methodology for comparing two Attack/Defense CTF exploit farm frameworks across five dimensions: memory usage, CPU usage, flag throughput, UI pagination latency, and qualitative architecture metrics.

---

## Table of Contents

1. [Environment Setup](#1-environment-setup)
2. [Test Parameters](#2-test-parameters)
3. [Metric Definitions](#3-metric-definitions)
4. [Benchmark Procedures](#4-benchmark-procedures)
   - [M1 — RAM Usage](#m1--ram-usage)
   - [M2 — CPU Usage](#m2--cpu-usage)
   - [M3 — Flag Store Throughput](#m3--flag-store-throughput)
   - [M4 — UI Pagination Latency](#m4--ui-pagination-latency)
   - [M5 — Qualitative Observations](#m5--qualitative-observations)
5. [Dummy Exploit Specification](#5-dummy-exploit-specification)
6. [Data Collection Scripts](#6-data-collection-scripts)
7. [Reporting Template](#7-reporting-template)
8. [Reproducibility Checklist](#8-reproducibility-checklist)

---

## 1. Environment Setup

### Hardware (Fixed Baseline)

All measurements are taken on the **same local machine** with no other heavy processes running.

| Parameter | Requirement |
|-----------|------------|
| Machine | Single physical host (no VM overhead differences) |
| CPU cores pinned | Use `taskset` or Docker `--cpuset-cpus` to pin to the **same cores** for both tools |
| RAM limit | No limit imposed — measure natural consumption |
| Disk | SSD recommended (flag DB I/O) |
| Network | Loopback only (`127.0.0.1`) — no real network traffic |
| OS | Record exact kernel version and distro |

### Isolation Rules

- **Kill all non-essential background processes** before each run (`htop` to verify).
- Run each tool in **a fresh terminal session** with no other services active.
- Reboot between tool comparisons if a full cold-start baseline is needed.
- CookieFarm is tested both **natively** and **via Docker Compose** — record separately.

### Versions to Record

```bash
# CookieFarm
ckc --version
docker --version
go version

# DestructiveFarm
python3 --version
pip show flask  # record Flask version
```

---

## 2. Test Parameters

These are **fixed across both tools** to ensure comparability.

| Parameter | Value | Rationale |
|-----------|-------|-----------|
| Simulated teams | **40** | Realistic A/D competition size |
| Flags per team per run | **30** | One dummy exploit run, 30 flags per target |
| Total flags per run | **1,200** | 40 × 30 |
| Runs (rounds) | **10** | Simulates 10 consecutive CTF rounds |
| Total flags ingested | **12,000** | Full stress volume |
| Round interval | **30 seconds** | Realistic CTF round cadence |
| Exploit concurrency | **40 workers** | One worker per team (max parallelism) |
| Pagination endpoint | `/api/v1/flags/{id}?cursor=<data>` | CookieFarm; adapt for DestructiveFarm |
| Pages fetched per latency test | **50 sequential requests** | Statistically meaningful p50/p95/p99 |

---

## 3. Metric Definitions

### M1 — RAM Usage

> **Definition:** RSS (Resident Set Size) of the server process(es) in MiB, sampled every 2 seconds.

- For **DestructiveFarm**: measure the Flask server process (`python3 server.py`).
- For **CookieFarm**: measure the Go server process + the Docker container (`docker stats`).
- Report: **idle baseline**, **peak during ingest**, **steady-state after 10 rounds**.

### M2 — CPU Usage

> **Definition:** CPU percentage of the server process(es), sampled every 2 seconds via `pidstat`.

- Normalize to **single-core equivalent** (divide by number of pinned cores).
- Report: **average during ingest**, **peak spike**, **idle after ingest**.

### M3 — Flag Store Throughput

> **Definition:** Wall-clock time (seconds) to successfully store all 1,200 flags per round.

- Measure from **first exploit invocation** to **last DB write confirmed**.
- Report: per-round time for all 10 rounds, then compute mean ± stddev.
- Secondary metric: **flags/second** ingestion rate.

### M4 — UI Pagination Latency

> **Definition:** HTTP response time (ms) for `GET /api/v1/flags/{teamID}?cursor=<token>` (or equivalent endpoint on DestructiveFarm), measured with 12,000 flags already in the database.

- Test with **cold cache** (restart server, no prior requests) and **warm cache** (after 10 sequential prefetch requests).
- Report: **p50**, **p95**, **p99** across 50 requests.
- Tool: `hyperfine`, `wrk`, or the Python script provided in §6.

### M5 — Qualitative Observations

> **Definition:** Non-numeric dimensions assessed by the operator during and after the benchmark.

| Dimension | What to Observe |
|-----------|----------------|
| Setup complexity | Time from git clone to first running exploit (minutes) |
| Config surface | Number of required config files / env variables |
| Error handling | Behavior when an exploit crashes mid-run |
| UI responsiveness | Subjective smoothness of the web dashboard under load |
| Docker support | Native Docker Compose vs manual setup effort |
| Architecture clarity | How easy it is to add a new exploit |

---

## 4. Benchmark Procedures

### M1 — RAM Usage

**Step 1:** Start the server under test (cold start, empty DB).

**Step 2:** Record idle RAM baseline for 30 seconds.

```bash
# Get PID of the server
PID=$(pgrep -f "python3 server.py")      # DestructiveFarm
PID=$(pgrep -f "cookiefarm-server")      # CookieFarm native

# Sample every 2s for the duration of the test
while true; do
  ps -o pid,rss,vsz --no-headers -p $PID
  sleep 2
done | tee ram_samples.txt
```

**Step 3:** Trigger the dummy exploit runner (see §5) for all 10 rounds while sampling continues.

**Step 4:** Let server idle for 60 seconds post-ingest. Continue sampling.

**Step 5:** Parse `ram_samples.txt` → extract idle/peak/steady-state values.

---

### M2 — CPU Usage

```bash
# Install sysstat if needed: sudo apt install sysstat
pidstat -u -p $PID 2 > cpu_samples.txt &
PIDSTAT_PID=$!

# Run the full 10-round benchmark here...

kill $PIDSTAT_PID
```

Extract `%usr` and `%sys` columns from `cpu_samples.txt`. Sum them for total CPU%.

---

### M3 — Flag Store Throughput

**Instrumentation approach:** The dummy exploit (§5) prints a timestamp before its first flag output and after its last. The server log records the DB write timestamp of the last flag.

```
Round N throughput = T_last_db_write − T_first_exploit_start
```

If the server does not expose DB write timestamps in logs, use an alternative:

```bash
# Poll flag count in DB every 0.5s
while true; do
  COUNT=$(curl -s http://localhost:8080/api/v1/flags/count | jq '.count')
  echo "$(date +%s%3N) $COUNT"
  sleep 0.5
done | tee flag_count_timeline.txt
```

Plot `flag_count_timeline.txt` to visualize ingest rate over time.

---

### M4 — UI Pagination Latency

Ensure the DB contains **all 12,000 flags** (after 10 full rounds) before running this test.

**Cold cache test:**

```bash
# Restart the server (DB intact, in-memory cache cleared)
# Then immediately run:
python3 scripts/measure_pagination.py \
  --url "http://localhost:8080/api/v1/flags/40" \
  --requests 50 \
  --mode cold \
  --output latency_cold.json
```

**Warm cache test:**

```bash
# Do 10 warmup requests first, then measure 50:
python3 scripts/measure_pagination.py \
  --url "http://localhost:8080/api/v1/flags/40" \
  --requests 50 \
  --warmup 10 \
  --mode warm \
  --output latency_warm.json
```

The measurement script is provided in §6.

**For DestructiveFarm:** adapt the URL to its equivalent flag listing endpoint. If DestructiveFarm does not support cursor-based pagination, measure the full flag list endpoint instead and annotate the difference in your report.

---

### M5 — Qualitative Observations

Use this structured observation form during testing. Fill it in **real-time**, not from memory.

```
Tool: _______________
Tester: Franco
Date: _______________

[ ] Setup time (git clone → first exploit running): ___ minutes
[ ] Config files required: ___
[ ] Environment variables required: ___
[ ] Exploit crashed during test? Y/N — server behavior: ___
[ ] Dashboard loaded without delay under 12k flags? Y/N
[ ] Any errors in server logs? List: ___
[ ] Time to add a second dummy exploit: ___ minutes
[ ] Docker Compose startup time: ___ seconds (CookieFarm only)
```

---

## 5. Dummy Exploit Specification

Both tools must use **identical exploit logic** to eliminate exploit-side variance.

### Specification

- Accepts one argument: target IP (ignored in dummy mode)
- Sleeps for a **random delay** between 50–200ms (simulates real network latency)
- Outputs exactly **30 flag-formatted strings** to stdout
- Exits with code 0

### Implementation

```python
#!/usr/bin/env python3
"""
dummy_exploit.py — Benchmark dummy exploit
Usage: python3 dummy_exploit.py <target_ip>
Outputs 30 fake flags to stdout, simulating real exploit timing.
"""
import sys
import time
import random
import string

TARGET = sys.argv[1] if len(sys.argv) > 1 else "10.0.0.1"

# Simulate network round-trip
time.sleep(random.uniform(0.05, 0.20))

FLAG_PREFIX = "FLAG"  # Replace with competition prefix if needed
FLAG_CHARS = string.ascii_uppercase + string.digits

for i in range(30):
    flag = FLAG_PREFIX + "{" + "".join(random.choices(FLAG_CHARS, k=31)) + "}"
    print(flag)
    time.sleep(random.uniform(0.001, 0.005))  # Simulate per-flag delay
```

> ⚠️ **Important:** Substitute `FLAG_PREFIX` and flag format regex to match whichever flag format the server under test validates. If neither server validates format during this benchmark, any string works.

### Invocation for CookieFarm

Wrap in a CookieFarm-compatible exploit module using the `ckc` client interface. The core logic is identical to the script above.

### Invocation for DestructiveFarm

Pass `dummy_exploit.py` to `start_sploit.py` with `--host` set to the loopback server and 40 team IPs.

---

## 6. Data Collection Scripts

### `scripts/measure_pagination.py`

```python
#!/usr/bin/env python3
"""
Measures HTTP GET latency for a paginated flags endpoint.
Outputs p50, p95, p99 response times.
"""
import argparse
import json
import time
import statistics
import urllib.request

parser = argparse.ArgumentParser()
parser.add_argument("--url", required=True)
parser.add_argument("--requests", type=int, default=50)
parser.add_argument("--warmup", type=int, default=0)
parser.add_argument("--mode", choices=["cold", "warm"], default="cold")
parser.add_argument("--output", default="latency.json")
args = parser.parse_args()

def fetch_latency(url):
    start = time.perf_counter()
    with urllib.request.urlopen(url, timeout=10) as r:
        _ = r.read()
    return (time.perf_counter() - start) * 1000  # ms

# Warmup
for _ in range(args.warmup):
    fetch_latency(args.url)
    time.sleep(0.05)

# Measure
latencies = []
for i in range(args.requests):
    latencies.append(fetch_latency(args.url))
    time.sleep(0.05)  # 50ms inter-request gap

latencies_sorted = sorted(latencies)
results = {
    "mode": args.mode,
    "url": args.url,
    "n": args.requests,
    "p50_ms":  round(latencies_sorted[int(len(latencies) * 0.50)], 2),
    "p95_ms":  round(latencies_sorted[int(len(latencies) * 0.95)], 2),
    "p99_ms":  round(latencies_sorted[int(len(latencies) * 0.99)], 2),
    "mean_ms": round(statistics.mean(latencies), 2),
    "min_ms":  round(min(latencies), 2),
    "max_ms":  round(max(latencies), 2),
    "raw_ms":  [round(x, 2) for x in latencies],
}

with open(args.output, "w") as f:
    json.dump(results, f, indent=2)

print(f"[{args.mode}] p50={results['p50_ms']}ms  p95={results['p95_ms']}ms  p99={results['p99_ms']}ms")
```

### `scripts/parse_ram.py`

```python
#!/usr/bin/env python3
"""
Parses output of: ps -o pid,rss --no-headers -p $PID
Computes idle baseline, peak, and steady-state RSS in MiB.
"""
import sys

lines = [l.strip() for l in open(sys.argv[1]) if l.strip()]
rss_values = [int(l.split()[1]) / 1024 for l in lines]  # KB → MiB

baseline = sum(rss_values[:15]) / 15  # first 30s at 2s interval
peak = max(rss_values)
steady = sum(rss_values[-15:]) / 15   # last 30s

print(f"Baseline: {baseline:.1f} MiB")
print(f"Peak:     {peak:.1f} MiB")
print(f"Steady:   {steady:.1f} MiB")
```

---

## 7. Reporting Template

After running all tests, fill in this results table. Replace `???` with measured values.

### Results Summary

| Metric | DestructiveFarm | CookieFarm (Native) | CookieFarm (Docker) |
|--------|----------------|---------------------|---------------------|
| Idle RAM (MiB) | ??? | ??? | ??? |
| Peak RAM (MiB) | ??? | ??? | ??? |
| Avg CPU% during ingest | ??? | ??? | ??? |
| Peak CPU% spike | ??? | ??? | ??? |
| Mean flags/sec (ingest) | ??? | ??? | ??? |
| Round ingest time — mean (s) | ??? | ??? | ??? |
| Round ingest time — stddev (s) | ??? | ??? | ??? |
| Pagination p50 cold (ms) | ??? | ??? | ??? |
| Pagination p95 cold (ms) | ??? | ??? | ??? |
| Pagination p99 cold (ms) | ??? | ??? | ??? |
| Pagination p50 warm (ms) | ??? | ??? | ??? |
| Pagination p95 warm (ms) | ??? | ??? | ??? |
| Setup time (min) | ??? | ??? | ??? |

### Raw Data Files

Commit these files alongside your README:

```
benchmark/
├── destructivefarm/
│   ├── ram_samples.txt
│   ├── cpu_samples.txt
│   ├── flag_count_timeline.txt
│   ├── latency_cold.json
│   └── latency_warm.json
└── cookiefarm/
    ├── ram_samples.txt
    ├── cpu_samples.txt
    ├── flag_count_timeline.txt
    ├── latency_cold.json
    └── latency_warm.json
```

### Visualization

Generate charts from raw data using the following:

```bash
python3 scripts/generate_charts.py \
  --df-ram  benchmark/destructivefarm/ram_samples.txt \
  --cf-ram  benchmark/cookiefarm/ram_samples.txt \
  --df-lat  benchmark/destructivefarm/latency_cold.json \
  --cf-lat  benchmark/cookiefarm/latency_cold.json \
  --output  benchmark/charts/
```

Embed charts directly in the README with:

```markdown
![RAM Usage Over Time](benchmark/charts/ram_timeline.png)
![CPU Usage During Ingest](benchmark/charts/cpu_ingest.png)
![Flags per Second per Round](benchmark/charts/flags_per_second.png)
![Pagination Latency Distribution](benchmark/charts/latency_boxplot.png)
```

---

## 8. Reproducibility Checklist

Before publishing results, verify every item below.

**Environment**
- [ ] Exact OS version, kernel, and hardware specs documented
- [ ] Both tool versions (git commit SHA) recorded
- [ ] No background processes during measurement (verified with `htop`)
- [ ] Same CPU cores pinned for both tools
- [ ] Same dummy exploit binary used for both tools
- [ ] Empty DB confirmed before each tool's test run

**Measurements**
- [ ] All 10 rounds completed for M1/M2/M3
- [ ] 50 latency samples collected for M4 (cold and warm)
- [ ] Raw data files committed to repository
- [ ] Scripts used to parse/aggregate data committed to repository

**Reporting**
- [ ] Results table fully populated (no `???` remaining)
- [ ] Charts generated from raw data (not mocked)
- [ ] Anomalous rounds (if any) documented with reason
- [ ] CookieFarm Docker overhead measured separately from native

---

*Methodology authored for the ByteTheCookies team benchmark. CookieFarm is developed by [ByteTheCookies](https://github.com/ByteTheCookies). DestructiveFarm is developed by [DestructiveVoice](https://github.com/DestructiveVoice/DestructiveFarm).*
