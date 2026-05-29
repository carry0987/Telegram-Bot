---
title: 快速開始
description: 用最少設定在本地跑起 Telegram bot 模板。
---

# 快速開始

最快的起步方式是使用 polling mode 搭配 in-memory session backend。你只需要 Telegram bot token 與本地 Go 環境。

## 前置需求

- Go 1.25 以上
- 從 BotFather 取得的 Telegram bot token
- 選用：Docker 與 Docker Compose，用於 Redis 或容器化執行
- 選用：Node.js 22 與 pnpm，用於本地建置這個 Docusaurus 使用手冊

## 1. 取得 bot token

在 Telegram 中對 BotFather 執行：

```text
/newbot
```

請保存 BotFather 回傳的 token。沒有 `TELEGRAM_BOT_TOKEN`，應用程式不會啟動。

## 2. 建立環境檔

```bash
cp .env.example .env
```

最小 polling mode 設定：

```dotenv
TELEGRAM_BOT_TOKEN=123456:replace-me
BOT_MODE=polling
```

## 3. 啟動應用程式

```bash
make run
```

或直接使用 Go：

```bash
go run ./cmd/server
```

啟動流程會依序進行：

1. 載入 `.env` 與 `.env.local`
2. 驗證設定
3. 建立 session store
4. 建立 Telegram bot service
5. 同步 Telegram command menu
6. 啟動 bot transport 與 HTTP server

## 4. 驗證 bot 是否可用

在 Telegram 對 bot 傳送：

```text
/start
/help
/ping
/echo hello world
/keyboard
/menu
/session
```

本地 HTTP server 也會提供：

- `GET http://localhost:3000/`
- `GET http://localhost:3000/healthz`
- `GET http://localhost:3000/readyz`

## 5. 需要持久化 session 時再加 Redis

```bash
docker compose -f docker-compose.dev.yml up -d
```

然後在 `.env` 補上：

```dotenv
REDIS_URL=redis://localhost:6381
SESSION_TTL=24h
```

若 `REDIS_URL` 沒有設定，程式會自動退回 in-memory store。

## 6. 部署時切換到 webhook mode

```dotenv
BOT_MODE=webhook
WEBHOOK_PUBLIC_URL=https://bot.example.com
WEBHOOK_PATH=/telegram/webhook
WEBHOOK_SECRET_TOKEN=replace-with-random-secret
```

啟動時會向 Telegram 註冊 `https://bot.example.com/telegram/webhook`。