---
title: Callback Actions
description: Callback data contract and runtime behavior for inline keyboard actions.
---

# Callback Actions

The inline keyboard demo uses callback queries to show button-driven bot behavior.

## Callback data format

The current buttons use the `menu:` prefix.

Supported actions:

- `menu:hello`
- `menu:session`
- `menu:help`

Using a prefix-based scheme is important because it keeps routing simple and avoids accidental overlap when more callback-driven features are added.

## Runtime behavior

When a callback query arrives:

1. The handler strips the `menu:` prefix
2. It resolves the action
3. It records `last_command` as `menu:<action>`
4. It answers the callback query with a short `Done` acknowledgement
5. It sends a regular chat message with the full response

## Session-aware action

`menu:session` demonstrates that callback flows can update the same session state used by slash commands. It increments `session_hits` and includes the active backend in the reply text.

## Unknown actions

Unrecognized callback suffixes produce `Unknown menu action.`

This gives you a clear failure mode while you are extending the keyboard surface.