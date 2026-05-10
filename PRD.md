# Budget Forge — MVP PRD

## Problem Statement

People struggle to plan where their money should go each month. Most budgeting tools are either too simple (just expense tracking) or too complex (requiring bank connections, double-entry accounting, or steep learning curves). Users need an easy way to assign every expected dollar to a specific purpose before the month begins — a zero-based envelope budget — without friction, without bank sync, and without rebuilding their budget from scratch every month.

## Solution

Budget Forge is a web-based zero-based budgeting application that helps users plan their monthly spending by assigning every dollar of expected income to named categories (envelopes) before the month begins. It is server-rendered with Go and Datastar, uses PostgreSQL for persistence, and follows a mobile-first responsive design. Users manually enter transactions, track account balances independently, and maintain their budget plan through simple reallocation when spending exceeds a category's allocation.

## User Stories

1. As a new user, I want to sign up with an email and password, so that I can create my own private budget.
2. As a new user, I want my first budget to be pre-populated with sensible default categories, so that I don't start from a blank slate.
3. As a user, I want to plan my expected income for the month, so that I know how much I have to allocate.
4. As a user, I want to allocate my expected income across spending categories at the start of the month, so that every dollar has a job.
5. As a user, I want to see my budget plan as my home screen, so that checking my available balances is the fastest path.
6. As a user, I want to log a spending transaction quickly, so that I can track where my money went.
7. As a user, I want to optionally link a transaction to a bank account, so that my account balances stay accurate.
8. As a user, I want to see three numbers for every category — Planned, Actual, and Available — so that I know how much I planned, how much I've spent, and how much is left.
9. As a user, I want to move money from one category to another when I overspend, so that my budget stays zero-based.
10. As a user, I want to see a warning when my total planned allocations don't match my total planned income, so that I know my budget is out of balance.
11. As a user, I want to create custom categories and category groups, so that my budget matches my life.
12. As a user, I want each new month to start as a copy of the previous month's categories and planned amounts, so that I don't rebuild my budget every month.
13. As a user, I want to review my month before closing it, so that I can see what I planned vs what I actually spent.
14. As a user, I want to decide what to do with unspent money in a category at month-end, so that I can carry it forward, move it to savings, or reallocate it.
15. As a user, I want closed months to be locked and immutable, so that my historical budget record stays accurate.
16. As a user, I want to view a list of all my transactions for the current month, so that I can review or correct entries.
17. As a user, I want to create and manage bank accounts with current balances, so that I can track my real-world financial containers.
18. As a user, I want to transfer money between accounts, so that my account balances reflect real-world movements without affecting my category budgets.
19. As a user, I want to edit or delete a transaction in the current month, so that I can fix mistakes.
20. As a user, I want to add income transactions when I receive a paycheck, so that my actual income tracks against my planned income.
21. As a user, I want to see a shortfall warning when my actual income is less than my planned income, so that I know I need to reallocate.
22. As a user, I want to manage my categories and groups from a settings page, so that I can rename, reorder, or delete them.
23. As a user, I want to access the app on my phone with a bottom tab bar, so that I can check my budget on the go.
24. As a user, I want to see the same app on desktop with a top navigation bar, so that I can manage my budget comfortably on a larger screen.
25. As a user, I want a floating action button to quickly add a transaction from any page, so that logging spending is always one tap away.
26. As a user, I want to see income categories and spending categories clearly distinguished, so that I don't confuse money coming in with money going out.
27. As a user, I want to add a transaction without linking it to an account, so that I can track spending even for cash or untracked accounts.
28. As a user, I want to change a category's planned amount during the month, so that I can adjust my plan when life changes.
29. As a user, I want to see my total planned income, total planned spending, and the difference at the top of my budget page, so that I can verify my budget is zero-based.
30. As a user, I want to view past closed months in read-only mode, so that I can review my budgeting history.

## Implementation Decisions

### Modules

The following modules will be built in `internal/`, each as a vertical slice containing its own domain logic, storage interface, HTTP handlers, and templates:

1. **Auth Module** (`internal/auth/`)
   - Handles user signup, login, logout, and session management.
   - Interface: `AuthService` with methods `Register(email, password) (User, error)`, `Login(email, password) (Session, error)`, `Logout(sessionID) error`.
   - Uses bcrypt for password hashing and secure HTTP-only cookies for sessions.
   - Depends on: `pkg/db` for persistence.

2. **User/Onboarding Module** (`internal/user/`)
   - Handles user profile and the warm-start onboarding flow.
   - On signup, provisions default categories for the user's first month.
   - Interface: `UserService` with methods `GetProfile(userID)`, `ProvisionDefaults(userID) error`.
   - Deep module: the onboarding logic (default category set, initial month creation) is complex but exposed through a single `ProvisionDefaults` call.
   - Depends on: Category Module, Budget Module.

3. **Budget/Month Module** (`internal/budget/`)
   - The core orchestration module for the monthly budget lifecycle.
   - Interface: `BudgetService` with methods:
     - `GetMonth(userID, year, month) (MonthView, error)`
     - `RollforwardMonth(userID, fromMonth, toMonth) error`
     - `CloseMonth(userID, year, month, surplusDecisions) error`
     - `UpdatePlan(userID, monthCategoryID, plannedAmount) error`
   - Deep module: encapsulates all zero-based math, rollforward logic, month closing, and surplus handling.
   - `MonthView` aggregates categories with their planned/actual/available balances.
   - Depends on: Category Module, Transaction Module.

4. **Category Module** (`internal/category/`)
   - CRUD for categories and category groups.
   - Interface: `CategoryService` with methods:
     - `Create(userID, name, direction, groupID) (Category, error)`
     - `Update(userID, categoryID, name, groupID) error`
     - `Delete(userID, categoryID) error`
     - `List(userID) ([]Category, error)`
     - `CreateGroup(userID, name) (Group, error)`
   - Categories have `direction`: `debit` (expense) or `credit` (income).
   - Depends on: `pkg/db`.

5. **Transaction Module** (`internal/transaction/`)
   - CRUD for transactions. Transactions are linked to a Category and optionally an Account.
   - Interface: `TransactionService` with methods:
     - `Create(userID, date, amount, categoryID, accountID, description) (Transaction, error)`
     - `Update(userID, txID, ...) error`
     - `Delete(userID, txID) error`
     - `ListByMonth(userID, year, month) ([]Transaction, error)`
     - `ListByCategory(userID, categoryID, year, month) ([]Transaction, error)`
   - Positive amount = income/refund; negative amount = expense.
   - Depends on: `pkg/db`, Account Module (for balance updates).

6. **Account Module** (`internal/account/`)
   - CRUD for accounts and balance tracking.
   - Interface: `AccountService` with methods:
     - `Create(userID, name, initialBalance) (Account, error)`
     - `UpdateBalance(userID, accountID, delta) error`
     - `List(userID) ([]Account, error)`
   - Depends on: `pkg/db`.

7. **Transfer Module** (`internal/transfer/`)
   - Handles two-phase account-to-account transfers.
   - Interface: `TransferService` with method `Execute(userID, fromAccountID, toAccountID, amount, date, description) error`.
   - Atomically creates two linked transactions (debit from source, credit to destination) with no category linkage.
   - Deep module: transfer atomicity is critical and should be encapsulated behind a simple interface.
   - Depends on: Transaction Module, Account Module.

8. **Web/HTTP Module** (`internal/web/`)
   - Routes, Datastar handlers, and HTML templates.
   - Not a deep module — it's a thin adapter layer that delegates to domain services.
   - Uses Go's `html/template` with Datastar for partial-page updates.
   - Pages: Budget (home), Transactions, Accounts, Settings, Login, Signup.

### Schema Decisions

- `users`: `id`, `email`, `password_hash`, `created_at`
- `sessions`: `id`, `user_id`, `token`, `expires_at`
- `categories`: `id`, `user_id`, `name`, `direction` (debit/credit), `group_id`, `sort_order`, `created_at`
- `category_groups`: `id`, `user_id`, `name`, `sort_order`
- `months`: `id`, `user_id`, `year`, `month`, `status` (open/closed), `closed_at`
- `month_categories`: `id`, `month_id`, `category_id`, `planned_amount` (positive integer, sign inferred from category direction)
- `accounts`: `id`, `user_id`, `name`, `current_balance`, `created_at`
- `transactions`: `id`, `user_id`, `date`, `amount` (positive = income, negative = expense), `category_id`, `account_id` (nullable), `description`, `transfer_pair_id` (nullable), `created_at`

### API Contracts

All endpoints are server-rendered HTML via Datastar. JSON APIs are not used for the MVP — the server returns HTML fragments for partial updates and full pages for navigation.

Key interactions:
- **Reallocation**: POST `/budget/reallocate` with `from_month_category_id`, `to_month_category_id`, `amount`. Returns updated budget fragment.
- **Add Transaction**: POST `/transactions` with form data. Returns updated transaction list + budget summary fragment.
- **Close Month**: POST `/month/close` with surplus decisions map. Returns confirmation or redirects to next month.

### Architecture Decisions

- **Monorepo with vertical slices**: Each domain module in `internal/` owns its own logic, storage queries, and handlers. This keeps the codebase navigable and lets agents work on one slice at a time.
- **Hypermedia-driven (Datastar)**: No JSON APIs, no client-side state management. The server owns all state. Datastar provides SPA-like interactivity via server-pushed HTML fragments.
- **Forecast budgeting**: Users plan full expected income at month start. The budget is zero-based when total planned income equals total planned expenses.
- **Immutable closed months**: Once a month is closed, its data is read-only. Historical accuracy is prioritized over edit flexibility.
- **Manual entry only**: No bank sync for MVP. Fast manual entry with smart defaults (today's date, last-used category).
- **Positive amounts with inferred sign**: Users enter positive numbers. The app applies sign based on category direction. This matches mental models.

## Testing Decisions

### What Makes a Good Test

Tests should verify **external behavior and invariants**, not implementation details. A good test asserts that given a set of inputs, a module produces the correct outputs or state changes — without caring about internal loops, variable names, or helper functions. Tests that mock database queries or assert on SQL strings are fragile and should be avoided in favor of integration-style tests against a real test database.

### Modules to Test

1. **Budget/Month Module** — Highest priority. Tests must cover:
   - Zero-based math: planned income = planned expenses → balanced; mismatch → unallocated warning.
   - Month rollforward: categories, groups, and planned amounts copied correctly.
   - Month closing: surplus handling decisions applied correctly; month becomes immutable.
   - Reallocation: moving planned amounts between categories preserves total.
   - Actual balance calculation: sum of transactions per category computed correctly.

2. **Transfer Module** — High priority. Tests must cover:
   - Atomic execution: both transactions created or neither.
   - Account balance updates: source decreased, destination increased by exact amount.
   - No category impact: category balances unchanged.
   - Failure rollback: if one side fails, no state changes.

3. **Transaction Module** — Medium priority. Tests must cover:
   - Amount sign conventions: positive increases category actual, negative decreases.
   - Optional account linkage: transaction without account does not affect account balance.
   - Account balance update: linked transaction updates account current_balance.
   - Month scoping: transactions filtered correctly by user + month.

4. **Category Module** — Medium priority. Tests must cover:
   - CRUD operations scoped to user.
   - Direction field behavior.
   - Group assignment and removal.

5. **Auth Module** — Medium priority. Tests must cover:
   - Password hashing (bcrypt rounds appropriate).
   - Session creation and validation.
   - Logout invalidates session.

### Prior Art

There is no existing test suite in the codebase. All tests will be net-new. Use Go's standard `testing` package with `testify/assert` for readability. Use `testcontainers-go` or a local ephemeral PostgreSQL database for integration tests against real SQL. Unit tests that don't touch the database can use in-memory stubs for the storage interface.

## Out of Scope

- Bank account sync or import (Plaid, CSV import, etc.)
- Recurring/scheduled transactions
- Split transactions (one transaction across multiple categories)
- Multiple budgets per user
- Paycheck-aligned budget cycles (bi-weekly, etc.)
- OAuth authentication
- Native mobile apps
- Push notifications or email alerts
- Reports, charts, or analytics beyond the basic budget view
- Multi-currency support
- Shared budgets between users
- Debt/savings goals or targets
- Attachments or receipt photos
- Budget templates beyond the single default set
- Automatic category suggestions or machine learning

## Further Notes

- The UI design system is documented in `DESIGN.md` and uses a Starbucks-inspired warm palette. The `ui.pen` file contains the design reference. All frontend work should align with these tokens.
- The floating action button (FAB) for "Add Transaction" should be visible on all four primary pages.
- The Budget page should show income categories at the top and spending categories below, with a summary bar showing total planned income, total planned spending, and the unallocated difference.
- Reallocation should be possible via inline controls on the Budget page (e.g., an "Adjust" button that opens a small form to move money between categories).
- Smart defaults for manual transaction entry: default date = today, default category = last used for this payee (if tracked), default account = last used.
- The first month for a new user is the current calendar month. Default planned amounts can be zero — the user fills them in during onboarding.
- The app does not enforce zero-based balance at the UI level (it warns but doesn't block). This reduces friction while still teaching the habit.
