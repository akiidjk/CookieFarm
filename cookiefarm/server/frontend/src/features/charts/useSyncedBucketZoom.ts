import { useCallback, useEffect, useMemo, useState } from "react";
import type { ChartEvents, KumoChartOption } from "@cloudflare/kumo/components/chart";

type DataZoomEvent = Parameters<NonNullable<ChartEvents["datazoom"]>>[0] & {
  batch?: Array<{
    start?: number;
    end?: number;
    startValue?: number;
    endValue?: number;
  }>;
  start?: number;
  end?: number;
  startValue?: number;
  endValue?: number;
};

function clamp(value: number, min: number, max: number): number {
  return Math.min(Math.max(value, min), max);
}

export function useSyncedBucketZoom(bucketCount: number): {
  dataZoom: NonNullable<KumoChartOption["dataZoom"]>;
  onDataZoom: ChartEvents["datazoom"];
  windowSize: number;
  setWindowSize: (value: number) => void;
} {
  const initialEnd = Math.max(bucketCount - 1, 0);
  const initialStart = Math.max(initialEnd - 23, 0);
  const [bucketRange, setBucketRange] = useState(() => ({
    start: initialStart,
    end: initialEnd,
  }));
  const maxBucketIndex = Math.max(bucketCount - 1, 0);
  const bucketWindowSize = Math.max(bucketRange.end - bucketRange.start + 1, 1);

  useEffect(() => {
    setBucketRange((current) => {
      const nextStart = clamp(current.start, 0, maxBucketIndex);
      const nextEnd = clamp(Math.max(current.end, nextStart), nextStart, maxBucketIndex);
      return { start: nextStart, end: nextEnd };
    });
  }, [maxBucketIndex]);

  const dataZoom = useMemo(
    () => [
      {
        type: "inside" as const,
        startValue: bucketRange.start,
        endValue: bucketRange.end,
        filterMode: "filter" as const,
      },
      {
        type: "slider" as const,
        startValue: bucketRange.start,
        endValue: bucketRange.end,
        filterMode: "filter" as const,
        height: 18,
        bottom: 10,
        borderColor: "#3F3F46",
        fillerColor: "rgba(173, 173, 184, 0.18)",
        handleStyle: { color: "#A1A1AA", borderColor: "#D4D4D8" },
        moveHandleStyle: { color: "#71717A" },
        textStyle: { color: "#A1A1AA" },
        brushSelect: false,
        showDetail: true,
      },
    ],
    [bucketRange.end, bucketRange.start],
  );

  const onDataZoom = useCallback(
    (params: DataZoomEvent) => {
      const zoom = params.batch?.[0] ?? params;
      let nextStart = bucketRange.start;
      let nextEnd = bucketRange.end;

      if (typeof zoom.startValue === "number") {
        nextStart = Math.round(zoom.startValue);
      } else if (typeof zoom.start === "number") {
        nextStart = Math.round((zoom.start / 100) * maxBucketIndex);
      }

      if (typeof zoom.endValue === "number") {
        nextEnd = Math.round(zoom.endValue);
      } else if (typeof zoom.end === "number") {
        nextEnd = Math.round((zoom.end / 100) * maxBucketIndex);
      }

      nextStart = clamp(nextStart, 0, maxBucketIndex);
      nextEnd = clamp(nextEnd, nextStart, maxBucketIndex);
      setBucketRange({ start: nextStart, end: nextEnd });
    },
    [bucketRange.end, bucketRange.start, maxBucketIndex],
  );

  const setWindowSize = useCallback(
    (value: number) => {
      const nextWindowSize = clamp(Math.round(value), 1, Math.max(bucketCount, 1));
      setBucketRange((current) => {
        const end = clamp(current.end, 0, maxBucketIndex);
        const start = clamp(end - nextWindowSize + 1, 0, end);
        return { start, end };
      });
    },
    [bucketCount, maxBucketIndex],
  );

  return { dataZoom, onDataZoom, windowSize: bucketWindowSize, setWindowSize };
}
