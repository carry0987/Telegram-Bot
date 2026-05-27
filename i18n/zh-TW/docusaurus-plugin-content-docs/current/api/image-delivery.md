---
title: Inline Query
description: 模板內建 inline query 行為的參考說明。
---

# Inline Query

模板內建一個以 article result 為基礎的 inline mode demo。

## 空 query 行為

當 inline query 為空或只有空白時，bot 會回傳三個預設結果：

- Bot Help
- Pong
- Greeting

## 非空 query 行為

當 query 有內容時，bot 會回傳：

- 用 query 內容當訊息文字的 echo article
- help article
- menu reminder article

## 回應特性

- `IsPersonal = true`
- `CacheTime = 1`
- 使用 `InlineQueryResultArticle`