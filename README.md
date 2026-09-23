# SplitMate

[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?logo=go)](https://go.dev/) [![CI](https://github.com/igorynos/SplitMate/actions/workflows/ci.yml/badge.svg)](https://github.com/igorynos/SplitMate/actions/workflows/ci.yml)

Shared-expense backend that records group purchases and calculates a compact settlement plan showing who should pay whom.

This is the Go successor to [SplitMate-python](https://github.com/igorynos/SplitMate-python).

## ✨ Features

- 💰 **Safe money calculations:** Uses integer minor units instead of floating-point values.
- ⚖️ **Fair expense splitting:** Distributes indivisible remainders deterministically.
- 🔄 **Debt minimization:** Produces a compact settlement plan between debtors and creditors.
- 👥 **Group accounting:** Manages users, ledgers, participants, and personal wallets.
- 🛒 **Financial documents:** Records purchases and direct payments separately.
- 🗑️ **Controlled deletion:** Allows users to remove their own expense and payment documents.
- 🗄️ **Persistent state:** Stores accounting data and relations in PostgreSQL.
- 🧪 **Testable domain:** Keeps calculation logic independent from HTTP and database packages.
- 📴 **Operational safety:** Supports graceful shutdown, health checks, containers, and CI.

## 🔄 Settlement Example

If Igor pays `3000` for a dinner shared by Igor, Anna, and Max:

```text
Anna ──1000──► Igor
Max  ──1000──► Igor
```

All amounts are represented in the smallest currency unit, so calculations never lose money through floating-point rounding.

## 🌐 API

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

## 🚀 Quick Start

```bash
go test ./...
go run ./cmd/splitmate

curl -X POST localhost:8080/ledgers/1/expenses \
  -H 'Content-Type: application/json' \
  -d '{"description":"Dinner","amount":3000,"paid_by":1,"participants":[1,2,3]}'

curl localhost:8080/ledgers/1/settlement
```

Start the complete environment:

```bash
docker compose -f deployments/compose.yml up --build
```

## 🏗️ Project Layout

```text
cmd/splitmate       application entry point
internal/model      accounting entities
internal/debts      settlement algorithm
internal/repository PostgreSQL persistence
internal/transport  HTTP handlers
configs             environment template
deployments         Docker Compose environment
```

## 🧪 Quality Checks

```bash
make test
make lint
make build
```

## 🐍 Previous Implementation

The original Telegram-oriented Python version is preserved in [SplitMate-python](https://github.com/igorynos/SplitMate-python).
