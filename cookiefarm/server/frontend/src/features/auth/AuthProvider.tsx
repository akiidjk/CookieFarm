import {
  createContext,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import useSWR from "swr";
import { Navigate, useLocation } from "react-router";
import {
  authVerifyKey,
  login as loginRequest,
  logout as logoutRequest,
  verifyAuth,
  type AuthSession,
} from "@/api/auth";
import { PageSkeleton } from "@/components/PageSkeleton";

type AuthStatus = "checking" | "authenticated" | "anonymous";

type AuthContextValue = {
  status: AuthStatus;
  user: AuthSession | null;
  login: (password: string, username?: string) => Promise<void>;
  logout: () => Promise<void>;
  refresh: () => Promise<boolean>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider(props: { children: ReactNode }) {
  const [status, setStatus] = useState<AuthStatus>("checking");
  const [user, setUser] = useState<AuthSession | null>(null);
  const authQuery = useSWR(authVerifyKey, verifyAuth, {
    refreshInterval: status === "authenticated" ? 60_000 : 0,
  });

  async function refresh(): Promise<boolean> {
    const session = await authQuery.mutate();
    setUser(session ?? null);
    setStatus(session ? "authenticated" : "anonymous");
    return Boolean(session);
  }

  useEffect(() => {
    if (authQuery.data === undefined) {
      return;
    }
    setUser(authQuery.data ?? null);
    setStatus(authQuery.data ? "authenticated" : "anonymous");
  }, [authQuery.data]);

  return (
    <AuthContext.Provider
      value={{
        status,
        user,
        login: async (password: string, username?: string) => {
          await loginRequest({
            ...(username ? { username } : {}),
            password,
          });
          const session = { username: username || "cookieguest" };
          void authQuery.mutate(session, { revalidate: false });
          setUser(session);
          setStatus("authenticated");
        },
        logout: async () => {
          try {
            await logoutRequest();
          } finally {
            void authQuery.mutate(null, { revalidate: false });
            setUser(null);
            setStatus("anonymous");
          }
        },
        refresh,
      }}
    >
      {props.children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return context;
}

export function RequireAuth(props: { children: ReactNode }) {
  const auth = useAuth();
  const location = useLocation();

  if (auth.status === "checking") {
    return <PageSkeleton />;
  }

  if (auth.status !== "authenticated") {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }

  return props.children;
}
