# CookieFarm Frontend

React + Vite dashboard for CookieFarm server operations.

## Requirements

- Node.js 20+ (or Bun if you prefer)
- A running CookieFarm server API (default: `http://127.0.0.1:8080`)

## Install

```/dev/null/install.sh#L1-2
npm install
# or: bun install
```

## Development

Run the frontend dev server on `localhost:5173`:

```/dev/null/dev.sh#L1-2
npm run dev
# or: bun run dev
```

In development, API requests should go to `/api/v1/*` and be proxied by Vite to the backend server.

### Dev API target (default)

- Proxy target: `http://127.0.0.1:8080`
- Frontend origin: `http://localhost:5173`

If you need to override API base directly, set:

```/dev/null/env.md#L1-1
VITE_API_BASE_URL=http://127.0.0.1:8080/api/v1
```

## Production Build

Build static assets:

```/dev/null/build.sh#L1-2
npm run build
# or: bun run build
```

Preview build locally:

```/dev/null/preview.sh#L1-2
npm run preview
# or: bun run preview
```

## Serving with Go Fiber

When the built frontend is served by the Go Fiber server on the same host, keep API calls relative:

- API base should resolve to `/api/v1`
- Frontend and API share origin in production

This avoids CORS issues and works cleanly in both modes:

1. Vite dev server (`localhost:5173`) via proxy
2. Fiber-served built frontend (same-origin `/api/v1`)

## Optional Mock API

To enable mock handlers in dev:

```/dev/null/mock-env.md#L1-1
VITE_USE_API_MOCKS=true
```

Use mocks only for local UI development. Keep it off for real backend integration.

## Scripts

- `dev` — start Vite dev server
- `build` — type-check + production build
- `preview` — preview production bundle
