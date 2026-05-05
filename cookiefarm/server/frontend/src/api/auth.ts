import { z } from "zod";
import { apiFetch } from "./client";

export const loginPayloadSchema = z.object({
  username: z.string().trim().optional(),
  password: z.string().trim().min(1, "Password is required."),
});

const emptyResponseSchema = z.object({}).passthrough();
const authSessionSchema = z.object({
  username: z.string().trim().min(1).default("cookieguest"),
});

export type LoginPayload = z.infer<typeof loginPayloadSchema>;
export type AuthSession = z.infer<typeof authSessionSchema>;

export const authVerifyKey = "/auth/verify";

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

export async function verifyAuth(): Promise<AuthSession | null> {
  try {
    return await apiFetch(authVerifyKey, {}, authSessionSchema);
  } catch {
    return null;
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
