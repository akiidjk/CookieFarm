#!/usr/bin/env python3
"""
generate_charts.py — Benchmark chart generator for CookieFarm vs DestructiveFarm.
"""

import argparse
import json
import os
import statistics
import sys

# ── Argument parsing ───────────────────────────────────────────────────────────

parser = argparse.ArgumentParser(
    description="Generate benchmark comparison charts for CookieFarm vs DestructiveFarm"
)
parser.add_argument(
    "--stats", required=False, help="Unified stats_samples.txt from get_cpu_ram.sh"
)
parser.add_argument(
    "--df-ram", required=False, help="DestructiveFarm ram_samples.txt (legacy)"
)
parser.add_argument(
    "--cf-ram", required=False, help="CookieFarm ram_samples.txt (legacy)"
)
parser.add_argument(
    "--df-cpu", required=False, help="DestructiveFarm cpu_samples.txt (legacy)"
)
parser.add_argument(
    "--cf-cpu", required=False, help="CookieFarm cpu_samples.txt (legacy)"
)
parser.add_argument(
    "--df-flags", required=True, help="DestructiveFarm flag_count_timeline.txt"
)
parser.add_argument(
    "--cf-flags", required=True, help="CookieFarm flag_count_timeline.txt"
)
parser.add_argument("--df-lat", required=True, help="DestructiveFarm latency_cold.json")
parser.add_argument("--cf-lat", required=True, help="CookieFarm latency_cold.json")
parser.add_argument(
    "--df-lat-warm",
    required=False,
    default=None,
    help="DestructiveFarm latency_warm.json",
)
parser.add_argument(
    "--cf-lat-warm", required=False, default=None, help="CookieFarm latency_warm.json"
)
parser.add_argument(
    "--cf-docker-ram",
    required=False,
    default=None,
    help="CookieFarm Docker ram_samples.txt (optional)",
)
parser.add_argument(
    "--cf-docker-lat",
    required=False,
    default=None,
    help="CookieFarm Docker latency_cold.json (optional)",
)
parser.add_argument("--output", required=True, help="Output directory for PNG charts")
parser.add_argument(
    "--flags-per-round",
    type=int,
    default=1200,
    help="Expected flags per round (default: 1200)",
)
args = parser.parse_args()

os.makedirs(args.output, exist_ok=True)

# ── Import plotly ──────────────────────────────────────────────────────────────
try:
    import plotly.graph_objects as go
    import plotly.io as pio
except ImportError:
    print(
        "ERROR: plotly not installed. Run: pip install plotly kaleido", file=sys.stderr
    )
    sys.exit(1)

COLORS = pio.templates[pio.templates.default].layout.colorway or [
    "#20b2aa",
    "#ff6b6b",
    "#ffa500",
    "#9370db",
    "#32cd32",
]

# ── File parsers ───────────────────────────────────────────────────────────────


def parse_unified_stats(path: str):
    df_ram, cf_ram = [], []
    df_cpu, cf_cpu = [], []

    with open(path) as f:
        for line in f:
            if line.startswith("FLASK:"):
                parts = line.split()
                if len(parts) >= 5:
                    df_cpu.append(float(parts[1]))
                    df_ram.append(float(parts[3]) / 1024.0)
            elif line.startswith("CKS:"):
                parts = line.split()
                if len(parts) >= 5:
                    cf_cpu.append(float(parts[1]))
                    cf_ram.append(float(parts[3]) / 1024.0)

    df_t = [i * 2 for i in range(len(df_ram))]
    cf_t = [i * 2 for i in range(len(cf_ram))]
    return (df_t, df_ram), (cf_t, cf_ram), (df_t, df_cpu), (cf_t, cf_cpu)


def parse_ram(path: str) -> tuple[list[float], list[float]]:
    lines = [l.strip() for l in open(path) if l.strip()]
    rss = []
    for line in lines:
        parts = line.split()
        try:
            rss.append(int(parts[-1]) / 1024.0)
        except (ValueError, IndexError):
            continue
    timestamps = [i * 2 for i in range(len(rss))]
    return timestamps, rss


def parse_cpu(path: str) -> tuple[list[float], list[float]]:
    cpu = []
    for line in open(path):
        parts = line.split()
        if len(parts) >= 8:
            try:
                usr = float(parts[2])
                sys_ = float(parts[3])
                cpu.append(usr + sys_)
            except ValueError:
                continue
    timestamps = [i * 2 for i in range(len(cpu))]
    return timestamps, cpu


def parse_flag_timeline(path: str) -> tuple[list[float], list[int]]:
    times, counts = [], []
    if not os.path.exists(path):
        return [], []
    for line in open(path):
        parts = line.strip().split()
        if len(parts) >= 2:
            try:
                times.append(int(parts[0]))
                counts.append(int(parts[1]))
            except ValueError:
                continue
    if not times:
        return [], []
    t0 = times[0]
    rel = [(t - t0) / 1000.0 for t in times]
    return rel, counts


def compute_flags_per_second(
    times: list[float], counts: list[int], flags_per_round: int
) -> list[float]:
    if not times or not counts:
        return []
    fps_list = []
    round_start_t = times[0]
    round_start_c = counts[0]

    for i in range(1, len(times)):
        delta_flags = counts[i] - round_start_c
        delta_t = times[i] - round_start_t

        if delta_flags >= flags_per_round:
            fps = delta_flags / delta_t if delta_t > 0 else 0
            fps_list.append(round(fps, 1))
            round_start_t = times[i]
            round_start_c = counts[i]

    if not fps_list and counts[-1] > counts[0] and times[-1] > times[0]:
        return [round((counts[-1] - counts[0]) / (times[-1] - times[0]), 1)]
    return fps_list


def parse_latency(path: str) -> dict:
    if not os.path.exists(path):
        return {"raw_ms": []}
    with open(path) as f:
        return json.load(f)


def save_chart(fig: go.Figure, name: str, caption: str, description: str):
    out_path = os.path.join(args.output, name)
    fig.write_image(out_path)
    with open(out_path + ".meta.json", "w") as f:
        json.dump({"caption": caption, "description": description}, f, indent=2)
    print(f"  ✓ {name}")


# ── Load data ──────────────────────────────────────────────────────────────────

print("Loading data files...")

if args.stats:
    (
        (df_ram_t, df_ram_v),
        (cf_ram_t, cf_ram_v),
        (df_cpu_t, df_cpu_v),
        (cf_cpu_t, cf_cpu_v),
    ) = parse_unified_stats(args.stats)
else:
    df_ram_t, df_ram_v = parse_ram(args.df_ram)
    cf_ram_t, cf_ram_v = parse_ram(args.cf_ram)
    df_cpu_t, df_cpu_v = parse_cpu(args.df_cpu)
    cf_cpu_t, cf_cpu_v = parse_cpu(args.cf_cpu)

df_flag_t, df_flag_c = parse_flag_timeline(args.df_flags)
cf_flag_t, cf_flag_c = parse_flag_timeline(args.cf_flags)
df_lat_cold = parse_latency(args.df_lat)
cf_lat_cold = parse_latency(args.cf_lat)
df_lat_warm = parse_latency(args.df_lat_warm) if args.df_lat_warm else None
cf_lat_warm = parse_latency(args.cf_lat_warm) if args.cf_lat_warm else None

docker_ram_t, docker_ram_v = ([], [])
if args.cf_docker_ram:
    docker_ram_t, docker_ram_v = parse_ram(args.cf_docker_ram)

docker_lat_cold = parse_latency(args.cf_docker_lat) if args.cf_docker_lat else None

# ── Chart 1: RAM Timeline ──────────────────────────────────────────────────────

print("\nGenerating charts...")
fig1 = go.Figure()
if df_ram_v:
    fig1.add_trace(
        go.Scatter(
            x=df_ram_t,
            y=df_ram_v,
            name="DestructiveFarm",
            line=dict(color=COLORS[0], width=2),
        )
    )
if cf_ram_v:
    fig1.add_trace(
        go.Scatter(
            x=cf_ram_t,
            y=cf_ram_v,
            name="CF Native",
            line=dict(color=COLORS[1], width=2),
        )
    )
if docker_ram_v:
    fig1.add_trace(
        go.Scatter(
            x=docker_ram_t,
            y=docker_ram_v,
            name="CF Docker",
            line=dict(color=COLORS[2], width=2, dash="dash"),
        )
    )

fig1.update_layout(
    title="RAM Usage Over Time",
    legend=dict(orientation="h", yanchor="bottom", y=1.02, xanchor="center", x=0.5),
)
fig1.update_xaxes(title_text="Time (s)")
fig1.update_yaxes(title_text="RAM (MiB)")
save_chart(
    fig1, "ram_timeline.png", "RAM Usage Over Time", "RSS memory usage over time"
)

# ── Chart 2: CPU Usage Bar ─────────────────────────────────────────────────────


def cpu_stats(vals):
    if not vals:
        return 0, 0, 0
    n = len(vals)
    baseline_n = min(15, n // 10)
    idle_n = min(15, n // 10)
    avg_ingest = (
        statistics.mean(vals[baseline_n : int(n * 0.9)])
        if len(vals[baseline_n : int(n * 0.9)]) > 0
        else statistics.mean(vals)
    )
    peak = max(vals)
    idle = statistics.mean(vals[-idle_n:]) if idle_n > 0 else 0
    return round(avg_ingest, 1), round(peak, 1), round(idle, 1)


df_avg, df_peak, df_idle = cpu_stats(df_cpu_v)
cf_avg, cf_peak, cf_idle = cpu_stats(cf_cpu_v)
states = ["Avg ingest", "Peak spike", "Idle"]

fig2 = go.Figure()
fig2.add_trace(
    go.Bar(
        name="DestructiveFarm",
        x=states,
        y=[df_avg, df_peak, df_idle],
        marker_color=COLORS[0],
    )
)
fig2.add_trace(
    go.Bar(
        name="CF Native", x=states, y=[cf_avg, cf_peak, cf_idle], marker_color=COLORS[1]
    )
)
fig2.update_layout(
    barmode="group",
    title="CPU Usage by State",
    legend=dict(orientation="h", yanchor="bottom", y=1.02, xanchor="center", x=0.5),
)
fig2.update_xaxes(title_text="State")
fig2.update_yaxes(title_text="CPU (%)")
fig2.update_traces(cliponaxis=False)
save_chart(
    fig2,
    "cpu_ingest.png",
    "CPU Usage by State",
    "Grouped bar chart of CPU% across states",
)

# ── Chart 3: Flags/sec per Round ──────────────────────────────────────────────

df_fps = compute_flags_per_second(df_flag_t, df_flag_c, args.flags_per_round)
cf_fps = compute_flags_per_second(cf_flag_t, cf_flag_c, args.flags_per_round)

if df_fps and cf_fps:
    rounds_x = list(range(1, max(len(df_fps), len(cf_fps)) + 1))
    fig3 = go.Figure()
    fig3.add_trace(
        go.Scatter(
            x=rounds_x[: len(df_fps)],
            y=df_fps,
            name="DestructiveFarm",
            mode="lines+markers",
            line=dict(color=COLORS[0], width=2),
        )
    )
    fig3.add_trace(
        go.Scatter(
            x=rounds_x[: len(cf_fps)],
            y=cf_fps,
            name="CF Native",
            mode="lines+markers",
            line=dict(color=COLORS[1], width=2),
        )
    )
    fig3.update_layout(
        title="Flags/sec per Round",
        legend=dict(orientation="h", yanchor="bottom", y=1.02, xanchor="center", x=0.5),
    )
    fig3.update_xaxes(title_text="Round")
    fig3.update_yaxes(title_text="Flags/sec")
    save_chart(
        fig3,
        "flags_per_second.png",
        "Flags/sec per Round",
        "Flag ingestion throughput per round",
    )

# ── Chart 4: Latency Box Plot ──────────────────────────────────────────────────

box_data = []
if df_lat_cold.get("raw_ms"):
    box_data.append(("DF Cold", df_lat_cold["raw_ms"], COLORS[0]))
if cf_lat_cold.get("raw_ms"):
    box_data.append(("CF Cold", cf_lat_cold["raw_ms"], COLORS[1]))
if docker_lat_cold and docker_lat_cold.get("raw_ms"):
    box_data.append(("Dk Cold", docker_lat_cold["raw_ms"], COLORS[2]))
if df_lat_warm and df_lat_warm.get("raw_ms"):
    box_data.append(("DF Warm", df_lat_warm["raw_ms"], COLORS[0]))
if cf_lat_warm and cf_lat_warm.get("raw_ms"):
    box_data.append(("CF Warm", cf_lat_warm["raw_ms"], COLORS[1]))

if box_data:
    fig4 = go.Figure()
    for label, vals, color in box_data:
        fig4.add_trace(
            go.Box(
                y=vals,
                name=label,
                marker_color=color,
                boxmean=True,
                line=dict(color=color),
            )
        )
    fig4.update_layout(title="Pagination Latency", showlegend=False)
    fig4.update_xaxes(title_text="Tool + Cache Mode")
    fig4.update_yaxes(title_text="Latency (ms)")
    save_chart(
        fig4,
        "latency_boxplot.png",
        "Pagination Latency",
        "Box plots comparing latencies",
    )

# ── Summary stats printout ─────────────────────────────────────────────────────
print("\n── Summary ─────────────────────────────────────────────────────────")
print(f"{'Metric':<30} {'DestructiveFarm':>14} {'CF Native':>12}")
print("-" * 58)


def ram_stats(vals):
    if not vals:
        return 0, 0, 0
    n = len(vals)
    base = statistics.mean(vals[: min(15, n)])
    peak = max(vals)
    steady = statistics.mean(vals[-min(15, n) :])
    return round(base, 1), round(peak, 1), round(steady, 1)


df_rb, df_rp, df_rs = ram_stats(df_ram_v)
cf_rb, cf_rp, cf_rs = ram_stats(cf_ram_v)
print(f"{'Idle RAM (MiB)':<30} {df_rb:>14} {cf_rb:>12}")
print(f"{'Peak RAM (MiB)':<30} {df_rp:>14} {cf_rp:>12}")
print(f"{'Avg CPU% (ingest)':<30} {df_avg:>14} {cf_avg:>12}")
print(
    f"{'Pagination p50 cold (ms)':<30} {df_lat_cold.get('p50_ms', '?'):>14} {cf_lat_cold.get('p50_ms', '?'):>12}"
)
print(
    f"{'Pagination p95 cold (ms)':<30} {df_lat_cold.get('p95_ms', '?'):>14} {cf_lat_cold.get('p95_ms', '?'):>12}"
)
print(
    f"{'Pagination p99 cold (ms)':<30} {df_lat_cold.get('p99_ms', '?'):>14} {cf_lat_cold.get('p99_ms', '?'):>12}"
)
if df_fps and cf_fps:
    df_mean_fps = round(statistics.mean(df_fps), 1)
    cf_mean_fps = round(statistics.mean(cf_fps), 1)
    print(f"{'Mean flags/sec':<30} {df_mean_fps:>14} {cf_mean_fps:>12}")
print(f"\nCharts saved to: {os.path.abspath(args.output)}")
