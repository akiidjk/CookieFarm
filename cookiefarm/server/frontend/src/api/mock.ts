import { shouldUseApiMocks } from "./client";

type MockFlag = {
  flag_code: string;
  service_name: string;
  port_service: number;
  submit_time: number;
  response_time: number;
  msg: string;
  status: 0 | 1 | 2 | 3 | 4;
  team_id: number;
  username: string;
  exploit_name: string;
};

const config = {
  configured: true,
  server: {
    url_flag_checker: "http://localhost:5001/flags",
    team_token: "mock-team-token",
    submit_flag_checker_time: 120,
    max_flag_batch_size: 1000,
    protocol: "cc_http",
    tick_time: 120,
    flag_ttl: 0,
    start_time: "2023-10-01T00:00:00Z",
    end_time: "2023-10-31T23:59:59Z",
  },
  shared: {
    services: {
      CookieService: 8081,
      vulnify: 1337,
      "app-nc": 1338,
    },
    regex_flag: "[A-Z0-9]{31}=",
    format_ip_teams: "10.10.{}.1",
    my_team_id: 1,
    url_flag_ids: "http://localhost:5001/flagIds",
    nop_team: 0,
    range_ip_teams: 9,
    configured: true,
  },
};

const flags: MockFlag[] = Array.from({ length: 200 }, (_, index) => ({
  flag_code: `FLAG{${String(index + 1).padStart(4, "0")}}`,
  service_name: index % 2 === 0 ? "CookieService" : "vulnify",
  port_service: index % 2 === 0 ? 8081 : 1337,
  submit_time: Math.floor(Date.now() / 1000) - index * 20,
  response_time: Math.floor(Date.now() / 1000) - index * 20 + (index % 3 === 0 ? 4 : 0),
  msg: index % 3 === 0 ? "Accepted by mock checker" : "Queued in mock collector",
  status: index % 4 === 0 ? 1 : index % 4 === 1 ? 0 : index % 4 === 2 ? 2 : 3,
  team_id: (index % 8) + 1,
  username: "mock",
  exploit_name: "manual",
}));

export async function startMocking(): Promise<void> {
  if (!shouldUseApiMocks || typeof window === "undefined") {
    return;
  }

  const [{ http, HttpResponse }, { setupWorker }] = await Promise.all([
    import("msw"),
    import("msw/browser"),
  ]);

  const worker = setupWorker(
    http.get("/api/v1/", async () =>
      HttpResponse.json({
        message: "The cookie is up!!",
        time: new Date().toISOString(),
      }),
    ),
    http.post("/api/v1/auth/login", async () => HttpResponse.json({})),
    http.post("/api/v1/auth/logout", async () => HttpResponse.json({})),
    http.get("/api/v1/auth/verify", async () => HttpResponse.json({})),
    http.get("/api/v1/protocols", async () =>
      HttpResponse.json({
        protocols: ["cc_http", "cc_tcp"],
      }),
    ),
    http.get("/api/v1/config", async () => HttpResponse.json(config)),
    http.post("/api/v1/config", async () => HttpResponse.json({ message: "ok" })),
    http.get("/api/v1/stats", async () =>
      HttpResponse.json({
        flags_stats: [
          {
            team_id: 2,
            total_flags: 315,
            accepted_flags: { Float64: 151, Valid: true },
            denied_flags: { Float64: 164, Valid: true },
            unsubmitted_flags: { Float64: 0, Valid: true },
            error_flags: { Float64: 0, Valid: true },
            not_valid_flags: { Float64: 0, Valid: true },
          },
          {
            team_id: 3,
            total_flags: 268,
            accepted_flags: { Float64: 135, Valid: true },
            denied_flags: { Float64: 133, Valid: true },
            unsubmitted_flags: { Float64: 0, Valid: true },
            error_flags: { Float64: 0, Valid: true },
            not_valid_flags: { Float64: 0, Valid: true },
          },
        ],
      }),
    ),
    http.get("/api/v1/flags", async () =>
      HttpResponse.json({
        flags,
        n_flags: flags.length,
      }),
    ),
    http.get("/api/v1/flags/:limit", async ({ request, params }) => {
      const url = new URL(request.url);
      const offset = Number(url.searchParams.get("offset") ?? "0");
      const status = url.searchParams.get("status");
      const service = url.searchParams.get("service");
      const team = url.searchParams.get("team");
      const search = (url.searchParams.get("search") ?? "").toLowerCase();
      const searchField = url.searchParams.get("search_field") ?? "flag_code";
      const limit = Number(params.limit);

      const filtered = flags.filter((flag) => {
        if (status && String(flag.status) !== status) {
          return false;
        }
        if (service && flag.service_name !== service) {
          return false;
        }
        if (team && String(flag.team_id) !== team) {
          return false;
        }
        if (search) {
          const value = String(flag[searchField as keyof MockFlag] ?? "").toLowerCase();
          return value.includes(search);
        }
        return true;
      });

      return HttpResponse.json({
        flags: filtered.slice(offset, offset + limit),
        n_flags: filtered.length,
      });
    }),
    http.post("/api/v1/submit-flag", async () => HttpResponse.json({ message: "ok" })),
    http.delete("/api/v1/delete-flag", async () => HttpResponse.json({ message: "ok" })),
  );

  await worker.start({
    onUnhandledRequest: "bypass",
  });
}
