---
title: 設定
description: Telegram bot 模板的 environment variables 與驗證規則。
---

# 設定

所有 runtime 行為都由 `config.Config` 內的 environment variables 控制。

## 必填變數

### `TELEGRAM_BOT_TOKEN`

所有模式都需要。若為空，程式會在啟動驗證階段直接結束。

## 一般 service 設定

- `PORT`：HTTP server port，預設 `3000`
- `DEBUG`：是否開啟 debug logging，預設 `false`
- `BOT_MODE`：`polling` 或 `webhook`，預設 `polling`
- `BOT_INIT_TIMEOUT`：bot 初始化 timeout，預設 `5s`
- `READ_HEADER_TIMEOUT`：HTTP read header timeout，預設 `5s`
- `SHUTDOWN_TIMEOUT`：graceful shutdown timeout，預設 `10s`

## Telegram 設定

- `TELEGRAM_BOT_USERNAME`：選填，前導 `@` 會自動去除

## Webhook 設定

- `WEBHOOK_PUBLIC_URL`：webhook mode 必填，必須是有效的 `http://` 或 `https://` URL
- `WEBHOOK_PATH`：預設 `/telegram/webhook`，必須以 `/` 開頭
- `WEBHOOK_SECRET_TOKEN`：選填但建議在正式環境設定
- `WEBHOOK_DROP_PENDING_UPDATES`：註冊或刪除 webhook 時是否丟棄待處理更新

## Session 設定

- `REDIS_URL`：選填，未設定時使用 in-memory store
- `SESSION_TTL`：session 寫入使用的 TTL，預設 `24h`