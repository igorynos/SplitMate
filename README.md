# SplitMate

Go rewrite of [SplitMate-python](https://github.com/igorynos/SplitMate-python), a service for tracking shared expenses and calculating who owes whom.

## Implemented

- Add expenses for a group through HTTP
- Split integer monetary amounts without rounding loss
- Produce a deterministic minimal settlement plan
- Concurrency-safe service and unit tests
- Small multi-stage Docker image

## Run

```bash
go test ./...
go run ./cmd/splitmate
```

```bash
curl -X POST localhost:8080/groups/trip/expenses \
  -H 'Content-Type: application/json' \
  -d '{"description":"Dinner","amount":3000,"paid_by":"Igor","participants":["Igor","Anna","Max"]}'
curl localhost:8080/groups/trip/settlement
```

Amounts use the smallest currency unit. PostgreSQL persistence and the Telegram adapter are the next migration steps.
