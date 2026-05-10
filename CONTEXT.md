# Architecture

- **Frontend:** Server-side rendered HTML with Go + Datastar (hypermedia-driven, not a SPA)
- **Templates:** `templ` — compile-time type-safe code generation. Not `html/template`.
- **API:** Content negotiation — same endpoint returns HTML fragments (Datastar SSE) or JSON (RFC 9457 Problem Details) based on `Accept` header.
- **Database:** PostgreSQL via `pgx/v5`, abstracted behind `shared_domain.DatabasePort`.
- **Repo structure:** Monorepo with vertical slices (hexagonal/ports-and-adapters).

```
budget-forge/
├── cmd/              # Service entry points (web, migrate)
│   ├── web/          # HTTP server binary
│   │   ├── main.go
│   │   └── routes/   # Route registration
│   └── migrate/      # Migration runner binary
│       └── main.go
├── internal/         # Domain modules (vertical slice)
│   ├── auth/         # Signup, login, logout, sessions
│   ├── user/         # Profile, onboarding, default categories
│   ├── budget/       # Month lifecycle: plan, rollforward, close, reallocate
│   ├── category/     # CRUD for categories and groups
│   ├── transaction/  # CRUD for transactions
│   ├── account/      # CRUD for accounts, balance tracking
│   └── transfer/     # Atomic account-to-account transfers
├── pkg/              # Shared packages (imported from omni-commerce, rewritten)
│   ├── server/       # HTTP server setup (mux, middleware, graceful shutdown)
│   ├── database/     # pgx pool, ContextRouter, PgWorkUnit (written from scratch)
│   ├── shared_domain/# Core interfaces: DatabasePort, Logger, WorkUnit, List[T]
│   ├── di/           # Thin wrapper over samber/do/v2
│   ├── sqlcraft/     # Fluent SQL query builder (rewritten, no oops)
│   ├── dafi/         # Dynamic filtering/pagination/sorting (rewritten, no oops)
│   ├── httpresponse/ # RFC 9457 Problem Details (rewritten for net/http)
│   ├── validation/   # Struct validation (ozzo, rewritten, no oops)
│   ├── null/         # Nullable types (guregu/null/v6 wrapper)
│   └── ui/           # Shared templ components (layout, button, shell)
├── docs/
│   └── adr/          # Architecture Decision Records
├── DESIGN.md
├── PRD.md
└── AGENTS.md
```

## Tech Stack

| Concern | Choice | Notes |
|---------|--------|-------|
| HTTP server | `net/http` (standard) | No Echo, no frameworks |
| Templates | `templ` | Compile-time, type-safe |
| Hypermedia | Datastar (official Go SDK) | SSE fragments |
| DI / wiring | `samber/do/v2` via `pkg/di` | Per-module `di.go`, per-request scope |
| Database driver | `pgx/v5` via `shared_domain.DatabasePort` | ContextRouter for tx routing |
| SQL builder | `pkg/sqlcraft` | Fluent builder, imported from omni-commerce |
| Query parsing | `pkg/dafi` | Filter/sort/pagination from HTTP params |
| Migrations | `golang-migrate/migrate` | SQL files, run via `cmd/migrate` |
| Logging | `slog` (standard) | Injected as `shared_domain.Logger` interface |
| Config | `os.Getenv` + mise.toml | No godotenv |
| Auth | bcrypt + `gorilla/sessions` + `gorilla/csrf` | HTTP-only cookies |
| Validation | `pkg/validation` (ozzo, without oops) | Struct validation on DTOs |
| Nullable types | `pkg/null` (guregu/null/v6) | Not pointer types |
| Error handling | Sentinel errors + `fmt.Errorf("%w")` wrapping | No samber/oops |
| Error → status | Callback function per module | Maps sentinel → HTTP code |
| CSS | Tailwind CSS | Compiled, embedded via `//go:embed` |
| Testing | Go `testing` + `testify/assert` + `testcontainers-go` | Integration against real PG |

## Module File Pattern

Each module follows the omni-commerce hexagonal pattern:

```
internal/{module}/
├── port.go          # Interfaces (Service, Repository)
├── domain.go        # Domain types, DTOs, validation
├── service.go       # Business logic
├── repository.go    # PostgreSQL implementation (sqlcraft)
├── handler.go       # HTTP handler (net/http)
├── di.go            # Wire into samber/do container
├── ui/              # Templ components specific to this module
├── *_test.go        # Tests co-located with source
└── wire.go          # (if needed) Module-specific wiring
```

- **Handler** checks content negotiation, renders templ or calls `httpresponse`.
- **Service** validates via `ozzo-validation`, delegates to repository, wraps errors.
- **Repository** uses `sqlcraft` to build queries, executes via `shared_domain.DatabasePort`.
- **DI** registers `Repository`, `Service`, `Handler` via `di.Provide()`.

## Content Negotiation

Middleware sets a context value based on `Accept` header:
- `Accept: application/json` → response format is JSON (Problem Details)
- Datastar header present → response format is HTML (templ + SSE framing)
- Neither → default to HTML (full page render)

Handlers check the context value and render accordingly. Error mapping uses a callback:

```go
type ProblemMapper func(err error) (status int, ok bool)
```

## Middleware Chain

1. Recovery (panic → 500 Problem Detail)
2. RequestID (unique ID per request → slog)
3. Logging (method, path, status, duration via slog)
4. CORS
5. Session (extract cookie, load session, attach user to context)
6. Content negotiation (detect Accept/Datastar header → context value)
7. Auth (require user on protected routes)
8. Handler

## Entry Points

- `cmd/web/main.go` — HTTP server, wires DI container, mounts routes, graceful shutdown
- `cmd/migrate/main.go` — Runs `golang-migrate` migrations against the database

## UI / Navigation

Four primary pages with a bottom tab bar on mobile and a top nav on desktop:

1. **Budget** (home) — the envelope plan for the current month. Default landing page.
2. **Transactions** — list of all transactions for the current month.
3. **Accounts** — list of accounts with current balances.
4. **Settings** — profile, category/group management, month-end actions.

A floating action button (+) provides the fast path to "Add Transaction." Datastar enables server-rendered HTML pages with partial updates for interactivity (reallocating funds, adding transactions inline).

## Domain Glossary

### Budget Cycle
The standard period for planning and tracking. **Monthly** — calendar month (1st to last day). Each month gets its own budget plan. Paycheck-aligned cycles are deferred to post-MVP.

### Zero-Based Budgeting
Every dollar of expected income is assigned to a specific purpose before the month begins. Income − Planned Expenses = $0 (with any leftover explicitly assigned to savings, buffer, or a "fun money" envelope).

### Reallocation
Moving planned or available funds from one Category to another to maintain the zero-based invariant. When actual spending exceeds a category's plan, the user must explicitly reallocate from another category. The app does not auto-compensate.

### Account
A real-world financial container (e.g., Chase Checking, Savings, Credit Card). Each Account has a current balance tracked independently from Categories. Transactions are optionally linked to an Account.

### Income
Money expected to arrive during the budget cycle (e.g., salary, freelance payments). Represented as Categories with direction = credit, each with a planned amount per month. Actual deposits are tracked as Transactions in these income categories. The budget view presents the aggregate total for allocation.

### Income Tracking
Each income Category has a **planned** amount (set at month start) and **actual** deposits (recorded as they arrive via Transactions). If actual income falls short of planned, the app surfaces a shortfall and prompts reallocation to maintain the zero-based invariant.

### Category / Envelope
A named bucket that receives a portion of monthly income or records an income source. Each Category has a **direction** (debit = expense, credit = income), a **planned** amount (set at month start), and an **actual** amount (sum of transactions during the month). Examples: "Rent" (debit), "Groceries" (debit), "Salary" (credit), "Emergency Fund" (debit).

### Month Rollforward
When a new budget cycle begins, the app creates a draft by copying the previous month's Categories, Groups, and planned amounts. The user reviews and adjusts before the month "starts." This is the primary onboarding path for an existing user — the budget is never empty.

### Closed Month
A month that has ended and been reviewed by the user. Once closed, its transactions and allocations become immutable. The user may view historical months but cannot edit them.

### User
An individual with their own set of budgets, categories, accounts, and transactions. Each user's data is fully isolated.

### Default Categories
On signup, the user is provisioned with a sensible set of pre-populated categories (e.g., Rent, Groceries, Utilities, Dining Out, Transportation, Savings) to provide a warm start. The user may edit, rename, or delete these.

### Surplus Handling
When a month ends with unspent funds in a Category (actual < planned), the user decides what to do with the remainder: carry it forward as extra available balance in that Category, move it to Savings, or reallocate it elsewhere. The app does not auto-carry surpluses.

### Month Review
The process a user performs before closing a month. The app presents a summary: total income vs planned, total spent vs planned, categories with surplus, and categories that went over. The user reviews, decides what to do with surpluses, and triggers Close Month.

### Category Group
An optional user-defined label that clusters related Categories for rollup display (e.g., "Needs", "Wants", "Savings"). Groups do not enforce rules — they are for presentation only. A Category may belong to zero or one Group.

### Transaction
A real-world spending or income event. Linked to a Category and optionally an Account. A positive amount represents income or a refund; a negative amount represents an expense. If linked to an Account, also updates that Account's current balance. Account linkage is not required for the transaction to affect the category balance.

### Transfer
An account-to-account movement of funds that affects two Accounts but zero Categories. Recorded as two linked Transactions (debit from source, credit to destination) executed atomically. Does not affect category balances because it is just moving money between containers.

## Example Dialogue

> **Dev:** "When a user closes a month and has a surplus in Groceries, does the app auto-carry it forward?"
> **Domain expert:** "No — surplus handling is explicit. The user decides during Month Review: carry forward, move to Savings, or reallocate elsewhere."
>
> **Dev:** "What happens if actual income is less than planned — does the app block transactions?"
> **Domain expert:** "The app warns about the shortfall but does not block. Users enter transactions manually; enforcing zero-based at the UI level would be friction. Education over enforcement."
>
> **Dev:** "A user transfers $500 from Checking to Savings. Does that affect their Groceries Available balance?"
> **Domain expert:** "No — a Transfer moves money between Accounts only. It creates two linked Transactions with no Category. Category balances are unchanged."

## Flagged Ambiguities

- "account" was used to mean both **Account** (financial container) and **User** (login identity) — resolved: these are distinct concepts with separate modules.
