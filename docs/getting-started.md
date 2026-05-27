---
title: Getting Started
description: Run the Telegram bot template locally with the minimum required setup.
---

# Getting Started

The fastest path to a working bot is polling mode with the in-memory session backend. You only need a Telegram bot token and Go installed locally.

## Prerequisites

- Go 1.25 or newer
- A Telegram bot token from BotFather
- Optional: Docker and Docker Compose if you want Redis or containerized runs
- Optional: Node.js 22 and pnpm if you want to build this Docusaurus site locally

## 1. Create a bot token

Use BotFather in Telegram:

```text
/newbot
```

Save the token that BotFather returns. The application will not start without `TELEGRAM_BOT_TOKEN`.

## 2. Create your environment file

Copy the example environment file and fill in the token:

```bash
cp .env.example .env
```

Minimum polling-mode configuration:

```dotenv
TELEGRAM_BOT_TOKEN=123456:replace-me
BOT_MODE=polling
```

## 3. Run the application

You can start the bot directly:

```bash
make run
```

Or use Go directly:

```bash
go run ./cmd/server
```

On startup the application will:

1. Load `.env` and `.env.local`
2. Validate configuration
3. Build the session store
4. Build the Telegram bot service
5. Sync the Telegram command menu
6. Start the bot transport and the HTTP server

## 4. Verify the bot is working

Open the bot in Telegram and try:

```text
/start
/help
/ping
/echo hello world
/keyboard
/menu
/session
```

The local HTTP server also exposes:

- `GET http://localhost:3000/`
- `GET http://localhost:3000/healthz`
- `GET http://localhost:3000/readyz`

## 5. Add Redis when you need durable session state

For local infrastructure only:

```bash
docker compose -f docker-compose.dev.yml up -d
```

Then point the app at Redis:

```dotenv
REDIS_URL=redis://localhost:6381
SESSION_TTL=24h
```

If `REDIS_URL` is unset, the bot falls back to the in-memory store automatically.

## 6. Switch to webhook mode when deploying

Webhook mode requires a public base URL:

```dotenv
BOT_MODE=webhook
WEBHOOK_PUBLIC_URL=https://bot.example.com
WEBHOOK_PATH=/telegram/webhook
WEBHOOK_SECRET_TOKEN=replace-with-random-secret
```

The bot will register `https://bot.example.com/telegram/webhook` with Telegram during startup.

## Common startup failures

### Missing token

If `TELEGRAM_BOT_TOKEN` is empty, the application exits during configuration validation.

### Invalid webhook configuration

When `BOT_MODE=webhook`, the app requires:

- A valid `WEBHOOK_PUBLIC_URL`
- A webhook path that starts with `/`
- A webhook path that is not exactly `/`

### Invalid Redis URL

If `REDIS_URL` is set, it must parse as a valid Redis connection string.

## Next steps

- Read the integration guide before adding new commands
- Read the architecture section if you plan to split business logic from transport
- Read the operations section before deploying webhook mode