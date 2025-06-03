"use client";

import React, { createContext, useContext, useMemo } from "react";
import { usePathname } from "next/navigation";
import { findRouteByPath, Route } from "@/config/routes";

interface RouteContextType {
  currentRoute: Route | undefined;
  breadcrumbs: Route[];
  isPublicRoute: boolean;
}

const RouteContext = createContext<RouteContextType | undefined>(undefined);

export function RouteProvider({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();

  const routeInfo = useMemo(() => {
    const currentRoute = findRouteByPath(pathname);
    const isPublicRoute = !!currentRoute?.isPublic;

    // Build breadcrumbs
    const breadcrumbs: Route[] = [];

    // If we're on a nested route, we need to build up the breadcrumb trail
    if (pathname) {
      const pathSegments = pathname.split('/').filter(Boolean);
      let currentPath = '';

      for (const segment of pathSegments) {
        currentPath += `/${segment}`;
        const route = findRouteByPath(currentPath);

        if (route) {
          breadcrumbs.push(route);
        }
      }
    }

    return {
      currentRoute,
      breadcrumbs,
      isPublicRoute
    };
  }, [pathname]);

  return (
    <RouteContext.Provider value={routeInfo}>
      {children}
    </RouteContext.Provider>
  );
}

export function useRoute() {
  const context = useContext(RouteContext);

  if (context === undefined) {
    throw new Error('useRoute must be used within a RouteProvider');
  }

  return context;
}
