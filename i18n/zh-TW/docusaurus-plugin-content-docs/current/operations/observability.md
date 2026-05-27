---
title: 可觀測性
description: 模板目前提供的 health、readiness、logging 與 runtime 檢查能力。
---

# 可觀測性

這個模板的 observability 很精簡，但不是空白。

## 結構化 logging

程式在啟動時會設定 `log/slog`：

- 預設層級是 `info`
- `DEBUG=true` 時會提升到 `debug`
- 啟動與關閉事件都會明確記錄

## Health 與 readiness

HTTP probe 端點：

- `/healthz`：liveness
- `/readyz`：Telegram connectivity readiness

差異在於：

- `/healthz` 只回答 process 是否活著
- `/readyz` 會透過 `GetMe()` 檢查 bot 與 Telegram 的可用性

## 目前尚未內建的項目

- Prometheus metrics
- OpenTelemetry traces
- Request ID
- Audit event persistence