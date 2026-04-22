import { useRef } from "react";
import { ClipboardText } from "@cloudflare/kumo/components/clipboard-text";
import { Empty } from "@cloudflare/kumo/components/empty";
import { Table } from "@cloudflare/kumo/components/table";
import { useVirtualizer } from "@tanstack/react-virtual";
import type { Flag } from "@/api/flags";
import { FlagStatusBadge } from "@/components/FlagStatusBadge";

function formatTimestamp(timestamp: number): string {
  if (timestamp === 0) {
    return "-";
  }

  return new Date(timestamp * 1000).toLocaleString([], {
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

function formatDuration(flag: Flag): string {
  if (flag.response_time === 0 || flag.response_time < flag.submit_time) {
    return "Pending";
  }

  return `${flag.response_time - flag.submit_time}s`;
}

export function FlagTable(props: {
  rows: Flag[];
  emptyTitle?: string;
  emptyDescription?: string;
}) {
  const parentRef = useRef<HTMLDivElement | null>(null);
  const shouldVirtualize = props.rows.length > 500;
  const virtualizer = useVirtualizer({
    count: props.rows.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 54,
    overscan: 12,
  });

  if (props.rows.length === 0) {
    return (
      <section className="rounded-2xl border border-kumo-line bg-kumo-base">
        <Empty
          size="sm"
          title={props.emptyTitle ?? "No flags match these filters"}
          description={
            props.emptyDescription ??
            "Adjust the current filters to widen the view or wait for new rows."
          }
        />
      </section>
    );
  }

  const virtualRows = shouldVirtualize ? virtualizer.getVirtualItems() : [];
  const paddingTop = shouldVirtualize && virtualRows[0] ? virtualRows[0].start : 0;
  const paddingBottom =
    shouldVirtualize && virtualRows.length > 0
      ? virtualizer.getTotalSize() - virtualRows[virtualRows.length - 1]!.end
      : 0;
  const visibleRows: Flag[] = shouldVirtualize
    ? virtualRows.flatMap((item) => {
      const row = props.rows[item.index];
      return row ? [row] : [];
    })
    : props.rows;

  return (
    <section className="overflow-hidden rounded-2xl border border-kumo-line bg-kumo-base">
      <div ref={parentRef} className="max-h-[70vh] overflow-auto">
        <Table layout="fixed">
          {/* Define column widths here. Adjust the className/style on each <col/> to change sizes. */}
          <colgroup>
            <col className="w-80" />
            <col />
            <col className="w-40" />
            <col className="w-30" />
            <col className="w-40" />
            <col className="w-40" />
            <col className="w-24" />
            <col className="w-24" />
          </colgroup>
          <Table.Header sticky>
            <Table.Row>
              <Table.Head>Flag</Table.Head>
              <Table.Head>Message</Table.Head>
              <Table.Head>Service</Table.Head>
              <Table.Head>Status</Table.Head>
              <Table.Head>Submit</Table.Head>
              <Table.Head>Response</Table.Head>
              <Table.Head>Duration</Table.Head>
              <Table.Head>Team</Table.Head>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {paddingTop > 0 ? (
              <Table.Row aria-hidden="true">
                <Table.Cell colSpan={8} style={{ height: paddingTop }} />
              </Table.Row>
            ) : null}

            {visibleRows.map((flag) => (
              <Table.Row key={`${flag.flag_code}-${flag.submit_time}`}>
                <Table.Cell>
                  <ClipboardText size="sm" text={flag.flag_code} />
                </Table.Cell>
                <Table.Cell > {flag.msg || "-"}</Table.Cell>
                <Table.Cell>{`${flag.service_name}:${flag.port_service}`}</Table.Cell>
                <Table.Cell>
                  <FlagStatusBadge status={flag.status} />
                </Table.Cell>
                <Table.Cell>{formatTimestamp(flag.submit_time)}</Table.Cell>
                <Table.Cell>{formatTimestamp(flag.response_time)}</Table.Cell>
                <Table.Cell>{formatDuration(flag)}</Table.Cell>
                <Table.Cell>{flag.team_id}</Table.Cell>
              </Table.Row>
            ))}

            {paddingBottom > 0 ? (
              <Table.Row aria-hidden="true">
                <Table.Cell colSpan={8} style={{ height: paddingBottom }} />
              </Table.Row>
            ) : null}
          </Table.Body>
        </Table>
      </div>
    </section>
  );
}
