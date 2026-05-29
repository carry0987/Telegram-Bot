---
title: Introduction
description: Overview of the Go Telegram-Bot and what the documentation covers.
---

# Telegram-Bot

This repository is a starter for building Telegram bots in Go without collapsing everything into a single `main.go` file. It gives you a working runtime shape, a small but useful command surface, session storage, health probes, and deployment modes that cover both local development and production delivery.

## What the template already includes

- A single process entrypoint in `cmd/server`
- Environment-based configuration with validation
- Telegram transport support for both polling and webhook delivery
- Command menu synchronization on startup
- Built-in commands for `/start`, `/help`, `/ping`, `/echo`, `/keyboard`, `/hidekeyboard`, `/menu`, and `/session`
- Inline query and callback query demos
- HTTP endpoints for `/`, `/healthz`, and `/readyz`
- Session storage backed by Redis or an in-memory fallback
- Docker Compose files for local infrastructure and containerized runs

## What this documentation site covers

- How the application boots, validates configuration, and shuts down cleanly
- How Telegram updates are routed into command handlers
- How session keys are stored in memory and Redis
- How to run the bot in polling mode locally and webhook mode in deployed environments
- How to extend the template without breaking the current architecture

## Core runtime shape

The application is deliberately split into a few stable layers:

- `cmd/server` loads environment variables, configures logging, validates config, and handles OS signals
- `internal/config` parses and validates environment variables
- `internal/app` wires the bot service, session backend, and HTTP server together
- `internal/bot` owns Telegram-specific transport and handler registration
- `internal/handler` exposes the HTTP endpoints used for health checks and webhook delivery
- `internal/session` abstracts session persistence behind a small store interface

This keeps Telegram SDK details local to one part of the codebase and leaves room for future application logic to live outside transport handlers.

## Default interaction model

The template is intentionally small but it demonstrates the most common Telegram interaction paths:

- Standard slash commands
- Reply keyboards
- Inline keyboards with callback data
- Inline query answers
- Per-chat session counters with TTL
- A default fallback response for unrecognized text messages

That combination is enough to turn the repository into a real bot instead of a transport proof of concept.

## When to use this template

Use it when you want:

- A Go bot project that can start locally with only a bot token
- A codebase that is already structured for tests and future growth
- A clean switch between polling and webhook deployment modes
- A lightweight session layer without committing to a database on day one

If you need background job orchestration, persistent domain data, or multiple external integrations, this repo is still a good starting point, but those parts are expected to be added on top of the current skeleton.