# Database Scaling Laboratory (Phase 0)

A high-performance backend engineering laboratory designed to evaluate database scalability, connection pooling efficiency, indexing strategies, and load-testing behavior under high-concurrency OLTP and analytical workloads.

---

# 🏗️ Architecture Overview

The platform is composed of isolated Docker containers communicating over an internal bridge network. The API uses PostgreSQL through `pgxpool`, while migrations and data seeding are executed as standalone containers.

```text
                           +----------------------+
                           |      k6 Load Test    |
                           |   ~5,300 req/sec    |
                           +----------+-----------+
                                      |
                                      | HTTP :8080
                                      v
                    +----------------------------------+
                    |        Go REST API              |
                    |      (pgxpool Connection Pool)  |
                    +----------------+----------------+
                                     |
                                     |
                     Docker Bridge Network
                                     |
          +--------------------------+--------------------------+
          |                                                     |
          v                                                     v
+----------------------+                       +----------------------+
|  postgres-primary    | <-------------------- | migrate/migrate      |
|                      |     Schema Updates    | Migration Container  |
|  8.1M+ Seeded Rows   |                       +----------------------+
+----------+-----------+
           ^
           |
           | Bulk Inserts (pgx.CopyFrom)
           |
+----------------------+
|    Seeder Runner     |
|  Data Generation     |
+----------------------+
```

## Infrastructure Stack

| Component         | Technology             |
| ----------------- | ---------------------- |
| Language Runtime  | Go 1.22+               |
| Database          | PostgreSQL             |
| Driver            | pgx/v5                 |
| Connection Pool   | pgxpool                |
| Migrations        | golang-migrate/migrate |
| Load Testing      | Grafana k6             |
| Metrics           | Prometheus             |
| Dashboards        | Grafana                |
| Container Runtime | Docker                 |

---

# 📊 Database Schema Matrix

The baseline dataset models a realistic SaaS workload consisting of approximately **8.1 million rows**.

| Table         | Columns                                     | Primary Key / Indexes               | Target Volume |
| ------------- | ------------------------------------------- | ----------------------------------- | ------------- |
| `users`       | `id`, `name`, `email`, `created_at`         | UUID PK, UNIQUE(email)              | 100,000       |
| `products`    | `id`, `name`, `price`                       | UUID PK                             | 10,000        |
| `orders`      | `id`, `user_id`, `total`, `created_at`      | UUID PK, idx_orders_user_id         | 1,000,000     |
| `order_items` | `id`, `order_id`, `product_id`, `quantity`  | UUID PK, idx_order_items_order_id   | 2,000,000+    |
| `events`      | `id`, `user_id`, `event_type`, `created_at` | BIGSERIAL PK, idx_events_created_at | 5,000,000     |

---

# 🛠️ Operational Guide

## 1. Start Infrastructure

Launch PostgreSQL, Prometheus, Grafana, and supporting services.

```bash
docker compose up -d
```

---

## 2. Apply Database Migrations

Execute migrations against the running PostgreSQL instance.

```bash
docker run --rm \
  -v $(pwd)/migrations:/migrations \
  --network host \
  migrate/migrate \
  -path=/migrations \
  -database "postgres://postgres:supersecretpassword@localhost:5432/scaling_lab?sslmode=disable" \
  up
```

---

## 3. Build and Run the Seeder

Compile the seeder for Linux and execute it inside the Docker network.

### Build

```bash
GOOS=linux GOARCH=amd64 go build -o seeder-runner cmd/seeder/main.go
```

### Seed Database

```bash
docker run --rm \
  -v $(pwd)/seeder-runner:/seeder-runner \
  -v $(pwd)/.env:/.env \
  --network database-scaling-lab_lab_network \
  ubuntu:latest \
  /seeder-runner
```

The seeder streams approximately **8.1 million rows** using PostgreSQL's high-performance `pgx.CopyFrom` bulk insertion mechanism.

---

## 4. Build and Run the API

### Build

```bash
GOOS=linux GOARCH=amd64 go build -o api-runner cmd/api/main.go
```

### Run

```bash
docker run --rm -d \
  -p 8080:8080 \
  -v $(pwd)/api-runner:/api-runner \
  -v $(pwd)/.env:/.env \
  --name scaling-api \
  --network database-scaling-lab_lab_network \
  ubuntu:latest \
  /api-runner
```

Verify the API:

```bash
curl http://localhost:8080/health
```

---

# ⚡ Benchmarking & Load Testing

A mixed read/write workload was executed using Grafana k6 for three minutes to establish a baseline before introducing optimization techniques.

## Run Benchmark

```bash
k6 run k6/mixed-load.js
```

## Test Configuration

```text
Duration: 3 minutes
Virtual Users: 100
Workload: Mixed Read / Write
Target: Go REST API + PostgreSQL
```

## Baseline Results

```text
checks_total.......: 900,328 (5383.55/s)
checks_succeeded...: 100.00%
checks_failed......: 0.00%

http_req_duration:
  avg=2.03ms
  med=1.33ms
  p(95)=4.82ms

http_req_failed....: 0.00%
http_reqs..........: 900,328 (5383.55/s)
```

---

# 🔍 Engineering Insights

### B-Tree Index Efficiency

Queries against `/users/:id` execute through the primary-key B-tree index, maintaining logarithmic lookup complexity and sub-5ms p95 latency even with millions of records.

### Buffer Cache Effectiveness

The benchmark repeatedly targets a fixed set of records, resulting in a high PostgreSQL buffer-cache hit ratio and minimizing physical disk I/O.

### Connection Pool Performance

`pgxpool` efficiently multiplexes database connections, allowing the API to sustain more than 5,000 requests per second without excessive connection creation overhead.

---

# 🚀 Next Phases

- Read replica deployment
- Query plan analysis (`EXPLAIN ANALYZE`)
- Partitioning strategies
- Materialized views
- PgBouncer integration
- Horizontal API scaling
- Advanced Prometheus metrics
- Grafana performance dashboards
