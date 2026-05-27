---
title: 測試
description: 模板目前的測試覆蓋與建議驗證流程。
---

# 測試

專案目前已經在幾個核心 slice 上提供聚焦單元測試。

## 目前的測試覆蓋

- `internal/config/config_test.go`：驗證 config parsing 與語意限制
- `internal/handler/http_test.go`：驗證 `/healthz`、`/readyz` 與 webhook route 掛載
- `internal/session/memory_test.go`：驗證 counter、讀寫與 TTL expiry
- `internal/bot/service_test.go`：驗證 webhook cleanup、command registration 與 inline query 行為

## 建議驗證指令

```bash
make test
```

或：

```bash
go test ./...
```

Docusaurus 站點可用：

```bash
cd docs
node_modules/.bin/tsc --noEmit
```