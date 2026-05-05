#!/usr/bin/env python3
"""
Parses output of: ps -o pid,command,%cpu,%mem,rss,vsz -p $PID
Computes idle baseline, peak, and steady-state RSS in MiB.
"""

import sys

with open(sys.argv[1]) as f:
    raw_lines = f.readlines()

# Skip header lines (those containing "PID")
data_lines = [line.strip() for line in raw_lines if line.strip() and "PID" not in line]

# Columns: PID COMMAND %CPU %MEM RSS VSZ
# RSS is column index 4 (KB → MiB)
rss_values = [int(line.split()[4]) / 1024 for line in data_lines]

baseline = sum(rss_values[:15]) / 15  # first 30s at 2s interval
peak = max(rss_values)
steady = sum(rss_values[-15:]) / 15  # last 30s

print(f"Baseline: {baseline:.1f} MiB")
print(f"Peak:     {peak:.1f} MiB")
print(f"Steady:   {steady:.1f} MiB")
