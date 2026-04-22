import {
  createContext,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";
import { Navigate, useLocation } from "react-router";
import { login as loginRequest, logout as logoutRequest, verifyAuth } from "@/api/auth";
import { PageSkeleton } from "@/components/PageSkeleton";
import { useInterval } from "@/hooks/useInterval";

type AuthStatus = "checking" | "authenticated" | "anonymous";

type AuthContextValue = {
  status: AuthStatus;
  login: (password: string, username?: string) => Promise<void>;
  logout: () => Promise<void>;
  refresh: () => Promise<boolean>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider(props: { children: ReactNode }) {
  const [status, setStatus] = useState<AuthStatus>("checking");

  async function refresh(): Promise<boolean> {
    const isAuthenticated = await verifyAuth();
    setStatus(isAuthenticated ? "authenticated" : "anonymous");
    return isAuthenticated;
  }

  useEffect(() => {
    void refresh();
  }, []);

  useInterval(
    () => {
      if (status !== "authenticated") {
        return;
      }
      void refresh();
    },
    60_000,
    { enabled: status === "authenticated" },
  );

  return (
    <AuthContext.Provider
      value={{
        status,
        login: async (password: string, username?: string) => {
          await loginRequest({
            ...(username ? { username } : {}),
            password,
          });
          setStatus("authenticated");
        },
        logout: async () => {
          try {
            await logoutRequest();
          } finally {
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
