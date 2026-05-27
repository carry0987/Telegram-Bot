---
title: Integration Guide
description: Extend the Telegram bot template while preserving the current architecture.
---

# Integration Guide

The current codebase is small enough to understand quickly, but it already has boundaries that are worth preserving. This guide describes how to integrate new bot behavior without pushing everything into the Telegram handler layer.

## Start from the owning layer

Use these responsibilities as your default rule:

- `internal/bot` owns Telegram update matching, input parsing, and Telegram API calls
- `internal/session` owns per-chat state persistence
- `internal/config` owns environment-derived behavior
- `internal/app` owns runtime composition and lifecycle

If you add business rules, they should usually become a new package that the bot layer calls into, instead of turning `internal/bot/service.go` into the entire application.

## Adding a new command

The current command registration happens in `registerHandlers()`.

Important detail: the `go-telegram/bot` library command pattern does not include the leading slash. Register `start`, not `/start`.

Use these match types intentionally:

- `MatchTypeCommandStartOnly` for commands that should only match at the start of the message
- `MatchTypeCommand` for commands that need arguments, such as `/echo <text>`
- `MatchTypePrefix` for callback data prefixes like `menu:`

Recommended flow for a new command:

1. Register the handler in `registerHandlers()`
2. Add a dedicated `handle...` method
3. Update `helpMessage()`
4. Add the command to `syncCommands()` if it should appear in Telegram's command menu
5. Write a focused test for the behavior or the routing rule

## Keeping transport and business logic separate

If a handler starts doing more than parsing input and formatting a response, extract the core behavior into a new package. For example:

- `internal/usecase` for application rules
- `internal/store` for domain-specific persistence
- `internal/service` for integration-oriented orchestration

The bot layer should remain the place where Telegram models are translated into application inputs and application outputs are translated back into Telegram responses.

## Working with session state

The session store is intentionally small:

- `Increment()` for counters
- `Set()` for string values
- `Get()` for reading string values
- `Backend()` for diagnostics and output

This is enough for lightweight bot state such as:

- Visit counts
- Last command tracking
- Menu flow checkpoints
- Short-lived conversational state

If you need more complex state, either wrap the existing store or introduce a dedicated persistence layer instead of expanding the store blindly.

## Choosing polling or webhook for integrations

Use polling when:

- You are developing locally
- You do not yet have a public endpoint
- Simplicity matters more than horizontal scale

Use webhook mode when:

- The bot runs behind a public HTTPS endpoint
- You want Telegram to push updates directly to your service
- Your deployment already exposes HTTP reliably

The rest of the application is intentionally similar across both modes. Only the delivery mechanism changes.

## Handling callback and inline workflows

The template already demonstrates two common Telegram interaction patterns:

- Inline keyboard buttons using callback data with the `menu:` prefix
- Inline query responses using article results

When you extend these flows:

- Keep callback data structured and prefix-based
- Validate and normalize any user-controlled data before using it
- Keep inline query results small and predictable
- Use session state only when the interaction is actually chat-scoped

## Recommended extension path

If you are evolving the template into a real product, this is a reasonable sequence:

1. Add new commands and extract reusable formatting
2. Introduce a use-case layer for business logic
3. Add external APIs or domain storage behind interfaces
4. Add metrics, tracing, or audit logging where needed
5. Move from polling to webhook mode when deployment requires it