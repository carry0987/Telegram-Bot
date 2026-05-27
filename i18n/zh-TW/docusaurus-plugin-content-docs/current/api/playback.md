---
title: Callback Action
description: Inline keyboard callback data 的格式與執行行為。
---

# Callback Action

inline keyboard demo 透過 callback query 展示按鈕驅動的 bot 行為。

## Callback data 格式

目前按鈕都使用 `menu:` 前綴。

支援的 action：

- `menu:hello`
- `menu:session`
- `menu:help`

## 執行流程

callback query 進來後會：

1. 去掉 `menu:` 前綴
2. 解析 action
3. 把 `last_command` 設成 `menu:<action>`
4. 用 `Done` 回應 callback query
5. 再送出一般聊天訊息作為完整結果

`menu:session` 也會更新 `session_hits`，示範 callback 流程與 slash command 共用 session state。