import { Meter } from "@cloudflare/kumo/components/meter";
import type { StatsSummary } from "@/api/stats";

function formatTimestamp(timestamp: string | null): string {
  if (!timestamp) {
    return "No submissions yet";
  }

  return new Date(timestamp).toLocaleString([], {
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

function freshnessValue(timestamp: string | null): number {
  if (!timestamp) {
    return 0;
  }

  const elapsedMs = Date.now() - new Date(timestamp).getTime();
  const freshness = 100 - Math.min((elapsedMs / (15 * 60_000)) * 100, 100);
  return Math.max(0, Math.round(freshness));
}

export function StatsBar(props: { summary: StatsSummary }) {
  const ratioValue =
    props.summary.total_flushes === 0
      ? 0
      : (props.summary.successful_flushes / props.summary.total_flushes) * 100;
  const collectorValue = props.summary.status.is_running
    ? 100
    : freshnessValue(props.summary.last_successful_flush);

  return (
    <div className="grid gap-4 lg:grid-cols-4">
      <div className="rounded-2xl border border-kumo-line bg-kumo-base p-4">
        <Meter
          label="Flags Received"
          value={100}
          customValue={props.summary.total_flags_received.toLocaleString()}
          indicatorClassName="bg-kumo-brand"
        />
      </div>

      <div className="rounded-2xl border border-kumo-line bg-kumo-base p-4">
        <Meter
          label="Flags Flushed"
          value={100}
          customValue={props.summary.total_flags_flushed.toLocaleString()}
          indicatorClassName="bg-kumo-info"
        />
      </div>

      <div className="rounded-2xl border border-kumo-line bg-kumo-base p-4">
        <Meter
          label="Flush Success"
          value={Number.isFinite(ratioValue) ? ratioValue : 0}
          customValue={`${props.summary.successful_flushes} / ${props.summary.total_flushes}`}
          indicatorClassName="bg-kumo-success"
        />
      </div>

      <div className="rounded-2xl border border-kumo-line bg-kumo-base p-4">
        <Meter
          label="Collector"
          value={collectorValue}
          customValue={
            props.summary.status.is_running
              ? `Running · last success ${formatTimestamp(props.summary.last_successful_flush)}`
              : `Stopped · buffer ${props.summary.buffer_size}`
          }
          indicatorClassName="bg-kumo-warning"
        />
      </div>
    </div>
  );
}
