import { use } from "react";
import { z } from "zod";
import { apiFetch, cached } from "./client";

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

export async function fetchStatsSummary(): Promise<StatsSummary> {
  return apiFetch("/stats", {}, statsSummarySchema);
}

export function readStatsSummary(): Promise<StatsSummary> {
  return cached("stats:summary", fetchStatsSummary);
}

export function useStatsSummary() {
  return use(readStatsSummary());
}
