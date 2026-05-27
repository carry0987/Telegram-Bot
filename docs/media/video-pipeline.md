---
title: Interactive Patterns
description: Reply keyboards, inline keyboards, callback queries, and inline queries in the template.
---

# Interactive Patterns

Beyond slash commands, the template includes the interaction patterns most bots need early.

## Reply keyboard demo

`/keyboard` sends a reply keyboard and `/hidekeyboard` removes it.

Use this pattern when:

- You want persistent shortcuts in the chat input area
- The choices are small and stable
- The user is expected to stay inside a narrow flow

## Inline keyboard demo

`/menu` sends an inline keyboard with buttons that emit callback data:

- `menu:hello`
- `menu:session`
- `menu:help`

The handler matches callback data using the `menu:` prefix and then switches on the suffix.

This is a good default pattern because it keeps callback parsing predictable and makes room for future prefixes such as `settings:` or `admin:`.

## Callback query response flow

When a button is clicked, the bot does two things:

1. Calls `answerCallbackQuery` with a short acknowledgement
2. Sends a regular chat message with the actual response text

That split keeps the Telegram client responsive while still giving the user a visible message in the chat history.

## Inline query demo

Inline mode is enabled through a match function that checks `update.InlineQuery != nil`.

Current behavior:

- Empty queries return help, pong, and greeting articles
- Non-empty queries return an echo article plus help and menu reminders
- Results are marked personal and use a short cache time

This keeps the demo safe for experimentation without pretending to be a real search system.

## Choosing the right interaction pattern

Use:

- Slash commands for discoverable, explicit actions
- Reply keyboards for repeated shortcuts inside a focused chat flow
- Inline keyboards for context-sensitive actions tied to a message
- Inline queries when the bot should insert content into another chat