import { useState } from "react";
import { mutate as swrMutate } from "swr";
import { Banner } from "@cloudflare/kumo/components/banner";
import { Breadcrumbs } from "@cloudflare/kumo/components/breadcrumbs";
import { Button } from "@cloudflare/kumo/components/button";
import { WarningCircleIcon } from "@phosphor-icons/react";
import { configKey, useConfig } from "@/api/config";
import { useChartStats, useStatsSummary } from "@/api/stats";
import { PageHeader } from "@/components/kumo/page-header/page-header";
import { ChartsSummary } from "./ChartsSummary";
import { ExploitCharts } from "./ExploitCharts";
import { TeamStatusChart } from "./TeamStatusChart";
import { TimeCharts } from "./TimeCharts";
import { useSyncedBucketZoom } from "./useSyncedBucketZoom";

export function ChartsPage() {
  const config = useConfig();
  const chartStatsQuery = useChartStats(config.server.tick_time);
  const summaryQuery = useStatsSummary();
  const chartStats = chartStatsQuery.data!;
  const summary = summaryQuery.data!;
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const { dataZoom, onDataZoom } = useSyncedBucketZoom(
    chartStats.tick_series.length,
  );

  async function refreshCharts(): Promise<void> {
    await Promise.all([swrMutate(configKey), chartStatsQuery.mutate(), summaryQuery.mutate()]);
    setErrorMessage(null);
  }

  const swrError = chartStatsQuery.error ?? summaryQuery.error;
  const visibleErrorMessage =
    errorMessage ?? (swrError instanceof Error ? swrError.message : null);

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
          icon={<WarningCircleIcon weight="fill" />}
          title="Unable to refresh charts"
          description={visibleErrorMessage}
        />
      ) : null}

      <ChartsSummary chartStats={chartStats} />
      <div className="grid gap-4 xl:grid-cols-2">
        <TimeCharts chartStats={chartStats} dataZoom={dataZoom} onDataZoom={onDataZoom} />
        <ExploitCharts chartStats={chartStats} />
        <TeamStatusChart summary={summary} />
      </div>
    </div>
  );
}
