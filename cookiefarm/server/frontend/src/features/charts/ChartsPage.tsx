import { useMemo, useState } from "react";
import * as echarts from "echarts/core";
import { mutate as swrMutate } from "swr";
import { BarChart, LineChart, PieChart } from "echarts/charts";
import { Banner } from "@cloudflare/kumo/components/banner";
import { Breadcrumbs } from "@cloudflare/kumo/components/breadcrumbs";
import { Button } from "@cloudflare/kumo/components/button";
import { Chart, ChartLegend, ChartPalette } from "@cloudflare/kumo/components/chart";
import {
  AriaComponent,
  AxisPointerComponent,
  GridComponent,
  LegendComponent,
  TooltipComponent,
} from "echarts/components";
import { CanvasRenderer } from "echarts/renderers";
import { WarningCircle } from "@phosphor-icons/react";
import { configKey, useConfig } from "@/api/config";
import { useAllFlags } from "@/api/flags";
import { useStatsSummary } from "@/api/stats";
import { PageHeader } from "@/components/kumo/page-header/page-header";
import { buildExploitShare, buildTickSeries, formatTickLabel } from "./chartData";

echarts.use([
  BarChart,
  LineChart,
  PieChart,
  AxisPointerComponent,
  AriaComponent,
  GridComponent,
  LegendComponent,
  TooltipComponent,
  CanvasRenderer,
]);

const isDarkMode = true;

export function ChartsPage() {
  const config = useConfig();
  const flagsQuery = useAllFlags();
  const summaryQuery = useStatsSummary();
  const flags = flagsQuery.data!.flags;
  const summary = summaryQuery.data!;
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  async function refreshCharts(): Promise<void> {
    await Promise.all([swrMutate(configKey), flagsQuery.mutate(), summaryQuery.mutate()]);
    setErrorMessage(null);
  }
  const swrError = flagsQuery.error ?? summaryQuery.error;
  const visibleErrorMessage =
    errorMessage ?? (swrError instanceof Error ? swrError.message : null);

  const tickSeries = useMemo(
    () => buildTickSeries(flags, config.server.tick_time),
    [flags, config.server.tick_time],
  );
  const exploitShare = useMemo(() => buildExploitShare(flags), [flags]);
  const topExploits = exploitShare.slice(0, 8);

  const totalPerTickOption = useMemo(() => {
    return {
      tooltip: {
        trigger: "axis" as const,
        axisPointer: { type: "line" as const },
      },
      grid: { top: 30, right: 20, bottom: 45, left: 55 },
      xAxis: {
        type: "category" as const,
        boundaryGap: false,
        data: tickSeries.map((point) => formatTickLabel(point.timestamp)),
        axisLabel: { color: "#A1A1AA", hideOverlap: true },
        axisLine: { lineStyle: { color: "#3F3F46" } },
      },
      yAxis: {
        type: "value" as const,
        minInterval: 1,
        axisLabel: { color: "#A1A1AA" },
        splitLine: { lineStyle: { color: "#3F3F46" } },
      },
      series: [
        {
          name: "Flags",
          type: "line" as const,
          smooth: true,
          symbolSize: 6,
          data: tickSeries.map((point) => point.total),
          lineStyle: { width: 3, color: ChartPalette.semantic("Neutral", isDarkMode) },
          itemStyle: { color: ChartPalette.semantic("Neutral", isDarkMode) },
          areaStyle: { color: "rgba(173, 173, 184, 0.16)" },
        },
      ],
    };
  }, [tickSeries]);

  const statusPerTickOption = useMemo(() => {
    return {
      tooltip: {
        trigger: "axis" as const,
        axisPointer: { type: "line" as const },
      },
      legend: {
        top: 0,
        textStyle: { color: "#D4D4D8" },
      },
      grid: { top: 50, right: 20, bottom: 45, left: 55 },
      xAxis: {
        type: "category" as const,
        boundaryGap: false,
        data: tickSeries.map((point) => formatTickLabel(point.timestamp)),
        axisLabel: { color: "#A1A1AA", hideOverlap: true },
        axisLine: { lineStyle: { color: "#3F3F46" } },
      },
      yAxis: {
        type: "value" as const,
        minInterval: 1,
        axisLabel: { color: "#A1A1AA" },
        splitLine: { lineStyle: { color: "#3F3F46" } },
      },
      series: [
        {
          name: "Accepted",
          type: "line" as const,
          smooth: true,
          data: tickSeries.map((point) => point.accepted),
          itemStyle: { color: ChartPalette.color(0, isDarkMode) },
          lineStyle: { color: ChartPalette.color(0, isDarkMode), width: 2 },
        },
        {
          name: "Denied",
          type: "line" as const,
          smooth: true,
          data: tickSeries.map((point) => point.denied),
          itemStyle: { color: ChartPalette.semantic("Attention", isDarkMode) },
          lineStyle: { color: ChartPalette.semantic("Attention", isDarkMode), width: 2 },
        },
        {
          name: "Error",
          type: "line" as const,
          smooth: true,
          data: tickSeries.map((point) => point.error),
          itemStyle: { color: ChartPalette.semantic("Warning", isDarkMode) },
          lineStyle: { color: ChartPalette.semantic("Warning", isDarkMode), width: 2 },
        },
        {
          name: "Queued",
          type: "line" as const,
          smooth: true,
          data: tickSeries.map((point) => point.queued),
          itemStyle: { color: ChartPalette.semantic("NeutralLight", isDarkMode) },
          lineStyle: { color: ChartPalette.semantic("NeutralLight", isDarkMode), width: 2 },
        },
      ],
    };
  }, [tickSeries]);

  const exploitPieOption = useMemo(() => {
    return {
      tooltip: {
        trigger: "item" as const,
        formatter: "{b}: {c} flags ({d}%)",
      },
      legend: {
        orient: "vertical" as const,
        right: 0,
        top: "center" as const,
        textStyle: { color: "#D4D4D8" },
      },
      series: [
        {
          name: "Exploit Share",
          type: "pie" as const,
          radius: ["45%", "72%"],
          center: ["35%", "50%"],
          avoidLabelOverlap: true,
          label: {
            color: "#F4F4F5",
            formatter: "{d}%",
          },
          data: exploitShare.map((item, index) => ({
            name: item.name,
            value: item.value,
            itemStyle: { color: ChartPalette.color(index, isDarkMode) },
          })),
        },
      ],
    };
  }, [exploitShare]);

  const flagDistributionOption = useMemo(() => {
    const stats = summary.flags_stats ?? [];

    return {
      tooltip: {
        trigger: "axis" as const,
        axisPointer: { type: "shadow" as const },
      },
      grid: { top: 30, right: 20, bottom: 30, left: 60 },
      xAxis: {
        type: "category" as const,
        data: stats.map((item) => `Team ${item.team_id}`),
        axisLabel: { color: "#A1A1AA" },
        axisLine: { lineStyle: { color: "#3F3F46" } },
      },
      yAxis: {
        type: "value" as const,
        axisLabel: { color: "#A1A1AA" },
        splitLine: { lineStyle: { color: "#3F3F46" } },
      },
      series: [
        {
          name: "Accepted",
          type: "bar" as const,
          stack: "total",
          data: stats.map((item) => item.accepted_flags),
          itemStyle: { color: ChartPalette.color(0, isDarkMode) },
        },
        {
          name: "Denied",
          type: "bar" as const,
          stack: "total",
          data: stats.map((item) => item.denied_flags),
          itemStyle: { color: ChartPalette.semantic("Attention", isDarkMode) },
        },
        {
          name: "Error",
          type: "bar" as const,
          stack: "total",
          data: stats.map((item) => item.error_flags),
          itemStyle: { color: ChartPalette.semantic("Warning", isDarkMode) },
        },
        {
          name: "Unsubmitted",
          type: "bar" as const,
          stack: "total",
          data: stats.map((item) => item.unsubmitted_flags),
          itemStyle: { color: ChartPalette.semantic("NeutralLight", isDarkMode) },
        },
      ],
    };
  }, [summary]);

  const totalFlags = flags.length;
  const latestTickCount = tickSeries[tickSeries.length - 1]?.total ?? 0;
  const leadingExploit = exploitShare[0];

  return (
    <div className="space-y-6">
      <PageHeader
        breadcrumbs={
          <Breadcrumbs className="px-3 py-2 text-sm">
            <Breadcrumbs.Link href="/">Operations</Breadcrumbs.Link>
            <Breadcrumbs.Separator />
            <Breadcrumbs.Current>Charts</Breadcrumbs.Current>
          </Breadcrumbs>
        }
        title="Charts"
        description={`Tick-based flag history using ${config.server.tick_time}s buckets, plus exploit distribution.`}
      >
        <Button
          variant="secondary"
          onClick={() => {
            void refreshCharts().catch((error: unknown) => {
              setErrorMessage(error instanceof Error ? error.message : "Chart refresh failed");
            });
          }}
        >
          Reload
        </Button>
      </PageHeader>

      {visibleErrorMessage ? (
        <Banner
          variant="error"
          icon={<WarningCircle weight="fill" />}
          title="Unable to refresh charts"
          description={visibleErrorMessage}
        />
      ) : null}

      <section className="flex flex-wrap gap-4 rounded-2xl border border-kumo-line bg-kumo-base p-4">
        <ChartLegend.LargeItem
          name="History"
          color={ChartPalette.semantic("Neutral", isDarkMode)}
          value={String(totalFlags)}
          unit="flags"
        />
        <ChartLegend.LargeItem
          name="Latest Tick"
          color={ChartPalette.color(0, isDarkMode)}
          value={String(latestTickCount)}
          unit="flags"
        />
        {leadingExploit ? (
          <ChartLegend.LargeItem
            name="Top Exploit"
            color={ChartPalette.color(1, isDarkMode)}
            value={`${leadingExploit.percentage.toFixed(1)}%`}
            unit={leadingExploit.name}
          />
        ) : (
          <ChartLegend.LargeItem
            name="Top Exploit"
            color={ChartPalette.color(1, isDarkMode)}
            value="0%"
          />
        )}
      </section>

      <div className="grid gap-4 xl:grid-cols-2">
        <section className="rounded-2xl border border-kumo-line bg-kumo-base p-4">
          <h2 className="mb-2 text-sm font-medium text-kumo-fg-primary">Flags Per Tick</h2>
          <p className="mb-4 text-sm text-kumo-fg-secondary">
            Total submitted flags grouped by collector tick.
          </p>
          <Chart
            echarts={echarts}
            options={totalPerTickOption}
            isDarkMode={isDarkMode}
            height={320}
          />
        </section>

        <section className="rounded-2xl border border-kumo-line bg-kumo-base p-4">
          <h2 className="mb-2 text-sm font-medium text-kumo-fg-primary">Status Per Tick</h2>
          <p className="mb-4 text-sm text-kumo-fg-secondary">
            Accepted, denied, error, and queued trends across time.
          </p>
          <Chart
            echarts={echarts}
            options={statusPerTickOption}
            isDarkMode={isDarkMode}
            height={320}
          />
        </section>

        <section className="rounded-2xl border border-kumo-line bg-kumo-base p-4">
          <h2 className="mb-2 text-sm font-medium text-kumo-fg-primary">Exploit Share</h2>
          <p className="mb-4 text-sm text-kumo-fg-secondary">
            Percentage of flags generated by each exploit.
          </p>
          <Chart
            echarts={echarts}
            options={exploitPieOption}
            isDarkMode={isDarkMode}
            height={320}
          />
        </section>

        <section className="rounded-2xl border border-kumo-line bg-kumo-base p-4">
          <h2 className="mb-2 text-sm font-medium text-kumo-fg-primary">Flag Status Distribution</h2>
          <p className="mb-4 text-sm text-kumo-fg-secondary">
            Team-by-team accepted, denied, error, and unsubmitted totals.
          </p>
          <Chart
            echarts={echarts}
            options={flagDistributionOption}
            isDarkMode={isDarkMode}
            height={320}
          />
        </section>
      </div>
    </div>
  );
}
