---
title: Bot Commands
description: Reference for the built-in Telegram commands exposed by the template.
---

# Bot Commands

The application currently exposes these built-in commands through Telegram.

## `/start`

Introduces the bot, shows the reply keyboard, and increments the `visits` session counter.

## `/help`

Returns a concise summary of the available commands.

## `/ping`

Returns `pong`.

## `/echo <text>`

Echoes the provided text back to the chat. When called without arguments, the bot returns `Usage: /echo <text>`.

## `/keyboard`

Shows the reply keyboard demo.

## `/hidekeyboard`

Hides the reply keyboard and instructs the user to run `/keyboard` to restore it.

## `/menu`

Shows the inline keyboard demo.

## `/session`

Returns the current session counter, the selected session backend, and the previous recorded command when available.

## Default fallback

For non-empty text messages that do not match a registered command, the bot sends the help message instead of ignoring the update.

That default makes the bot friendlier during local development and gives you a clear place to change the unmatched-message behavior later.