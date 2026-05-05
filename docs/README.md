# CookieFarm Documentation

This directory contains architecture and protocol documentation for CookieFarm.

## Reference Documents

- [System Understanding Document](../tech-docs/system-understanding.md)
- [Dependency Graph](../tech-docs/dependency-graph.md)
- [CKP Protocol](../tech-docs/ckp-protocol.md)

## Documentation App

The `docs` project can also be used as a Fumadocs/Next.js documentation app.

Run the development server:

```bash
npm run dev
# or
pnpm dev
# or
yarn dev
```

Open `http://localhost:3000` in your browser.

## Project Notes

- `lib/source.ts` provides the content source adapter.
- `lib/layout.shared.tsx` stores shared layout options.
- `app/docs` contains the documentation layout and pages.
- `app/api/search/route.ts` contains the search route handler.
