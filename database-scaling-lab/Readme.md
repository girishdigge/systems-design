# Database Scaling Laboratory (Phase 0)

A high-performance backend engineering laboratory designed to test database scalability, connection pooling optimizations, indexing strategies, and load testing under high-concurrency transactional (OLTP) and analytical (OLAP) workloads.

---

## 🏗️ Architecture Overview

The system is designed as a modular, containerized multi-service infrastructure environment optimizing Go backend performance against a heavy PostgreSQL data layer.

+-----------------------+
| k6 Load Tester |
| (5,300+ req/sec) |
+-----------+-----------+
| (HTTP via Port 8080)
v
+-----------------------+
| Go REST API Container|
| (pgxpool Connection) |
+-----------+-----------+
|
+--------------------+--------------------+
| (Internal Docker Network Routing) |
v v
+------------------+ +------------------+
| postgres-primary | <------------------- | migrate/migrate |
| (8.1M+ Rows) | (Schema Up) | (Schema Tool) |
+------------------+ +------------------+

### Infrastructure Stack

- **Language Runtime:** Go 1.22+ (utilizing the native high-performance `pgx/v5` driver)
- **Database Engine:** PostgreSQL 18 (OLTP/OLAP target layer)
- **Migration Layer:** `golang-migrate/migrate` (Docker runtime)
- **Load Automation:** Grafana k6 (JavaScript-driven high-concurrency engine)
- **Observability Stack:** Prometheus (Metrics Scraper) & Grafana (Performance Dashboards)

---

## 📊 Database Schema Matrix

The baseline layout models a real-world SaaS infrastructure consisting of over **8.1 Million rows** partitioned conceptually into active business operations and heavy analytical event logging:

| Table Name        | Column Layout                               | Primary Key / Indexing                            | Estimated Target Volume |
| :---------------- | :------------------------------------------ | :------------------------------------------------ | :---------------------- |
| **`users`**       | `id`, `name`, `email`, `created_at`         | `id` (UUID B-Tree PK), `email` (UNIQUE)           | 100,000 Rows            |
| **`products`**    | `id`, `name`, `price`                       | `id` (UUID B-Tree PK)                             | 10,000 Rows             |
| **`orders`**      | `id`, `user_id`, `total`, `created_at`      | `id` (UUID B-Tree PK), `idx_orders_user_id`       | 1,000,000 Rows          |
| **`order_items`** | `id`, `order_id`, `product_id`, `quantity`  | `id` (UUID B-Tree PK), `idx_order_items_order_id` | ~2,000,000 Rows         |
| **`events`**      | `id`, `user_id`, `event_type`, `created_at` | `id` (BIGSERIAL PK), `idx_events_created_at`      | 5,000,000 Rows          |

---

## 🛠️ Operational Guide

### 1. Provisioning Infrastructure

Spin up the decoupled operational services inside the isolated virtual bridge network:

```bash
docker compose up -d
```

2. Executing Migrations
   Apply structural database layout schemas cleanly using the internal network host mapping:

docker run --rm -v $(pwd)/migrations:/migrations --network host migrate/migrate \
 -path=/migrations -database "postgres://postgres:supersecretpassword@localhost:5432/scaling_lab?sslmode=disable" up

3. Compiling & Executing Mass Seeder
   To bypass host-layer 127.0.0.1 network socket collisions (caused by ghost background host Postgres services), compile the seeder to target Linux and execute it safely inside the isolated Docker network:

# Compile binary target

GOOS=linux GOARCH=amd64 go build -o seeder-runner cmd/seeder/main.go

# Stream 8.1M+ rows via pgx.CopyFrom Bulk Insertion

docker run --rm \
 -v $(pwd)/seeder-runner:/seeder-runner \
 -v $(pwd)/.env:/.env \
 --network database-scaling-lab_lab_network \
 ubuntu:latest /seeder-runner

4. Running the API Web Server
   Compile and launch the Go REST API container exposed on port 8080:

GOOS=linux GOARCH=amd64 go build -o api-runner cmd/api/main.go

docker run --rm -d \
 -p 8080:8080 \
 -v $(pwd)/api-runner:/api-runner \
 -v $(pwd)/.env:/.env \
 --name scaling-api \
 --network database-scaling-lab_lab_network \
 ubuntu:latest /api-runner

⚡ Benchmarking & Load Testing Metrics (Phase 0 Baseline)
A heavy 3-minute mixed read/write scenario was simulated using k6 to establish the performance baseline of the unoptimized infrastructure.

Execution Command
k6 run k6/mixed-load.js

Verified Baseline Performance Metrics
scenarios: 100 max VUs looping for 3 minutes over 3 stages

checks_total.......: 900328 5383.55/s
checks_succeeded...: 100.00% (900328 out of 900328)
checks_failed......: 0.00%

HTTP Metrics:
http_req_duration: avg=2.03ms med=1.33ms p(95)=4.82ms
http_req_failed..: 0.00%
http_reqs........: 900328 5383.55/s

Engineering InsightsB-Tree Efficiency: Lookups on /users/:id executed in $O(\log N)$ via the primary key index, maintaining sub-5ms $p(95)$ latencies despite the database housing millions of rows.Buffer Cache Impact: High request volumes achieved a high cache-hit ratio due to repeated hits on the static benchmark target user UUID string, minimizing physical disk read cycles.
