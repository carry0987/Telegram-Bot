---
title: 介紹
description: Go Telegram-Bot 的整體介紹與文件範圍。
---

# Telegram-Bot

這個專案是一個用 Go 開發 Telegram-Bot 的起手式模板，不把所有邏輯都塞進單一 `main.go`。它提供可直接運作的 runtime 骨架、實用的指令面、session 儲存、健康檢查，以及同時涵蓋本地開發與正式環境的傳輸模式。

## 目前已內建的能力

- `cmd/server` 單一程序入口
- 以環境變數驅動的設定與驗證
- 同時支援 polling 與 webhook 的 Telegram transport
- 啟動時自動同步 Telegram command menu
- 內建 `/start`、`/help`、`/ping`、`/echo`、`/keyboard`、`/hidekeyboard`、`/menu`、`/session`
- inline query 與 callback query 範例
- `/`、`/healthz`、`/readyz` HTTP 端點
- Redis session store 與 in-memory fallback
- 本地基礎設施與容器化執行的 Docker Compose 設定

## 這份文件會涵蓋什麼

- 應用程式如何啟動、驗證設定並做 graceful shutdown
- Telegram update 如何被路由到 command handler
- session key 在 memory 與 Redis 中的儲存方式
- 如何在本地用 polling、在部署環境用 webhook 運行
- 如何在不破壞現有架構的前提下擴充模板

## 核心分層

- `cmd/server` 載入環境變數、設定 logger、驗證 config、處理 signal
- `internal/config` 解析與驗證 environment variables
- `internal/app` 負責 bot service、session backend 與 HTTP server 的 wiring
- `internal/bot` 負責 Telegram transport 與 handler 註冊
- `internal/handler` 提供 health check 與 webhook 入口
- `internal/session` 抽象 session persistence

這樣的分層把 Telegram SDK 細節收斂在單一區塊，也為之後抽出 use case 或 domain logic 留出空間。

## 內建互動模式

這個模板雖然小，但已經涵蓋常見的 Telegram bot 互動形式：

- slash commands
- reply keyboard
- inline keyboard 與 callback data
- inline query 回覆
- 具 TTL 的 chat-scoped session counter
- 未匹配文字訊息的預設 fallback 回覆

這些能力足夠讓它成為真正可延伸的 bot 起點，而不只是 transport demo。