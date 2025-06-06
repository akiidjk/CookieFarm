
export interface Route {
  path: string;
  label: string;
  icon?: string;
  children?: Route[];
  isPublic?: boolean;
  hideFromSidebar?: boolean;
}

// Define the applications's route structure
export const routes: Route[] = [
  {
    path: '/login',
    label: 'Login',
    isPublic: true,
    hideFromSidebar: true,
  },
  {
    path: '/dashboard',
    label: 'Dashboard',
  },
];

// Helper function to get all public routes
export function getPublicRoutes(): string[] {
  return routes
    .filter(route => route.isPublic)
    .map(route => route.path);
}

// Helper function to get protected routes
export function getProtectedRoutes(): string[] {
  const protectedRoutes: string[] = [];

  function collectRoutes(routeList: Route[]) {
    routeList.forEach(route => {
      if (!route.isPublic) {
        protectedRoutes.push(route.path);
      }

      if (route.children) {
        collectRoutes(route.children);
      }
    });
  }

  collectRoutes(routes);
  return protectedRoutes;
}

// Helper function to find a route by path
export function findRouteByPath(path: string): Route | undefined {
  function findRoute(routeList: Route[]): Route | undefined {
    for (const route of routeList) {
      if (route.path === path) {
        return route;
      }

      if (route.children) {
        const found = findRoute(route.children);
        if (found) return found;
      }
    }
    return undefined;
  }

  return findRoute(routes);
}
