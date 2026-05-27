---
title: 儲存佈局
description: session 在 memory 與 Redis 中的 key 形式與 TTL 行為。
---

# 儲存佈局

這個模板目前沒有關聯式資料庫或大型 domain store。唯一的持久化 runtime state 是 chat-scoped session data。

## Session backend 選擇

- `REDIS_URL` 為空時使用 in-memory store
- `REDIS_URL` 有設定時使用 Redis

其餘程式只依賴 `session.Store` 介面。

## 目前儲存的資料

- `visits`：`/start` 的計數器
- `session_hits`：`/session` 與 session menu action 的計數器
- `last_command`：記錄上一個 command 或 menu action

## Memory store key

```text
{chatID}:{key}
```

例如：

```text
42:last_command
```

每個 entry 都會保存值與過期時間，過期項目在讀取時惰性清除。

## Redis key

```text
telegram-bot-template:session:{chatID}:{key}
```

例如：

```text
telegram-bot-template:session:42:last_command
telegram-bot-template:session:42:session_hits
```

這個 namespace 能避免共用 Redis 時與其他應用程式衝突。