---
title: Session 狀態
description: 模板如何儲存輕量 chat state，以及何時選擇 memory 或 Redis。
---

# Session 狀態

這個模板包含 session layer，因為即使是簡單 bot，也常常需要少量 chat-scoped state。

## 為什麼這裡需要 session

沒有 session 的話，模板只能示範無狀態 command，這對真實 bot 開發太狹窄。

目前 session layer 支援：

- counter
- string value
- 透過設定選擇 backend
- 以 TTL 為基礎的過期機制

## 目前的使用方式

- `/start` 追蹤互動次數
- `/session` 追蹤 session demo 使用次數
- 大多數 command 都會寫入 `last_command`
- inline menu 的 `session` action 也會更新 `session_hits`