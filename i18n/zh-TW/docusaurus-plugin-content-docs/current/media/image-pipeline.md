---
title: 指令面
description: 模板內建的 Telegram commands 與它們示範的行為。
---

# 指令面

模板內建一組小而完整的 command，足以涵蓋幾種常見 Telegram pattern。

## 內建 commands

### `/start`

- 打招呼
- 重新開啟 reply keyboard
- 增加 `visits` session counter

### `/help`

- 回傳目前的 command 清單
- 更新 `last_command`

### `/ping`

- 回傳 `pong`

### `/echo <text>`

- 解析 command 後方參數
- 沒有文字時回傳 usage

### `/keyboard`

- 顯示 reply keyboard demo

### `/hidekeyboard`

- 移除 reply keyboard

### `/menu`

- 顯示 inline keyboard demo

### `/session`

- 增加 session counter
- 顯示目前 backend 與上一個 command