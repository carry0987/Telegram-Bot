---
title: Bot 指令
description: 模板目前透過 Telegram 暴露的內建 commands 參考。
---

# Bot 指令

目前應用程式透過 Telegram 提供以下內建 command。

## `/start`

介紹 bot、顯示 reply keyboard，並增加 `visits` session counter。

## `/help`

回傳可用 command 摘要。

## `/ping`

回傳 `pong`。

## `/echo <text>`

把輸入文字回送到聊天室。若沒有參數，會回傳 `Usage: /echo <text>`。

## `/keyboard`

顯示 reply keyboard demo。

## `/hidekeyboard`

隱藏 reply keyboard，並提示使用者用 `/keyboard` 再次開啟。

## `/menu`

顯示 inline keyboard demo。

## `/session`

回傳目前的 session counter、session backend，以及上一個記錄到的 command。