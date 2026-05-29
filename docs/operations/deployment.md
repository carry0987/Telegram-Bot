---
title: Deployment
description: Polling and webhook deployment patterns for the Telegram bot.
---

# Deployment

The repository supports two runtime delivery modes. Choose the one that matches your environment.

## Polling deployment

Polling mode is the simplest path.

Use it when:

- You are developing locally
- You do not want to expose a public HTTPS endpoint yet
- A single bot instance is enough

Minimum configuration:

```dotenv
TELEGRAM_BOT_TOKEN=123456:replace-me
BOT_MODE=polling
```

Important behavior:

- Startup inspects current webhook state
- If Telegram still has an active webhook URL, the app deletes it before polling begins
- If no webhook URL is configured, no delete call is made

## Webhook deployment

Webhook mode is the better fit when your application already runs behind a reachable HTTPS endpoint.

Example configuration:

```dotenv
TELEGRAM_BOT_TOKEN=123456:replace-me
BOT_MODE=webhook
WEBHOOK_PUBLIC_URL=https://bot.example.com
WEBHOOK_PATH=/telegram/webhook
WEBHOOK_SECRET_TOKEN=replace-with-random-secret
REDIS_URL=redis://redis:6379
```

Deployment requirements:

- Public HTTPS reachability
- Reverse proxy or load balancer forwarding to the app
- Health probes wired to `/healthz` and `/readyz`
- A stable webhook URL that Telegram can reach consistently

## Docker Compose options

### Infrastructure only

Use the development compose file when you want Redis locally but still run the app with `go run`:

```bash
docker compose -f docker-compose.dev.yml up -d
```

### Full containerized run

Use the main compose file when you want the app and Redis inside containers:

```bash
docker compose up --build
```

### Tunnel-assisted webhook development

The compose setup also includes an optional Cloudflare Tunnel profile. Use it when you want a public ingress path during development without managing your own domain directly.

## Recommended production baseline

- Run webhook mode behind HTTPS
- Set `WEBHOOK_SECRET_TOKEN`
- Use Redis for session state
- Expose `/healthz` and `/readyz` to your platform's probe system
- Keep logs aggregated so Telegram API failures are visible quickly