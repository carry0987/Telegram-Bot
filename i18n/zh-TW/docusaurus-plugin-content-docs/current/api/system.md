---
title: 系統端點
description: 應用程式目前暴露的 HTTP 端點參考。
---

# 系統端點

HTTP server 在 polling 與 webhook mode 都會存在。

## `GET /`

回傳基本服務資訊：

- `name`
- `description`
- `mode`
- `session`

## `GET /healthz`

回傳 liveness 狀態，不依賴 Telegram 是否可達。

## `GET /readyz`

透過 Telegram API probe 決定 readiness。handler 會建立短 timeout context，呼叫 bot service 的 `Ping()`，而 `Ping()` 會再呼叫 `GetMe()`。

成功時回傳 `200`，失敗時回傳 `503` 與錯誤訊息。

## `POST /telegram/webhook`

只會在 webhook mode 掛上，polling mode 不提供這條 route。