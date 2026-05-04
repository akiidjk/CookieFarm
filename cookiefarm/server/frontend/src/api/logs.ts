import useSWR, { type SWRConfiguration } from "swr";
import { z } from "zod";
import { apiFetch } from "./client";

export const logLevelSchema = z.enum(["DEBUG", "INFO", "WARN", "ERROR"]);

export const logEntrySchema = z.object({
  id: z.string(),
  level: logLevelSchema,
  timestamp: z.string(),
  message: z.string(),
  ansi: z.string(),
});

export type LogLevel = z.infer<typeof logLevelSchema>;
export type LogEntry = z.infer<typeof logEntrySchema>;

const logsResponseSchema = z.object({
  items: z.array(logEntrySchema),
});

export function logsKey(limit = 200) {
  return ["/logs", limit] as const;
}

export async function fetchRecentLogs(limit = 200): Promise<LogEntry[]> {
  const response = await apiFetch(`/logs?limit=${limit}`, {}, logsResponseSchema);
  return response.items;
}

export function useRecentLogs(limit = 200, options: SWRConfiguration = {}) {
  return useSWR(logsKey(limit), () => fetchRecentLogs(limit), {
    suspense: true,
    ...options,
  });
}
