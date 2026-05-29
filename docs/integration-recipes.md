---
title: Integration Recipes
description: Practical patterns for extending the bot with new commands, state, and delivery modes.
---

# Integration Recipes

This page collects common extension tasks you are likely to perform first.

## Add a simple command

Use this pattern when the bot only needs to parse input and send a reply.

1. Register the command in `registerHandlers()`
2. Add a `handleYourCommand()` method
3. Add a help line in `helpMessage()`
4. Add a command menu entry in `syncCommands()` if appropriate

This is the right shape for commands like `/about`, `/version`, or `/status`.

## Add a command with arguments

The `/echo` handler already shows the simplest version of this pattern. Reuse the same shape when you need free-form text after a command.

Typical rules:

- Trim whitespace before parsing
- Return explicit usage text when arguments are missing
- Avoid doing domain work directly inside the parsing branch

## Track a short conversational step

Use the session store when the state is:

- Scoped to one chat
- Small enough to fit into string values or counters
- Safe to expire after a fixed TTL

Examples:

- Last selected menu action
- Current onboarding step
- Temporary confirmation state

The template already tracks counters and `last_command`, so you can extend the same convention.

## Switch from in-memory to Redis sessions

No code change is required. Set `REDIS_URL` and restart the application.

```dotenv
REDIS_URL=redis://localhost:6381
SESSION_TTL=24h
```

Use Redis when you need session continuity across restarts or across multiple running instances.

## Add a new inline keyboard action

The existing callback handling uses the `menu:` prefix and switches on the suffix.

Recommended pattern:

1. Add a new button in `inlineMenuMarkup()`
2. Add a new case in `handleMenuAction()`
3. Keep the callback data prefix stable
4. Update any session bookkeeping if the new action changes user flow

This keeps callback handling easy to read and easy to test.

## Add inline query results

The template uses article results for inline mode because they are simple and predictable.

To extend the behavior:

1. Decide whether empty queries should return defaults or nothing
2. Build results in `inlineQueryResults()`
3. Keep result IDs stable per logical item
4. Prefer message text that works well in arbitrary chats

If inline mode becomes a real search surface, move result construction into a dedicated package.

## Prepare a production webhook deployment

Set the required environment variables:

```dotenv
BOT_MODE=webhook
WEBHOOK_PUBLIC_URL=https://bot.example.com
WEBHOOK_PATH=/telegram/webhook
WEBHOOK_SECRET_TOKEN=replace-with-random-secret
WEBHOOK_DROP_PENDING_UPDATES=false
```

Before switching environments, make sure:

- The public URL is reachable by Telegram
- HTTPS termination is in place
- Your reverse proxy forwards requests to the application
- Health probes are wired to `/healthz` and `/readyz`

## Add domain logic cleanly

When a command becomes more than a demo, extract the core logic first and keep the handler responsible only for:

- Reading Telegram input
- Calling the application layer
- Formatting the response

That keeps the template maintainable as the command list grows.