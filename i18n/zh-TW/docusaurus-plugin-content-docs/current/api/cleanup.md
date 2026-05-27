---
title: Webhook 端點
description: Telegram webhook endpoint 在 webhook mode 中如何註冊與暴露。
---

# Webhook 端點

只有在 `BOT_MODE=webhook` 時才會啟用 webhook delivery。

## 路徑

預設路徑是：

```text
POST /telegram/webhook
```

你可以透過 `WEBHOOK_PATH` 變更，但它必須以 `/` 開頭，且在 webhook mode 中不能等於 `/`。

## 註冊流程

啟動時 bot service 會用以下兩個值組出完整 URL：

- `WEBHOOK_PUBLIC_URL`
- `WEBHOOK_PATH`

接著呼叫 Telegram `setWebhook`，並帶入完整 URL、secret token 與 `WEBHOOK_DROP_PENDING_UPDATES`。