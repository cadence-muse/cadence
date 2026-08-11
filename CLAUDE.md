# Repository Guidelines

## Overview

Cadence is a Go modular music band repertoire management system. It exposes OpenAPI public API. The service uses
PostgreSQL 16 for persistence.

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

# Scaffold a new migration pair
bin/create-db-migration <name>
```

## Databases

### Migrations

Migrations live in `data/migrations/` as timestamped SQL file pairs:

```
20240101120000_<name>.up.sql
20240101120000_<name>.down.sql
```

Use `bin/create-db-migration <name>` to scaffold a new pair.

## Routing

- `gorilla/mux` router, modules register via `router.PathPrefix("/api")`
- Public API: `/api/*`
- Health: `/resilience/ready`

## Constraints

- Never edit existing migration files - only scaffold new ones via `bin/create-db-migration`. Modifying an
  already-applied migration corrupts DB state in deployed environments.
- No hardcoded credentials or secrets - all configuration must flow through `CADENCE_`-prefixed environment variables.
  Nothing sensitive in source code.
- API changes must be backward-compatible - do not remove or rename existing fields; add new fields instead.

## Definition of Done

A task is considered done only when all the following pass without errors:

1. Full build, unit tests and lint - run:

```shell
mise run
```
