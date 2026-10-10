# Go Backend Engineering Lab

A hands-on project for learning Go backend engineering by implementing
infrastructure concepts ourselves. The business domain is deliberately
small: a Product API. The learning goals are the cache, rate limiter,
worker pool, reverse proxy/load balancer, networking, concurrency,
testing, and observability around it.

> This is a learning lab, not a production replacement for Redis, NGINX,
> or a mature message broker.

## 1. Technology and dependency rules

-   **Go** for the backend.
-   **Gorilla Mux** for routing.
-   Go standard library for HTTP server primitives, networking,
    concurrency, JSON, logging, and tests.
-   PostgreSQL in the database phase; use a suitable driver rather than
    implementing the PostgreSQL wire protocol.
-   No third-party libraries for the infrastructure algorithms we are
    specifically learning to implement.

Gorilla Mux uses Go's `net/http` interfaces. We will use Gorilla Mux for
route matching and `net/http` for the server, request/response types,
and HTTP primitives. We are not implementing HTTP itself.

## 2. What we are building

A small Product REST API that gradually gains:

1.  A clean handler/service/repository structure.
2.  PostgreSQL persistence.
3.  A custom in-memory cache with TTL and expiration.
4.  A token-bucket rate limiter.
5.  A bounded Go worker pool.
6.  A Go reverse proxy/load balancer with round-robin selection and
    health checks.
7.  A simplified Redis-like TCP server supporting a small subset of
    RESP.
8.  Tests, logs, metrics, graceful shutdown, Docker, and load testing.

### Target architecture

``` text
Client (curl/Postman)
        |
        v
Our Go reverse proxy/load balancer
        |
        +-------------------+
        |                   |
        v                   v
   Product API 1       Product API 2
        |                   |
        +---------+---------+
                  |
          +-------+-------+
          |               |
          v               v
      PostgreSQL     Our cache service
      (persistent)   (custom implementation)

Job submission API --> bounded Go worker pool
```

This is the eventual architecture. We will add one component at a time.
The early cache is in-memory and volatile; a process restart clears it.

## 3. Planned repository structure

Create directories when the relevant phase begins; do not create empty
packages in advance.

``` text
go-backend-lab/
├── cmd/
│   ├── api/main.go
│   ├── worker/main.go
│   ├── loadbalancer/main.go
│   └── cache-server/main.go
├── internal/
│   ├── config/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   ├── model/
│   ├── middleware/
│   ├── cache/
│   ├── ratelimit/
│   ├── worker/
│   ├── proxy/
│   ├── protocol/
│   ├── database/
│   └── observability/
├── migrations/
├── tests/
├── scripts/
├── compose.yaml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

### Responsibilities

-   `cmd/`: process entry points and dependency wiring.
-   `handler/`: HTTP input, validation, status codes, and response
    encoding.
-   `service/`: business rules and orchestration.
-   `repository/`: storage interfaces and implementations.
-   `model/`: domain types.
-   `middleware/`: logging, recovery, request IDs, and rate limiting.
-   `cache/`: concurrency-safe key-value store and TTL behavior.
-   `ratelimit/`: token-bucket algorithm.
-   `worker/`: jobs, bounded queue, worker lifecycle, and retry policy.
-   `proxy/`: backend registry, selection, health state, and proxy
    wiring.
-   `protocol/`: parser/encoder for the small RESP subset.
-   `database/`: PostgreSQL connection setup.
-   `observability/`: metrics and health information.

## 4. API endpoints

### Health

  -----------------------------------------------------------------------
  Method                  Path                    Purpose
  ----------------------- ----------------------- -----------------------
  GET                     `/health/live`          Process is alive

  GET                     `/health/ready`         Application is ready;
                                                  later checks
                                                  dependencies
  -----------------------------------------------------------------------

### Products

  Method   Path                      Purpose
  -------- ------------------------- ------------------------
  POST     `/api/v1/products`        Create a product
  GET      `/api/v1/products`        List products
  GET      `/api/v1/products/{id}`   Fetch one product
  PUT      `/api/v1/products/{id}`   Update editable fields
  DELETE   `/api/v1/products/{id}`   Delete a product

List query parameters: - `limit`: page size with a configured maximum. -
`offset`: number of records to skip. - `q`: optional case-insensitive
product-name search. - `sort`: allow-listed field, such as `name`,
`price_cents`, or `created_at`. - `order`: `asc` or `desc`.

Never concatenate arbitrary client input into SQL.

Example request:

``` json
{
  "name": "Mechanical Keyboard",
  "description": "Compact wired keyboard",
  "price_cents": 5999,
  "currency": "INR",
  "stock_quantity": 25
}
```

Example response:

``` json
{
  "id": "generated-uuid",
  "name": "Mechanical Keyboard",
  "description": "Compact wired keyboard",
  "price_cents": 5999,
  "currency": "INR",
  "stock_quantity": 25,
  "created_at": "2026-10-10T10:00:00Z",
  "updated_at": "2026-10-10T10:00:00Z"
}
```

The server generates the ID and timestamps. `price_cents` is an integer
minor-unit amount, avoiding floating-point money calculations. This is a
practice API, not a payment system.

## 5. Database schema

We may use an in-memory repository for the first HTTP iteration.
PostgreSQL becomes the persistent source of truth in Phase 2.

### `products`

  Column             PostgreSQL type   Constraints / meaning
  ------------------ ----------------- -------------------------------------------------------
  `id`               UUID              Primary key, generated by the application or database
  `name`             VARCHAR(200)      Not null
  `description`      TEXT              Nullable
  `price_cents`      BIGINT            Not null, must be \>= 0
  `currency`         CHAR(3)           Not null, default `INR`
  `stock_quantity`   INTEGER           Not null, must be \>= 0
  `created_at`       TIMESTAMPTZ       Not null
  `updated_at`       TIMESTAMPTZ       Not null

Indexes: - Primary key index on `id`. - Index on `created_at` for common
ordering. - Add a name/search index only if query patterns justify it. -
Product names are not unique by default.

### `jobs` (added when durable jobs are introduced)

The first worker pool uses an in-memory channel. Later, use PostgreSQL
to persist jobs so queued work can survive process restarts.

  ------------------------------------------------------------------------
  Column                  PostgreSQL type         Constraints / meaning
  ----------------------- ----------------------- ------------------------
  `id`                    UUID                    Primary key

  `type`                  VARCHAR(100)            Not null;
                                                  e.g. `product_reindex`

  `payload`               JSONB                   Not null; job input

  `status`                VARCHAR(20)             Not null; `pending`,
                                                  `running`, `succeeded`,
                                                  `failed`

  `attempts`              INTEGER                 Not null, default 0

  `max_attempts`          INTEGER                 Not null, default 5

  `available_at`          TIMESTAMPTZ             Not null; earliest
                                                  eligible run/retry time

  `locked_at`             TIMESTAMPTZ             Nullable; when a worker
                                                  claimed it

  `last_error`            TEXT                    Nullable

  `created_at`            TIMESTAMPTZ             Not null

  `updated_at`            TIMESTAMPTZ             Not null
  ------------------------------------------------------------------------

Indexes: - `(status, available_at)` for finding runnable jobs. -
Optional `locked_at` index for investigating stale claims.

Multiple workers claiming jobs safely requires a transaction and
row-locking strategy; we will implement and test that before running
multiple worker processes.

### Cache entry (not a database table)

Our early cache stores entries in memory:

``` text
key -> { value, expires_at }
```

It supports `SET`, `GET`, `DELETE`, then TTL and expiration. It is
volatile initially. Persistence, replication, eviction, and full Redis
compatibility are separate advanced topics.

We will not store every HTTP access log in PostgreSQL. Logs go to
stdout/stderr; a separate audit table can be designed later if the
product needs queryable audit history.

## 6. Phases and implementation steps

### Phase 0 --- Setup

1.  Create the repository and initialize the Go module.
2.  Add Gorilla Mux.
3.  Create `cmd/api/main.go`.
4.  Register `GET /health/live`.
5.  Configure the port and start the server.
6.  Run `go test ./...`.

**Done when:** the server starts and the health endpoint works.

### Phase 1 --- Clean HTTP API

1.  Define the product model.
2.  Define a `ProductRepository` interface.
3.  Implement an in-memory repository.
4.  Implement service validation and business rules.
5.  Implement product CRUD handlers.
6.  Register routes with Gorilla Mux.
7.  Add consistent JSON responses and errors.
8.  Add logging, panic-recovery middleware, and request IDs.
9.  Propagate request contexts.
10. Test success cases, invalid input, missing products, and errors.

**Learn:** interfaces, dependency injection, pointers, errors,
middleware, context, HTTP status codes, and layering.

### Phase 2 --- PostgreSQL

1.  Create versioned SQL migrations.
2.  Connect PostgreSQL and configure connection pooling.
3.  Implement a PostgreSQL repository behind the same interface.
4.  Add constraints and indexes.
5.  Implement pagination and safe filtering/sorting.
6.  Handle database errors without leaking internals.
7.  Add transactions and integration tests.

**Learn:** SQL, transactions, connection pools, constraints, migrations,
and query plans.

### Phase 3 --- Our cache

1.  Implement a concurrency-safe map-backed cache.
2.  Add `Set(key, value, ttl)`, `Get(key)`, and `Delete(key)`.
3.  Define missing-key and expired-key behavior.
4.  Implement lazy expiration.
5.  Add cleanup for expired entries.
6.  Test expiration with controlled time.
7.  Integrate cache-aside product reads.
8.  Invalidate/update cache entries after writes.
9.  Test concurrency and stale-data cases.
10. Benchmark cached and uncached reads.

**Learn:** mutexes, race conditions, TTL, cache-aside, invalidation,
memory growth, and benchmarking.

### Phase 4 --- Rate limiter

1.  Implement a token-bucket algorithm.
2.  Configure capacity and refill rate.
3.  Protect shared limiter state.
4.  Define how client identity is determined safely.
5.  Add middleware returning `429 Too Many Requests`.
6.  Test bursts, refill, concurrency, and limits.
7.  Bound and clean up limiter state so arbitrary client IDs cannot grow
    memory forever.

**Learn:** algorithms, time calculations, synchronization, admission
control, and middleware.

### Phase 5 --- Worker pool

1.  Define job and result types.
2.  Implement a bounded job channel.
3.  Start a configurable number of goroutines.
4.  Add an endpoint to submit jobs and inspect results.
5.  Handle cancellation and graceful shutdown.
6.  Add bounded retries with exponential backoff.
7.  Define queue-full/backpressure behavior.
8.  Test worker limits, failures, cancellation, and shutdown.
9.  Later, persist jobs in PostgreSQL.

**Learn:** goroutines, channels, WaitGroups, cancellation, retries,
backpressure, and durable vs. in-memory jobs.

### Phase 6 --- Our load balancer

1.  Start multiple API instances on different ports.
2.  Create a separate load-balancer executable.
3.  Maintain a concurrency-safe backend registry.
4.  Implement round-robin selection ourselves.
5.  Forward requests with standard-library HTTP primitives.
6.  Add health checks and backend health states.
7.  Skip unhealthy backends.
8.  Handle timeouts and upstream failures.
9.  Define backend recovery behavior.
10. Test distribution, failed backends, concurrency, and shutdown.

**Learn:** reverse proxying, networking, synchronization, health states,
and failure handling.

We implement the balancing algorithm and health logic ourselves; we do
not reimplement HTTP parsing and forwarding.

### Phase 7 --- Simplified Redis-like TCP server

1.  Listen for TCP connections using `net`.
2.  Accept multiple clients.
3.  Parse a deliberately small RESP subset.
4.  Support `PING`, `SET`, `GET`, and `DEL`.
5.  Add expiration, `EXPIRE`, and `TTL`.
6.  Handle malformed commands and disconnects.
7.  Make shared storage concurrency-safe.
8.  Write a small client for our protocol.
9.  Test simultaneous clients.
10. Optionally test supported commands with `redis-cli` if compatibility
    is sufficient.

This is a learning server, not a complete Redis replacement. It will not
initially implement persistence, replication, transactions, eviction
policies, or the full command set.

### Phase 8 --- Observability and resilience

1.  Add structured logs with `log/slog`.
2.  Log request IDs and request durations.
3.  Implement liveness and readiness checks.
4.  Record request counts, errors, and latency.
5.  Add outbound timeouts.
6.  Add graceful shutdown for APIs, proxy, cache, and workers.
7.  Test dependency failures.
8.  Run the race detector and load tests.

### Phase 9 --- Docker and repeatable execution

1.  Add Dockerfiles for the required processes.
2.  Add Docker Compose for PostgreSQL and implemented services.
3.  Configure ports and environment variables.
4.  Add health checks and startup handling.
5.  Add build and test checks in CI.
6.  Document local development and reset commands.
7.  Run multiple API instances behind our load balancer.

### Phase 10 --- Advanced experiments

-   Compare round robin with least-connections and weighted round robin.
-   Compare token bucket with fixed-window and sliding-window rate
    limiting.
-   Add an LRU cache eviction policy.
-   Explore snapshots/persistence for the cache.
-   Build a bounded connection-pool exercise.
-   Add circuit breaking and retry budgets.
-   Compare in-memory queues with durable jobs.
-   Benchmark throughput, latency percentiles, memory, and lock
    contention.
-   Profile with `pprof`.
-   Explore distributed coordination and locking.

## 7. Engineering rules

-   Use `context.Context` for cancellation and bounded work.
-   Pass dependencies explicitly; avoid hidden global state.
-   Protect shared mutable state.
-   Keep handlers thin and business logic in services.
-   Avoid holding locks during slow network or database operations.
-   Bound queues, retries, and resource usage.
-   Do not retry non-idempotent operations blindly.
-   Never log passwords, tokens, secrets, or sensitive payloads.
-   Test failure paths and concurrency, not only happy paths.
-   Run `go test ./...` regularly and `go test -race ./...` during
    concurrency work.

## 8. Definition of done

We should be able to demonstrate and explain:

-   A layered Go REST API using Gorilla Mux and PostgreSQL.
-   A concurrency-safe cache with TTL.
-   A working rate limiter.
-   A bounded worker pool with cancellation and retries.
-   A reverse proxy/load balancer with health-aware selection.
-   A simplified TCP/RESP server.
-   Automated tests, race checks, benchmarks, logs, and health checks.
-   Multiple API instances running behind our own load balancer.
-   The trade-offs, limitations, and failure modes of each component.

The outcome is not simply a list of features. It is the ability to
explain how each component works, why it exists, what can fail, and how
we verified it.

## 9. Initial commands

``` bash
mkdir go-backend-lab
cd go-backend-lab

go mod init go-backend-lab
go get github.com/gorilla/mux

mkdir -p cmd/api
mkdir -p internal/handler
mkdir -p internal/service
mkdir -p internal/repository
mkdir -p internal/model
mkdir -p internal/middleware
```

Start with Phase 0 only. Create a minimal `cmd/api/main.go`, register
`GET /health/live` with Gorilla Mux, and start the HTTP server. Do not
create every future package or begin Redis/load balancing before the
first API is working and tested.
