import { useEffect, useEffectEvent } from "react";

type IntervalOptions = {
  enabled?: boolean;
  immediate?: boolean;
};

export function useInterval(
  callback: () => void,
  delayMs: number | null,
  options: IntervalOptions = {},
) {
  const onTick = useEffectEvent(callback);
  const enabled = options.enabled ?? true;

  useEffect(() => {
    if (!enabled || delayMs === null) {
      return;
    }

    if (options.immediate) {
      onTick();
    }

    const timer = window.setInterval(() => {
      onTick();
    }, delayMs);

    return () => {
      window.clearInterval(timer);
    };
  }, [delayMs, enabled, onTick, options.immediate]);
}
