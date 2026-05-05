#!/usr/bin/env python3
"""
Measures HTTP/GET/POST latency for a paginated flags endpoint.
Outputs p50, p95, p99 response times.
"""

import argparse
import base64
import json
import statistics
import time
import urllib.request

parser = argparse.ArgumentParser()
parser.add_argument("--url", required=True)
parser.add_argument("--requests", type=int, default=50)
parser.add_argument("--warmup", type=int, default=0)
parser.add_argument("--mode", choices=["cold", "warm"], default="cold")
parser.add_argument("--output", default="latency.json")
parser.add_argument(
    "--cookie", default=None, help="Cookie string to pass, e.g. token=123"
)
parser.add_argument("--token", default=None, help="Token for Authorization header")
parser.add_argument(
    "--basic-auth", default=None, help="Basic auth credentials, e.g. username:password"
)
parser.add_argument(
    "--method", default="GET", choices=["GET", "POST"], help="HTTP Method"
)
parser.add_argument("--data", default=None, help="Data payload for POST requests")
args = parser.parse_args()


def fetch_latency(url):
    req = urllib.request.Request(url, method=args.method)
    if args.data:
        req.data = args.data.encode("utf-8")
        req.add_header("Content-Type", "application/x-www-form-urlencoded")
    if args.cookie:
        req.add_header("Cookie", args.cookie)
    if args.token:
        req.add_header("Authorization", args.token)
    if args.basic_auth:
        auth_bytes = args.basic_auth.encode("utf-8")
        base64_auth = base64.b64encode(auth_bytes).decode("utf-8")
        req.add_header("Authorization", f"Basic {base64_auth}")

    start = time.perf_counter()
    with urllib.request.urlopen(req, timeout=10) as r:
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
    "p50_ms": round(latencies_sorted[int(len(latencies) * 0.50)], 2),
    "p95_ms": round(latencies_sorted[int(len(latencies) * 0.95)], 2),
    "p99_ms": round(latencies_sorted[int(len(latencies) * 0.99)], 2),
    "mean_ms": round(statistics.mean(latencies), 2),
    "min_ms": round(min(latencies), 2),
    "max_ms": round(max(latencies), 2),
    "raw_ms": [round(x, 2) for x in latencies],
}

with open(args.output, "w") as f:
    json.dump(results, f, indent=2)

print(
    f"[{args.mode}] p50={results['p50_ms']}ms  p95={results['p95_ms']}ms  p99={results['p99_ms']}ms"
)
