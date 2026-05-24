A RESTful budget tracking API built with Go and PostgreSQL. Manage accounts, log income and expenses, and query balances — all over HTTP with no frameworks, just Go's standard library.

---

## Features

- Create and manage multiple accounts (cash, bank, credit)
- Log income and expense transactions per account
- Query account summaries with total income, total expenses, and current balance
- Input validation with proper HTTP error responses
- PostgreSQL persistence with auto-created tables on startup

---

## Tech Stack

- **Go** — standard library `net/http` only, no frameworks
- **PostgreSQL** — relational storage with foreign key constraints
- **godotenv** — environment variable loading from `.env`

---

## Installation

**Requirements:** Go 1.22+, PostgreSQL 15+

```bash
git clone https://github.com/XU4NLOC/go_http-rest-api.git
cd go_http-rest-api
```

Set up your database:

```bash
psql postgres
```

```sql
CREATE DATABASE budgettracker;
CREATE USER budgetuser WITH PASSWORD 'yourpassword';
GRANT ALL PRIVILEGES ON DATABASE budgettracker TO budgetuser;
\c budgettracker
GRANT ALL ON SCHEMA public TO budgetuser;
\q
```

Create your `.env` file:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=budgetuser
DB_PASSWORD=yourpassword
DB_NAME=budgettracker
```

Run the server:

```bash
go run main.go
# Server running on http://localhost:8080
```

Tables are created automatically on first run.

---

## API Reference

### Accounts

| Method | Endpoint | Description |
|---|---|---|
| GET | `/accounts` | List all accounts |
| POST | `/accounts` | Create an account |
| GET | `/accounts/{id}` | Get one account |
| DELETE | `/accounts/{id}` | Delete an account |

### Transactions

| Method | Endpoint | Description |
|---|---|---|
| GET | `/accounts/{id}/transactions` | List transactions for an account |
| POST | `/accounts/{id}/transactions` | Add a transaction |
| GET | `/accounts/{id}/summary` | Get balance and totals |

---

## Usage Examples

```bash
# Create an account
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{"name": "My Bank", "type": "bank"}'

# Add an income transaction
curl -X POST http://localhost:8080/accounts/1/transactions \
  -H "Content-Type: application/json" \
  -d '{"amount": 5000000, "type": "income", "description": "Salary"}'

# Add an expense transaction
curl -X POST http://localhost:8080/accounts/1/transactions \
  -H "Content-Type: application/json" \
  -d '{"amount": 200000, "type": "expense", "description": "Groceries"}'

# Get account summary
curl http://localhost:8080/accounts/1/summary
```

**Summary response:**
```json
{
  "account_id": 1,
  "total_income": 5000000,
  "total_expenses": 200000,
  "balance": 4800000
}
```

---

## Project Structure

```
budgettracker/
├── main.go                  # Entry point, server startup
├── .env                     # Environment variables (not committed)
├── db/
│   └── db.go                # PostgreSQL connection + table creation
├── models/
│   ├── account.go           # Account struct
│   └── transaction.go       # Transaction + Summary structs
├── handlers/
│   ├── account.go           # Account HTTP handlers
│   └── transaction.go       # Transaction + Summary HTTP handlers
└── router/
    └── router.go            # Route definitions
```

---

## Validation

The API enforces the following rules:

- `name` and `type` are required when creating an account
- Account `type` must be one of: `cash`, `bank`, `credit`
- Transaction `amount` must be greater than 0
- Transaction `type` must be `income` or `expense`
- Transactions cannot be added to non-existent accounts

---

## Motivation

Second project in my Go learning path, focused on backend engineering. This project covers database/sql, connection management, SQL aggregations, JSON encoding, manual HTTP routing, and input validation — all without reaching for a framework.

Previous: [taskmanager](https://github.com/yourusername/taskmanager) — CLI task manager built with Cobra.
Next: Port scanner using goroutines and channels.
