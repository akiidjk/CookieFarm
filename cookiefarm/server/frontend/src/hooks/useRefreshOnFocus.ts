import { useEffect, useEffectEvent, useRef } from "react";

type RefreshOnFocusOptions = {
  enabled?: boolean;
  immediate?: boolean;
  dedupeMs?: number;
};

export function useRefreshOnFocus(
  callback: () => void,
  options: RefreshOnFocusOptions = {},
) {
  const onRefresh = useEffectEvent(callback);
  const enabled = options.enabled ?? true;
  const dedupeMs = options.dedupeMs ?? 500;
  const lastRunAtRef = useRef(0);

  useEffect(() => {
    if (!enabled) {
      return;
    }

    function runRefresh() {
      const now = Date.now();
      if (now - lastRunAtRef.current < dedupeMs) {
        return;
      }

      lastRunAtRef.current = now;
      onRefresh();
    }

    function handleVisibilityChange() {
      if (document.visibilityState === "visible") {
        runRefresh();
      }
    }

    if (options.immediate) {
      runRefresh();
    }

    window.addEventListener("focus", runRefresh);
    document.addEventListener("visibilitychange", handleVisibilityChange);

    return () => {
      window.removeEventListener("focus", runRefresh);
      document.removeEventListener("visibilitychange", handleVisibilityChange);
    };
  }, [dedupeMs, enabled, onRefresh, options.immediate]);
}
