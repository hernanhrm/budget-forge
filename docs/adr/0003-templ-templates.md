# templ for templates (not html/template)

Use `templ` — a compile-time type-safe Go code generator for HTML templates — instead of Go's standard `html/template`.

## Decision

All HTML rendering uses `templ` components. Components are Go functions generated from `.templ` files via `templ generate`. Shared components live in `pkg/ui/`; module-specific components live in `internal/{module}/ui/`. Handlers pass domain types directly to components (no view-model mapping).

## Why not html/template

- `html/template` is runtime-parsed with no type safety — a typo in a template variable name is a runtime error.
- `templ` catches type mismatches at compile time. If a budget view struct changes, templates that reference renamed fields fail the build.
- `templ` components compose as Go functions with typed parameters — IDE autocomplete, refactoring, and navigation all work.
- Datastar SSE fragments are natural in `templ`: render a component to a buffer, wrap in SSE frame.

## Why not a JS framework or WASM

- The app is hypermedia-driven (Datastar). Server-rendered templates are the right abstraction.
- `templ` keeps all rendering in Go, aligned with the "server owns all state" philosophy.

## Consequences

- Build step required: `templ generate` must run before `go build`.
- Developers need the `templ` CLI and LSP plugin for the best experience.
- Switch from `html/template` means no `{{define}}` / `{{template}}` inheritance — components compose via Go function calls instead.
