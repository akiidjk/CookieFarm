import useSWR from "swr";
import { z } from "zod";
import { apiFetch } from "./client";

export const serviceSchema = z.object({
  name: z.string(),
  port: z.number().int().min(1).max(65535),
  active: z.boolean(),
});

export type Service = z.infer<typeof serviceSchema>;

const servicesSchema = z.array(serviceSchema);

export const servicesKey = "/services";

export async function fetchServices(): Promise<Service[]> {
  return apiFetch(servicesKey, {}, servicesSchema);
}

export function useServices() {
  return useSWR(servicesKey, fetchServices, { suspense: true });
}
