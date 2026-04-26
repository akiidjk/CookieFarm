import { Suspense, lazy, type ReactNode } from "react";
import { Navigate, createBrowserRouter } from "react-router";
import { PageSkeleton } from "@/components/PageSkeleton";
import { RouteError } from "@/components/RouteError";
import { RequireAuth } from "@/features/auth/AuthProvider";
import { AppLayout } from "@/layouts/AppLayout";

const LoginPage = lazy(async () => {
  const module = await import("@/features/auth/LoginPage");
  return { default: module.LoginPage };
});

const DashboardPage = lazy(async () => {
  const module = await import("@/features/dashboard/DashboardPage");
  return { default: module.DashboardPage };
});

const FlagsPage = lazy(async () => {
  const module = await import("@/features/flags/FlagsPage");
  return { default: module.FlagsPage };
});

const ExploitsPage = lazy(async () => {
  const module = await import("@/features/exploits/ExploitsPage");
  return { default: module.ExploitsPage };
});

const ChartsPage = lazy(async () => {
  const module = await import("@/features/charts/ChartsPage");
  return { default: module.ChartsPage };
});

const ConfigPage = lazy(async () => {
  const module = await import("@/features/config/ConfigPage");
  return { default: module.ConfigPage };
});

function suspenseElement(node: ReactNode) {
  return <Suspense fallback={<PageSkeleton />}>{node}</Suspense>;
}

export const router = createBrowserRouter([
  {
    path: "/login",
    element: suspenseElement(<LoginPage />),
    errorElement: <RouteError />,
  },
  {
    path: "/",
    element: (
      <RequireAuth>
        <AppLayout />
      </RequireAuth>
    ),
    errorElement: <RouteError />,
    children: [
      {
        index: true,
        element: suspenseElement(<DashboardPage />),
      },
      {
        path: "flags",
        element: suspenseElement(<FlagsPage />),
      },
      {
        path: "exploits",
        element: suspenseElement(<ExploitsPage />),
      },
      {
        path: "charts",
        element: suspenseElement(<ChartsPage />),
      },
      {
        path: "config",
        element: suspenseElement(<ConfigPage />),
      },
      {
        path: "dashboard",
        element: <Navigate to="/" replace />,
      },
    ],
  },
]);
