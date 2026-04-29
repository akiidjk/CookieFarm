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
