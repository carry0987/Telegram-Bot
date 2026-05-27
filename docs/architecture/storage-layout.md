---
title: Storage Layout
description: Session storage behavior, key layout, and TTL handling in memory and Redis.
---

# Storage Layout

The template does not ship with a relational database or a domain event store. The only persisted runtime state today is chat-scoped session data.

## Session storage model

The application chooses the session backend at startup:

- If `REDIS_URL` is empty, it uses the in-memory store
- If `REDIS_URL` is set, it uses Redis

The rest of the application only sees the `session.Store` interface.

## Data currently stored

The built-in handlers store three logical values:

- `visits` as a counter for `/start`
- `session_hits` as a counter for `/session` and the session menu action
- `last_command` as a string recording the previous command or menu action

This is intentionally minimal. It demonstrates how to store conversational state without implying a larger persistence model.

## Memory store layout

The in-memory backend uses a `map[string]memoryEntry` guarded by a mutex.

Key shape:

```text
{chatID}:{key}
```

Example:

```text
42:last_command
```

Each entry tracks:

- The stored string value
- An expiration timestamp derived from the configured TTL

Expired items are removed lazily when they are accessed.

## Redis key layout

The Redis backend uses namespaced string keys.

Key shape:

```text
telegram-bot-template:session:{chatID}:{key}
```

Examples:

```text
telegram-bot-template:session:42:last_command
telegram-bot-template:session:42:session_hits
```

This namespace prevents collisions when the same Redis instance is shared with other applications.

## TTL behavior

`SESSION_TTL` defaults to `24h` and is applied whenever the session store is written.

Current behavior:

- `Increment()` refreshes the TTL after incrementing a counter
- `Set()` stores the value with the provided TTL
- `Get()` does not extend TTL by itself

That means session state behaves like short-lived activity state, not like permanent user data.

## Operational implications

Use the in-memory backend when:

- You are developing locally
- Restarts can safely reset session state
- You only run one process

Use Redis when:

- Session continuity matters across restarts
- You may run more than one application instance
- You want the `/session` demo to survive deploy cycles