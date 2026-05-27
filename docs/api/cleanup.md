---
title: Webhook Endpoint
description: How the Telegram webhook endpoint is registered and exposed in webhook mode.
---

# Webhook Endpoint

Webhook delivery is only active when `BOT_MODE=webhook`.

## Route

The default route is:

```text
POST /telegram/webhook
```

You can change the path with `WEBHOOK_PATH`, but it must start with `/` and cannot be exactly `/` in webhook mode.

## Registration flow

During startup, the bot service builds the full public webhook URL from:

- `WEBHOOK_PUBLIC_URL`
- `WEBHOOK_PATH`

It then calls Telegram `setWebhook` with:

- The full URL
- The optional secret token
- The `WEBHOOK_DROP_PENDING_UPDATES` flag

## Secret token behavior

If `WEBHOOK_SECRET_TOKEN` is configured, the Telegram client is initialized with webhook secret validation enabled.

This is recommended for production because it helps distinguish legitimate Telegram requests from arbitrary traffic reaching the same endpoint.

## Polling-mode cleanup

When the application starts in polling mode, it first checks webhook state through `getWebhookInfo`. It only calls `deleteWebhook` if Telegram reports that a webhook URL is already configured.

That prevents unnecessary deletion calls and avoids the common local-development problem where a leftover webhook blocks polling.