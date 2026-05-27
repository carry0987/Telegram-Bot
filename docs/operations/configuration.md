---
title: Configuration
description: Environment variables and configuration rules for the Telegram bot template.
---

# Configuration

All runtime behavior is controlled through environment variables parsed into `config.Config`.

## Required variables

### `TELEGRAM_BOT_TOKEN`

Required in all modes. The application exits during startup if it is empty.

## General service variables

### `PORT`

HTTP server port. Default: `3000`.

### `DEBUG`

Enables debug logging when `true`. Default: `false`.

### `BOT_MODE`

Supported values:

- `polling`
- `webhook`

Default: `polling`.

### `BOT_INIT_TIMEOUT`

Timeout used during bot client initialization checks. Default: `5s`.

### `READ_HEADER_TIMEOUT`

HTTP server read-header timeout. Default: `5s`.

### `SHUTDOWN_TIMEOUT`

Maximum graceful shutdown window. Default: `10s`.

## Telegram variables

### `TELEGRAM_BOT_USERNAME`

Optional bot username. If present, leading `@` is stripped automatically. It is used to make help text more accurate for group-chat command usage.

## Webhook variables

### `WEBHOOK_PUBLIC_URL`

Required only in webhook mode. Must be a valid `http://` or `https://` URL with a host.

### `WEBHOOK_PATH`

Default: `/telegram/webhook`. Must start with `/`. In webhook mode it cannot be exactly `/`.

### `WEBHOOK_SECRET_TOKEN`

Optional but recommended in webhook mode.

### `WEBHOOK_DROP_PENDING_UPDATES`

Controls whether pending updates are dropped when registering or deleting the webhook. Default: `false`.

## Session variables

### `REDIS_URL`

Optional Redis connection string. When unset, the application uses the in-memory session store.

### `SESSION_TTL`

TTL applied to session writes. Default: `24h`.

## Validation rules worth noting

- `PORT` must be between `1` and `65535`
- All configured timeouts and TTLs must be greater than `0`
- `BOT_MODE` must be `polling` or `webhook`
- `REDIS_URL` must parse as a valid Redis URL if provided

Keep configuration decisions in this layer. Do not spread environment lookups across handlers or use-case code.