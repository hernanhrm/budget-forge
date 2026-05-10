# Architecture

- **Frontend:** Server-side rendered HTML with Go + Datastar (hypermedia-driven, not a SPA)
- **Database:** PostgreSQL
- **Repo structure:** Monorepo

```
budget-forge/
├── cmd/              # Service entry points
├── internal/         # Application modules (vertical slice, clean architecture, optional hexagonal)
└── pkg/              # Utility packages (db connection, logger, config from .env)
```

# UI / Navigation

Four primary pages with a bottom tab bar on mobile and a top nav on desktop:

1. **Budget** (home) — the envelope plan for the current month. Default landing page.
2. **Transactions** — list of all transactions for the current month.
3. **Accounts** — list of accounts with current balances.
4. **Settings** — profile, category/group management, month-end actions.

A floating action button (+) provides the fast path to "Add Transaction." Datastar enables server-rendered HTML pages with partial updates for interactivity (reallocating funds, adding transactions inline).

# Domain Glossary

## Budget Cycle
The standard period for planning and tracking. **Monthly** — calendar month (1st to last day). Each month gets its own budget plan. Paycheck-aligned cycles are deferred to post-MVP.

## Zero-Based Budgeting
Every dollar of expected income is assigned to a specific purpose before the month begins. Income − Planned Expenses = $0 (with any leftover explicitly assigned to savings, buffer, or a "fun money" envelope).

## Reallocation
Moving planned or available funds from one Category to another to maintain the zero-based invariant. When actual spending exceeds a category's plan, the user must explicitly reallocate from another category. The app does not auto-compensate.

## Account
A real-world financial container (e.g., Chase Checking, Savings, Credit Card). Each Account has a current balance tracked independently from Categories. Transactions are optionally linked to an Account.

## Income
Money expected to arrive during the budget cycle (e.g., salary, freelance payments). Represented as Categories with direction = credit, each with a planned amount per month. Actual deposits are tracked as Transactions in these income categories. The budget view presents the aggregate total for allocation.

## Income Tracking
Each income Category has a **planned** amount (set at month start) and **actual** deposits (recorded as they arrive via Transactions). If actual income falls short of planned, the app surfaces a shortfall and prompts reallocation to maintain the zero-based invariant.

## Category / Envelope
A named bucket that receives a portion of monthly income or records an income source. Each Category has a **direction** (debit = expense, credit = income), a **planned** amount (set at month start), and an **actual** amount (sum of transactions during the month). Examples: "Rent" (debit), "Groceries" (debit), "Salary" (credit), "Emergency Fund" (debit).

## Month Rollforward
When a new budget cycle begins, the app creates a draft by copying the previous month's Categories, Groups, and planned amounts. The user reviews and adjusts before the month "starts." This is the primary onboarding path for an existing user — the budget is never empty.

## Closed Month
A month that has ended and been reviewed by the user. Once closed, its transactions and allocations become immutable. The user may view historical months but cannot edit them.

## User
An individual with their own set of budgets, categories, accounts, and transactions. Each user's data is fully isolated.

## Default Categories
On signup, the user is provisioned with a sensible set of pre-populated categories (e.g., Rent, Groceries, Utilities, Dining Out, Transportation, Savings) to provide a warm start. The user may edit, rename, or delete these.

## Surplus Handling
When a month ends with unspent funds in a Category (actual < planned), the user decides what to do with the remainder: carry it forward as extra available balance in that Category, move it to Savings, or reallocate it elsewhere. The app does not auto-carry surpluses.

## Month Review
The process a user performs before closing a month. The app presents a summary: total income vs planned, total spent vs planned, categories with surplus, and categories that went over. The user reviews, decides what to do with surpluses, and triggers Close Month.

## Category Group
An optional user-defined label that clusters related Categories for rollup display (e.g., "Needs", "Wants", "Savings"). Groups do not enforce rules — they are for presentation only. A Category may belong to zero or one Group.

## Transaction
A real-world spending or income event. Linked to a Category and optionally an Account. A positive amount represents income or a refund; a negative amount represents an expense. If linked to an Account, also updates that Account's current balance. Account linkage is not required for the transaction to affect the category balance.

## Transfer
An account-to-account movement of funds that affects two Accounts but zero Categories. Recorded as two linked Transactions (debit from source, credit to destination) executed atomically. Does not affect category balances because it is just moving money between containers.
