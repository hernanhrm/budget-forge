# AGENTS.md

## Project

Zero-based monthly budgeting app. Server-rendered Go + Datastar (hypermedia, not SPA). PostgreSQL. Mobile-first.

## Key Docs

- `CONTEXT.md` — architecture, domain glossary (read this first)
- `PRD.md` — user stories, module definitions, schema, API contracts, testing strategy
- `DESIGN.md` — Starbucks-inspired design system tokens, colors, typography, components
- `ui.pen` — visual design reference (Pencil file)

## Tech Stack

- **Backend:** Go, `templ` (compile-time type-safe templates), Datastar for partial-page updates
- **Database:** PostgreSQL via `pgx/v5`, abstracted behind `shared_domain.DatabasePort`
- **Auth:** bcrypt + `gorilla/sessions` + `gorilla/csrf` + secure HTTP-only cookies (no OAuth for MVP)
- **Content negotiation:** Same endpoints return HTML fragments (Datastar SSE) or JSON (RFC 9457 Problem Details) based on `Accept` header. No JSON APIs — but JSON error responses exist.
- **No bank sync** — manual entry only for MVP
- **No oops** — sentinel errors with `fmt.Errorf("%w")` wrapping at each layer

## Architecture

Monorepo with vertical slices. Each module in `internal/` owns its domain logic, storage interface, HTTP handlers, and templates.

```
cmd/              # Service entry points
  web/            # HTTP server + routes
  migrate/        # DB migration runner
internal/         # Domain modules (vertical slice)
  auth/           # Signup, login, logout, sessions
  user/           # Profile, onboarding, default categories
  budget/         # Month lifecycle: plan, rollforward, close, reallocate
  category/       # CRUD for categories and groups
  transaction/    # CRUD for transactions
  account/        # CRUD for accounts, balance tracking
  transfer/       # Atomic account-to-account transfers
pkg/              # Shared packages
  server/         # HTTP server setup (mux, middleware)
  shared_domain/  # Core interfaces (DatabasePort, Logger, WorkUnit, List[T])
  database/       # pgx pool, ContextRouter, PgWorkUnit
  di/             # Thin wrapper over samber/do/v2
  sqlcraft/       # Fluent SQL query builder
  dafi/           # Dynamic filtering/pagination/sorting
  httpresponse/   # RFC 9457 Problem Details (net/http)
  validation/     # Struct validation (ozzo, without oops)
  null/           # Nullable types (guregu/null/v6)
  ui/             # Shared templ components (layout, button, shell)
```

## Key Conventions

- **Zero-based invariant:** Income − Planned Expenses = $0. App warns but does not block at UI level.
- **Amounts:** Positive = income/refund, negative = expense. Users enter positive numbers; sign is inferred from category direction (`debit`/`credit`).
- **Month lifecycle:** Open → draft → active → closed. Closed months are immutable.
- **Month rollforward:** New month copies previous month's categories, groups, and planned amounts.
- **Reallocation:** Explicit user action to move funds between categories. No auto-compensation.
- **Transfers:** Two linked transactions (debit + credit), atomic, no category impact.
- **Categories have direction:** `debit` (expense) or `credit` (income). Income categories appear at top of budget view.

## Schema (key tables)

`users`, `sessions`, `categories`, `category_groups`, `months`, `month_categories`, `accounts`, `transactions`

`month_categories` links a category to a month with a `planned_amount`. `transactions.amount` is signed. `transactions.transfer_pair_id` links transfer pairs. `transactions.account_id` is nullable (cash/untracked spending).

## API Pattern

All endpoints support content negotiation. Handlers check context value set by middleware:
- `Accept: application/json` → JSON (RFC 9457 Problem Details)
- `Datastar-Request` header → HTML fragments via SSE
- Neither → full HTML page (browser navigation)
- POST without JS → 303 redirect (PRG fallback)

Key routes:
- `POST /budget/reallocate` — move funds between categories
- `POST /transactions` — add transaction, returns updated list + budget summary
- `POST /month/close` — close month with surplus decisions

## Testing Strategy

Use Go `testing` + `testify/assert`. Integration tests against real PostgreSQL (testcontainers-go). Unit tests can use in-memory stubs for storage interfaces.

Priority order: Budget/Month > Transfer > Transaction > Category > Auth.

Test external behavior and invariants, not implementation details. Don't mock SQL queries.

## Design System

Starbucks-inspired warm palette. Key tokens:
- Page canvas: `#f2f0eb` (warm cream, not white)
- Cards: `#ffffff`, `12px` radius
- Primary CTA: `#00754A` (Green Accent), `50px` pill, `scale(0.95)` active
- Feature bands: `#1E3932` (House Green)
- Body text: `rgba(0,0,0,0.87)` (not pure black)

See `DESIGN.md` §9 "Agent Prompt Guide" for component-level specs and example prompts.

## Out of Scope (MVP)

Bank sync, recurring transactions, split transactions, multiple budgets per user, OAuth, native apps, reports/charts, multi-currency, shared budgets.
