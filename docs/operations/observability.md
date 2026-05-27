---
title: Observability
description: What the template currently exposes for health, readiness, logging, and runtime checks.
---

# Observability

The template keeps observability simple but not empty.

## Structured logging

The process configures `log/slog` at startup.

- Default level is `info`
- `DEBUG=true` raises the level to `debug`
- Startup and shutdown events are logged explicitly
- Bot transport mode and HTTP listener address are logged during runtime start

This gives you enough visibility for local development and for first deployments.

## Health and readiness

Two HTTP probe endpoints are available:

- `/healthz` for liveness
- `/readyz` for Telegram connectivity readiness

The distinction matters:

- `/healthz` answers whether the process is alive
- `/readyz` answers whether the bot can currently reach Telegram through `GetMe()`

That makes readiness more useful than a trivial always-healthy endpoint.

## Session backend visibility

The selected session backend appears in:

- The root endpoint response
- The `/session` command output

This is small, but it is useful when verifying whether Redis configuration actually took effect.

## What is not included yet

The template does not currently include:

- Prometheus metrics
- OpenTelemetry traces
- Request IDs
- Audit-event persistence
- Alerting integrations

Those are reasonable next additions once the bot grows beyond the starter stage.

## Recommended next steps for production observability

1. Add metrics around Telegram API failures and handler latency.
2. Add tracing if commands call external services.
3. Add structured fields that identify handler name and chat scope when operationally safe.