---
title: 部署
description: Telegram bot 模板的 polling 與 webhook 部署方式。
---

# 部署

這個 repository 支援兩種 delivery mode，請依照環境選擇。

## Polling mode

適合：

- 本地開發
- 還沒有公開 HTTPS endpoint
- 單一 bot instance 即可

最小設定：

```dotenv
TELEGRAM_BOT_TOKEN=123456:replace-me
BOT_MODE=polling
```

啟動時會先檢查目前 webhook 狀態；如果 Telegram 還有設定 webhook，就會先刪掉，避免阻塞 polling。

## Webhook mode

適合已經有穩定 HTTPS 入口的部署環境。

```dotenv
TELEGRAM_BOT_TOKEN=123456:replace-me
BOT_MODE=webhook
WEBHOOK_PUBLIC_URL=https://bot.example.com
WEBHOOK_PATH=/telegram/webhook
WEBHOOK_SECRET_TOKEN=replace-with-random-secret
REDIS_URL=redis://redis:6379
```

部署時至少要確保：

- 外部可達的 HTTPS
- 代理層能把請求轉給應用程式
- `/healthz` 與 `/readyz` 已接入 probe 系統