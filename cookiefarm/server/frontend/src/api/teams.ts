import { use } from "react";
import { z } from "zod";
import { apiFetch, cached } from "./client";

export const teamSchema = z.object({
  ip: z.string(),
  name: z.string(),
  active: z.boolean(),
});

export type Team = z.infer<typeof teamSchema>;

const teamsSchema = z.array(teamSchema);

export async function fetchTeams(): Promise<Team[]> {
  return apiFetch("/teams", {}, teamsSchema);
}

export function readTeams(): Promise<Team[]> {
  return cached("teams:list", fetchTeams);
}

export function useTeams() {
  return use(readTeams());
}
