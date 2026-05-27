---
title: 整合指南
description: 在保留現有架構的前提下擴充 Telegram bot 模板。
---

# 整合指南

目前程式碼量不大，但已經有值得保留的邊界。這份指南說明如何加入新功能，同時避免把所有邏輯都塞進 Telegram handler。

## 先從擁有行為的層開始

- `internal/bot` 負責 Telegram update matching、輸入解析與 Telegram API 呼叫
- `internal/session` 負責 chat-scoped state persistence
- `internal/config` 負責來自環境的行為設定
- `internal/app` 負責 runtime composition 與 lifecycle

若你要加入商業邏輯，通常應該新增 package 讓 bot 層呼叫，而不是直接把 `internal/bot/service.go` 變成整個應用程式。

## 新增 command 的方式

目前 command 註冊都在 `registerHandlers()`。

重要細節：`go-telegram/bot` 的 command pattern 不能帶前導 `/`。註冊時要用 `start`，不是 `/start`。

常用 match type：

- `MatchTypeCommandStartOnly` 用在只應該匹配訊息開頭的 command
- `MatchTypeCommand` 用在需要參數的 command，例如 `/echo <text>`
- `MatchTypePrefix` 用在 callback data 前綴，例如 `menu:`

建議流程：

1. 在 `registerHandlers()` 註冊 handler
2. 新增對應的 `handle...` 方法
3. 更新 `helpMessage()`
4. 如果需要顯示在 Telegram command menu，更新 `syncCommands()`
5. 補上聚焦的測試

## 保持 transport 與 business logic 分離

如果 handler 開始負責太多業務判斷，就應該把核心邏輯抽到新 package，例如：

- `internal/usecase`
- `internal/store`
- `internal/service`

bot 層應該專注在 Telegram model 與應用程式輸入輸出之間的轉換。