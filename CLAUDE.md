# Repository Guidelines

## Overview

Cadence: Go modular music band repertoire management system. Exposes OpenAPI public API. PostgreSQL 16 for persistence,
Redis for session storage.

Domain: users create/join bands (via invite code), bands have members (owner/member roles), tracks (title, artist,
optional duration/tempo/key/notes) and setlists (grouping tracks into a set).

## Coding Style & Naming Conventions

- Language: Go 1.26
- Indentation: tabs (enforced by `gofmt`)
- Imports: formatted `goimports`; local module prefix `cadence/`, strict order enforced with `gci`. Import order:
    - Standard library
    - External packages (starting with `github.com/` or others)
    - Local packages (starting with `cadence/`)
- Linting: `golangci-lint` with config in `.golangci.yml`
- Naming: standard Go conventions — `PascalCase` exported, `camelCase` unexported; repository
  types use `Repository` suffix, services use `Service`
- Code declarations:
    - Types, structs, funcs, receivers sorted Public (top) → private (bottom)
    - Public/private constants, top-level var always at top
    - Constructors (`func NewService()` etc) always at top, after consts and top-level vars

## Build & Development Commands

Never plain `go build` or `golangci-lint run` — always mise.

```shell
# Generate OpenAPI, build service binary and lint project
mise run

# Generate API Code
mise run generate

# Build only service binary
mise run build

# Lint
mise run lint

# Unit tests
mise run test:unit

# Lint + test
mise run check

# End-to-end tests (spins up real Postgres/Redis via testcontainers)
mise run test:e2e

# Scaffold a new migration pair
bin/create-db-migration <name>
```

`compose.yml` runs local Postgres + Redis (+ app) for manual testing; copy `compose.override.example.yml` to override.

## Tests

E2E tests under `test/e2e/`, `e2e` build tag — `go test ./...` and `mise run test:unit` skip them;
run explicit via `mise run test:e2e` (needs running Docker daemon).

Unit tests (testify) colocated with code (`domain/`, `common/`). E2E: one full-journey test
(`test/e2e/user_journey_test.go`) covering register/login/band/track/setlist flows against real Postgres/Redis
via testcontainers — extend that flow with new assertions rather than adding separate E2E test files.

## Architecture

### Module structure

Main module under `pkg/cadence/`:

```
pkg/cadence/
  app/
    service/     - write-side application services (Create/Update/Remove), one per aggregate
    query/       - read-only query-service interfaces (CQRS-style split from service/)
  domain/        - domain entities
  infrastructure/
    transport/   - ogen server, API handlers, auth security handler, domain<->API conversion
    persistence/postgresql/
      repo/      - hand-written SQL writes (sqlx-style), soft delete, advisory-lock transactions
      query/     - hand-written SQL reads for app/query interfaces
```

### Key packages

- `pkg/common/` — shared utilities: `auth` (context helpers), `jsonlog`/`log` (structured logging), `maybe`
  (3-state optional type for PATCH semantics), `ogenerrors` (app error taxonomy: not_found, invalid_input,
  permission_denied, already_exists, operation_rejected), `ogenmiddleware` (request logging), `postgresql`
  (client/DSN/migrator), `redis` (session store client), `transactional` (Unit-of-Work with optional
  Postgres advisory-lock support), `uuid` (UUIDv7)

### CQRS split

Write and read sides fully separate stacks, one per aggregate, never share models:

- Write side: `app/service` (business logic, mutates via `domain` entities) ->
  `infrastructure/persistence/postgresql/repo`
  (`domain.XRepository` impl, `sqlx`-style, returns/accepts domain entities, e.g. `domain.LoadBand(...)`).
- Read side: `app/query` (`XQueryService` interfaces) -> `infrastructure/persistence/postgresql/query`
  (hand-written SQL, `SelectContext`/`GetContext` straight into flat DTOs like `TrackListItem`, `UserSetlistListItem`).
  Read-side code never imports `domain` — no entities, no invariants, just projections for API responses.
- Both sides get own `authorized` decorator (see below), own constructor wired independently
  in `cmd/cadence/dependencycontainer.go` (e.g. `trackService` and `trackQueryService` separate objects,
  both backed by same `track` table but through different SQL and different Go types).
- Adding new read model: define DTO + interface method in `app/query`, hand-write SQL in
  `infrastructure/persistence/postgresql/query`, add `authorized` pass-through/`requireMember` check,
  wire constructor. Don't reuse `domain` types or `repo/` queries for read paths, even if SQL would
  look similar — keeps read side free to diverge (joins, aggregates, search) without touching write-side
  invariants.

### Authorized services pattern

`app/service/authorized` and `app/query/authorized` decorate plain `app/service`/`app/query`
implementations with request-level authorization, keeping inner services free of auth concerns:

- Each authorized type wraps inner `service.XService`/`query.XQueryService` plus a
  `transactional.Executor[app.RepoProvider]`. Wired for every aggregate in `cmd/cadence/dependencycontainer.go`.
- `requesterIDFromContext(ctx)` pulls caller's user ID from session auth
  (`pkg/common/auth`), returns `permission_denied` if unauthenticated.
- `requireMember(ctx, executor, bandID, requesterID)` (in `app/service/authorized/band.go` and mirrored in
  `app/query/authorized`) loads band and checks `band.HasMember` before delegating — use for any
  op scoped to "any band member may do this".
- Owner-only operations (e.g. `Band.TransferOwnership`, `Band.RegenerateInviteCode`) enforce `IsOwner` inside
  domain method itself, return `domain.ErrNotBandOwner`; authorized wrapper just passes call
  through — don't duplicate owner check at authorized layer.
- Ops implicitly scoped to caller (`ListUserTracks`, `ListUserSetlists`, `JoinByInviteCode` — keyed by
  `userID` or invite code, not band membership) also pass through unchecked.
- New service/query methods on band-scoped aggregates: same split — "is requester
  even allowed near this band" in authorized wrapper via `requireMember`, finer-grained role checks
  (owner vs member) in domain type.

## Databases

### Migrations

Migrations live in `data/migrations/` as timestamped SQL file pairs:

```
20240101120000_<name>.up.sql
20240101120000_<name>.down.sql
```

- Use `bin/create-db-migration <name>` to scaffold new pair.
- One query per file only.
- Timestamp at start must be unique — add `sleep 1` calls when scaffolding several migrations in a row.

## Routing

- `gorilla/mux` router
- Public API: `/api/*` registered via `PathPrefix`, CORS middleware (`CADENCE_CORS_ALLOWED_ORIGINS`)
- Auth: session-token based (`SessionAuth` security scheme backed by Redis), see `pkg/common/auth`
- Health: `/resilience/live` (liveness, process up), `/resilience/ready` (readiness, checks DB/Redis)
- OpenAPI spec: `api/server/publicapi.yml`; generated client/server code lands in `api/server/publicapi/`
  (never hand-edit, regenerate via `mise run generate`)

## Constraints

- Never edit existing migration files — only scaffold new ones via `bin/create-db-migration`. Modifying an
  already-applied migration corrupts DB state in deployed environments.
- No hardcoded credentials or secrets — all config must flow through `CADENCE_`-prefixed environment variables.
  Nothing sensitive in source code.
- API changes must be backward-compatible — don't remove or rename existing fields; add new fields instead.
- Deletes are soft (`deleted_at`/`deleted_by` columns, filtered out of reads) — never hard-`DELETE` domain rows.
- Mutations racing on shared aggregate (band/track update or remove) take Postgres advisory lock via
  `transactional.Executor` — follow existing `repo/transactionfactory.go` pattern for new mutating operations.

## Definition of Done

Task done only when all below pass without errors:

1. Full build, unit tests and lint — run:

```shell
mise run
```

2. E2E tests, if touching public API behavior, transport wiring, or persistence — run:

```shell
mise run test:e2e
```
