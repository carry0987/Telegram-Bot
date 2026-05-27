---
title: 互動模式
description: 模板中的 reply keyboard、inline keyboard、callback query 與 inline query。
---

# 互動模式

除了 slash commands，模板也內建了早期最常見的互動型態。

## Reply keyboard

`/keyboard` 會送出 reply keyboard，`/hidekeyboard` 會把它移除。

適合用在：

- 想提供固定快捷操作
- 選項少而穩定
- 使用者停留在單一小流程中

## Inline keyboard

`/menu` 會送出帶 callback data 的 inline keyboard：

- `menu:hello`
- `menu:session`
- `menu:help`

handler 會先匹配 `menu:` 前綴，再根據 suffix 決定動作。

## Inline query

inline mode 透過檢查 `update.InlineQuery != nil` 啟用。

目前行為：

- 空 query 會回傳 help、pong、greeting 三個 article
- 非空 query 會回傳 echo article，再附上 help 與 menu 提示
- 結果標記為 personal，cache time 很短