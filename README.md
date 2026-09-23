# SplitMate

[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://go.dev/) [![CI](https://github.com/igorynos/SplitMate/actions/workflows/ci.yml/badge.svg)](https://github.com/igorynos/SplitMate/actions/workflows/ci.yml)

Shared-expense backend that records group purchases and calculates a compact settlement plan showing who should pay whom.

This is the Go successor to [SplitMate-python](https://github.com/igorynos/SplitMate-python).

## Highlights

- Integer monetary arithmetic without floating-point rounding
- Fair remainder distribution when an expense cannot be divided evenly
- Deterministic debt settlement between participants
- Concurrency-safe group state
- PostgreSQL persistence for users, ledgers, wallets, expenses, and payments
- Participant management and document deletion
- Transport-independent business logic with unit tests
- HTTP API, graceful shutdown, and minimal Docker image

## API

```text
POST /users
POST /ledgers
POST /ledgers/{id}/members
POST /ledgers/{id}/expenses
POST /ledgers/{id}/payments
GET  /ledgers/{id}/settlement
DELETE /documents/{kind}/{id}?user_id={id}
GET  /health
```

Amounts are expressed in the smallest currency unit—for example, cents rather than floating-point dollars.

```bash
go test ./...
go run ./cmd/splitmate

curl -X POST localhost:8080/ledgers/1/expenses \
  -H 'Content-Type: application/json' \
  -d '{"description":"Dinner","amount":3000,"paid_by":1,"participants":[1,2,3]}'

curl localhost:8080/ledgers/1/settlement
```

The service uses the project layout conventions relevant to an application: `cmd`, `internal`, `configs`, and `deployments`. Start the complete environment with `docker compose -f deployments/compose.yml up --build`.
