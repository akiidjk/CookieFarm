# Benchmark Notes

## 30-04-2025

## Goal

Check the perfomance differences between the current implementation of the Database and the new one using mattn + some query tweks. We will run the same exploit on both implementations and compare the results.

### Summary

- 40 hosts
- sha: 9c2ed70be7e5248f3f7fd2c4ff9ca35b7c6a17e5
- branch: improve/database
- cks version: 1.3.0
- ckc version: 1.3.0

The original is the dev branch with hash 2d9a804d4bb855d93c7f4eb68dd20b6ee8d7c6da

## CK config

### 9c2ed70be7e5248f3f7fd2c4ff9ca35b7c6a17e5

32 GB RAM
32 vCPU
LAN

- Build command: `just server-build-prod`, `just server-build-plugins-prod`
- Command run: `../bin/cks -c`
- Build command: `just client-build-prod`
- Command run: `./bin/ckc exploit run -e benchmark -n CookieService -t 5 -T 10`

### 2d9a804d4bb855d93c7f4eb68dd20b6ee8d7c6da

32 GB RAM
16 vCPU
LAN

- Build command: `just client-build-prod`
- Command run: `./bin/ckc exploit run -e benchmark -n CookieService -t 5 -T 10`
- Build command: `just client-build-prod`
- Command run: `./bin/ckc exploit run -e benchmark -n CookieService -t 5 -T 10`

### Cks config used:

```yaml
configured: true

# Server
server:
  url_flag_checker: "http://localhost:5001/flags"
  team_token: "pippo"
  submit_flag_checker_time: 30
  max_flag_batch_size: 5000
  protocol: "cc_http"
  tick_time: 120
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

### Ckc config used:

```yaml
host: 127.0.0.1
username: cookieguest
port: 8080
https: false
```

### Exploit used

```python
#!/usr/bin/env python3
import requests
from cookiefarm import exploit_manager

@exploit_manager
def exploit(ip, port, name_service, flag_ids: list):
    for _ in range(30):
        r = requests.get(f"http://{ip}:{port}/get-flag")
        print(r.text)
```

### Results

#### Perfomance metrics:
