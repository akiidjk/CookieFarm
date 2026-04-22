import { useMemo } from "react";
import * as echarts from "echarts/core";
import { BarChart, PieChart } from "echarts/charts";
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
  PieChart,
  AriaComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
  CanvasRenderer,
]);

export function StatsBar(props: { summary: StatsSummary }) {

  const flagsOption = useMemo(() => {
    return {
      tooltip: {
        trigger: "axis" as const,
        axisPointer: { type: "shadow" as const },
      },
      grid: { top: 30, right: 20, bottom: 30, left: 60 },
      xAxis: {
        type: "category" as const,
        data: ["Received", "Flushed"],
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
          data: [
            {
              value: props.summary.total_flags_received,
              itemStyle: { color: ChartPalette.semantic("Neutral", true) },
            },
            {
              value: props.summary.total_flags_flushed,
              itemStyle: { color: ChartPalette.semantic("NeutralLight", true) },
            },
          ],
          type: "bar" as const,
          barWidth: "40%",
        },
      ],
    };
  }, [props.summary]);

  const flushesOption = useMemo(() => {
    const failed = Math.max(0, props.summary.total_flushes - props.summary.successful_flushes);
    return {
      tooltip: { trigger: "item" as const },
      series: [
        {
          type: "pie" as const,
          radius: ["40%", "70%"],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 4,
            borderColor: "#18181B",
            borderWidth: 2,
          },
          label: { show: false },
          data: [
            {
              value: props.summary.successful_flushes,
              name: "Successful",
              itemStyle: { color: ChartPalette.color(0, true) },
            },
            {
              value: failed,
              name: "Failed",
              itemStyle: { color: ChartPalette.semantic("Attention", true) },
            },
          ],
        },
      ],
    };
  }, [props.summary]);

  return (
    <div className="grid gap-4 lg:grid-cols-2">
      <div className="rounded-2xl border border-kumo-line bg-kumo-base p-4 flex flex-col">
        <h3 className="mb-4 text-sm font-medium text-kumo-fg-primary">Flags Volume</h3>
        <div className="flex-1">
          <Chart
            echarts={echarts}
            options={flagsOption}
            isDarkMode={true}
            height={250}
          />
        </div>
        <div className="mt-4 flex flex-wrap gap-4">
          <ChartLegend.SmallItem
            name="Received"
            color={ChartPalette.semantic("Neutral", true)}
            value={String(props.summary.total_flags_received)}
          />
          <ChartLegend.SmallItem
            name="Flushed"
            color={ChartPalette.semantic("NeutralLight", true)}
            value={String(props.summary.total_flags_flushed)}
          />
        </div>
      </div>

      <div className="rounded-2xl border border-kumo-line bg-kumo-base p-4 flex flex-col">
        <h3 className="mb-4 text-sm font-medium text-kumo-fg-primary">Flush Success Rate</h3>
        <div className="flex-1">
          <Chart
            echarts={echarts}
            options={flushesOption}
            isDarkMode={true}
            height={250}
          />
        </div>
        <div className="mt-4 flex flex-wrap gap-4">
          <ChartLegend.SmallItem
            name="Successful"
            color={ChartPalette.color(0, true)}
            value={String(props.summary.successful_flushes)}
          />
          <ChartLegend.SmallItem
            name="Failed"
            color={ChartPalette.semantic("Attention", true)}
            value={String(Math.max(0, props.summary.total_flushes - props.summary.successful_flushes))}
          />
        </div>
      </div>
    </div>
  );
}
