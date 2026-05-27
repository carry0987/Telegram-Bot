---
title: Extending the Template
description: Guidelines for evolving the starter into a larger Telegram bot application.
---

# Extending the Template

The starter is intentionally small, but it is not meant to stay frozen. The important part is extending it without erasing the boundaries that already help the codebase stay maintainable.

## Preserve the current roles

Try to keep these rules intact:

- `cmd/server` stays responsible for process startup
- `internal/config` remains the only place that knows how environment variables are parsed and validated
- `internal/app` remains the wiring layer
- `internal/bot` remains the Telegram-facing transport layer
- `internal/session` remains a small chat-state abstraction

## Introduce a use-case layer when needed

The next architectural step for most real bots is a package that owns business rules separate from Telegram delivery.

That package might:

- Validate domain inputs
- Call external APIs
- Talk to domain persistence
- Return domain results that handlers format for Telegram

Do that before `internal/bot/service.go` becomes the home of everything.

## Evolve storage deliberately

The current session store is appropriate for ephemeral chat state. If your bot needs:

- Permanent user profiles
- Reports
- Admin data
- Search indexes

Add dedicated persistence abstractions rather than stretching the session API beyond its current purpose.

## Keep operational behavior first-class

As you add features, preserve these service qualities:

- Configuration validation should fail fast
- Health and readiness behavior should stay meaningful
- Transport switching should remain explicit
- Graceful shutdown should still close external resources cleanly

## A practical growth path

1. Add new handlers and extract reusable formatting helpers.
2. Move business logic into a separate use-case layer.
3. Add external integrations behind interfaces.
4. Add richer observability once the bot becomes operationally important.
5. Split the bot package internally if the number of handlers grows substantially.

The goal is not to keep the starter tiny forever. The goal is to let it grow without losing coherence.