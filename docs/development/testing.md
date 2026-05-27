---
title: Testing
description: Current test coverage and recommended validation workflow for extending the template.
---

# Testing

The project already includes focused unit tests around the current core slices.

## Existing test coverage

### Configuration

`internal/config/config_test.go` verifies:

- Minimal valid configuration
- Rejection of invalid values
- Webhook-mode validation requirements
- Username normalization

### HTTP handler

`internal/handler/http_test.go` verifies:

- `/healthz` returns `200`
- `/readyz` returns `503` on readiness failure
- `/readyz` returns `200` on readiness success
- The webhook route is mounted when configured

### Session storage

`internal/session/memory_test.go` verifies:

- Counter increment behavior
- Value set and get behavior
- TTL expiry in the in-memory backend

### Bot service

`internal/bot/service_test.go` verifies:

- Polling-mode webhook cleanup behavior
- Command registration without leading slashes
- Inline query result behavior

## Default validation commands

Use:

```bash
make test
```

You can also run Go tests directly:

```bash
go test ./...
```

For the Docusaurus site:

```bash
cd docs
node_modules/.bin/tsc --noEmit
```

## Testing recommendations when extending the template

- Add focused tests near the package you change
- Prefer handler-level tests for routing rules
- Keep Telegram SDK behavior behind stubs or interfaces where practical
- Validate configuration semantics with table-driven tests when adding new env vars

If a feature crosses packages, start by testing the smallest slice that proves the behavior rather than reaching immediately for full integration tests.