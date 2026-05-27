---
title: Command Surface
description: Built-in Telegram commands and what they demonstrate.
---

# Command Surface

The template deliberately ships with a small command set that exercises several Telegram patterns without making the code hard to follow.

## Built-in commands

### `/start`

- Greets the user
- Reopens the reply keyboard
- Increments the `visits` session counter

### `/help`

- Returns the current command reference
- Records `last_command`

### `/ping`

- Returns `pong`
- Provides a minimal command path for connectivity checks inside Telegram

### `/echo <text>`

- Parses arguments after the command name
- Returns a usage string when no text is provided
- Demonstrates a command that needs free-form input

### `/keyboard`

- Sends a reply keyboard
- Shows how to keep a fixed set of chat-visible shortcuts open

### `/hidekeyboard`

- Removes the reply keyboard
- Demonstrates reply keyboard cleanup

### `/menu`

- Sends an inline keyboard with callback actions
- Demonstrates button-driven interaction

### `/session`

- Increments a session counter
- Shows the active session backend
- Shows the previous recorded command when present

## Command menu sync

During startup, the bot calls `setMyCommands` so Telegram clients can show the command menu automatically.

If you add a new user-facing command, update both:

- `helpMessage()`
- `syncCommands()`

Otherwise the runtime behavior and Telegram's command menu will drift apart.

## Matching behavior

The project uses two command match strategies:

- `MatchTypeCommandStartOnly` for commands such as `/start`
- `MatchTypeCommand` for commands such as `/echo <text>`

Register command names without a leading slash. That is a library requirement, not a style choice.