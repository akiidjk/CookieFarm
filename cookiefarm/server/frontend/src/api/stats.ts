import useSWR, { type SWRConfiguration } from "swr";
import { z } from "zod";
import { apiFetch } from "./client";

const sqlFloatSchema = z.object({
  Float64: z.number(),
  Valid: z.boolean(),
}).transform((val) => val.Valid ? val.Float64 : 0);

export const teamStatsSchema = z.object({
  team_id: z.number().int(),
  total_flags: z.number().int(),
  accepted_flags: sqlFloatSchema,
  denied_flags: sqlFloatSchema,
  unsubmitted_flags: sqlFloatSchema,
  error_flags: sqlFloatSchema,
  not_valid_flags: sqlFloatSchema,
});

export type TeamStats = z.infer<typeof teamStatsSchema>;

export const statsSummarySchema = z.object({
  flags_stats: z.array(teamStatsSchema),
});

export type StatsSummary = z.infer<typeof statsSummarySchema>;

export const chartTickPointSchema = z.object({
  timestamp: z.number().int(),
  total: z.number().int(),
  queued: z.number().int(),
  accepted: z.number().int(),
  denied: z.number().int(),
  error: z.number().int(),
  invalid: z.number().int(),
});

export const exploitShareSchema = z.object({
  name: z.string(),
  value: z.number().int(),
  percentage: z.number(),
});

export const exploitTickPointSchema = z.object({
  timestamp: z.number().int(),
  value: z.number().int(),
});

export const exploitTickSeriesSchema = z.object({
  name: z.string(),
  total: z.number().int(),
  data: z.array(exploitTickPointSchema),
});

export const exploitStatusPercentageSchema = z.object({
  name: z.string(),
  total: z.number().int(),
  queued: z.number(),
  accepted: z.number(),
  denied: z.number(),
  error: z.number(),
  invalid: z.number(),
  queued_count: z.number().int(),
  accepted_count: z.number().int(),
  denied_count: z.number().int(),
  error_count: z.number().int(),
  invalid_count: z.number().int(),
});

export const chartStatsSchema = z.object({
  tick_series: z.array(chartTickPointSchema),
  exploit_share: z.array(exploitShareSchema),
  exploit_tick_series: z.array(exploitTickSeriesSchema),
  exploit_status_percentage: z.array(exploitStatusPercentageSchema),
  total_flags: z.number().int(),
});

export type ChartStats = z.infer<typeof chartStatsSchema>;

export const statsSummaryKey = "/stats";

export async function fetchStatsSummary(): Promise<StatsSummary> {
  return apiFetch(statsSummaryKey, {}, statsSummarySchema);
}

export function useStatsSummary(options: SWRConfiguration = {}) {
  return useSWR(statsSummaryKey, fetchStatsSummary, {
    suspense: true,
    ...options,
  });
}

export function chartStatsKey(tickSeconds: number) {
  return ["/stats/charts", tickSeconds] as const;
}

export async function fetchChartStats(tickSeconds: number): Promise<ChartStats> {
  const query = new URLSearchParams({
    tick_seconds: String(tickSeconds),
  });

  return apiFetch(`/stats/charts?${query.toString()}`, {}, chartStatsSchema);
}

export function useChartStats(tickSeconds: number, options: SWRConfiguration = {}) {
  return useSWR(chartStatsKey(tickSeconds), () => fetchChartStats(tickSeconds), {
    suspense: true,
    ...options,
  });
}
