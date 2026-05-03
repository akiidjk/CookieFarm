import { Input } from "@cloudflare/kumo/components/input";

export function ChartWindowControls(props: {
  value: number;
  max: number;
  onChange: (value: number) => void;
}) {
  const max = Math.max(props.max, 1);

  return (
    <section className="flex flex-wrap items-end gap-3 rounded-2xl border border-kumo-line bg-kumo-base p-4">
      <div className="min-w-40">
        <label
          htmlFor="chart-window-size"
          className="mb-1 block text-xs font-medium text-kumo-fg-secondary"
        >
          Window
        </label>
        <Input
          id="chart-window-size"
          aria-label="Visible bucket window size"
          type="number"
          min={1}
          max={max}
          step={1}
          size="sm"
          value={props.value}
          onChange={(event) => {
            props.onChange(event.currentTarget.valueAsNumber);
          }}
        />
      </div>
      <span className="pb-2 text-sm text-kumo-fg-secondary">
        Resize from the bottom chart window, or enter 1-{max} buckets.
      </span>
    </section>
  );
}
