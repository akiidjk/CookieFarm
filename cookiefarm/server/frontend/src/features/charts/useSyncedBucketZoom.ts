import { useCallback, useEffect, useMemo, useRef, useState } from "react";
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

function rangesEqual(
  left: { start: number; end: number },
  right: { start: number; end: number },
): boolean {
  return left.start === right.start && left.end === right.end;
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
  const bucketRangeRef = useRef(bucketRange);
  const pendingRangeRef = useRef(bucketRange);
  const frameRef = useRef<number | null>(null);
  const maxBucketIndex = Math.max(bucketCount - 1, 0);
  const bucketWindowSize = Math.max(bucketRange.end - bucketRange.start + 1, 1);

  useEffect(() => {
    bucketRangeRef.current = bucketRange;
    pendingRangeRef.current = bucketRange;
  }, [bucketRange]);

  useEffect(() => {
    return () => {
      if (frameRef.current !== null) {
        cancelAnimationFrame(frameRef.current);
      }
    };
  }, []);

  useEffect(() => {
    setBucketRange((current) => {
      const nextStart = clamp(current.start, 0, maxBucketIndex);
      const nextEnd = clamp(Math.max(current.end, nextStart), nextStart, maxBucketIndex);
      const nextRange = { start: nextStart, end: nextEnd };
      return rangesEqual(current, nextRange) ? current : nextRange;
    });
  }, [maxBucketIndex]);

  const scheduleRangeUpdate = useCallback((nextRange: { start: number; end: number }) => {
    pendingRangeRef.current = nextRange;

    if (frameRef.current !== null) {
      return;
    }

    frameRef.current = requestAnimationFrame(() => {
      frameRef.current = null;
      const pendingRange = pendingRangeRef.current;

      setBucketRange((current) => {
        if (rangesEqual(current, pendingRange)) {
          return current;
        }

        return pendingRange;
      });
    });
  }, []);

  const dataZoom = useMemo(
    () => [
      {
        type: "inside" as const,
        startValue: bucketRange.start,
        endValue: bucketRange.end,
        filterMode: "filter" as const,
        throttle: 32,
      },
      {
        type: "slider" as const,
        startValue: bucketRange.start,
        endValue: bucketRange.end,
        filterMode: "filter" as const,
        throttle: 32,
        height: 18,
        bottom: 10,
        borderColor: "#3F3F46",
        fillerColor: "rgba(173, 173, 184, 0.18)",
        handleStyle: { color: "#A1A1AA", borderColor: "#D4D4D8" },
        moveHandleStyle: { color: "#71717A" },
        textStyle: { color: "#A1A1AA" },
        brushSelect: false,
        showDetail: true,
        realtime: true,
      },
    ],
    [bucketRange.end, bucketRange.start],
  );

  const onDataZoom = useCallback(
    (params: DataZoomEvent) => {
      const zoom = params.batch?.[0] ?? params;
      let nextStart = bucketRangeRef.current.start;
      let nextEnd = bucketRangeRef.current.end;

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

      const nextRange = { start: nextStart, end: nextEnd };
      if (rangesEqual(bucketRangeRef.current, nextRange)) {
        return;
      }

      scheduleRangeUpdate(nextRange);
    },
    [maxBucketIndex, scheduleRangeUpdate],
  );

  const setWindowSize = useCallback(
    (value: number) => {
      const nextWindowSize = clamp(Math.round(value), 1, Math.max(bucketCount, 1));
      setBucketRange((current) => {
        const end = clamp(current.end, 0, maxBucketIndex);
        const start = clamp(end - nextWindowSize + 1, 0, end);
        const nextRange = { start, end };
        return rangesEqual(current, nextRange) ? current : nextRange;
      });
    },
    [bucketCount, maxBucketIndex],
  );

  return { dataZoom, onDataZoom, windowSize: bucketWindowSize, setWindowSize };
}
