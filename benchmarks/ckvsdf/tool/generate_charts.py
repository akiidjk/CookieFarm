#!/usr/bin/env python3
"""
generate_charts.py — Benchmark chart generator for CookieFarm vs DestructiveFarm.
"""

import argparse
import json
import multiprocessing
import os
import statistics
import sys

try:
    CPU_COUNT = multiprocessing.cpu_count()
except NotImplementedError:
    CPU_COUNT = 1

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
    "--df-lat-break",
    required=False,
    default=None,
    help="DestructiveFarm latency_cold_breakQUERY.json",
)
parser.add_argument(
    "--cf-lat-break",
    required=False,
    default=None,
    help="CookieFarm latency_cold_breakQUERY.json",
)
parser.add_argument(
    "--df-lat-warm-break",
    required=False,
    default=None,
    help="DestructiveFarm latency_warm_breakQUERY.json",
)
parser.add_argument(
    "--cf-lat-warm-break",
    required=False,
    default=None,
    help="CookieFarm latency_warm_breakQUERY.json",
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

# ── Import matplotlib ────────────────────────────────────────────────────────
try:
    import matplotlib.colors as mcolors
    import matplotlib.pyplot as plt
except ImportError:
    print(
        "ERROR: matplotlib not installed. Run: pip install matplotlib", file=sys.stderr
    )
    sys.exit(1)

COLORS = [
    "#20b2aa",
    "#ff6b6b",
    "#ffa500",
    "#9370db",
    "#32cd32",
]

# ── File parsers ───────────────────────────────────────────────────────────────


def parse_unified_stats(path: str):
    data = {
        "df_server": {"cpu": [], "ram": []},
        "cf_server": {"cpu": [], "ram": []},
        "df_client": {"cpu": [], "ram": []},
        "cf_client": {"cpu": [], "ram": []},
    }

    with open(path) as f:
        for line in f:
            parts = line.split()
            if len(parts) >= 5:
                try:
                    cpu = float(parts[1]) / CPU_COUNT
                    ram = float(parts[3]) / 1024.0
                except ValueError:
                    continue
                if line.startswith("FLASK:"):
                    data["df_server"]["cpu"].append(cpu)
                    data["df_server"]["ram"].append(ram)
                elif line.startswith("CKS:"):
                    data["cf_server"]["cpu"].append(cpu)
                    data["cf_server"]["ram"].append(ram)
                elif line.startswith("CLIENTS:"):
                    data["df_client"]["cpu"].append(cpu)
                    data["df_client"]["ram"].append(ram)
                elif line.startswith("CKC:"):
                    data["cf_client"]["cpu"].append(cpu)
                    data["cf_client"]["ram"].append(ram)

    for k in data:
        data[k]["t"] = [i * 2 for i in range(len(data[k]["cpu"]))]
    return data


def parse_ram(path: str) -> tuple[list[float], list[float]]:
    lines = [l.strip() for l in open(path) if l.strip()]
    rss = []
    for line in lines:
        parts = line.split()
        try:
            rss.append(int(parts[-1]) / 1024.0)
        except (ValueError, IndexError):
            continue
    timestamps = [float(i * 2) for i in range(len(rss))]
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
    timestamps = [float(i * 2) for i in range(len(cpu))]
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


def save_chart(fig, name: str, caption: str, description: str):
    out_path = os.path.join(args.output, name)
    fig.savefig(out_path, bbox_inches="tight", dpi=150)
    plt.close(fig)
    with open(out_path + ".meta.json", "w") as f:
        json.dump({"caption": caption, "description": description}, f, indent=2)
    print(f"  ✓ {name}")


# ── Load data ──────────────────────────────────────────────────────────────────

print("Loading data files...")

if args.stats:
    stats_data = parse_unified_stats(args.stats)
else:
    df_ram_t, df_ram_v = parse_ram(args.df_ram)
    cf_ram_t, cf_ram_v = parse_ram(args.cf_ram)
    df_cpu_t, df_cpu_v = parse_cpu(args.df_cpu)
    cf_cpu_t, cf_cpu_v = parse_cpu(args.cf_cpu)
    stats_data = {
        "df_server": {"t": df_ram_t, "ram": df_ram_v, "cpu": df_cpu_v},
        "cf_server": {"t": cf_ram_t, "ram": cf_ram_v, "cpu": cf_cpu_v},
        "df_client": {"t": [], "ram": [], "cpu": []},
        "cf_client": {"t": [], "ram": [], "cpu": []},
    }

df_flag_t, df_flag_c = parse_flag_timeline(args.df_flags)
cf_flag_t, cf_flag_c = parse_flag_timeline(args.cf_flags)
df_lat_cold = parse_latency(args.df_lat)
cf_lat_cold = parse_latency(args.cf_lat)
df_lat_warm = parse_latency(args.df_lat_warm) if args.df_lat_warm else None
cf_lat_warm = parse_latency(args.cf_lat_warm) if args.cf_lat_warm else None

df_lat_break = parse_latency(args.df_lat_break) if args.df_lat_break else None
cf_lat_break = parse_latency(args.cf_lat_break) if args.cf_lat_break else None
df_lat_warm_break = (
    parse_latency(args.df_lat_warm_break) if args.df_lat_warm_break else None
)
cf_lat_warm_break = (
    parse_latency(args.cf_lat_warm_break) if args.cf_lat_warm_break else None
)

docker_ram_t, docker_ram_v = ([], [])
if args.cf_docker_ram:
    docker_ram_t, docker_ram_v = parse_ram(args.cf_docker_ram)

docker_lat_cold = parse_latency(args.cf_docker_lat) if args.cf_docker_lat else None

# ── Stats Helpers ──────────────────────────────────────────────────────────────


def get_ram_stats(vals):
    if not vals:
        return 0, 0, 0
    n = len(vals)
    base = statistics.mean(vals[: min(15, n)])
    peak = max(vals)
    steady = statistics.mean(vals[-min(15, n) :]) if n > 0 else 0
    return round(base, 1), round(peak, 1), round(steady, 1)


def get_cpu_stats(vals):
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


# ── Chart 1: RAM Usage Timeline ────────────────────────────────────────────────

print("\nGenerating charts...")
fig_ram_t, ax_ram_t = plt.subplots(figsize=(30, 16))

for key, name, color, ls in [
    ("df_server", "DF Server", COLORS[0], "-"),
    ("df_client", "DF Client", COLORS[0], "--"),
    ("cf_server", "CF Server", COLORS[1], "-"),
    ("cf_client", "CF Client", COLORS[1], "--"),
]:
    if stats_data[key]["ram"]:
        ax_ram_t.plot(
            stats_data[key]["t"],
            stats_data[key]["ram"],
            label=name,
            color=color,
            linestyle=ls,
            linewidth=2,
        )

if docker_ram_v:
    ax_ram_t.plot(
        docker_ram_t,
        docker_ram_v,
        label="CF Docker",
        color=COLORS[2],
        linestyle=":",
        linewidth=2,
    )

ax_ram_t.set_title("RAM Usage Over Time")
ax_ram_t.set_xlabel("Time (s)")
ax_ram_t.set_ylabel("RAM (MiB)")
ax_ram_t.legend(loc="upper center", bbox_to_anchor=(0.5, -0.15), ncol=3)
ax_ram_t.grid(True, alpha=0.3)
save_chart(
    fig_ram_t, "ram_timeline.png", "RAM Usage Over Time", "RSS memory usage over time"
)

# ── Chart 2: CPU Usage Timeline ────────────────────────────────────────────────

fig_cpu_t, ax_cpu_t = plt.subplots(figsize=(30, 16))

for key, name, color, ls in [
    ("df_server", "DF Server", COLORS[0], "-"),
    ("df_client", "DF Client", COLORS[0], "--"),
    ("cf_server", "CF Server", COLORS[1], "-"),
    ("cf_client", "CF Client", COLORS[1], "--"),
]:
    if stats_data[key]["cpu"]:
        ax_cpu_t.plot(
            stats_data[key]["t"],
            stats_data[key]["cpu"],
            label=name,
            color=color,
            linestyle=ls,
            linewidth=2,
        )

ax_cpu_t.set_title("CPU Usage Over Time")
ax_cpu_t.set_xlabel("Time (s)")
ax_cpu_t.set_ylabel("CPU (%)")
ax_cpu_t.legend(loc="upper center", bbox_to_anchor=(0.5, -0.15), ncol=3)
ax_cpu_t.grid(True, alpha=0.3)
save_chart(
    fig_cpu_t, "cpu_timeline.png", "CPU Usage Over Time", "CPU % utilization over time"
)

# ── Chart 3: RAM Overall (Bar) ─────────────────────────────────────────────────

fig_ram_bar, ax_ram_bar = plt.subplots(figsize=(10, 6))
states_ram = ["Idle (Start)", "Peak", "Steady (End)"]
x = [0, 1, 2]
width = 0.2

series_list_ram = [
    ("df_server", "DF Server", COLORS[0], 1.0),
    ("cf_server", "CF Server", COLORS[1], 1.0),
    ("df_client", "DF Client", COLORS[0], 0.5),
    ("cf_client", "CF Client", COLORS[1], 0.5),
]

active_series_ram = [s for s in series_list_ram if stats_data[s[0]]["ram"]]
num_series_ram = len(active_series_ram)

for idx, (key, name, color, opacity) in enumerate(active_series_ram):
    r_base, r_peak, r_steady = get_ram_stats(stats_data[key]["ram"])
    offset = (idx - num_series_ram / 2 + 0.5) * width
    ax_ram_bar.bar(
        [xi + offset for xi in x],
        [r_base, r_peak, r_steady],
        width,
        label=name,
        color=mcolors.to_rgba(color, alpha=opacity),
    )

ax_ram_bar.set_title("RAM Usage Summary by Role")
ax_ram_bar.set_xticks(x)
ax_ram_bar.set_xticklabels(states_ram)
ax_ram_bar.set_ylabel("RAM (MiB)")
ax_ram_bar.legend(loc="upper center", bbox_to_anchor=(0.5, -0.15), ncol=4)
ax_ram_bar.grid(True, alpha=0.3, axis="y")
save_chart(
    fig_ram_bar,
    "ram_summary.png",
    "RAM Usage Summary",
    "Grouped bar chart for RAM stats",
)

# ── Chart 4: CPU Overall (Bar) ─────────────────────────────────────────────────

fig_cpu_bar, ax_cpu_bar = plt.subplots(figsize=(10, 6))
states_cpu = ["Avg ingest", "Peak spike", "Idle"]

series_list_cpu = [
    ("df_server", "DF Server", COLORS[0], 1.0),
    ("cf_server", "CF Server", COLORS[1], 1.0),
    ("df_client", "DF Client", COLORS[0], 0.5),
    ("cf_client", "CF Client", COLORS[1], 0.5),
]

active_series_cpu = [s for s in series_list_cpu if stats_data[s[0]]["cpu"]]
num_series_cpu = len(active_series_cpu)

for idx, (key, name, color, opacity) in enumerate(active_series_cpu):
    c_avg, c_peak, c_idle = get_cpu_stats(stats_data[key]["cpu"])
    offset = (idx - num_series_cpu / 2 + 0.5) * width
    ax_cpu_bar.bar(
        [xi + offset for xi in x],
        [c_avg, c_peak, c_idle],
        width,
        label=name,
        color=mcolors.to_rgba(color, alpha=opacity),
    )

ax_cpu_bar.set_title("CPU Usage Summary by Role")
ax_cpu_bar.set_xticks(x)
ax_cpu_bar.set_xticklabels(states_cpu)
ax_cpu_bar.set_ylabel("CPU (%)")
ax_cpu_bar.legend(loc="upper center", bbox_to_anchor=(0.5, -0.15), ncol=4)
ax_cpu_bar.grid(True, alpha=0.3, axis="y")
save_chart(
    fig_cpu_bar,
    "cpu_summary.png",
    "CPU Usage Summary",
    "Grouped bar chart for CPU stats",
)

# ── Chart 5: Cumulative Flag Growth Over Time ──────────────────────────────────

if df_flag_t and df_flag_c or cf_flag_t and cf_flag_c:
    fig_flags, ax_flags = plt.subplots(figsize=(10, 6))
    if df_flag_t and df_flag_c:
        ax_flags.plot(
            df_flag_t, df_flag_c, label="DestructiveFarm", color=COLORS[0], linewidth=2
        )
        ax_flags.fill_between(df_flag_t, 0, df_flag_c, color=COLORS[0], alpha=0.1)
    if cf_flag_t and cf_flag_c:
        ax_flags.plot(
            cf_flag_t, cf_flag_c, label="CF Native", color=COLORS[1], linewidth=2
        )
        ax_flags.fill_between(cf_flag_t, 0, cf_flag_c, color=COLORS[1], alpha=0.1)

    ax_flags.set_title("Flag Ingestion Growth Over Time")
    ax_flags.set_xlabel("Time (s)")
    ax_flags.set_ylabel("Total Flags")
    ax_flags.legend(loc="upper center", bbox_to_anchor=(0.5, -0.15), ncol=2)
    ax_flags.grid(True, alpha=0.3)
    save_chart(
        fig_flags,
        "flags_growth.png",
        "Flag Ingestion Growth Over Time",
        "Cumulative flag count over time",
    )

df_fps = compute_flags_per_second(df_flag_t, df_flag_c, args.flags_per_round)
cf_fps = compute_flags_per_second(cf_flag_t, cf_flag_c, args.flags_per_round)

# ── Chart 6 & 7: Latency Box Plot & Timeline ───────────────────────────────────

box_data = []
if df_lat_cold.get("raw_ms"):
    box_data.append(("DF Cold", df_lat_cold["raw_ms"], COLORS[0], "-"))
if cf_lat_cold.get("raw_ms"):
    box_data.append(("CF Cold", cf_lat_cold["raw_ms"], COLORS[1], "-"))
if docker_lat_cold and docker_lat_cold.get("raw_ms"):
    box_data.append(("Dk Cold", docker_lat_cold["raw_ms"], COLORS[2], "-"))
if df_lat_warm and df_lat_warm.get("raw_ms"):
    box_data.append(("DF Warm", df_lat_warm["raw_ms"], COLORS[0], "--"))
if cf_lat_warm and cf_lat_warm.get("raw_ms"):
    box_data.append(("CF Warm", cf_lat_warm["raw_ms"], COLORS[1], "--"))

if df_lat_break and df_lat_break.get("raw_ms"):
    box_data.append(("DF Cold Brk", df_lat_break["raw_ms"], COLORS[0], ":"))
if cf_lat_break and cf_lat_break.get("raw_ms"):
    box_data.append(("CF Cold Brk", cf_lat_break["raw_ms"], COLORS[1], ":"))
if df_lat_warm_break and df_lat_warm_break.get("raw_ms"):
    box_data.append(("DF Warm Brk", df_lat_warm_break["raw_ms"], COLORS[0], "-."))
if cf_lat_warm_break and cf_lat_warm_break.get("raw_ms"):
    box_data.append(("CF Warm Brk", cf_lat_warm_break["raw_ms"], COLORS[1], "-."))

if box_data:
    # Box Plot
    fig_box, ax_box = plt.subplots(figsize=(10, 6))
    labels = []
    vals_list = []
    colors_list = []

    for label, vals, color, _ in box_data:
        labels.append(label)
        vals_list.append(vals)
        colors_list.append(color)

    bplot = ax_box.boxplot(
        vals_list,
        labels=labels,
        patch_artist=True,
        showmeans=True,
        boxprops=dict(facecolor="white"),
        medianprops=dict(color="black", linewidth=1.5),
        meanprops=dict(marker="^", markerfacecolor="black", markeredgecolor="black"),
    )

    for patch, color in zip(bplot["boxes"], colors_list):
        patch.set_facecolor(mcolors.to_rgba(color, alpha=0.5))
        patch.set_edgecolor(color)
        patch.set_linewidth(1.5)

    ax_box.set_title("Pagination Latency Boxplot")
    ax_box.set_ylabel("Latency (ms)")
    ax_box.grid(True, alpha=0.3, axis="y")
    plt.setp(ax_box.get_xticklabels(), rotation=45, ha="right")
    # Adjust layout to fit rotated labels
    fig_box.subplots_adjust(bottom=0.2)
    save_chart(
        fig_box,
        "latency_boxplot.png",
        "Pagination Latency",
        "Box plots comparing latencies",
    )

    # Sequence / Throughput line graph
    fig_seq, ax_seq = plt.subplots(figsize=(10, 6))
    for label, vals, color, ls in box_data:
        x_vals = list(range(1, len(vals) + 1))
        ax_seq.plot(
            x_vals,
            vals,
            label=label,
            color=color,
            linestyle=ls,
            marker="o",
            markersize=4,
            alpha=0.8,
        )

    ax_seq.set_title("Latency Across Request Sequence")
    ax_seq.set_xlabel("Request Number")
    ax_seq.set_ylabel("Latency (ms)")
    ax_seq.legend(loc="upper center", bbox_to_anchor=(0.5, -0.15), ncol=3)
    ax_seq.grid(True, alpha=0.3)
    save_chart(
        fig_seq,
        "latency_sequence.png",
        "Latency Throughput Graph",
        "Latency changes over subsequent API requests",
    )

# ── Summary stats printout ─────────────────────────────────────────────────────
print("\n── Summary ─────────────────────────────────────────────────────────")
print(f"{'Metric':<30} {'DestructiveFarm':>14} {'CF Native':>12}")
print("-" * 58)

df_sb, df_sp, df_ss = get_ram_stats(stats_data["df_server"]["ram"])
cf_sb, cf_sp, cf_ss = get_ram_stats(stats_data["cf_server"]["ram"])
print(f"{'Server Idle RAM (MiB)':<30} {df_sb:>14} {cf_sb:>12}")
print(f"{'Server Peak RAM (MiB)':<30} {df_sp:>14} {cf_sp:>12}")

df_cb, df_cp, df_cs = get_ram_stats(stats_data["df_client"]["ram"])
cf_cb, cf_cp, cf_cs = get_ram_stats(stats_data["cf_client"]["ram"])
if stats_data["df_client"]["ram"] or stats_data["cf_client"]["ram"]:
    print(f"{'Client Peak RAM (MiB)':<30} {df_cp:>14} {cf_cp:>12}")

df_ca, df_cpeak, df_ci = get_cpu_stats(stats_data["df_server"]["cpu"])
cf_ca, cf_cpeak, cf_ci = get_cpu_stats(stats_data["cf_server"]["cpu"])
print(f"{'Server Avg CPU%':<30} {df_ca:>14} {cf_ca:>12}")

df_cla, df_clpeak, df_cli = get_cpu_stats(stats_data["df_client"]["cpu"])
cf_cla, cf_clpeak, cf_cli = get_cpu_stats(stats_data["cf_client"]["cpu"])
if stats_data["df_client"]["cpu"] or stats_data["cf_client"]["cpu"]:
    print(f"{'Client Avg CPU%':<30} {df_cla:>14} {cf_cla:>12}")

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

print("\n── 2M Flags Scale Test (Break vs Normal Difference) ────────────────")
print(
    f"{'Metric (Diff = Break - Normal)':<35} {'DestructiveFarm':>15} {'CF Native':>15}"
)
print("-" * 68)


def diff_str(brk_dict, norm_dict, key):
    if brk_dict and norm_dict and key in brk_dict and key in norm_dict:
        try:
            diff = float(brk_dict[key]) - float(norm_dict[key])
            return f"{diff:+.2f}"
        except (ValueError, TypeError):
            return "?"
    return "?"


print(
    f"{'Cold p50 diff (ms)':<35} {diff_str(df_lat_break, df_lat_cold, 'p50_ms'):>15} {diff_str(cf_lat_break, cf_lat_cold, 'p50_ms'):>15}"
)
print(
    f"{'Cold p95 diff (ms)':<35} {diff_str(df_lat_break, df_lat_cold, 'p95_ms'):>15} {diff_str(cf_lat_break, cf_lat_cold, 'p95_ms'):>15}"
)
print(
    f"{'Cold p99 diff (ms)':<35} {diff_str(df_lat_break, df_lat_cold, 'p99_ms'):>15} {diff_str(cf_lat_break, cf_lat_cold, 'p99_ms'):>15}"
)
print(
    f"{'Warm p50 diff (ms)':<35} {diff_str(df_lat_warm_break, df_lat_warm, 'p50_ms'):>15} {diff_str(cf_lat_warm_break, cf_lat_warm, 'p50_ms'):>15}"
)
print(
    f"{'Warm p95 diff (ms)':<35} {diff_str(df_lat_warm_break, df_lat_warm, 'p95_ms'):>15} {diff_str(cf_lat_warm_break, cf_lat_warm, 'p95_ms'):>15}"
)
print(
    f"{'Warm p99 diff (ms)':<35} {diff_str(df_lat_warm_break, df_lat_warm, 'p99_ms'):>15} {diff_str(cf_lat_warm_break, cf_lat_warm, 'p99_ms'):>15}"
)

print("\n── 📋 Full Summary (Markdown) ──────────────────────────────────────")
print("| Metric | DestructiveFarm | CookieFarm (Native) |  Winner |")
print("|--------|:--------------:|:-------------------:|:------:|")


def get_winner(d, c, lower_is_better=True):
    try:
        d_val, c_val = float(d), float(c)
        if d_val == c_val:
            return "Tie"
        if lower_is_better:
            return "DestructiveFarm" if d_val < c_val else "CookieFarm"
        return "DestructiveFarm" if d_val > c_val else "CookieFarm"
    except (ValueError, TypeError):
        return "???"


df_fps_val = round(statistics.mean(df_fps), 1) if df_fps else "?"
cf_fps_val = round(statistics.mean(cf_fps), 1) if cf_fps else "?"

df_ingest_time = (
    round(args.flags_per_round / float(df_fps_val), 2)
    if df_fps_val != "?" and float(df_fps_val) > 0
    else "?"
)
cf_ingest_time = (
    round(args.flags_per_round / float(cf_fps_val), 2)
    if cf_fps_val != "?" and float(cf_fps_val) > 0
    else "?"
)

df_p50 = df_lat_cold.get("p50_ms", "?")
cf_p50 = cf_lat_cold.get("p50_ms", "?")
df_p99 = df_lat_cold.get("p99_ms", "?")
cf_p99 = cf_lat_cold.get("p99_ms", "?")

print(
    f"| Server Idle RAM | {df_sb} MiB | {cf_sb} MiB |  {get_winner(df_sb, cf_sb, True)} |"
)
print(
    f"| Server Peak RAM | {df_sp} MiB | {cf_sp} MiB |  {get_winner(df_sp, cf_sp, True)} |"
)
print(
    f"| Client Peak RAM | {df_cp} MiB | {cf_cp} MiB |  {get_winner(df_cp, cf_cp, True)} |"
)
print(f"| Server Avg CPU% | {df_ca}% | {cf_ca}% |  {get_winner(df_ca, cf_ca, True)} |")
print(
    f"| Server Peak CPU% | {df_cpeak}% | {cf_cpeak}% | {get_winner(df_cpeak, cf_cpeak, True)} |"
)
print(
    f"| Client Avg CPU% | {df_cla}% | {cf_cla}% | {get_winner(df_cla, cf_cla, True)} |"
)
print(
    f"| Client Peak CPU% | {df_clpeak}% | {cf_clpeak}% | {get_winner(df_clpeak, cf_clpeak, True)} |"
)
print(
    f"| Mean ingest time/round | {df_ingest_time} s | {cf_ingest_time} s | {get_winner(df_ingest_time, cf_ingest_time, True)} |"
)
print(
    f"| Flags/sec | {df_fps_val} | {cf_fps_val} | {get_winner(df_fps_val, cf_fps_val, False)} |"
)
print(
    f"| Pagination p50 (cold) | {df_p50} ms | {cf_p50} ms |  {get_winner(df_p50, cf_p50, True)} |"
)
print(
    f"| Pagination p99 (cold) | {df_p99} ms | {cf_p99} ms | {get_winner(df_p99, cf_p99, True)} |"
)

df_brk_p50_diff = diff_str(df_lat_break, df_lat_cold, "p50_ms")
cf_brk_p50_diff = diff_str(cf_lat_break, cf_lat_cold, "p50_ms")
df_brk_p99_diff = diff_str(df_lat_break, df_lat_cold, "p99_ms")
cf_brk_p99_diff = diff_str(cf_lat_break, cf_lat_cold, "p99_ms")

print(
    f"| 2M Flags Scale p50 Diff (cold) | {df_brk_p50_diff} ms | {cf_brk_p50_diff} ms |  {get_winner(df_brk_p50_diff, cf_brk_p50_diff, True)} |"
)
print(
    f"| 2M Flags Scale p99 Diff (cold) | {df_brk_p99_diff} ms | {cf_brk_p99_diff} ms |  {get_winner(df_brk_p99_diff, cf_brk_p99_diff, True)} |"
)

print(f"\nCharts saved to: {os.path.abspath(args.output)}")
