#!/usr/bin/env just

[group("frontend")]
[working-directory("cookiefarm/server/frontend")]
fe-build:
    bun install && bun run build

[group("frontend")]
[working-directory("cookiefarm/server/frontend")]
fe-dev:
    bun install && bun run dev
