import { z, type ZodType } from "zod";

const cachedRequests = new Map<string, Promise<unknown>>();

export const API_BASE = (import.meta.env.VITE_API_BASE_URL ?? "/api/v1").trim();
export const shouldUseApiMocks =
  import.meta.env.DEV && import.meta.env.VITE_USE_API_MOCKS === "true";

export class ApiError extends Error {
  readonly status: number;
  readonly details?: string | undefined;
  readonly body?: unknown;
  readonly fieldErrors?: Record<string, string> | undefined;

  constructor(
    status: number,
    message: string,
    options: {
      body?: unknown;
      details?: string;
      fieldErrors?: Record<string, string>;
    } = {},
  ) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.body = options.body;
    this.details = options.details;
    this.fieldErrors = options.fieldErrors;
  }
}

function buildAbsoluteUrl(path: string): string {
  if (/^https?:\/\//.test(path)) {
    return path;
  }

  const normalizedBase = API_BASE.endsWith("/") ? API_BASE.slice(0, -1) : API_BASE;
  const normalizedPath = path.startsWith("/") ? path : `/${path}`;
  return `${normalizedBase}${normalizedPath}`;
}

async function parseResponseBody(response: Response): Promise<unknown> {
  const text = await response.text();
  if (!text) {
    return undefined;
  }

  try {
    return JSON.parse(text) as unknown;
  } catch {
    return text;
  }
}

export async function apiFetch<T>(
  path: string,
  init: RequestInit = {},
  schema?: ZodType<T>,
): Promise<T> {
  const headers = new Headers(init.headers);
  headers.set("Accept", "application/json");
  if (init.body && !(init.body instanceof FormData) && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  const response = await fetch(buildAbsoluteUrl(path), {
    credentials: "include",
    ...init,
    headers,
  });

  const body = await parseResponseBody(response);
  if (!response.ok) {
    const errorPayload =
      body && typeof body === "object" ? (body as Record<string, unknown>) : undefined;
    const errorOptions: {
      body?: unknown;
      details?: string;
      fieldErrors?: Record<string, string>;
    } = {};

    if (body !== undefined) {
      errorOptions.body = body;
    }

    if (typeof errorPayload?.details === "string") {
      errorOptions.details = errorPayload.details;
    }

    if (
      errorPayload?.fieldErrors &&
      typeof errorPayload.fieldErrors === "object" &&
      !Array.isArray(errorPayload.fieldErrors)
    ) {
      errorOptions.fieldErrors = Object.fromEntries(
        Object.entries(errorPayload.fieldErrors).flatMap(([key, value]) =>
          typeof value === "string" ? [[key, value]] : [],
        ),
      );
    }

    throw new ApiError(
      response.status,
      typeof errorPayload?.message === "string"
        ? errorPayload.message
        : typeof errorPayload?.error === "string"
          ? errorPayload.error
          : response.statusText || "Request failed",
      errorOptions,
    );
  }

  if (!schema) {
    return body as T;
  }

  return schema.parse(body);
}

export function cached<T>(key: string, loader: () => Promise<T>): Promise<T> {
  const existing = cachedRequests.get(key);
  if (existing) {
    return existing as Promise<T>;
  }

  const promise = loader();
  cachedRequests.set(key, promise);
  promise.catch(() => {
    cachedRequests.delete(key);
  });
  return promise;
}

export function invalidateCached(prefix: string) {
  for (const key of cachedRequests.keys()) {
    if (key.startsWith(prefix)) {
      cachedRequests.delete(key);
    }
  }
}

export function buildWebSocketUrl(path: string): string {
  if (shouldUseApiMocks) {
    return `mock://${path.replace(/^\/+/, "")}`;
  }

  const target = new URL(buildAbsoluteUrl(path), window.location.origin);
  const protocol = target.protocol === "https:" ? "wss:" : "ws:";
  return `${protocol}//${target.host}${target.pathname}${target.search}`;
}

export function buildEventSourceUrl(path: string): string {
  if (shouldUseApiMocks) {
    return `mock://${path.replace(/^\/+/, "")}`;
  }

  return buildAbsoluteUrl(path);
}

export const messageResponseSchema = z.object({
  message: z.string(),
});
