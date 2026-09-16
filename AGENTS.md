# Repository Guidelines

## Overview

Cadence is a Go service for managing band repertoires.

## Project Structure

The executable and dependency wiring live in `cmd/cadence/`. Core code is under `pkg/cadence/`: `domain/` contains
entities and invariants, `app/service/` handles writes, `app/query/`handles read models, and `infrastructure/` provides
HTTP transport and persistence adapters. Shared packages live in `pkg/common/`.

The public OpenAPI contract is `api/server/publicapi.yml`; generated code in `api/server/publicapi/` must not be edited
manually. SQL migrations are embedded from `data/migrations/`. Unit tests are colocated with source; end-to-end tests
are in `test/e2e/`.

## Build, Test, and Development

Use mise for all project commands; it selects the required Go, ogen, and golangci-lint versions.

```shell
mise run                 # generate API code, build, lint, and run unit tests

mise run generate        # regenerate OpenAPI server code
mise run build           # build ./bin/cadence
mise run lint            # run golangci-lint and format checks
mise run test:unit       # run unit tests
mise run test:e2e        # run Docker-backed E2E tests


mise run dev             # build and start Docker Compose, waiting for health
mise run dev:reload      # rebuild and restart the local service

bin/create-db-migration <name>  # create an up/down migration pair
```

For local services, copy `compose.override.example.yml` to `compose.override.yml`, adjust its volume path, then run
`docker compose up -d`.

## Coding Style and Architecture

Format Go code with tabs, `gofmt`, `goimports`, and `gci`; `mise run lint` enforces all four. Follow Go naming: exported
identifiers use `PascalCase`, unexported identifiers use `camelCase`; repository and service types end in `Repository`
and `Service`.

Keep write and read paths separate: mutations go through `app/service/`, domain entities, and
`persistence/postgresql/repo/`; projections go through `app/query/` and `persistence/postgresql/query/`. Domain code
must not import infrastructure. Place public declarations before private ones, with constants, variables, and
constructors at the top.

## Testing and Database Changes

Use testify for unit tests and name files `*_test.go`. Add focused tests beside changed domain or service code. Run
`mise run test:e2e` when changing API behavior, transport wiring, or persistence; it requires a running Docker daemon
and uses testcontainers.

Never modify an existing migration. Create timestamped `.up.sql` and `.down.sql` pairs with the helper, use one query
per file, and add new migrations only. Keep credentials out of source; configuration uses `CADENCE_`-prefixed
environment variables.

## Commits and Pull Requests

Recent history uses concise, imperative subjects such as `Adjust mise tasks` and `Fix compile flags`. Keep each commit
scoped to one change. PRs should explain the behavior change, link relevant issues, note migration or configuration
effects, include API examples or screenshots where applicable, and confirm `mise run` (plus relevant E2E tests) passes.
