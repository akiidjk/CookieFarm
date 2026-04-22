import { useMemo } from "react";
import * as echarts from "echarts/core";
import { BarChart } from "echarts/charts";
import {
  AriaComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
} from "echarts/components";
import { CanvasRenderer } from "echarts/renderers";
import { Chart, ChartPalette, ChartLegend } from "@cloudflare/kumo/components/chart";
import type { StatsSummary } from "@/api/stats";

echarts.use([
  BarChart,
  AriaComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
  CanvasRenderer,
]);

export function StatsBar(props: { summary: StatsSummary }) {
  const isDarkMode = true;

  const stats = props.summary?.flags_stats || [];

  const totalFlagsOption = useMemo(() => {
    const categories = stats.map((s) => `Team ${s.team_id}`);
    const data = stats.map((s) => s.total_flags);

    return {
      tooltip: {
        trigger: "axis" as const,
        axisPointer: { type: "shadow" as const },
      },
      grid: { top: 30, right: 20, bottom: 30, left: 60 },
      xAxis: {
        type: "category" as const,
        data: categories,
        axisLabel: { color: isDarkMode ? "#A1A1AA" : "#52525B" },
        axisLine: { lineStyle: { color: isDarkMode ? "#3F3F46" : "#E4E4E7" } },
      },
      yAxis: {
        type: "value" as const,
        axisLabel: { color: isDarkMode ? "#A1A1AA" : "#52525B" },
        splitLine: { lineStyle: { color: isDarkMode ? "#3F3F46" : "#E4E4E7" } },
      },
      series: [
        {
          name: "Total Flags",
          data: data,
          type: "bar" as const,
          itemStyle: { color: ChartPalette.semantic("Neutral", isDarkMode) },
        },
      ],
    };
  }, [stats, isDarkMode]);

  const statusOption = useMemo(() => {
    const categories = stats.map((s) => `Team ${s.team_id}`);

    return {
      tooltip: {
        trigger: "axis" as const,
        axisPointer: { type: "shadow" as const },
      },
      grid: { top: 30, right: 20, bottom: 30, left: 60 },
      xAxis: {
        type: "category" as const,
        data: categories,
        axisLabel: { color: isDarkMode ? "#A1A1AA" : "#52525B" },
        axisLine: { lineStyle: { color: isDarkMode ? "#3F3F46" : "#E4E4E7" } },
      },
      yAxis: {
        type: "value" as const,
        axisLabel: { color: isDarkMode ? "#A1A1AA" : "#52525B" },
        splitLine: { lineStyle: { color: isDarkMode ? "#3F3F46" : "#E4E4E7" } },
      },
      series: [
        {
          name: "Accepted",
          type: "bar" as const,
          stack: "total",
          data: stats.map((s) => s.accepted_flags),
          itemStyle: { color: ChartPalette.color(0, isDarkMode) },
        },
        {
          name: "Denied",
          type: "bar" as const,
          stack: "total",
          data: stats.map((s) => s.denied_flags),
          itemStyle: { color: ChartPalette.semantic("Attention", isDarkMode) },
        },
        {
          name: "Error",
          type: "bar" as const,
          stack: "total",
          data: stats.map((s) => s.error_flags),
          itemStyle: { color: ChartPalette.semantic("Warning", isDarkMode) },
        },
        {
          name: "Unsubmitted",
          type: "bar" as const,
          stack: "total",
          data: stats.map((s) => s.unsubmitted_flags),
          itemStyle: { color: ChartPalette.semantic("NeutralLight", isDarkMode) },
        },
      ],
    };
  }, [stats, isDarkMode]);

  const totalFlagsCount = stats.reduce((acc, s) => acc + s.total_flags, 0);
  const totalAcceptedCount = stats.reduce((acc, s) => acc + s.accepted_flags, 0);
  const totalDeniedCount = stats.reduce((acc, s) => acc + s.denied_flags, 0);
  const totalErrorCount = stats.reduce((acc, s) => acc + s.error_flags, 0);
  const totalUnsubmittedCount = stats.reduce((acc, s) => acc + s.unsubmitted_flags, 0);

  return (
    <div className="grid gap-4 lg:grid-cols-2">
      <div className="rounded-2xl border border-kumo-line bg-kumo-base p-4 flex flex-col">
        <h3 className="mb-4 text-sm font-medium text-kumo-fg-primary">Total Flags per Team</h3>
        <div className="flex-1">
          <Chart
            echarts={echarts}
            options={totalFlagsOption}
            isDarkMode={isDarkMode}
            height={250}
          />
        </div>
        <div className="mt-4 flex flex-wrap gap-4">
          <ChartLegend.SmallItem
            name="Total Flags"
            color={ChartPalette.semantic("Neutral", isDarkMode)}
            value={String(totalFlagsCount)}
          />
        </div>
      </div>

      <div className="rounded-2xl border border-kumo-line bg-kumo-base p-4 flex flex-col">
        <h3 className="mb-4 text-sm font-medium text-kumo-fg-primary">Flag Status Distribution</h3>
        <div className="flex-1">
          <Chart
            echarts={echarts}
            options={statusOption}
            isDarkMode={isDarkMode}
            height={250}
          />
        </div>
        <div className="mt-4 flex flex-wrap gap-4">
          <ChartLegend.SmallItem
            name="Accepted"
            color={ChartPalette.color(0, isDarkMode)}
            value={String(totalAcceptedCount)}
          />
          <ChartLegend.SmallItem
            name="Denied"
            color={ChartPalette.semantic("Attention", isDarkMode)}
            value={String(totalDeniedCount)}
          />
          <ChartLegend.SmallItem
            name="Error"
            color={ChartPalette.semantic("Warning", isDarkMode)}
            value={String(totalErrorCount)}
          />
          <ChartLegend.SmallItem
            name="Unsubmitted"
            color={ChartPalette.semantic("NeutralLight", isDarkMode)}
            value={String(totalUnsubmittedCount)}
          />
        </div>
      </div>
    </div>
  );
}
