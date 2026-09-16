# Cadence API

The backend service for [Cadence](https://github.com/cadence-muse), a shared repertoire manager for musical bands.
It exposes the public HTTP API, applies domain rules, stores application data in PostgreSQL, and keeps sessions in Redis.

## 🎯 Responsibilities

- User registration, login, profiles, passwords, and sessions
- Band membership, invitations, ownership, and member management
- Track catalogs with search and musical metadata
- Setlist creation, editing, track ordering, and search
- Readiness and liveness endpoints for deployment health checks

## 🏗️ Architecture

The service uses a layered architecture with separate write and read paths:

```text
cmd/cadence                       executable, configuration, and dependency wiring
pkg/cadence/domain                entities and business invariants
pkg/cadence/app/service           commands and write workflows
pkg/cadence/app/query             read models and queries
pkg/cadence/infrastructure        HTTP and persistence adapters
pkg/common                        shared authentication, metrics, and middleware
api/server/publicapi.yml          public OpenAPI contract
data/migrations                   embedded PostgreSQL migrations
```

HTTP requests enter through the generated ogen transport. Commands pass through application services and domain entities
before reaching PostgreSQL repositories; reads use dedicated query implementations. The domain layer does not depend on infrastructure.

## 🛠️ Local development

### Prerequisites

- [mise](https://mise.jdx.dev/)
- Docker with the Compose plugin

Clone the repository and create the local Compose override. Set its bind-mount source to the absolute path of the binary in your checkout.

```shell
git clone https://github.com/cadence-muse/cadence.git
cd cadence
cp compose.override.example.yml compose.override.yml
$EDITOR compose.override.yml
mise run dev
```

The API listens on `http://localhost:8080`; PostgreSQL and Redis listen on their standard ports on localhost. Common commands are:

```shell
mise run                 # generate, build, lint, and run unit tests

mise run generate        # regenerate the ogen server from OpenAPI
mise run build           # build bin/cadence
mise run lint            # run formatting and static checks

mise run test:unit       # run unit tests
mise run test:e2e        # run Docker-backed end-to-end tests

mise run dev:reload      # rebuild and restart the API container
mise run dev:logs        # follow local service logs
mise run dev:down        # stop the local stack
```

## 🔌 Working with the API

[`api/server/publicapi.yml`](api/server/publicapi.yml) is the source of truth. It covers users, bands, tracks, and setlists under `/api`.
Registration and login are public; other operations expect the token returned by login in the `Authorization` header.

```shell
curl -X POST http://localhost:8080/api/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"demo","password":"change-me"}'

curl http://localhost:8080/api/band/list \
  -H 'Authorization: TOKEN_FROM_LOGIN'
```

Change the YAML contract first when adding or modifying an endpoint, then run `mise run generate`.
Generated files under `api/server/publicapi/` must not be edited by hand.

## 🗃️ Migrations

Migrations are embedded into the service and applied by running the binary in migration mode.
Local Compose and the production platform use the same migration set.

```shell
bin/create-db-migration add_track_field
```

The helper creates a timestamped `.up.sql` and `.down.sql` pair in `data/migrations/`. Add new pairs
instead of changing migrations that may already have run, and keep one SQL statement in each file.

## ⚙️ Configuration

Configuration is read from environment variables prefixed with `CADENCE_`.

| Variable | Purpose | Default |
| --- | --- | --- |
| `CADENCE_SERVE_REST_ADDRESS` | HTTP listen address | `:8080` |
| `CADENCE_DB_HOST`, `CADENCE_DB_PORT` | PostgreSQL address | Required |
| `CADENCE_DB_NAME`, `CADENCE_DB_USER`, `CADENCE_DB_PASSWORD` | PostgreSQL credentials | Required |
| `CADENCE_DB_MAX_CONN` | Maximum database connections | `10` |
| `CADENCE_DB_CONN_LIFETIME` | Connection lifetime in seconds | `60` |
| `CADENCE_REDIS_HOST`, `CADENCE_REDIS_PORT` | Redis address | Required |
| `CADENCE_REDIS_PASSWORD`, `CADENCE_REDIS_DB` | Redis credentials and database | Required |
| `CADENCE_SESSION_MAX_PER_USER` | Concurrent sessions per user | `5` |
| `CADENCE_SESSION_TTL` | Session lifetime | `24h` |
| `CADENCE_CORS_ALLOWED_ORIGINS` | Allowed browser origins | Required for web clients |

See [`compose.yml`](compose.yml) for a complete local configuration.

## 🧪 Testing

Unit tests live beside the code and run with `mise run test:unit`.
End-to-end tests live in `test/e2e/`, use testcontainers, and require Docker.
Run `mise run test:e2e` after API, transport, or persistence changes.

## 📜 License

Distributed under the [MIT License](LICENSE).
