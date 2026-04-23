import { use } from "react";
import { z } from "zod";
import { apiFetch, cached } from "./client";

export const serviceSchema = z.object({
  name: z.string(),
  port: z.number().int().min(1).max(65535),
  active: z.boolean(),
});

export type Service = z.infer<typeof serviceSchema>;

const servicesSchema = z.array(serviceSchema);

export async function fetchServices(): Promise<Service[]> {
  return apiFetch("/services", {}, servicesSchema);
}

export function readServices(): Promise<Service[]> {
  return cached("services:list", fetchServices);
}

export function useServices() {
  return use(readServices());
}
