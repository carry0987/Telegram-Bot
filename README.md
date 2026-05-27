# Telegram-Bot

A starter template for building Telegram bots in Go.

The repository includes a runnable application shape with polling and webhook support, command menu sync, health probes, inline and callback demos, and session storage backed by Redis or an in-memory fallback.

## Quick Start

1. Create a bot with BotFather and get the token.
2. Copy the environment file.
3. Set `TELEGRAM_BOT_TOKEN`.
4. Run the bot in polling mode.

```bash
cp .env.example .env
make run
```

Minimum `.env` example:

```dotenv
TELEGRAM_BOT_TOKEN=123456:replace-me
BOT_MODE=polling
```

Then open the bot in Telegram and try:

```text
/start
/help
/ping
/echo hello world
/keyboard
/menu
/session
```

## What You Get

- `cmd/server` as the single process entrypoint
- Environment-based configuration and validation
- Polling and webhook transport modes
- Telegram command menu synchronization on startup
- Built-in command and interaction demos
- `GET /healthz` and `GET /readyz`
- Redis session storage with automatic in-memory fallback
- Docker Compose files for local infrastructure and containerized runs

## Documentation

Use the documentation site for full setup and operational details:

- Getting started
- Integration and extension guidance
- Architecture and request flow
- Session storage behavior
- Configuration and deployment

The documentation site is built from the `gh-pages/` folder and deployed with GitHub Pages.

## Development

```bash
make fmt
make lint
make test
make tidy
```
