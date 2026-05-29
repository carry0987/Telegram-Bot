---
title: Architecture Overview
description: Runtime structure and package boundaries of the Telegram bot.
---

# Architecture Overview

The repository is structured as a small service, not as a one-file bot script. That choice matters because it gives you explicit boundaries for runtime composition, transport logic, and state persistence.

## Package map

```text
cmd/server            process entrypoint
internal/config       env loading, normalization, validation
internal/app          composition root and runtime lifecycle
internal/bot          Telegram transport, handlers, menu sync
internal/handler      HTTP endpoints and webhook wiring
internal/session      session store abstraction and backends
```

## Startup sequence

At startup the process follows this order:

1. Load `.env` and `.env.local`
2. Parse environment variables into `config.Config`
3. Validate semantic constraints such as mode, webhook URL, and TTLs
4. Initialize the session store
5. Initialize the Telegram bot service
6. Build the HTTP handler and server
7. Prepare the bot transport
8. Start the bot transport and the HTTP server under a shared `errgroup`

That ordering keeps configuration failures cheap and prevents the application from opening the HTTP listener before the bot layer is ready.

## Composition root

`internal/app` is the composition root. It is responsible for:

- Choosing the session backend
- Building the bot service
- Building the HTTP handler
- Running the HTTP server and bot service together
- Coordinating graceful shutdown

It does not contain Telegram command logic or HTTP response shaping.

## Telegram transport layer

`internal/bot` owns the Telegram client lifecycle and update routing. Its responsibilities include:

- Creating the `go-telegram/bot` client with the correct options
- Registering handlers for text commands, callback data, and inline queries
- Syncing Telegram commands via `setMyCommands`
- Preparing polling or webhook transport before runtime starts
- Providing a readiness probe through `GetMe()`

The package also contains the current demo interaction flows. As your application grows, this package should increasingly delegate business logic outward instead of absorbing it.

## HTTP surface

`internal/handler` exposes four concerns:

- A root endpoint with service metadata
- A liveness endpoint
- A readiness endpoint that checks Telegram connectivity
- The Telegram webhook endpoint in webhook mode

The HTTP server exists even in polling mode because operational probes still matter.

## Session persistence

`internal/session` hides whether chat-scoped state lives in memory or Redis. The interface is intentionally narrow so the rest of the application only relies on the operations it actually uses today.

That tradeoff keeps the template simple while still allowing:

- Local development without extra infrastructure
- Durable session data when Redis is available
- Clear backend reporting in `/session` responses and root metadata

## Shutdown model

The bot transport and HTTP server run under a shared context. When the context is canceled:

- The HTTP server receives a bounded shutdown window
- The session backend is closed
- The process exits through the shared `errgroup`

This keeps shutdown logic centralized instead of scattering cleanup in multiple packages.