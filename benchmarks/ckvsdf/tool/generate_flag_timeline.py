#!/usr/bin/env python3
"""
generate_flag_timeline.py — Records flag count growth over time during a benchmark run.

Polls the server's flag count endpoint every 500ms and writes:
    <epoch_ms> <count>
to a file, which is later consumed by generate_charts.py.

Usage:
    python3 generate_flag_timeline.py \
        --url   http://localhost:8080/api/v1/flags/count \
        --field count \
        --output benchmark/cookiefarm/flag_count_timeline.txt \
        --duration 600

    # Or query a SQLite database directly:
    python3 generate_flag_timeline.py \
        --db cookiefarm.db \
        --query "SELECT COUNT(*) FROM flags" \
        --output benchmark/cookiefarm/flag_count_timeline.txt \
        --duration 600

    # Or stop manually with Ctrl+C — partial data is always saved.

Supported response shapes:
    {"count": 1234}
    {"data": {"total": 1234}}
    {"flags": [...]}          ← will count array length if --field not resolvable
    Plain integer text: 1234
"""

import argparse
import json
import os
import signal
import sqlite3
import sys
import time
import urllib.error
import urllib.request

# ── Args ───────────────────────────────────────────────────────────────────────

parser = argparse.ArgumentParser(
    description="Poll a flag count endpoint or database and record a timeline for benchmark charting."
)
parser.add_argument("--url", required=False, help="Full URL of the flag count endpoint")
parser.add_argument(
    "--db", required=False, help="Path to SQLite database to query directly"
)
parser.add_argument(
    "--query",
    required=False,
    default="SELECT COUNT(*) FROM flags",
    help="SQL query to execute if --db is used",
)
parser.add_argument(
    "--field",
    default="count",
    help="JSON field name for the count value (default: count). Supports dot notation: data.total",
)
parser.add_argument(
    "--output",
    required=True,
    help="Output file path, e.g. benchmark/cookiefarm/flag_count_timeline.txt",
)
parser.add_argument(
    "--interval",
    type=float,
    default=0.5,
    help="Polling interval in seconds (default: 0.5)",
)
parser.add_argument(
    "--duration",
    type=int,
    default=0,
    help="Stop after N seconds (default: 0 = run until Ctrl+C)",
)
parser.add_argument(
    "--timeout",
    type=int,
    default=5,
    help="HTTP request or DB timeout in seconds (default: 5)",
)
parser.add_argument(
    "--header",
    action="append",
    default=[],
    metavar="KEY:VALUE",
    help="Extra HTTP headers, e.g. --header 'Authorization: Bearer TOKEN'. Repeatable.",
)
args = parser.parse_args()

if not args.url and not args.db:
    parser.error("Must provide either --url or --db")

# ── Helpers ────────────────────────────────────────────────────────────────────


def resolve_field(data: dict, dotpath: str):
    """Navigate dot-separated keys into a nested dict."""
    parts = dotpath.split(".")
    node = data
    for p in parts:
        if isinstance(node, dict) and p in node:
            node = node[p]
        else:
            return None
    return node


def extract_count(body: bytes) -> int | None:
    text = body.decode("utf-8", errors="replace").strip()

    # Try plain integer response first
    try:
        return int(text)
    except ValueError:
        pass

    # Try JSON
    try:
        data = json.loads(text)
    except json.JSONDecodeError:
        return None

    # Dot-path field resolution
    val = resolve_field(data, args.field)
    if val is not None:
        try:
            return int(val)
        except (TypeError, ValueError):
            pass

    # Fallback: count array if top-level is list
    if isinstance(data, list):
        return len(data)

    # Fallback: look for any key named 'count', 'total', 'size', 'length'
    for key in ("count", "total", "size", "length", "n"):
        if key in data:
            try:
                return int(data[key])
            except (TypeError, ValueError):
                pass

    return None


def build_request(url: str) -> urllib.request.Request:
    req = urllib.request.Request(url)
    for h in args.header:
        if ":" in h:
            k, v = h.split(":", 1)
            req.add_header(k.strip(), v.strip())
    return req


def poll_db() -> int | None:
    try:
        # Use URI connection with timeout
        conn = sqlite3.connect(args.db, timeout=args.timeout)
        cursor = conn.cursor()
        cursor.execute(args.query)
        row = cursor.fetchone()
        conn.close()
        if row is not None:
            return int(row[0])
    except Exception as e:
        print(f"  DB Error: {e}")
    return None


# ── Setup output file ──────────────────────────────────────────────────────────

os.makedirs(
    os.path.dirname(args.output) if os.path.dirname(args.output) else ".", exist_ok=True
)

samples: list[tuple[int, int]] = []
start_time = time.time()


def flush_output():
    with open(args.output, "w") as f:
        for epoch_ms, count in samples:
            f.write(f"{epoch_ms} {count}\n")
    print(f"\n✓ Saved {len(samples)} samples → {args.output}")


def handle_interrupt(sig, frame):
    flush_output()
    sys.exit(0)


signal.signal(signal.SIGINT, handle_interrupt)
signal.signal(signal.SIGTERM, handle_interrupt)

# ── Polling loop ───────────────────────────────────────────────────────────────

if args.db:
    print(f"Polling DB: {args.db}")
    print(f"Query:   {args.query}")
else:
    print(f"Polling URL: {args.url}")
    print(f"Field:   {args.field}")

print(f"Interval: {args.interval}s  |  Duration: {args.duration or '∞'}s")
print(f"Output:  {args.output}")
print("Press Ctrl+C to stop early.\n")

req = build_request(args.url) if args.url else None
consecutive_errors = 0
MAX_ERRORS = 10
last_count = 0

while True:
    now_ms = int(time.time() * 1000)
    elapsed = time.time() - start_time

    count = None
    try:
        if args.db:
            count = poll_db()
        else:
            with urllib.request.urlopen(req, timeout=args.timeout) as resp:
                body = resp.read()
            count = extract_count(body)

        if count is None:
            if not args.db:
                print(
                    f"  [{elapsed:6.1f}s] WARNING: Could not parse count from response."
                )
            consecutive_errors += 1
        else:
            samples.append((now_ms, count))
            delta = count - last_count
            bar = "█" * min(
                40,
                int(
                    count
                    / max(
                        1,
                        (
                            args.flags_per_round
                            if hasattr(args, "flags_per_round")
                            else 1200
                        ),
                    )
                    * 40
                ),
            )
            print(
                f"  [{elapsed:6.1f}s] flags={count:6d}  Δ={delta:+5d}  {bar}", end="\r"
            )
            last_count = count
            consecutive_errors = 0

    except urllib.error.URLError as e:
        consecutive_errors += 1
        print(
            f"  [{elapsed:6.1f}s] HTTP error ({consecutive_errors}/{MAX_ERRORS}): {e.reason}"
        )
        if consecutive_errors >= MAX_ERRORS:
            print(f"\nERROR: {MAX_ERRORS} consecutive failures. Stopping.")
            flush_output()
            sys.exit(1)

    except Exception as e:
        consecutive_errors += 1
        print(f"  [{elapsed:6.1f}s] Unexpected error: {e}")

    # Check duration limit
    if args.duration > 0 and elapsed >= args.duration:
        print(f"\nDuration limit reached ({args.duration}s).")
        break

    time.sleep(args.interval)

flush_output()
