# 📊 Benchmark: CookieFarm vs DestructiveFarm
  
> **Date:** 02-05-2026    
> **CookieFarm commit:** `2b359d52d4791df23de653a6580490162dd441c5`  
> **DestructiveFarm commit:** `69cc5821a16fc38e6670e666bc0c8d5ee311e57a`  
> **Host OS:** Linux 6.18.25-1-lts 
> **CPU:** AMD Ryzen 7 3700X (16) @ 4.43 GHz  
> **RAM:** 32 GB DDR4 @ 3200 MHz
> **GO Version:** go1.26.2
> **Python Version:** Python 3.14.4

---

## ⚙️ Test Parameters

| Parameter | Value |
|-----------|-------|
| Simulated teams | 40 |
| Flags per team per run | 30 |
| Total flags per round | 1,1170 |
| Total rounds | 240 (full 8h A/D) |
| Total flags ingested | 288,000 |
| Round interval | 5s (not realistic but we want MAX) |
| Exploit concurrency | ? workers |
| Pagination requests sampled | 50 |

## Config used

Cookiefarm's config:

```yaml
configured: true

# Server
server:
  url_flag_checker: "http://localhost:5001/flags"
  team_token: "2b359d52d4791df23de653a6580490162dd441c5"
  submit_flag_checker_time: 30
  max_flag_batch_size: 5000
  protocol: "cc_http"
  tick_time: 30
  flag_ttl: 0 # in ticks
  start_time: "2023-10-01T00:00:00Z"
  end_time: "2023-10-31T23:59:59Z"

# Client
shared:
  services:
    CookieService: 8081
    vulnify: 1337
    app-nc: 1338
  range_ip_teams: 40
  format_ip_teams: "10.10.{}.1"
  my_team_id: 1
  regex_flag: "[A-Z0-9]{31}="
  nop_team: 0
  url_flag_ids: "http://localhost:5001/flagIds"
```

DestructiveFarm's config:

```python
CONFIG = {
    "TEAMS": {"Team #{}".format(i): "10.10.{}.1".format(i) for i in range(0, 40 + 1)},
    "FLAG_FORMAT": r"[A-Z0-9]{31}=",
    "SYSTEM_PROTOCOL": "ructf_http",
    "SYSTEM_URL": "http://localhost:5001/flags",
    "SYSTEM_TOKEN": "password",
    "SUBMIT_FLAG_LIMIT": 100,
    "SUBMIT_PERIOD": 5,
    "FLAG_LIFETIME": 5 * 60,
    "SERVER_PASSWORD": "password",
    "ENABLE_API_AUTH": False,
    "API_TOKEN": "00000000000000000000",
}
```

### Exploit used

Cookiefarm's exploit

```python
#!/usr/bin/env python3
import requests
from cookiefarm import exploit_manager

# "ip" are the IP address of the target team (example: 10.10.X.1)
# "port" is the port of the target service (example: 1337)
# "name_service" is the name of the service to exploit (example: "CookieService")
# "flag_ids" is the flag IDs of the target team and target service (example: [{"username": "psQSDAasd", "password": "qweqweqwe"}, {"username": "sdafjhAS", "password": "HIUOasdb"}])


@exploit_manager
def exploit(ip, port, name_service, flag_ids: list):
    for _ in range(30):
        r = requests.get(f"http://{ip}:{port}/get-flag")
        print(r.text)

```

DestructiveFarm's exploit

```python
#!/usr/bin/env python3
import sys

import requests


def exploit(ip, port, name_service, flag_ids: list):
    for _ in range(30):
        r = requests.get(f"http://{ip}:{port}/get-flag")
        print(r.text, flush=True)


exploit(sys.argv[1], 8081, None, [])
```

---

## 🧠 M1 — RAM Usage

> Sampled every 2s via `ps -o rss`. Values in **MiB**.

| State | DestructiveFarm | CookieFarm (Native) | CookieFarm (Docker) |
|-------|:--------------:|:-------------------:|:-------------------:|
| Idle baseline | ??? MiB | ??? MiB | ??? MiB |
| Peak (during ingest) | ??? MiB | ??? MiB | ??? MiB |
| Steady-state (post-ingest) | ??? MiB | ??? MiB | ??? MiB |

<!--
  Replace ??? with values from: python3 scripts/parse_ram.py benchmark/<tool>/ram_samples.txt
-->

### RAM Timeline

![RAM Usage Over Time](benchmark/charts/ram_timeline.png)

---

## ⚡ M2 — CPU Usage

> Sampled every 2s via `pidstat`. Normalized to single-core equivalent.

| State | DestructiveFarm | CookieFarm (Native) | CookieFarm (Docker) |
|-------|:--------------:|:-------------------:|:-------------------:|
| Average during ingest | ???% | ???% | ???% |
| Peak spike | ???% | ???% | ???% |
| Idle after ingest | ???% | ???% | ???% |

![CPU Usage During Ingest](benchmark/charts/cpu_ingest.png)

---

## 🚩 M3 — Flag Store Throughput

> Wall-clock time per round to store all 1,200 flags. 10 rounds total.

| Round | DestructiveFarm (s) | CookieFarm Native (s) | CookieFarm Docker (s) |
|:-----:|:-------------------:|:---------------------:|:---------------------:|
| 1 | ??? | ??? | ??? |
| 2 | ??? | ??? | ??? |
| 3 | ??? | ??? | ??? |
| 4 | ??? | ??? | ??? |
| 5 | ??? | ??? | ??? |
| 6 | ??? | ??? | ??? |
| 7 | ??? | ??? | ??? |
| 8 | ??? | ??? | ??? |
| 9 | ??? | ??? | ??? |
| 10 | ??? | ??? | ??? |
| **Mean** | **???** | **???** | **???** |
| **Stddev** | **±???** | **±???** | **±???** |
| **Flags/sec** | **???** | **???** | **???** |

![Flags per Second per Round](benchmark/charts/flags_per_second.png)

---

## 🌐 M4 — UI Pagination Latency

> `GET /api/v1/flags/40?cursor=<token>` — 50 requests, 12,000 flags in DB.  
> DestructiveFarm endpoint: `_______________` *(adapt if different)*

### Cold Cache (server restarted, no prior requests)

| Percentile | DestructiveFarm | CookieFarm (Native) | CookieFarm (Docker) |
|:----------:|:--------------:|:-------------------:|:-------------------:|
| p50 | ??? ms | ??? ms | ??? ms |
| p95 | ??? ms | ??? ms | ??? ms |
| p99 | ??? ms | ??? ms | ??? ms |
| Mean | ??? ms | ??? ms | ??? ms |
| Min | ??? ms | ??? ms | ??? ms |
| Max | ??? ms | ??? ms | ??? ms |

### Warm Cache (10 warmup requests prior)

| Percentile | DestructiveFarm | CookieFarm (Native) | CookieFarm (Docker) |
|:----------:|:--------------:|:-------------------:|:-------------------:|
| p50 | ??? ms | ??? ms | ??? ms |
| p95 | ??? ms | ??? ms | ??? ms |
| p99 | ??? ms | ??? ms | ??? ms |
| Mean | ??? ms | ??? ms | ??? ms |

![Pagination Latency Distribution](benchmark/charts/latency_boxplot.png)

---

## 🔍 M5 — Qualitative Observations

### DestructiveFarm

| Dimension | Observation |
|-----------|-------------|
| Setup time (clone → first exploit) | ??? min |
| Config files required | ??? |
| Env variables required | ??? |
| Exploit crash behavior | ??? |
| Dashboard feel under 12k flags | ??? |
| Time to add a second exploit | ??? min |
| Server errors during test | ??? |
| Notes | ??? |

### CookieFarm

| Dimension | Observation |
|-----------|-------------|
| Setup time — Native (clone → first exploit) | ??? min |
| Setup time — Docker Compose | ??? min |
| Config files required | ??? |
| Env variables required | ??? |
| Exploit crash behavior | ??? |
| Dashboard feel under 12k flags | ??? |
| Time to add a second exploit | ??? min |
| Server errors during test | ??? |
| Docker Compose startup time | ??? s |
| Notes | ??? |

---

## 📋 Full Summary

| Metric | DestructiveFarm | CookieFarm (Native) | CookieFarm (Docker) | Winner |
|--------|:--------------:|:-------------------:|:-------------------:|:------:|
| Idle RAM | ??? MiB | ??? MiB | ??? MiB | ??? |
| Peak RAM | ??? MiB | ??? MiB | ??? MiB | ??? |
| Avg CPU% (ingest) | ???% | ???% | ???% | ??? |
| Mean ingest time/round | ??? s | ??? s | ??? s | ??? |
| Flags/sec | ??? | ??? | ??? | ??? |
| Pagination p50 (cold) | ??? ms | ??? ms | ??? ms | ??? |
| Pagination p99 (cold) | ??? ms | ??? ms | ??? ms | ??? |
| Setup time | ??? min | ??? min | ??? min | ??? |

---

## 📁 Raw Data

All raw measurement files are committed under `benchmark/`:

```
benchmark/
├── destructivefarm/
│   ├── ram_samples.txt
│   ├── cpu_samples.txt
│   ├── flag_count_timeline.txt
│   ├── latency_cold.json
│   └── latency_warm.json
├── cookiefarm/
│   ├── ram_samples.txt
│   ├── cpu_samples.txt
│   ├── flag_count_timeline.txt
│   ├── latency_cold.json
│   └── latency_warm.json
└── charts/
    ├── ram_timeline.png
    ├── cpu_ingest.png
    ├── flags_per_second.png
    └── latency_boxplot.png
```

To regenerate charts from raw data:

```bash
python3 scripts/generate_charts.py \
  --df-ram  benchmark/destructivefarm/ram_samples.txt \
  --cf-ram  benchmark/cookiefarm/ram_samples.txt \
  --df-lat  benchmark/destructivefarm/latency_cold.json \
  --cf-lat  benchmark/cookiefarm/latency_cold.json \
  --output  benchmark/charts/
```

---

*Full methodology: [BENCHMARK_METHODOLOGY.md](./benchmark-methodology.md)*  
*CookieFarm — [github.com/ByteTheCookies/CookieFarm](https://github.com/ByteTheCookies/CookieFarm)*  
*DestructiveFarm — [github.com/DestructiveVoice/DestructiveFarm](https://github.com/DestructiveVoice/DestructiveFarm)*
