# SplitMate

[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://go.dev/) [![CI](https://github.com/igorynos/SplitMate/actions/workflows/ci.yml/badge.svg)](https://github.com/igorynos/SplitMate/actions/workflows/ci.yml)

Shared-expense backend that records group purchases and calculates a compact settlement plan showing who should pay whom.

This is the Go successor to [SplitMate-python](https://github.com/igorynos/SplitMate-python).

## Highlights

- Integer monetary arithmetic without floating-point rounding
- Fair remainder distribution when an expense cannot be divided evenly
- Deterministic debt settlement between participants
- Concurrency-safe group state
- Transport-independent business logic with unit tests
- HTTP API and minimal Docker image

## API

```text
POST /groups/{group}/expenses
GET  /groups/{group}/settlement
GET  /health
```

Amounts are expressed in the smallest currency unit—for example, cents rather than floating-point dollars.

```bash
go test ./...
go run ./cmd/splitmate

curl -X POST localhost:8080/groups/trip/expenses \
  -H 'Content-Type: application/json' \
  -d '{"description":"Dinner","amount":3000,"paid_by":"Igor","participants":["Igor","Anna","Max"]}'

curl localhost:8080/groups/trip/settlement
```

## Roadmap

- PostgreSQL event and balance persistence
- Telegram bot adapter
- Expense correction and deletion
- Monthly summaries and history export
- Idempotency keys for repeated Telegram updates
