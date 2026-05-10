# Import omni-commerce packages into pkg/ (not go module dependencies)

`sqlcraft`, `dafi`, `di`, `httpresponse`, `validation`, `null`, and `shared_domain` are copied from `omni-commerce/backend/pkg/` into this repo's `pkg/` directory, then rewritten to remove `samber/oops` and replace `echo` with `net/http`.

## Decision

Copy source files into `pkg/` rather than depending on `omni-commerce/backend/...` as a Go module. Each package is then rewritten:
- `httpresponse`: remove `echo` dependency, rewrite for `net/http`
- `sqlcraft`, `dafi`, `validation`: remove `samber/oops`, keep core logic
- `database`: written from scratch (no oops, no omni_errors)
- `di`, `null`, `shared_domain`: minimal changes (already clean)

## Why not go module dependency

- `omni-commerce` is a separate product with its own lifecycle. Changes there could break Budget Forge unexpectedly.
- Copying lets us remove `samber/oops` everywhere — a cross-cutting change that would be awkward as patches.
- The packages are small (~a few hundred lines each) and stable — maintenance burden of forks is low.
- Eliminates the risk of `omni-commerce`'s transitive dependencies (OpenTelemetry, etc.) leaking in unintentionally.

## Consequences

- Must manually port bug fixes from omni-commerce (if any).
- Initial import is a one-time copy; after that, packages evolve independently.
- `database` package is written from scratch — no shared history with omni-commerce.
