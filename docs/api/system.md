---
title: System Endpoints
description: Reference for the HTTP endpoints exposed by the application.
---

# System Endpoints

The HTTP server runs in both polling and webhook mode.

## `GET /`

Returns basic service metadata:

- `name`
- `description`
- `mode`
- `session`

Example response:

```json
{
  "name": "telegram-bot",
  "description": "Go Telegram-Bot with health checks, polling, webhook, and session demos",
  "mode": "polling",
  "session": "memory"
}
```

## `GET /healthz`

Returns a liveness response and does not depend on Telegram being reachable.

Example response:

```json
{
  "ok": true,
  "status": "healthy",
  "service": "telegram-bot"
}
```

## `GET /readyz`

Returns readiness based on a Telegram API probe.

The handler creates a short timeout context and calls the bot service `Ping()` method, which delegates to `GetMe()`.

Successful response:

```json
{
  "ok": true,
  "status": "ready",
  "message": "telegram bot client is ready"
}
```

Failure response uses `503 Service Unavailable` and includes the probe error message.

## `POST /telegram/webhook`

Only mounted when webhook mode is active. The route is provided by the Telegram client's webhook handler and is not available in polling mode.