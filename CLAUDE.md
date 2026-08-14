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
mise run check:lint

# Unit tests
mise run check:test

# Lint + test
mise run check

# Scaffold a new migration pair
bin/create-db-migration <name>
```

`compose.yml` runs local Postgres + Redis (+ app) for manual testing; copy `compose.override.example.yml` to override.

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

## Databases

### Migrations

Migrations live in `data/migrations/` as timestamped SQL file pairs:

```
20240101120000_<name>.up.sql
20240101120000_<name>.down.sql
```

Use `bin/create-db-migration <name>` to scaffold a new pair.

## Routing

- `gorilla/mux` router
- Public API: `/api/*` registered via `PathPrefix`, CORS middleware (`CADENCE_CORS_ALLOWED_ORIGINS`)
- Auth: session-token based (`SessionAuth` security scheme backed by Redis), see `pkg/common/auth`
- Health: `/resilience/ready`
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
