import { use } from "react";
import { z } from "zod";
import { apiFetch, cached } from "./client";

const statsTimestampSchema = z.union([z.string(), z.null()]).transform((value) => {
  if (!value || value.startsWith("0001-01-01")) {
    return null;
  }
  return value;
});

export const statsSummarySchema = z.object({
  buffer_size: z.number().int().nonnegative(),
  total_flags_received: z.number().int().nonnegative(),
  total_flags_flushed: z.number().int().nonnegative(),
  total_flushes: z.number().int().nonnegative(),
  successful_flushes: z.number().int().nonnegative(),
  failed_flushes: z.number().int().nonnegative(),
  last_flush_time: statsTimestampSchema,
  last_successful_flush: statsTimestampSchema,
  efficiency_ratio: z.number().nonnegative(),
  status: z.object({
    is_running: z.boolean(),
  }),
});

export type StatsSummary = z.infer<typeof statsSummarySchema>;

export async function fetchStatsSummary(): Promise<StatsSummary> {
  return apiFetch("/stats", {}, statsSummarySchema);
}

export function readStatsSummary(): Promise<StatsSummary> {
  return cached("stats:summary", fetchStatsSummary);
}

export function useStatsSummary() {
  return use(readStatsSummary());
}
