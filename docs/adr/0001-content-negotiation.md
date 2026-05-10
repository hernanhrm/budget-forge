# Content Negotiation: HTML (Datastar) + JSON (RFC 9457 Problem Details)

Same endpoint returns HTML fragments (Datastar SSE) or JSON (RFC 9457 Problem Details) based on `Accept` and Datastar headers. No separate `/api/` prefix — content negotiation, not URL partitioning.

## Decision

Middleware inspects incoming headers and sets a context value indicating response format: `"html"`, `"json"`, or `"full_page"` (no header → browser navigation). Handlers check the context value and render accordingly. Errors always produce Problem Details JSON when the `Accept` header requests JSON; otherwise they produce Datastar SSE fragments with inline error messaging. Middleware handles the fallback: if the handler panics or encounters an error before rendering, Problem Details is returned when safe.

## Why not pure Datastar/HTML

- Future mobile app or external tools need machine-readable errors.
- RFC 9457 Problem Details is a standard the ecosystem already understands.
- The app is hypermedia-first, but JSON errors are a pragmatic escape hatch for non-browser consumers.

## Why not separate `/api/` routes

- `accept` header and `datastar-request` header are both present so content negotiation is cleaner than URL prefixes.
- Content negotiation keeps the route tree flat and avoids duplication.
- Same handler, same validation, same business logic — only rendering forks.
