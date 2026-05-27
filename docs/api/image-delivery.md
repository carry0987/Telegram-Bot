---
title: Inline Queries
description: Reference for the inline query behavior implemented by the template.
---

# Inline Queries

The template includes a lightweight inline mode demo built around article results.

## Empty-query behavior

When the inline query text is empty or whitespace, the bot returns three default results:

- Bot Help
- Pong
- Greeting

This gives users something useful even before they type a search phrase.

## Non-empty query behavior

When a query string is present, the bot returns:

- An echo article using the query as the inserted message text
- A help article
- A menu reminder article

## Response characteristics

- Results are marked `IsPersonal = true`
- `CacheTime` is set to `1`
- Results use `InlineQueryResultArticle`

This favors correctness and fast iteration over aggressive caching.

## Typical extension points

Replace the current article construction when you need:

- Search against your own data
- User-specific suggestions
- Richer content insertion

If inline mode becomes important, move result generation into a dedicated package so the Telegram transport layer stays focused on delivery.