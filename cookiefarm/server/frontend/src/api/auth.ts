import { z } from "zod";
import { apiFetch } from "./client";

export const loginPayloadSchema = z.object({
  username: z.string().trim().optional(),
  password: z.string().trim().min(1, "Password is required."),
});

const emptyResponseSchema = z.object({}).passthrough();

export type LoginPayload = z.infer<typeof loginPayloadSchema>;

export async function login(payload: LoginPayload): Promise<void> {
  await apiFetch(
    "/auth/login",
    {
      method: "POST",
      body: JSON.stringify(payload),
    },
    emptyResponseSchema,
  );
}

export async function verifyAuth(): Promise<boolean> {
  try {
    await apiFetch("/auth/verify", {}, emptyResponseSchema);
    return true;
  } catch {
    return false;
  }
}

export async function logout(): Promise<void> {
  await apiFetch(
    "/auth/logout",
    {
      method: "POST",
    },
    emptyResponseSchema,
  );
}
