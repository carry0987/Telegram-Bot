---
title: 整合範例
description: 擴充 command、session 與部署模式的實用做法。
---

# 整合範例

這一頁整理了最常見的擴充工作。

## 新增簡單 command

適用於只需要解析輸入並回覆訊息的情況。

1. 在 `registerHandlers()` 註冊 command
2. 新增 `handleYourCommand()`
3. 在 `helpMessage()` 補上說明
4. 需要顯示在 Telegram menu 時，更新 `syncCommands()`

例如 `/about`、`/version`、`/status` 都適合這個型態。

## 新增帶參數的 command

`/echo` 已經示範了這種模式。通常要注意：

- 先 trim whitespace
- 缺參數時回傳明確 usage
- 不要把業務邏輯直接寫在解析分支內

## 紀錄短期對話狀態

session store 適合這種狀態：

- 作用域是單一 chat
- 資料量小，可用 string 或 counter 表示
- 過一段時間過期是可接受的

範例：

- 上次選到的 menu action
- onboarding 的目前步驟
- 臨時確認狀態

## 從 memory 切到 Redis

不需要改程式，只要設定：

```dotenv
REDIS_URL=redis://localhost:6381
SESSION_TTL=24h
```

需要跨重啟或跨多實例保留 session 時，就應該改用 Redis。