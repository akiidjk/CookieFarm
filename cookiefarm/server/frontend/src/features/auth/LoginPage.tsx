import { useActionState, useEffect } from "react";
import { Banner } from "@cloudflare/kumo/components/banner";
import { Button } from "@cloudflare/kumo/components/button";
import { Input } from "@cloudflare/kumo/components/input";
import { WarningCircle } from "@phosphor-icons/react";
import { useLocation, useNavigate } from "react-router";
import { ApiError } from "@/api/client";
import { useAuth } from "./AuthProvider";

type LoginState = {
  errorMessage: string | null;
  completed: boolean;
};

const initialState: LoginState = {
  errorMessage: null,
  completed: false,
};

export function LoginPage() {
  const auth = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const redirectTarget =
    typeof location.state === "object" &&
    location.state &&
    "from" in location.state &&
    typeof location.state.from === "string"
      ? location.state.from
      : "/";

  const [state, submitAction, pending] = useActionState(
    async (_previousState: LoginState, formData: FormData): Promise<LoginState> => {
      const password = String(formData.get("password") ?? "").trim();
      const username = String(formData.get("username") ?? "").trim();

      try {
        await auth.login(password, username || undefined);
        return {
          errorMessage: null,
          completed: true,
        };
      } catch (error) {
        return {
          errorMessage:
            error instanceof ApiError
              ? error.message
              : error instanceof Error
                ? error.message
                : "Login failed.",
          completed: false,
        };
      }
    },
    initialState,
  );

  useEffect(() => {
    if (auth.status === "authenticated" || state.completed) {
      navigate(redirectTarget, { replace: true });
    }
  }, [auth.status, navigate, redirectTarget, state.completed]);

  return (
    <main className="flex min-h-screen items-center justify-center px-4 py-10">
      <section className="dashboard-surface w-full max-w-md rounded-3xl border border-kumo-line p-6 shadow-sm">
        <div className="space-y-2">
          <p className="text-sm uppercase tracking-[0.2em] text-kumo-fg-secondary">
            CookieFarm
          </p>
          <h1 className="text-3xl font-semibold text-kumo-fg-primary">Operator Login</h1>
          <p className="text-sm text-kumo-fg-secondary">
            Authenticate with the server password to access the dashboard.
          </p>
        </div>

        <form action={submitAction} className="mt-6 space-y-4">
          {state.errorMessage ? (
            <Banner
              variant="error"
              icon={<WarningCircle weight="fill" />}
              title="Authentication failed"
              description={state.errorMessage}
            />
          ) : null}

          <Input
            name="username"
            label="Username"
            placeholder="cookieguest"
            autoComplete="username"
          />

          <Input
            name="password"
            type="password"
            label="Password"
            placeholder="Enter server password"
            autoComplete="current-password"
            required
          />

          <Button type="submit" className="w-full" loading={pending}>
            Sign in
          </Button>
        </form>
      </section>
    </main>
  );
}
