---
title: Session State
description: How the template stores lightweight chat state and how to decide between memory and Redis.
---

# Session State

The template includes a session layer because even simple bots usually need a small amount of chat-scoped state.

## Why session state exists here

Without session storage, the template would only demonstrate stateless commands. That is too narrow for real bot development.

The current session layer supports:

- Counters
- Simple string values
- Backend selection through configuration
- TTL-based expiry

## Current built-in uses

- `/start` tracks how many times the current chat started interacting with the bot
- `/session` tracks how many times the session demo has been used
- Most commands write the `last_command` key
- The inline menu `session` action updates the same `session_hits` counter

## Backend selection

The choice is automatic at startup:

- No `REDIS_URL` means the in-memory store
- A configured `REDIS_URL` means the Redis store

The bot exposes the selected backend in two places:

- The root HTTP endpoint
- The `/session` command output

## TTL and expiry

Every session write uses `SESSION_TTL`, which defaults to `24h`.

This makes the store suitable for:

- Short-lived menu state
- Temporary onboarding progress
- Diagnostic counters

It is not designed to be your permanent source of user records.

## When to keep memory only

Memory is enough when:

- You are prototyping locally
- Restarting the bot can safely clear state
- You are running a single instance

## When to move to Redis

Redis is the better choice when:

- Restart-safe session continuity matters
- You want consistent behavior across multiple instances
- You are deploying webhook mode behind real infrastructure

If your bot outgrows string values and counters, keep the session store small and introduce a separate domain persistence layer rather than turning this package into a general database abstraction.