---
title: 架構總覽
description: Telegram bot 模板的 runtime 結構與 package 邊界。
---

# 架構總覽

這個 repository 被設計成一個小型 service，而不是單檔 bot script。這讓 runtime composition、transport 邏輯與狀態儲存有清楚邊界。

## Package 地圖

```text
cmd/server            程序入口
internal/config       env 載入、正規化與驗證
internal/app          composition root 與 runtime lifecycle
internal/bot          Telegram transport、handler、menu sync
internal/handler      HTTP 端點與 webhook wiring
internal/session      session abstraction 與 backend
```

## 啟動順序

1. 載入 `.env` 與 `.env.local`
2. 解析 environment variables 成 `config.Config`
3. 驗證 mode、webhook URL、TTL 等語意限制
4. 初始化 session store
5. 初始化 Telegram bot service
6. 建立 HTTP handler 與 server
7. 執行 bot transport 的 prepare
8. 用共享 `errgroup` 啟動 bot transport 與 HTTP server

這個順序讓設定錯誤能提早失敗，也避免 HTTP server 先開起來但 bot 還沒準備好。

## Composition root

`internal/app` 負責：

- 選擇 session backend
- 建立 bot service
- 建立 HTTP handler
- 同時執行 bot 與 HTTP server
- 協調 graceful shutdown

它不應該承擔 Telegram command 邏輯或 HTTP response shaping。