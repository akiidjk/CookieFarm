import useSWR from "swr";
import { z } from "zod";
import { apiFetch } from "./client";

export const teamSchema = z.object({
  ip: z.string(),
  name: z.string(),
  active: z.boolean(),
});

export type Team = z.infer<typeof teamSchema>;

const teamsSchema = z.array(teamSchema);

export const teamsKey = "/teams";

export async function fetchTeams(): Promise<Team[]> {
  return apiFetch(teamsKey, {}, teamsSchema);
}

export function useTeams() {
  return useSWR(teamsKey, fetchTeams, { suspense: true });
}
