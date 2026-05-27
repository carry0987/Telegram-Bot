---
title: 擴充模板
description: 把 starter 演進成較大型 Telegram bot 應用程式的指引。
---

# 擴充模板

這個 starter 雖然小，但不是要永遠保持不變。重要的是在擴充時不要破壞已經存在的良好邊界。

## 保留現有責任分工

- `cmd/server` 負責程序啟動
- `internal/config` 負責 environment parsing 與 validation
- `internal/app` 負責 wiring
- `internal/bot` 負責 Telegram-facing transport
- `internal/session` 維持為小型 chat-state abstraction

## 需要時導入 use-case layer

大多數真實 bot 的下一步，都會是把商業邏輯從 Telegram delivery 中抽離出來。

這個新 layer 可以負責：

- 驗證 domain input
- 呼叫外部 API
- 操作持久化層
- 回傳讓 handler 格式化成 Telegram response 的結果

請在 `internal/bot/service.go` 變成萬用檔案之前先做這一步。