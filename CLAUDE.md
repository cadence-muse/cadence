# Repository Guidelines

## Overview

Cadence is a Go modular music band repertoire management system. It exposes OpenAPI public API. The service uses
PostgreSQL 16 for persistence and Redis for session storage.

Domain: users create/join bands (via invite code), bands have members (owner/member roles), tracks (title, artist,
optional duration/tempo/key/notes) and setlists (grouping tracks into a set).

## Coding Style & Naming Conventions

- Language: Go 1.26
- Indentation: tabs (enforced by `gofmt`)
- Imports: formatted `goimports`; local module prefix is `cadence/`, strict order is enforced with `gci`. Order of
  imports:
    - Standard library
    - External packages (starting with `github.com/` or others)
    - Local packages (starting with `cadence/`)
- Linting: `golangci-lint` with the config in `.golangci.yml`
- Naming: follow standard Go conventions - `PascalCase` for exported identifiers, `camelCase` for unexported; repository
  types use the `Repository` suffix, services use `Service`
- Code declarations:
    - Declarations of types, structs, funcs, receivers sorted from Public (at top) and private (at bottom)
    - Public and private constants, top-level var always at top
    - Constructors (funcs like `func NewService()`) always at top, after consts and top-level vars

## Build & Development Commands

Do not use plain `go build` or `golangci-lint run`, always use mise.

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

E2E tests live under `test/e2e/` behind the `e2e` build tag, so `go test ./...` and `mise run test:unit` skip them;
run them explicitly via `mise run test:e2e` (needs a running Docker daemon).

Unit tests (testify) are colocated with the code (`domain/`, `common/`). E2E is a single full-journey test
(`test/e2e/user_journey_test.go`) covering register/login/band/track/setlist flows against real Postgres/Redis
via testcontainers - extend that flow with new assertions rather than adding separate E2E test files.

## Architecture

### Module structure

Main module under `pkg/cadence/` follows:

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

- `pkg/common/` - shared utilities: `auth` (context helpers), `jsonlog`/`log` (structured logging), `maybe`
  (3-state optional type for PATCH semantics), `ogenerrors` (app error taxonomy: not_found, invalid_input,
  permission_denied, already_exists, operation_rejected), `ogenmiddleware` (request logging), `postgresql`
  (client/DSN/migrator), `redis` (session store client), `transactional` (Unit-of-Work with optional
  Postgres advisory-lock support), `uuid` (UUIDv7)

### CQRS split

Write and read sides are fully separate stacks, one per aggregate, never sharing models:

- Write side: `app/service` (business logic, mutates via `domain` entities) -> `infrastructure/persistence/postgresql/repo`
  (`domain.XRepository` impl, `sqlx`-style, returns/accepts domain entities, e.g. `domain.LoadBand(...)`).
- Read side: `app/query` (`XQueryService` interfaces) -> `infrastructure/persistence/postgresql/query`
  (hand-written SQL, `SelectContext`/`GetContext` straight into flat DTOs like `TrackListItem`, `UserSetlistListItem`).
  Read-side code never imports `domain` - no entities, no invariants, just projections for API responses.
- Both sides get their own `authorized` decorator (see below) and their own constructor wired independently
  in `cmd/cadence/dependencycontainer.go` (e.g. `trackService` and `trackQueryService` are separate objects,
  both backed by the same `track` table but through different SQL and different Go types).
- Adding a new read model: define the DTO + interface method in `app/query`, hand-write the SQL in
  `infrastructure/persistence/postgresql/query`, add the `authorized` pass-through/`requireMember` check,
  wire the constructor. Don't reuse `domain` types or `repo/` queries for read paths, even if the SQL would
  look similar - keeps the read side free to diverge (joins, aggregates, search) without touching write-side
  invariants.

### Authorized services pattern

`app/service/authorized` and `app/query/authorized` decorate the plain `app/service`/`app/query`
implementations with request-level authorization, keeping the inner services free of auth concerns:

- Each authorized type wraps an inner `service.XService`/`query.XQueryService` plus a
  `transactional.Executor[app.RepoProvider]`, e.g. `authorizedservice.NewBandService(appservice.NewBandService(executor), executor)`.
  Wired for every aggregate in `cmd/cadence/dependencycontainer.go`.
- `requesterIDFromContext(ctx)` pulls the caller's user ID from session auth
  (`pkg/common/auth`), returning `permission_denied` if unauthenticated.
- `requireMember(ctx, executor, bandID, requesterID)` (in `app/service/authorized/band.go` and mirrored in
  `app/query/authorized`) loads the band and checks `band.HasMember` before delegating - use this for any
  op scoped to "any band member may do this".
- Owner-only operations (e.g. `Band.TransferOwnership`, `Band.RegenerateInviteCode`) enforce `IsOwner` inside
  the domain method itself and return `domain.ErrNotBandOwner`; the authorized wrapper just passes the call
  through - don't duplicate the owner check at the authorized layer.
- Ops implicitly scoped to the caller (`ListUserTracks`, `ListUserSetlists`, `JoinByInviteCode` - keyed by
  `userID` or by invite code, not by band membership) also pass through unchecked.
- New service/query methods on band-scoped aggregates should follow this same split: put "is this requester
  even allowed near this band" in the authorized wrapper via `requireMember`, and finer-grained role checks
  (owner vs member) in the domain type.

## Databases

### Migrations

Migrations live in `data/migrations/` as timestamped SQL file pairs:

```
20240101120000_<name>.up.sql
20240101120000_<name>.down.sql
```

- Use `bin/create-db-migration <name>` to scaffold a new pair.
- Write only one query per file.
- Timestamp in the beginning should be unique, so add `sleep 1` calls when scaffolding several migrations in a row.

## Routing

- `gorilla/mux` router
- Public API: `/api/*` registered via `PathPrefix`, CORS middleware (`CADENCE_CORS_ALLOWED_ORIGINS`)
- Auth: session-token based (`SessionAuth` security scheme backed by Redis), see `pkg/common/auth`
- Health: `/resilience/live` (liveness, process up), `/resilience/ready` (readiness, checks DB/Redis)
- OpenAPI spec: `api/server/publicapi.yml`; generated client/server code lands in `api/server/publicapi/`
  (never hand-edit, regenerate via `mise run generate`)

## Constraints

- Never edit existing migration files - only scaffold new ones via `bin/create-db-migration`. Modifying an
  already-applied migration corrupts DB state in deployed environments.
- No hardcoded credentials or secrets - all configuration must flow through `CADENCE_`-prefixed environment variables.
  Nothing sensitive in source code.
- API changes must be backward-compatible - do not remove or rename existing fields; add new fields instead.
- Deletes are soft (`deleted_at`/`deleted_by` columns, filtered out of reads) - never hard-`DELETE` domain rows.
- Mutations that race on a shared aggregate (band/track update or remove) take a Postgres advisory lock via
  `transactional.Executor` - follow the existing `repo/transactionfactory.go` pattern for new mutating operations.

## Definition of Done

A task is considered done only when all the following pass without errors:

1. Full build, unit tests and lint - run:

```shell
mise run
```

2. E2E tests, if touching public API behavior, transport wiring, or persistence - run:

```shell
mise run test:e2e
```
